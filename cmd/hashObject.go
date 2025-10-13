package cmd

import (
	"fmt"
	"os"
	"senpai/core"

	"github.com/spf13/cobra"
)

var (
	hashObjectType string
	writeFlag      bool
)
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [file]",
	Short: "Compute object ID and optionally create an object from a file",
	Long:  `Computes the object ID value for an object with specified type with the contents of the named file (which can be outside of the worktree), and optionally writes the resulting object into the object database. Reports its object ID to its standard output. When <type> is not specified, it defaults to "blob".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("didnt not provide file in command line args, specify the file")
		}
		filePath := args[0]
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file: %v\n", err)
		}

		write, _ := cmd.Flags().GetBool("write")
		objectType, _ := cmd.Flags().GetString("type")

		hash, err := core.HashObject(data, objectType, write)
		if err != nil {
			fmt.Printf("Error hashing object: %v\n", err)
		}
		fmt.Println(hash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)
	hashObjectCmd.Flags().BoolVarP(&writeFlag, "write", "w", false, "Write object into the object database")
	hashObjectCmd.Flags().StringVarP(&hashObjectType, "type", "t", "blob", "Specify object type (blob, tree, commit, tag)")
}
