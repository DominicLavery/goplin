package ui

import (
	"dominiclavery/goplin/commands"
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/logs"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
	"github.com/spf13/cobra"
	"log"

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
		record, err := r.Read()
		if err != nil {
			log.Println("Unexpected issue reading a command", err)
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
	cmdLine.rootCmd = &cobra.Command{}
	cmdLine.rootCmd.SetOut(logs.ConsoleView)
	cmdLine.rootCmd.SetErr(logs.ConsoleView)

	mkbookCommand := commands.NewMakeBookCommand(source)
	cmdLine.rootCmd.AddCommand(mkbookCommand)
	mknoteCommand := commands.NewMakeNoteCommand(source)
	cmdLine.rootCmd.AddCommand(mknoteCommand)
	cmdLine.InputField.SetFinishedFunc(func(key tcell.Key) { cmdLine.finishedFunc(key) })
	return &cmdLine
}

func (c CmdLine) SetFinishedFunc(handler func(key tcell.Key)) {
	c.InputField.SetFinishedFunc(func(key tcell.Key) {
		c.finishedFunc(key)
		handler(key)
	})
}
