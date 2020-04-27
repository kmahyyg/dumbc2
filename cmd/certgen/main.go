package main

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
	"github.com/kmahyyg/dumbc2/config"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"time"
)

func main() {
	
}


func generateCertificate(conf config.UserConfig) ([]*config.SSCertificate, error) {
	// build certificate
	var ssCACert = &config.SSCertificate{
		CertData:        nil,
		PrivKeyData:     nil,
		CertFingerprint: nil,
	}
	var ssCert = &config.SSCertificate{
		CertData:        nil,
		PrivKeyData:     nil,
		CertFingerprint: nil,
	}
	bits := 4096
	// generate rsa key pairs
	cAprivKey, err := rsa.GenerateKey(rand.Reader, bits)
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return []*config.SSCertificate{ssCACert, ssCert}, errors.Wrap(err, "Generate PRIVKEY Error")
	}
	// generate random things
	serialNo, _ := rand.Int(rand.Reader, big.NewInt(1<<64))
	randDomain := utils.RandString(8) + ".com"
	// generate self-signed certificate
	ssCASignedCert := &x509.Certificate{
		SerialNumber: serialNo,
		Subject: pkix.Name{
			CommonName: "dumbyc2." + randDomain,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageCRLSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageCodeSigning,
		},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	ssSignedCert := &x509.Certificate{
		SerialNumber: serialNo,
		Subject: pkix.Name{
			CommonName: "dumbyc2." + randDomain,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	derCACert, err := x509.CreateCertificate(rand.Reader, ssCASignedCert, ssCASignedCert, &cAprivKey.PublicKey, cAprivKey)
	derCert, err := x509.CreateCertificate(rand.Reader, ssSignedCert, ssCASignedCert, &privKey.PublicKey, privKey)

	if err != nil {
		return []*config.SSCertificate{ssCACert, ssCert}, errors.Wrap(err, "Create Certificate Error")
	}

	err = writeCert(derCACert, &conf, cAprivKey, ssCACert, true)
	if err != nil {
		panic(err)
	}
	err = writeCert(derCert, &conf, privKey, ssCert, false)
	if err != nil {
		panic(err)
	}

	return []*config.SSCertificate{ssCACert, ssCert}, nil
}

func writeCert(cert []byte, conf *config.UserConfig, certKey *rsa.PrivateKey, ssCert *config.SSCertificate, isCA bool) error {
	// pem encode and write out
	buf := &bytes.Buffer{}
	err := pem.Encode(buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	if err != nil {
		return errors.Wrap(err, "PEM Cert Encoding Error")
	}
	ssCert.CertData = buf.Bytes()
	// write rsa private key and write out
	buf = &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certKey),
	})
	if err != nil {
		return errors.Wrap(err, "PEM PRIVKEY Encoding Error")
	}
	ssCert.PrivKeyData = buf.Bytes()
	// load generated cert
	parsedCert, err := x509.ParseCertificate(cert)
	if err != nil {
		return errors.Wrap(err, "Parsing Cert Error!")
	}

	// from loaded cert, get public key and ssl pin
	pubDer := x509.MarshalPKCS1PublicKey(parsedCert.PublicKey.(*rsa.PublicKey))
	pubsum := sha256.Sum256(pubDer)
	pubPin := make([]byte, base64.StdEncoding.EncodedLen(len(pubsum)))
	base64.StdEncoding.Encode(pubPin, pubsum[:])
	ssCert.CertFingerprint = pubPin

	if isCA {
		err = ioutil.WriteFile(conf.CAPrivateKeyPath, ssCert.PrivKeyData, 0600)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Generated Private Key"))
		}
		err = ioutil.WriteFile(conf.CAPath, ssCert.CertData, 0755)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Generated Cert"))
		}
		err = ioutil.WriteFile(conf.CACertPinPath, ssCert.CertFingerprint, 0644)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Certificate Pin"))
		}
	} else {
		err = ioutil.WriteFile(conf.ClientPrivateKeyPath, ssCert.PrivKeyData, 0600)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Generated Private Key"))
		}
		err = ioutil.WriteFile(conf.ClientPath, ssCert.CertData, 0755)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Generated Cert"))
		}
		err = ioutil.WriteFile(conf.ClientPinPath, ssCert.CertFingerprint, 0644)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Certificate Pin"))
		}
	}
	return nil
}