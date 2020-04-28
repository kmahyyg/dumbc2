package remoteop

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
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
	DataLength byte    // 0 means status only
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


func (pbb *PingBack) BuildnSend(conn net.Conn) error {
	var data []byte
	data = append(data, pbb.StatusCode, Delimiter, pbb.DataLength, Delimiter)
	data = append(data[:], pbb.DataPart[:]...)
	lendt, err := conn.Write(data)
	if lendt != len(data) || err != nil {
		log.Println("Error in PingBack BuildnSend")
		return errors.New("Failure in buildNSend.")
	}
	return nil
}

func (pbb *PingBack) Parse(conn net.Conn) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	if err != nil {
		return err
	} else {
		pbb.StatusCode = buf.Bytes()[0]
		pbb.DataLength = buf.Bytes()[2]
		if pbb.DataLength != byte(0) {
			pbb.DataPart = buf.Bytes()[4:]
		}
		return nil
	}
}

func (rtc *RTCommand) BuildnSend(conn net.Conn) error {
	var data []byte
	data = append(data[:], rtc.Command[:]...)
	data = append(data, Delimiter, rtc.HasData, Delimiter, rtc.FileLength, Delimiter)
	data = append(data[:], rtc.FilePathRemote[:]...)
	data = append(data, Delimiter)
	if rtc.FilePathLocal != nil {
		var err error
		rtc.RealData, err = ioutil.ReadFile(string(rtc.FilePathLocal))
		if err != nil {
			log.Println("Read Error for local file.")
			return err
		}
	}
	data = append(data[:], rtc.RealData...)
	lendt, err := conn.Write(data)
	if lendt != len(data) || err != io.EOF {
		log.Println("Error in RtCommand BuildnSend")
		return errors.New("Failure in buildNSend.")
	}
	return nil
}
