package config

type NetworkService int
type WorkingAs int

const (
	Reverse NetworkService = 0
	Bind NetworkService = 1
	Client WorkingAs = 2
	Server WorkingAs = 3
)

type RemoteConf struct {
	Host string
	Port int
	CertificateData UserConfig
	ConnectionType NetworkService
}
