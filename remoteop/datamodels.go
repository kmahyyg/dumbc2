package remoteop

import (
	"log"
	"net"
)

const (
	StatusOK = byte(0)
	StatusFailed = byte(1)
	Delimiter = byte('$')
	StatusEEXIST = byte(2)
)

type PingBack struct {
	StatusCode byte
	FileSize byte    // 0 means status only
	DataPart []byte
}

type RTCommand struct {
	Command []byte
	FilePathLocal []byte
	FilePathRemote []byte
	FileLength byte     // Megabyte
	HasData	 byte		// 0 for false, 1 for true
	RealData  []byte
}


func (pbb *PingBack) BuildnSend(conn net.Conn){
	var data []byte
	data = append(data, pbb.StatusCode, Delimiter, pbb.FileSize, Delimiter)
	data = append(data[:], pbb.DataPart[:]...)
	lendt, err := conn.Write(data)
	if lendt != len(data) || err != nil {
		log.Println("Error in PingBack BuildnSend")
	}
}

func (rtc *RTCommand) BuildnSend(conn net.Conn) {
	var data []byte
	data = append(data[:], rtc.Command[:]...)
	data = append(data, Delimiter, rtc.HasData, Delimiter, rtc.FileLength, Delimiter)
	data = append(data[:], rtc.FilePathRemote[:]...)
	data = append(data, Delimiter)
	data = append(data[:], rtc.RealData...)
	lendt, err := conn.Write(data)
	if lendt != len(data) || err != nil {
		log.Println("Error in RtCommand BuildnSend")
	}
}
