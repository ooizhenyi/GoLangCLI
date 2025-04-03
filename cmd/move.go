package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "mv [sourceFolder] [destinationFolder]",
	Short: "To move folder",
	Long:  "Move a folder from one location to another.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		sourceFolder := filepath.Join(dir, args[0])
		destinationFolder := filepath.Join(dir, args[1])

		if _, err := os.Stat(sourceFolder); os.IsNotExist(err) {
			fmt.Printf("Folder '%s' does not exists\n", sourceFolder)
			return err
		}

		destParent := filepath.Dir(destinationFolder)
		if _, err := os.Stat(destParent); os.IsNotExist(err) {
			fmt.Printf("Error destination directory: %v\n", err)
			return err
		}

		err = os.Rename(sourceFolder, destinationFolder)
		if err != nil {
			fmt.Println("Unable to move the folder to the destination")
			return err
		}

		fmt.Printf("Folder moved from '%s' to '%s' successfully\n", args[0], args[1])
		return nil

	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
