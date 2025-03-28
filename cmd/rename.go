package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "rename [old folder name] [new folder name]",
	Short: "Rename a folder",
	Long:  `Rename a folder to a new name.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		oldName := args[0]
		newName := args[1]

		oldPath := filepath.Join(dir, oldName)
		newPath := filepath.Join(dir, newName)

		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			fmt.Printf("Folder '%s' does not exist\n", oldName)
			return err
		}

		if _, err := os.Stat(newPath); err == nil {
			fmt.Printf("Folder '%s' already exists\n", newName)
			return err
		}

		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("Error renaming folder: %v\n", err)
			return err
		}

		fmt.Printf("Folder renamed from '%s' to '%s' successfully\n", oldName, newName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
