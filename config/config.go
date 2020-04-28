package config

import (
	"github.com/kmahyyg/dumbc2/utils"
	"os"
)

const CurrentVersion string = "v0.1.0-git"
const certPathPrefix = "/.dumbyc2"
const certCC string = "/.dumbyc2/clientcert.pem"
const certCCPK string = "/.dumbyc2/clientpk.pem"
const certCCPin string = "/.dumbyc2/clientpin.txt"
const certFC string = "/.dumbyc2/cacert.pem"
const certPK string = "/.dumbyc2/caprivkey.pem"
const certPin string = "/.dumbyc2/cacertpin.txt"

var (
	GlobalCert *UserConfig
	GlobalOP *UserOperation
)

type UserConfig struct {
	OutputPath			 string
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
	IsServer bool
	Host string
	Port int
	CertLocation string
}


func BuildCertPath(dataDir string) *UserConfig {
	usr := utils.GetAbsolutePath(dataDir)
	GlobalCert = &UserConfig{
		OutputPath: 		  usr + certPathPrefix,
		ClientPath:           usr + certCC,
		ClientPrivateKeyPath: usr + certCCPK,
		ClientPinPath:        usr + certCCPin,
		CAPrivateKeyPath:     usr + certPK,
		CACertPinPath:        usr + certPin,
		CAPath:               usr + certFC,
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

func BuildUserOperation(server bool, client bool, lhost string, lport int, certstor string) *UserOperation{
	var isServer bool
	if server && client {
		panic("Conflic Settings.")
	} else if server {
		isServer = true
	} else {
		isServer = false
	}
	GlobalOP = &UserOperation{
		IsServer:     isServer,
		Host:         lhost,
		Port:         lport,
		CertLocation: certstor,
	}
	return GlobalOP
}
