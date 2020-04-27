package config

import (
	"os"
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


func BuildCertPath() *UserConfig {
	usr, _ := user.Current()
	GlobalCert = &UserConfig{
		ClientPath:           usr.HomeDir + certCC,
		ClientPrivateKeyPath: usr.HomeDir + certCCPK,
		ClientPinPath:        usr.HomeDir + certCCPin,
		CAPrivateKeyPath:     usr.HomeDir + certPK,
		CACertPinPath:        usr.HomeDir + certPin,
		CAPath:               usr.HomeDir + certFC,
	}
	return GlobalCert
}

func CheckCert(isAgent bool) bool {
	check1 := checkFileExists(GlobalCert.ClientPath) && checkFileExists(GlobalCert.ClientPinPath) && checkFileExists(GlobalCert.ClientPrivateKeyPath)
	check2 := checkFileExists(GlobalCert.CAPath) && checkFileExists(GlobalCert.CACertPinPath) && checkFileExists(GlobalCert.CACertPinPath)
	if isAgent {
		if !check1 {
			return false
		}
	} else {
		if !check1 || !check2 {
			return false
		}
	}
	return true
}

func checkFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	} else {
		return true
	}
}

func BuildUserOperation(){
	//todo: build user operation
}
