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

type AfterConnectOperation struct {
	operation  string // shell, filemgr
	muxEnabled bool   // mux or out-of-band
}

func buildConf() *UserConfig {
	usr, _ := user.Current()
	var globalConf = &UserConfig{
		certPath:       usr.HomeDir + certFC,
		privateKeyPath: usr.HomeDir + certPK,
		certPinPath:    usr.HomeDir + certPin,
	}
	var err2 error
	_, err3 := os.Stat(globalConf.privateKeyPath)
	if _, err := os.Stat(globalConf.certPath); err != nil || err3 != nil {
		// file not exists, call generate
		_, err2 = generateCertificate(*globalConf)
		if err2 != nil {
			panic("Errors in Generate Certificate")
		}
	}
	return globalConf
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
