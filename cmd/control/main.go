package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/getsentry/sentry-go"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/utils"
	"github.com/kmahyyg/dumbc2/useri"
	"log"
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
	config.CheckCert()
	config.BuildUserOperation()
	printVersion()
	printBanner()
	printIPAddr()
	//todo: argparse

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
