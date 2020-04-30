package transport

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"github.com/kmahyyg/dumbc2/config"
	"io/ioutil"
	"log"
	"net"
	"time"
)

func TLSDialer(pinnedFGP []byte, clientCert []byte, clientKey []byte, addr string) (net.Conn, error) {
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		panic(err)
	}
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
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
		_ = conn.SetDeadline(time.Now().Add(time.Minute * 10))
		return conn, nil
}

// TLSServerBuilder: Just give the server listen addr, we do next.
// You do need to check the client certificate if you use Bind Shell.
func TLSServerBuilder(laddr string, verifyClient bool) (net.Listener, error) {
	var certLoca = config.GlobalCert
	cert, err := tls.LoadX509KeyPair(certLoca.ClientPath, certLoca.ClientPrivateKeyPath)
	if err != nil {
		log.Fatalln(err)
	}

	var tlsConf *tls.Config
	if !verifyClient {
		tlsConf = &tls.Config{Certificates: []tls.Certificate{cert}}
	} else {
		caAuthCert, err := ioutil.ReadFile(certLoca.CAPath)
		if err != nil {
			log.Fatalln(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caAuthCert)
		tlsConf = &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
		}
	}

	now := time.Now()
	tlsConf.Time = func() time.Time {
		return now
	}
	tlsConf.Rand = rand.Reader
	curLis, err := tls.Listen("tcp", laddr, tlsConf)
	if err != nil || curLis == nil {
		log.Fatalln(curLis)
	}
	return curLis, err
	// the current bind listener will not verify client (which is control side here)
	// you need to accept and check certificates.
}
