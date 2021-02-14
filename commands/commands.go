package commands

import (
	"dominiclavery/goplin/data"
	"github.com/spf13/cobra"
	"strings"
)

func NewMakeBookCommand(source data.Source) *cobra.Command {
	helpText := `
Usage: makebook NAME
  Creates a new notebook under the the currently
  selected notebook with the given name.

  Quotes can be used to create books with spaces 
  in the path. For example: "Parent Book/New Book"

Options:
  [None]
`
	return &cobra.Command{
		Use:     "makeBook name",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"mkbook", "mkb", "makebook"},
		Short:   "Creates a new notebook",
		Long:    strings.TrimSpace(helpText),
		RunE: func(cmd *cobra.Command, args []string) error {
			return source.MakeBook(args[0])
		},
	}
}

func NewMakeNoteCommand(source data.Source) *cobra.Command {
	helpText := `
Usage: makenote NAME
  Creates a new note at the currently selected
  notebook with the given name

  Quotes can be used to create notes with spaces 
  in the path. For example: "New Note"

Options:
  [None]
`
	return &cobra.Command{
		Use:     "makeNote name",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"mknote", "mkn", "makenote"},
		Short:   "Creates a new note",
		Long:    strings.TrimSpace(helpText),
		RunE: func(cmd *cobra.Command, args []string) error {
			return source.MakeNote(args[0])
		},
	}
}
