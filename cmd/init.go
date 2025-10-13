package cmd

import (
	"fmt"
	"senpai/core"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initialBranch string
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an empty git repository in the current directory",
	Long:  `This command creates an empty Git repository - basically a .git directory with subdirectories for objects, refs/heads, refs/tags, and template files. An initial branch without any commits will be created`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := core.InitRepo(".", initialBranch); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Initialized empty repository on branch", initialBranch)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringVar(&initialBranch, "initial-branch", "master", "Name of initial branch")
}
