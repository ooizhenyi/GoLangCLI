package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ooizhenyi/GoLangCLI/utils"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract [zipfile] [destination]",
	Short: "Extract a zip file",
	Long:  `Extract the contents of a zip file to a destination folder.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		zipPath := filepath.Join(dir, args[0])

		var extractPath string
		if len(args) == 2 {
			extractPath = filepath.Join(dir, args[1])
		} else {
			zipBase := filepath.Base(zipPath)
			extractPath = filepath.Join(dir, strings.TrimSuffix(zipBase, filepath.Ext(zipBase)))
		}

		if _, err := os.Stat(zipPath); os.IsNotExist(err) {
			fmt.Printf("Zip file '%s' does not exist\n", args[0])
			return err
		}

		if _, err := os.Stat(extractPath); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("Destination folder '%s' already exists. Files may be overwritten. Continue? (y/n): ", extractPath)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Extraction cancelled")
					return nil
				}
			}
		}

		fmt.Printf("Extracting '%s' to '%s'...\n", args[0], extractPath)
		err := utils.UnzipFile(zipPath, extractPath)
		if err != nil {
			fmt.Printf("Error extracting zip file: %v\n", err)
			return err
		}

		fmt.Println("Extraction completed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.Flags().BoolP("force", "f", false, "Force extraction without confirmation")
}
