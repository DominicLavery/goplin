package main

import (
	"dominiclavery/goplin/data"
	_ "dominiclavery/goplin/logs"
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
		//fileSource = data.NewDummySource()
	} else {
		fileSource = data.NewFilesystemSource(path)
	}

	app := ui.MakeApp(fileSource)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
