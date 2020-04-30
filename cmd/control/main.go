package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/common-nighthawk/go-figure"
	"github.com/getsentry/sentry-go"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/useri"
	"github.com/kmahyyg/dumbc2/utils"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://48c87cf56e1f4354805e08b4fb5b8b06@o132236.ingest.sentry.io/5212008",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func printVersion() {
	log.Printf("Current Version: %s \n", config.CurrentVersion)
	log.Println("For more information: https://github.com/kmahyyg/dumbc2")
	log.Println("This Program is licensed under AGPLv3.")
}

func main() {
	printBanner()
	printVersion()
	parser := argparse.NewParser(os.Args[0], "Dumb C2")
	certStor := parser.String("c", "cert", &argparse.Options{
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
	printIPAddr()
	if !config.CheckCert(false) {
		panic("Certificate not exists. Generate first.")
	}
	useri.StartServer(*config.GlobalOP)
}

func printIPAddr() {
	fmt.Printf("Interface IPs: %s \n", utils.GetAllIPs())
}

func printBanner() {
	const banner = "DumbY-C2"
	myFigure := figure.NewFigure(banner, "", true)
	myFigure.Print()
	fmt.Println("\t\t   " + config.CurrentVersion + " - master - 20200430")
}
