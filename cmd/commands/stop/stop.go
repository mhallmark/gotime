package stop

import (
	"errors"
	"fmt"
	"gotime/internal/timeline"
	"time"

	"github.com/spf13/cobra"
)

func New(tl *timeline.Timeline) *cobra.Command {
	cmd := &cobra.Command{
		Use: "stop",
		RunE: func(cmd *cobra.Command, args []string) error {
			running, err := tl.RunningEvent()
			if err != nil {
				if errors.Is(err, timeline.ErrNotExist) {
					return errors.New("you have no running time - exiting")
				}

				return err
			}

			n := time.Now()
			running.End = &n

			dur := n.Sub(running.Start)

			fmt.Printf("Time stoped at %v\n", dur)

			return nil
		},
	}

	return cmd
}
