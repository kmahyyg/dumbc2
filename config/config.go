package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"github.com/kmahyyg/dumbc2/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"os"
	"os/user"
	"time"
)

const certFC string = "/.dumbyc2/certs.pem"
const certPK string = "/.dumbyc2/privkey.pem"
const certPin string = "/.dumbyc2/certpin.txt"

type UserConfig struct {
	certPath       string
	privateKeyPath string
	certPinPath    string
}

type SSCertificate struct {
	certData        []byte
	privKeyData     []byte
	certFingerprint []byte
}

type parsedSSCert struct {
	certkey  rsa.PrivateKey
	pubk     rsa.PublicKey
	selfcert x509.Certificate
	sslpin   []byte
}

func (ssc *SSCertificate) parse() (*parsedSSCert, error) {
	//todo: read to proper format and type object
}

type UserOperation struct {
	operation        string // generate or connect or serve
	agentOptConnType string // reverse or bind
	agentOptPath     string
	connEndp         string                 // endpoint ip address
	connEndpPort     int                    // endpoint port
	connCert         *SSCertificate         // must not be nil
	afterCo          *AfterConnectOperation // can be nil
}

type AfterConnectOperation struct {
	operation  string // shell, filemgr
	muxEnabled bool   // mux or out-of-band
}

func ParseCLIInput() *UserOperation {
	_, certData := buildConf()
	//todo: receive user input from command line and do mapping.
}

func buildConf() (*UserConfig, *SSCertificate) {
	usr, _ := user.Current()
	var globalConf = &UserConfig{
		certPath:       usr.HomeDir + certFC,
		privateKeyPath: usr.HomeDir + certPK,
		certPinPath:    usr.HomeDir + certPin,
	}
	var certData *SSCertificate
	var err2 error
	_, err3 := os.Stat(globalConf.privateKeyPath)
	if _, err := os.Stat(globalConf.certPath); err != nil || err3 != nil {
		// file not exists, call generate
		certData, err2 = generateCertificate(*globalConf)
		if err2 != nil {
			panic("Errors in Generate Certificate")
		} else {
		}
	} else {
		certData, err = readCertificate(*globalConf)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Reading Cert And Key from PEM File"))
		}
	}
	return globalConf, certData
}

func readCertificate(config UserConfig) (*SSCertificate, error) {
	rawcert, err := ioutil.ReadFile(config.certPath)
	if err != nil {
		panic(err)
	}
	rawcertpin, err := ioutil.ReadFile(config.certPinPath)
	if err != nil {
		panic(err)
	}
	rawcertpk, err := ioutil.ReadFile(config.privateKeyPath)
	if err != nil {
		panic(err)
	}
	return &SSCertificate{
		certData:        rawcert,
		privKeyData:     rawcertpk,
		certFingerprint: rawcertpin,
	}, nil
}

func generateCertificate(conf UserConfig) (*SSCertificate, error) {
	// build return result
	var ssCert = &SSCertificate{
		certData:        nil,
		privKeyData:     nil,
		certFingerprint: nil,
	}
	bits := 4096
	// generate rsa key pairs
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return ssCert, errors.Wrap(err, "Generate PRIVKEY Error")
	}
	// generate random things
	serialNo, _ := rand.Int(rand.Reader, big.NewInt(1<<64))
	randDomain := utils.RandString(8) + ".com"
	// generate self-signed certificate
	ssSignedCert := x509.Certificate{
		SerialNumber: serialNo,
		Subject: pkix.Name{
			CommonName: "dumbyc2." + randDomain,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	}
	derCert, err := x509.CreateCertificate(rand.Reader, &ssSignedCert, &ssSignedCert, &privKey.PublicKey, privKey)
	if err != nil {
		return ssCert, errors.Wrap(err, "Create Certificate Error")
	}
	// pem encode and write out
	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derCert,
	})
	if err != nil {
		return ssCert, errors.Wrap(err, "PEM Encoding Error")
	}
	ssCert.certData = buf.Bytes()
	// write rsa private key and write out
	buf = &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	if err != nil {
		return ssCert, errors.Wrap(err, "PEM PRIVKEY Encoding Error")
	}
	ssCert.privKeyData = buf.Bytes()
	// load generated cert
	parsedCert, err := x509.ParseCertificate(derCert)
	if err != nil {
		return ssCert, errors.Wrap(err, "Parsing Cert Error!")
	}

	// from loaded cert, get public key and ssl pin
	pubDer := x509.MarshalPKCS1PublicKey(parsedCert.PublicKey.(*rsa.PublicKey))
	pubsum := sha256.Sum256(pubDer)
	pubPin := make([]byte, base64.StdEncoding.EncodedLen(len(pubsum)))
	base64.StdEncoding.Encode(pubPin, pubsum[:])
	ssCert.certFingerprint = pubPin

	err = ioutil.WriteFile(conf.privateKeyPath, ssCert.privKeyData, 0600)
	if err != nil {
		panic(errors.Wrap(err, "Errors in Write Generated Private Key"))
	}
	err = ioutil.WriteFile(conf.certPath, ssCert.certData, 0755)
	if err != nil {
		panic(errors.Wrap(err, "Errors in Write Generated Cert"))
	}
	err = ioutil.WriteFile(conf.certPinPath, ssCert.certFingerprint, 0644)
	if err != nil {
		panic(errors.Wrap(err, "Errors in Write Certificate Pin"))
	}

	return ssCert, nil
}
