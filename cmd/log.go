package cmd

import (
	"fmt"
	"senpai/core"
	"time"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "show commit logs",
	Long:  "show commit logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		commits, err := core.Log(".")
		if err != nil {
			return fmt.Errorf("error reading log: %w", err)
		}

		for _, commit := range commits {
			fmt.Printf("commit %s\n", commit.Hash)

			if commit.Author != "" && commit.Email != "" {
				fmt.Printf("Author: %s <%s>\n", commit.Author, commit.Email)
			}

			if commit.Timestamp > 0 {
				t := time.Unix(commit.Timestamp, 0)
				dateStr := t.Format("Mon Jan 2 15:04:05 2006")
				fmt.Printf("Date:   %s %s\n", dateStr, commit.Timezone)
			}

			fmt.Printf("\n    %s\n\n", commit.Message)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
