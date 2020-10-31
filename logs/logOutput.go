package logs

import (
	"log"
	"os"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Couldn't find a home directory to log to", err)
	}

	logDir := home + "/.config/goplin/"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatal("Couldn't make a log directory", err)
	}

	file, err := os.OpenFile(logDir+"logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}
