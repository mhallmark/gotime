package add

import (
	"gotime/cmd/commands/add/note"
	"gotime/internal/timeline"

	"github.com/spf13/cobra"
)

func New(tl *timeline.Timeline) *cobra.Command {
	cmd := &cobra.Command{
		Use: "add <item>",
	}

	cmd.AddCommand(note.New(tl))

	return cmd
}
