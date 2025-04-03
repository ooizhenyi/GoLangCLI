package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ooizhenyi/GoLangCLI/utils"
	"github.com/spf13/cobra"
)

var compressCmd = &cobra.Command{
	Use:   "compress [folder] [zipfile]",
	Short: "Compress a folder into a zip file",
	Long:  `Compress a folder and its contents into a zip file.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		sourcePath := filepath.Join(dir, args[0])

		var zipPath string
		if len(args) == 2 {
			zipPath = filepath.Join(dir, args[1])
		} else {
			zipPath = sourcePath + ".zip"
		}

		sourceInfo, err := os.Stat(sourcePath)
		if os.IsNotExist(err) {
			fmt.Printf("Source folder '%s' does not exist\n", args[0])
			return err
		}

		if !sourceInfo.IsDir() {
			fmt.Printf("'%s' is not a folder\n", args[0])
			return nil
		}

		if _, err := os.Stat(zipPath); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("Zip file '%s' already exists. Overwrite? (y/n): ", zipPath)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Compression cancelled")
					return nil
				}
			}
		}

		fmt.Printf("Compressing '%s' to '%s'...\n", args[0], zipPath)
		err = utils.ZipFolder(sourcePath, zipPath)
		if err != nil {
			fmt.Printf("Error compressing folder: %v\n", err)
			return err
		}
		zipInfo, err := os.Stat(zipPath)
		if err != nil {
			fmt.Println("Compression completed successfully")
			return err
		}

		fmt.Printf("Compression completed successfully. Zip file size: %d bytes\n", zipInfo.Size())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)
	compressCmd.Flags().BoolP("force", "f", false, "Force overwrite without confirmation")
}
