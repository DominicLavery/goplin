package main

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/ui"
	"log"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var fileSource data.Source
	if len(os.Args) > 1 {
		if os.Args[1] != "devmode" {
			log.Fatal("Unknown command", os.Args[1])
		}
		fileSource = data.NewDummySource()
	} else {
		fileSource = data.NewFilesystemSource(path)
	}

	setUpLogs()
	app := ui.MakeApp(fileSource)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func setUpLogs() {
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
