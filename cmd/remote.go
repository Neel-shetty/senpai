package cmd

import (
	"fmt"
	"os"

	"senpai/core"

	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage set of tracked repositories",
	Long:  `Manage the set of repositories ("remotes") whose branches you track.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		remotes, err := core.ListRemotes(repoPath)
		if err != nil {
			return fmt.Errorf("listing remotes: %w", err)
		}

		for _, remote := range remotes {
			fmt.Println(remote.Name)
		}

		return nil
	},
}

var remoteAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add a new remote",
	Long:  `Add a remote named <name> for the repository at <url>.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url := args[1]

		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		if err := core.AddRemote(repoPath, name, url); err != nil {
			return fmt.Errorf("adding remote: %w", err)
		}

		fmt.Printf("Added remote '%s' with URL '%s'\n", name, url)
		return nil
	},
}

var remoteRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a remote",
	Long:  `Remove the remote named <name>. All remote-tracking branches and configuration settings for the remote are removed.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		if err := core.RemoveRemote(repoPath, name); err != nil {
			return fmt.Errorf("removing remote: %w", err)
		}

		fmt.Printf("Removed remote '%s'\n", name)
		return nil
	},
}

var remoteGetURLCmd = &cobra.Command{
	Use:   "get-url <name>",
	Short: "Print the URL for a remote",
	Long:  `Retrieve the URLs for a remote.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		url, err := core.GetRemoteURL(repoPath, name)
		if err != nil {
			return fmt.Errorf("getting remote URL: %w", err)
		}

		fmt.Println(url)
		return nil
	},
}

var remoteSetURLCmd = &cobra.Command{
	Use:   "set-url <name> <newurl>",
	Short: "Change the URL for a remote",
	Long:  `Change the URL for a remote. Sets first URL for remote <name> that matches regex <newurl>.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		newURL := args[1]

		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		if err := core.SetRemoteURL(repoPath, name, newURL); err != nil {
			return fmt.Errorf("setting remote URL: %w", err)
		}

		fmt.Printf("Changed URL for remote '%s' to '%s'\n", name, newURL)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteRemoveCmd)
	remoteCmd.AddCommand(remoteGetURLCmd)
	remoteCmd.AddCommand(remoteSetURLCmd)
}
