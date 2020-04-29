package config

import (
	"github.com/kmahyyg/dumbc2/utils"
	"os"
)

const (
	CurrentVersion string = "v0.1.0-git"
	certPathPrefix        = "/.dumbyc2"
	certCC         string = "/.dumbyc2/servercert.pem"
	certCCPK       string = "/.dumbyc2/serverpk.pem"
	certCCPin      string = "/.dumbyc2/serverpin.txt"
	certFC         string = "/.dumbyc2/cacert.pem"
	certPK         string = "/.dumbyc2/caprivkey.pem"
	certPin        string = "/.dumbyc2/cacertpin.txt"
)

var (
	GlobalCert *UserConfig
	GlobalOP   *UserOperation
)

type UserConfig struct {
	OutputPath           string
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
	// only support reverse connection
	ListenAddr   string
	CertLocation string
}

func BuildCertPath(dataDir string) *UserConfig {
	usr := utils.GetAbsolutePath(dataDir)
	GlobalCert = &UserConfig{
		OutputPath:           usr + certPathPrefix,
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

func BuildUserOperation(laddr string, certstor string) *UserOperation {
	GlobalOP = &UserOperation{
		ListenAddr:   laddr,
		CertLocation: certstor,
	}
	return GlobalOP
}
