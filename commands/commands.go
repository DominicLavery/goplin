package commands

import (
	"dominiclavery/goplin/data"
	"github.com/spf13/cobra"
	"strings"
)

func NewMakeBookCommand(update func()) *cobra.Command {
	helpText := `
Usage: makebook NAME
  Creates a new notebook under the the currently
  selected notebook with the given name.

  Quotes can be used to create books with spaces
  in the name. For example: "New Book"

Options:
  [None]
`
	return &cobra.Command{
		Use:     "makeBook name",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"mb", "mkbook", "mkb", "makebook"},
		Short:   "Creates a new notebook",
		Long:    strings.TrimSpace(helpText),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := data.MakeBook(args[0])
			if err != nil {
				return err
			}
			update()
			return nil
		},
	}
}

func NewMakeNoteCommand(update func()) *cobra.Command {
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
		Aliases: []string{"mn", "mknote", "mkn", "makenote"},
		Short:   "Creates a new note",
		Long:    strings.TrimSpace(helpText),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := data.MakeNote(args[0])
			if err != nil {
				return err
			}
			update()
			return nil
		},
	}
}
