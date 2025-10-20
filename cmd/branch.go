package cmd

import (
	"fmt"
	"os"
	"senpai/core"

	"github.com/spf13/cobra"
)

var (
	deleteFlag  bool
	allFlag     bool
	verboseFlag bool
)

var branchCmd = &cobra.Command{
	Use:   "branch [flags] <name>",
	Short: "List, create, or delete branches",
	Long: `Manage branches in your repository.

Without arguments, this command lists all existing branches,
marking the current branch with an asterisk (*).

You can also create or delete branches:

  senpai branch new-feature     # Create a new branch
  senpai branch -d old-feature  # Delete a branch
  senpai branch                 # List all branches

Branches are simple references stored under .senpai/refs/heads/.
Each branch points to a specific commit hash.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		if len(args) == 0 {
			branches, err := core.ListBranches(repoPath)
			if err != nil {
				return err
			}

			current, err := core.GetCurrentBranch(repoPath)
			if err != nil {
				current = ""
			}

			for _, b := range branches {
				prefix := "  "
				if b == current {
					prefix = "* "
				}
				if verboseFlag {
					hash, err := core.ResolveBranchCommit(repoPath, b)
					if err == nil && len(hash) >= 7 {
						fmt.Printf("%s%s %s\n", prefix, b, hash[:7])
					} else {
						fmt.Printf("%s%s\n", prefix, b)
					}
				} else {
					fmt.Printf("%s%s\n", prefix, b)
				}
			}
			return nil
		}

		name := args[0]
		if deleteFlag {
			return core.DeleteBranch(repoPath, name)
		}
		return core.CreateBranch(repoPath, name)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "Delete the specified branch")
	branchCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "List all branches (local + remote)")
	branchCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Show commit hashes next to branch names")
}
