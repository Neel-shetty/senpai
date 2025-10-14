package cmd

import (
	"fmt"
	"io"
	"os"
	"senpai/core"
	"strings"

	"github.com/spf13/cobra"
)

var (
	parentHashes []string
	messages     []string
)

// commitTreeCmd represents the commitTree command
var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree <tree hash> [flags]",
	Short: "Create a new commit object",
	Long:  "Creates a new commit object based on the provided tree object and emits the new commit object id on stdout. The log message is read from the standard input, unless -m option is given.",
	RunE: func(cmd *cobra.Command, args []string) error {
		treeHash := args[0]

		var commitMsg string
		if len(messages) > 0 {
			commitMsg = strings.Join(messages, "\n\n")
		} else {
			fmt.Println("Enter commit message, end with Ctrl+D:")
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			commitMsg = string(data)
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

		commitHash, err := core.CommitTree(treeHash, parentHashes, commitMsg, authorName, authorEmail)
		if err != nil {
			return fmt.Errorf("failed to create commit: %w", err)
		}

		fmt.Println(commitHash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitTreeCmd)

	commitTreeCmd.Flags().StringArrayVarP(&parentHashes, "parent", "p", []string{}, "Parent commit hash (can be given multiple times)")
	commitTreeCmd.Flags().StringArrayVarP(&messages, "message", "m", []string{}, "Commit message paragraph (can be given multiple times)")
}
