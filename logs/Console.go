package logs

import (
	"fmt"
	"github.com/derailed/tview"
	"log"
)

var ConsoleView = tview.NewTextView().
	SetDynamicColors(true).
	SetRegions(true)

func SetApp(app *tview.Application) *tview.TextView {
	ConsoleView.SetChangedFunc(func() {
		app.QueueUpdateDraw(func() {
			ConsoleView.ScrollToEnd()
		})
	}).SetBorder(true)

	return ConsoleView
}

func TeeLog(v ...interface{}) {
	log.Println(v...)
	if _, err := fmt.Fprintln(ConsoleView, v...); err != nil {
		log.Println("Couldn't write to console", err)
	}
}
