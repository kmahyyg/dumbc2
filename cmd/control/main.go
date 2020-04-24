package main

import (
	"github.com/kmahyyg/dumbc2/config"
	"log"
	"github.com/getsentry/sentry-go"
	"github.com/manifoldco/promptui"
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
	config.BuildConf()
	printVersion()

}
