package main

import (
	"github.com/getsentry/sentry-go"
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


func main() {
	
}
