package main

import (
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/common-nighthawk/go-figure"
	"github.com/getsentry/sentry-go"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/useri"
	"github.com/kmahyyg/dumbc2/utils"
	"log"
	"net"
	"os"
	"strconv"
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
	server := parser.Flag("s","server", &argparse.Options{
		Required: false,
		Help:     "Run as Server",
		Default:  true,
	})
	client := parser.Flag("c", "client", &argparse.Options{
		Required: false,
		Help:     "Run as Client",
		Default:  false,
	})
	if *client && *server {
		panic("Conflict Flags.")
	}
	certStor := parser.String("C","cert", &argparse.Options{
		Required: false,
		Help:     "Certificate Location",
		Default:  "~",
	})
	lhost := parser.String("H","host", &argparse.Options{
		Required: true,
		Validate: func(args []string) error {
			ip := net.ParseIP(args[0])
			if ip == nil {
				return errors.New("IP Invalid.")
			}
			return nil
		},
		Help:     "The IP You are gonna listen or connect, default is your interface local IP.",
		Default:  utils.GetLocalIP(),
	})
	lport := parser.Int("P", "port", &argparse.Options{
		Required: true,
		Validate: func(args []string) error {
			si, err := strconv.Atoi(args[0])
			if err != nil || si < 1 || si > 65535{
				return errors.New("Illegal Port.")
			}
			return nil
		},
		Help:     "The Port You are gonna listen or connect, default is 25985.",
		Default:  25985,
	})
	err := parser.Parse(os.Args)
	if err != nil {
		panic(err)
	}
	config.BuildCertPath(*certStor)
	config.BuildUserOperation(*server, *client, *lhost, *lport, *certStor)
	printIPAddr()
	if !config.CheckCert(*client) {
		panic("Certificate not exists. Generate first.")
	}
	switch result {
	case "Server":
		useri.StartServer()
	case "Generate":
		useri.GenerateAgent()
	default:
		log.Fatalln("Illegal Operation. ")
	}

}

func printIPAddr() {
	fmt.Println("Interface IP: " + utils.GetLocalIP())
}

func printBanner() {
	const banner = "DumbY-C2"
	myFigure := figure.NewFigure(banner, "", true)
	myFigure.Print()
	fmt.Println(config.CurrentVersion + " - master - 20200426")
}
