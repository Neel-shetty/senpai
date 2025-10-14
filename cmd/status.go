package cmd

import (
	"fmt"
	"senpai/core"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Long: `Displays paths that have differences between the index file and the current HEAD commit, paths that have differences between the working
       tree and the index file, and paths in the working tree that are not tracked by Git (and are not ignored by gitignore(5)). The first are
       what you would commit by running git commit; the second and third are what you could commit by running git add before running git commit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		statues, err := core.Status(".")
		if err != nil {
			return fmt.Errorf("could not check status")
		}
		core.PrettyPrint(statues)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
