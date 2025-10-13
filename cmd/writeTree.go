package cmd

import (
	"fmt"
	"os"
	"senpai/core"

	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create a tree object from the current index",
	Long:  "Creates a tree object using the current index. The name of the new tree object is printed to standard output.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		treeHash, err := core.WriteTree(cwd)
		if err != nil {
			return fmt.Errorf("write-tree failed: %w", err)
		}

		fmt.Println(treeHash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
