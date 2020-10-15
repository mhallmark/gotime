package note

import (
	"errors"
	"gotime/internal/timeline"

	"github.com/spf13/cobra"
)

func New(tl *timeline.Timeline) *cobra.Command {
	cmd := &cobra.Command{
		Use: "note <note>",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				err := tl.AddNote(arg)
				if err != nil {
					if errors.Is(err, timeline.ErrNotExist) {
						return errors.New("you have no running time - exiting")
					}

					return err
				}
			}

			return nil
		},
	}

	return cmd
}
