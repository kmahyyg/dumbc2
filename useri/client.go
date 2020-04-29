package useri

import (
	"bytes"
	"encoding/base64"
	"github.com/kmahyyg/dumbc2/buildtime"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/transport"
	"log"
)

func StartAgent(userOP *config.UserOperation) {
	servFGP := make([]byte, base64.StdEncoding.DecodedLen(len(*buildtime.RemoteFingerprint)))
	_, err := base64.StdEncoding.Decode(servFGP, *buildtime.RemoteFingerprint)
	if err != nil {
		log.Fatalln("Server Pinned Key Error.")
	}
	pinnedDialer := transport.TLSDialerBuilder(servFGP, buildtime.ClientCertificateKey, buildtime.ClientCertificatePEM)
	curConn, err := pinnedDialer("tcp", userOP.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	var buf bytes.Buffer

}
