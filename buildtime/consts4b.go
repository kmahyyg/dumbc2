package buildtime

import (
	_ "github.com/kmahyyg/dumbc2/statik"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
)

var (
	ClientCertificatePEM []byte
	ClientCertificateKey []byte
	RemoteFingerprint    []byte
	CACertificate        []byte
)

func GetCertificates() error {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatalln(err)
	}
	CertFD, err := statikFS.Open("/clientcert.pem")
	if err != nil {
		log.Fatalln(err)
	}
	ClientCertificatePEM = bufferAndRead(CertFD)
	CertKey, err := statikFS.Open("/clientpk.pem")
	if err != nil {
		log.Fatalln(err)
	}
	ClientCertificateKey = bufferAndRead(CertKey)
	SPKFd, err := statikFS.Open("/serverpin.txt")
	if err != nil {
		log.Fatalln(err)
	}
	RemoteFingerprint = bufferAndRead(SPKFd)
	CACertFD, err := statikFS.Open("/cacert.pem")
	if err != nil {
		log.Fatalln(err)
	}
	CACertificate = bufferAndRead(CACertFD)
	return nil
}

func bufferAndRead(fd http.File) []byte {
	fdlen1, err := fd.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	buf := make([]byte, fdlen1.Size())
	_, err = fd.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	return buf
}
