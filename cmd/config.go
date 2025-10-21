package cmd

import (
	"fmt"
	"os"
	"strings"

	"senpai/core"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set repository or global options",
	Long:  `You can query/set/replace/unset options with this command. The name is actually the section and the key separated by a dot, and the value will be escaped.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Usage: senpai config [list|get|set] [options]")
		fmt.Println("Use 'senpai config --help' for more information")
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration variables",
	Long:  `List all variables set in config file, along with their values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		if err := core.ListConfig(repoPath); err != nil {
			return fmt.Errorf("listing config: %w", err)
		}

		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get the value for a given key",
	Long:  `Get the value for a given key. The key should be in the format section.key (e.g., user.name or core.bare).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		// Parse section and key from the dot notation
		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			return fmt.Errorf("key must be in format section.key (e.g., user.name)")
		}

		section := parts[0]
		keyName := parts[1]

		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		value, err := core.GetConfig(repoPath, section, keyName)
		if err != nil {
			return fmt.Errorf("getting config: %w", err)
		}

		fmt.Println(value)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration variable",
	Long:  `Set a configuration variable. The key should be in the format section.key (e.g., user.name or core.bare).`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			return fmt.Errorf("key must be in format section.key (e.g., user.name)")
		}

		section := parts[0]
		keyName := parts[1]

		repoPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		if err := core.SetConfig(repoPath, section, keyName, value); err != nil {
			return fmt.Errorf("setting config: %w", err)
		}

		fmt.Printf("Set %s to '%s'\n", key, value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
