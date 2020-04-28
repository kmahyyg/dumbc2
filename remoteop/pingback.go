package remoteop

import (
	"log"
	"net"
)


func SendData2Conn(data []byte,conn net.Conn){
	wrbytes, err := conn.Write(data)
	if err != nil || wrbytes != len(data) {
		log.Println("Error in sending pingback...")
	}
}