package useri

import (
	"github.com/go-errors/errors"
	"github.com/hashicorp/yamux"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/kmahyyg/dumbc2/utils"
	"github.com/manifoldco/promptui"
	"log"
	"net"
	"strconv"
	"time"
)

func StartServer(){
	allIPs := utils.GetAllIPs()
	prpt_ip := promptui.Select{
		Label: "Select Bind IP: ",
		Items: allIPs,
	}
	_, lres, err := prpt_ip.Run()

	if err != nil {
		log.Fatal("Internal Error.")
	}

	prpt_port := promptui.Prompt{
		Label: "Bind Port: ",
		Validate: func(portipt string) error {
			data, err := strconv.Atoi(portipt)
			if err != nil || data < 1 || data > 65534 {
				return errors.New("Input Port Error.")
			}
			return nil
		},
	}

	lport, err := prpt_port.Run()
	if err != nil {
		log.Fatalln("Internal Error.")
	}

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
			sess, err := yamux.Server(conn, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			stream, err := sess.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			handleClient(stream)
		} else {
			log.Fatalln("Failed to bind.")
		}
	}
}

func handleClient(stream net.Conn){
	// todo: full interactive pty
}
