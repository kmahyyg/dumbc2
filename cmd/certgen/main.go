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
	"flag"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	var fileLoca string
	flag.StringVar(&fileLoca, "o", "~", "Output Directory")
	flag.Parse()
	fileLoca = utils.GetAbsolutePath(fileLoca)
	log.Println("Program started.")
	currentConf := config.BuildCertPath(fileLoca)
	_ = os.Mkdir(currentConf.OutputPath, 0755)
	if config.CheckCert(false) {
		log.Fatalln("Certificate Already Exists. Delete ~/.dumbyc2 then re-run this program.")
	}
	log.Println("Start Generation, Please hang on.")
	_, err := generateCertificate(currentConf)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Generate finished.")
	}
}

func generateCertificate(conf *config.UserConfig) ([]*config.SSCertificate, error) {
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
	var ssClientCert = &config.SSCertificate{
		CertData:        nil,
		PrivKeyData:     nil,
		CertFingerprint: nil,
	}
	bits := 4096
	// generate rsa key pairs
	cAprivKey, err := rsa.GenerateKey(rand.Reader, bits)
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	cLientprivKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return []*config.SSCertificate{ssCACert, ssCert, ssClientCert}, errors.Wrap(err, "Generate PRIVKEY Error")
	}
	// generate random things
	serialNo, _ := rand.Int(rand.Reader, big.NewInt(utils.RandMathInt64(65535)))
	serialNo2, _ := rand.Int(rand.Reader, big.NewInt(utils.RandMathInt64(65535)))
	serialNo3, _ := rand.Int(rand.Reader, big.NewInt(utils.RandMathInt64(65535)))
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
		SerialNumber: serialNo2,
		Subject: pkix.Name{
			CommonName: "dumbyc2-S." + randDomain,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	ssClientSignedCert := &x509.Certificate{
		SerialNumber: serialNo3,
		Subject: pkix.Name{
			CommonName: "dumbyc2-C." + randDomain,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	derCACert, err := x509.CreateCertificate(rand.Reader, ssCASignedCert, ssCASignedCert, &cAprivKey.PublicKey, cAprivKey)
	derCert, err := x509.CreateCertificate(rand.Reader, ssSignedCert, ssCASignedCert, &privKey.PublicKey, privKey)
	derClientCert, err := x509.CreateCertificate(rand.Reader, ssClientSignedCert, ssCASignedCert, &privKey.PublicKey, privKey)

	if err != nil {
		return []*config.SSCertificate{ssCACert, ssCert, ssClientCert}, errors.Wrap(err, "Create Certificate Error")
	}

	err = writeCert(derCACert, conf, cAprivKey, ssCACert, true, false)
	if err != nil {
		panic(err)
	}
	err = writeCert(derCert, conf, privKey, ssCert, false, true)
	if err != nil {
		panic(err)
	}
	err = writeCert(derClientCert, conf, cLientprivKey, ssClientCert, false, false)
	if err != nil {
		panic(err)
	}

	return []*config.SSCertificate{ssCACert, ssCert, ssClientCert}, nil
}

func writeCert(cert []byte, conf *config.UserConfig, certKey *rsa.PrivateKey, ssCert *config.SSCertificate, isCA bool, isServer bool) error {
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
	} else if isServer {
		err = ioutil.WriteFile(conf.ServerPrivateKeyPath, ssCert.PrivKeyData, 0600)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Generated Private Key"))
		}
		err = ioutil.WriteFile(conf.ServerPath, ssCert.CertData, 0755)
		if err != nil {
			panic(errors.Wrap(err, "Errors in Write Generated Cert"))
		}
		err = ioutil.WriteFile(conf.ServerPinPath, ssCert.CertFingerprint, 0644)
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
