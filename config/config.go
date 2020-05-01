package config

import (
	"github.com/kmahyyg/dumbc2/utils"
	"os"
)

const (
	CurrentVersion string = "v0.4.3-git"
	certPathPrefix string = "/.dumbyc2"
	certSCC        string = "/.dumbyc2/servercert.pem"
	certSCCPK      string = "/.dumbyc2/serverpk.pem"
	certSCCPin     string = "/.dumbyc2/serverpin.txt"
	certFC         string = "/.dumbyc2/cacert.pem"
	certPK         string = "/.dumbyc2/caprivkey.pem"
	certPin        string = "/.dumbyc2/cacertpin.txt"
	certCCC        string = "/.dumbyc2/clientcert.pem"
	certCCCPK      string = "/.dumbyc2/clientpk.pem"
	certCCCPin     string = "/.dumbyc2/clientpin.txt"
)

var (
	GlobalCert *UserConfig
	GlobalOP   *UserOperation
)

type UserConfig struct {
	OutputPath           string
	ServerPath           string
	ServerPrivateKeyPath string
	ServerPinPath        string
	CAPath               string
	CAPrivateKeyPath     string
	CACertPinPath        string
	ClientPath           string
	ClientPrivateKeyPath string
	ClientPinPath        string
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
		ServerPath:           usr + certSCC,
		ServerPrivateKeyPath: usr + certSCCPK,
		ServerPinPath:        usr + certSCCPin,
		CAPrivateKeyPath:     usr + certPK,
		CACertPinPath:        usr + certPin,
		CAPath:               usr + certFC,
		ClientPrivateKeyPath: usr + certCCCPK,
		ClientPath:           usr + certCCC,
		ClientPinPath:        usr + certCCCPin,
	}
	return GlobalCert
}

func CheckCert(isAgent bool) bool {
	check1 := checkFileExists(GlobalCert.ServerPath) && checkFileExists(GlobalCert.ServerPinPath) && checkFileExists(GlobalCert.ServerPrivateKeyPath)
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
