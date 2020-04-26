package useri

import (
	"github.com/manifoldco/promptui"
	"log"
)

func GenerateAgent(){
	pmpt_conntype := promptui.Select{
		Label: "Connection Type: ",
		Items: []string{"Reverse", "Bind"},
	}
	_, result, err := pmpt_conntype.Run()
	
	if err != nil {
		log.Fatalln(err)
	}

	switch result {
	case "Reverse":

	case "Bind":
		
	}
}
