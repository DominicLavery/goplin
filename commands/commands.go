package commands

import (
	"dominiclavery/goplin/data"
	"github.com/spf13/cobra"
	"strings"
)

func NewMakeBookCommand(source data.Source) *cobra.Command {
	helpText := `
Usage: makebook PATH
  Creates a new notebook at the given path PATH.

  PATH should be in the form of notebook names
  separated by a forward slash '/'. For example,
  Parent Book/New Book 

Options:
  [None]
`
	return &cobra.Command{
		Use:     "makeBook",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"mkbook", "mkb", "makebook"},
		Short:   "Creates a new notebook",
		Long:    strings.TrimSpace(helpText),
		Run: func(cmd *cobra.Command, args []string) {
			source.MakeBook(args[0])
		}}
}
