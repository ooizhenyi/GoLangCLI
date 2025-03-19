package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//To create file/folder -> mkdir

var createCmd = &cobra.Command{
	Use:     "cf [folderName]",
	Short:   "To Create Folder/File",
	Long:    "Create a new folder in the specified or current directory.",
	Args:    cobra.MinimumNArgs(1),
	Example: "cf myGithub",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}
		folderName := args[0]
		fullPath := filepath.Join(dir, folderName)

		//check whether folder exist in the current directory
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("Folder '%s' already exists\n", folderName)
			return err
		}

		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
			return err
		}

		fmt.Printf("Folder '%s' created successfully\n", folderName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
