package main

import (
	"dominiclavery/goplin/data"
	_ "dominiclavery/goplin/logs"
	"dominiclavery/goplin/ui"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	data.RegisterSource("Local", data.NewFilesystemSource(path))
	//data.RegisterSource("inMem", data.NewInMemorySource())

	app := ui.MakeApp()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
