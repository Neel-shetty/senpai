package cmd

import (
	"fmt"
	"senpai/core"

	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file <object hash> [flags]",
	Short: "Provide contents or details of repository objects",
	Long:  "Output the contents or other properties such as size, type or delta information of one or more objects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no object specified")
		}
		hash := args[0]

		showType, err := cmd.Flags().GetBool("type")
		if err != nil {
			return fmt.Errorf("failed to read flag 'type': %w", err)
		}

		showSize, err := cmd.Flags().GetBool("size")
		if err != nil {
			return fmt.Errorf("failed to read flag 'size': %w", err)
		}

		pretty, err := cmd.Flags().GetBool("pretty")
		if err != nil {
			return fmt.Errorf("failed to read flag 'pretty': %w", err)
		}

		exists, err := cmd.Flags().GetBool("exists")
		if err != nil {
			return fmt.Errorf("failed to read flag 'exists': %w", err)
		}

		return core.CatFile(hash, showType, showSize, pretty, exists)
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	catFileCmd.Flags().BoolP("type", "t", false, "Show object type")
	catFileCmd.Flags().BoolP("size", "s", false, "Show object size")
	catFileCmd.Flags().BoolP("pretty", "p", false, "Pretty-print contents of object")
	catFileCmd.Flags().BoolP("exists", "e", false, "Check if object exists (exit code 0 if true)")
}
