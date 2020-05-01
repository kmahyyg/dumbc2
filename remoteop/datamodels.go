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
	CommandBOOM       = "boom"
	CommandBASH       = "bash"
	CommandUPLD       = "upld"
	CommandDWLD       = "dwld"
	CommandEXIT       = "exit"
	StatusOK          = 0
	StatusFinishTrans = 1
	StatusNEXISTS     = 2
	StatusFAILED      = 3
	//StatusOtherErr    = 4
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
`)
}

func ParseUserInput(uipt []string) (*UserCmd, error) {
	var ucmd = &UserCmd{}
	if len(uipt) == 1 || len(uipt) > 3 {
		switch uipt[0] {
		case CommandBASH:
			ucmd.Cmd = CommandBASH
			return ucmd, nil
		case CommandBOOM:
			ucmd.Cmd = CommandBOOM
			return ucmd, nil
		case CommandEXIT:
			return nil, errors.New("USER_EXIT")
		case "help":
			fallthrough
		default:
			printHelp()
			return nil, nil
		}
	} else if len(uipt) == 2 {
		switch uipt[0] {
		default:
			printHelp()
			return nil, errors.New("Unknown Err.")
		}
	} else {
		switch uipt[0]{
		case "download":
			ucmd.Cmd = CommandDWLD
			ucmd.OptionRMT = uipt[1]
			ucmd.OptionLCL = uipt[2]
			return ucmd, nil
		case "upload":
			ucmd.Cmd = CommandUPLD
			ucmd.OptionLCL = uipt[1]
			ucmd.OptionRMT = uipt[2]
			return ucmd, nil
		default:
			printHelp()
			return nil, nil
		}
	}
}
