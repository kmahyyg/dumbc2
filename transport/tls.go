package transport

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"github.com/kmahyyg/dumbc2/config"
	"log"
	"net"
	"time"
)

type TLSPinnedDialer func(network, addr string)(net.Conn, error)

func TLSDialerBuilder(pinnedFGP []byte) TLSPinnedDialer {
	return func(network, addr string) (net.Conn,error) {
		conn, err := tls.Dial(network, addr, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			return conn, err
		}
		connState := conn.ConnectionState()
		for _, peerCert := range connState.PeerCertificates {
			der, err := x509.MarshalPKIXPublicKey(peerCert.PublicKey)
			hash := sha256.Sum256(der)
			if err != nil {
				log.Fatalln(err)
			}
			if bytes.Compare(hash[0:], pinnedFGP) == 0 {
				log.Println("Certificate Checked!")
			} else {
				log.Fatalln("Certificate check error.")
			}
		}
		return conn, nil
	}
}

// TLSServerBuilder: Just give the server listen addr, we do next.
// You do need to check the client certificate if you use Bind Shell.
func TLSServerBuilder(laddr string) (net.Listener, error){
	var certLoca = config.GlobalConf
	cert, err := tls.LoadX509KeyPair(certLoca.CertPath, certLoca.PrivateKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	tlsConf := tls.Config{Certificates: []tls.Certificate{cert}}
	now := time.Now()
	tlsConf.Time = func() time.Time {
		return now
	}
	tlsConf.Rand = rand.Reader
	curLis, err := tls.Listen("tcp", laddr, &tlsConf)
	if err != nil {
		log.Fatalln(curLis)
	}
	return curLis, err
	// the current bind listener will not verify client (which is control side here)
	// you need to accept and check certificates.
}
