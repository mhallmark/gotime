package list

import (
	"fmt"
	"gotime/internal/timeline"

	"github.com/spf13/cobra"
)

func New(tl *timeline.Timeline) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			f := tl.ReportFormat()
			fmt.Println(f)
		},
	}
}
