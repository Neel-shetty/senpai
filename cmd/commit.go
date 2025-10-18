package cmd

import (
	"fmt"
	"os"
	"senpai/core"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Long: `Create a new commit containing the current contents of the index and the given log message describing the changes.
The new commit is a direct child of HEAD, usually the tip of the current branch, and the branch is updated to
point to it`,
	RunE: func(cmd *cobra.Command, args []string) error {
		msg, err := cmd.Flags().GetString("message")
		if err != nil {
			return fmt.Errorf("failed to read message flag: %w", err)
		}
		if msg == "" {
			return fmt.Errorf("commit message required (use -m)")
		}
		authorName := os.Getenv("GIT_AUTHOR_NAME")
		authorEmail := os.Getenv("GIT_AUTHOR_EMAIL")

		if authorName == "" {
			authorName = os.Getenv("GIT_COMMITTER_NAME")
		}
		if authorEmail == "" {
			authorEmail = os.Getenv("GIT_COMMITTER_EMAIL")
		}

		if authorName == "" || authorEmail == "" {
			return fmt.Errorf("missing author info: set GIT_AUTHOR_NAME and GIT_AUTHOR_EMAIL (or committer variants)")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		commitHash, err := core.Commit(cwd, msg, authorName, authorEmail)
		if err != nil {
			return fmt.Errorf("failed to create commit: %w", err)
		}

		fmt.Printf("[%s] %s\n", commitHash[:7], msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "commit message")
}
