package config

import (
	"os/user"
)

const CurrentVersion string = "v0.1.0-git"
const certCC string = "/.dumbyc2/clientcert.pem"
const certCCPK string = "/.dumbyc2/clientpk.pem"
const certCCPin string = "/.dumbyc2/clientpin.txt"
const certFC string = "/.dumbyc2/cacert.pem"
const certPK string = "/.dumbyc2/caprivkey.pem"
const certPin string = "/.dumbyc2/cacertpin.txt"

var (
	GlobalCert *UserConfig
)

type UserConfig struct {
	ClientPath           string
	ClientPrivateKeyPath string
	ClientPinPath        string
	CAPath               string
	CAPrivateKeyPath     string
	CACertPinPath        string
}

type SSCertificate struct {
	CertData        []byte
	PrivKeyData     []byte
	CertFingerprint []byte
}

type UserOperation struct {
	IsServer int
	Host string
	Port int
}

func CheckCert() {
	usr, _ := user.Current()
	GlobalCert = &UserConfig{
		ClientPath:           usr.HomeDir + certCC,
		ClientPrivateKeyPath: usr.HomeDir + certCCPK,
		ClientPinPath:        usr.HomeDir + certCCPin,
		CAPrivateKeyPath:     usr.HomeDir + certPK,
		CACertPinPath:        usr.HomeDir + certPin,
		CAPath:               usr.HomeDir + certFC,
	}
	//todo: check if all file exists
}

func BuildUserOperation(){
	//todo: build user operation
}
