package start

import (
	"errors"
	"fmt"
	"gotime/internal/timeline"
	"time"

	"github.com/spf13/cobra"
)

func New(tl *timeline.Timeline) *cobra.Command {
	var notes []string
	cmd := &cobra.Command{
		Use: "start",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := tl.RunningEvent()
			if errors.Is(err, timeline.ErrNotExist) {
				tl.Start(time.Now(), notes)
				fmt.Println("Time entry created!")
				return nil
			}

			if err == nil {
				return errors.New("cannot start a new entry - current time running")
			}

			return err
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVar(&notes, "notes", []string{}, "add notes while starting your time entry")

	return cmd
}
