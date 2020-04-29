//go:generate go get -u github.com/rakyll/statik
//go:generate bash -l -c "echo 'Please Copy Your Certs to buildtime/certs. '"
//go:generate bash -l -c "cd ../../; statik -f -include=*.pem, *.txt -src=buildtime/certs/"

package main

import (
	"github.com/akamensky/argparse"
	"github.com/kmahyyg/dumbc2/buildtime"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/utils"
	"os"
)

func main() {
	parser := argparse.NewParser(os.Args[0], "Dumb C2")
	certStor := parser.String("C", "cert", &argparse.Options{
		Required: false,
		Help:     "Certificate Location",
		Default:  "~",
	})
	laddr := parser.String("l", "listen", &argparse.Options{
		Required: true,
		Help:     "The IP You are gonna listen or connect, default is your interface local IP.",
		Default:  utils.GetLocalIP(),
	})
	err := parser.Parse(os.Args)
	if err != nil {
		panic(err)
	}
	*certStor = utils.GetAbsolutePath(*certStor)
	config.BuildCertPath(*certStor)
	config.BuildUserOperation(*laddr, *certStor)
	if buildtime.GetCertificates() != nil{
		panic("Certificate not exists. Generate first.")
	}
}
