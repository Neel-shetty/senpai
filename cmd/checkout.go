package cmd

import (
	"fmt"
	"senpai/core"

	"github.com/spf13/cobra"
)

var (
	newBranch bool
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout [flags] <branch | commit>",
	Short: "Switch branches or restore working tree files",
	Long: `Switches to a specified branch or commit.

Examples:
  senpai checkout main
  senpai checkout -b feature/api
  senpai checkout a1b2c3d
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		if newBranch {
			fmt.Printf("Creating and switching to new branch: %s\n", target)
			return core.CheckoutNewBranch(".", target)
		}

		fmt.Printf("Switching to branch or commit: %s\n", target)
		return core.Checkout(".", target)
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
	checkoutCmd.Flags().BoolVarP(&newBranch, "branch", "b", false, "Create and switch to a new branch")
}
