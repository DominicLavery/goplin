package main

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/ui"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var notebooks, notes = data.NewFilesystemSource(path).Dataset()
	app := ui.MakeApp(notebooks, notes)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
