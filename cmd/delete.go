package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "dlt [folderName]",
	Short:   "To Delete Folder/File",
	Long:    "Delete a folder and all its contents.",
	Args:    cobra.ExactArgs(1),
	Example: "dlt myGithub",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}
		folderName := args[0]
		fullPath := filepath.Join(dir, folderName)

		//check folder exist or not
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Printf("Folder '%s' does not exists\n", folderName)
			return err
		}

		err = os.RemoveAll(fullPath)
		if err != nil {
			fmt.Printf("Error deleting folder: %v\n", err)
			return err
		}
		fmt.Printf("Folder '%s' deleted successfully\n", folderName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
