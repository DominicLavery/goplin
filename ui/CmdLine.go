package ui

import (
	"dominiclavery/goplin/commands"
	"dominiclavery/goplin/data"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
	"github.com/spf13/cobra"

	"encoding/csv"
	"strings"
)

type CmdLine struct {
	source  data.Source
	rootCmd *cobra.Command
	*tview.InputField
}

func (c CmdLine) finishedFunc(key tcell.Key) {
	if key == tcell.KeyEnter {
		command := c.GetText()
		command = strings.TrimLeft(command, ":")

		r := csv.NewReader(strings.NewReader(command))
		r.Comma = ' '
		record, err := r.Read() //TODO Why might this error?
		if err != nil {
			_ = c.rootCmd.Help() // Roots Help() can't return an error
		}
		c.rootCmd.SetArgs(record)
		_ = c.rootCmd.Execute() //TODO handle me too
		c.SetText("")
	}
}

func MakeCmdLine(source data.Source) *CmdLine {
	cmdLine := CmdLine{InputField: tview.NewInputField(), source: source}
	cmdLine.SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	cmdLine.rootCmd = &cobra.Command{Use: "goplin"}
	mkbookCommand := commands.NewMakeBookCommand(source)
	cmdLine.rootCmd.AddCommand(mkbookCommand)
	cmdLine.InputField.SetFinishedFunc(func(key tcell.Key) { cmdLine.finishedFunc(key) })
	return &cmdLine
}

func (c CmdLine) SetFinishedFunc(handler func(key tcell.Key)) {
	c.InputField.SetFinishedFunc(func(key tcell.Key) {
		c.finishedFunc(key)
		handler(key)
	})
}
