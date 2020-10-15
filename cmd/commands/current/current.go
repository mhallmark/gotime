package current

import (
	"errors"
	"fmt"
	"gotime/internal/timeline"
	"time"

	"github.com/spf13/cobra"
)

func New(tl *timeline.Timeline) *cobra.Command {
	return &cobra.Command{
		Use: "current",
		RunE: func(cmd *cobra.Command, args []string) error {
			running, err := tl.RunningEvent()
			if err != nil {
				if errors.Is(err, timeline.ErrNotExist) {
					fmt.Println("No timer running.")
					return nil
				}

				return err
			}

			dur := time.Since(running.Start)
			fmt.Printf("Current Duration: %v\n", dur.Round(time.Second))

			for _, note := range running.Notes {
				fmt.Printf("-- %q\n", note)
			}

			return nil
		},
	}
}
