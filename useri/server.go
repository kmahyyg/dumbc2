package useri

import (
	"github.com/go-errors/errors"
	"github.com/kmahyyg/dumbc2/utils"
	"github.com/kmahyyg/dumbc2/transport"
	"github.com/manifoldco/promptui"
	"log"
	"strconv"
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

	fladdr := lres + lport
    lbserver, err := transport.TLSServerBuilder(fladdr, false)
    for true{
    	conn, err := lbserver.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

	}
}
