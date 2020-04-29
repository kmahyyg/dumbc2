package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetAbsolutePath(dataDir string) string {
	var usr string
	var err error
	if len(dataDir) != 0 {
		homeDir, _ := os.UserHomeDir()
		usr = dataDir
		if strings.HasPrefix(dataDir, "~") {
			usr = strings.Replace(dataDir, "~", homeDir, 1)
		}
		usr, _ = filepath.Abs(usr)
	} else {
		usr, err = os.UserHomeDir()
		if err != nil {
			log.Fatalln(err)
		}
	}
	return usr
}
