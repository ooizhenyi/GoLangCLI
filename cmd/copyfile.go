package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copyfile [source file] [destination]",
	Short: "Copy a file",
	Long:  `Copy a file from source to destination.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		sourceFile := filepath.Join(dir, args[0])
		destFile := filepath.Join(dir, args[1])

		srcInfo, err := os.Stat(sourceFile)
		if os.IsNotExist(err) {
			fmt.Printf("Source file '%s' does not exist\n", args[0])
			return err
		}

		if srcInfo.IsDir() {
			fmt.Printf("'%s' is a directory, not a file. Use 'mv' command for directories.\n", args[0])
			return errors.New("use 'mv' command for directories")
		}

		destInfo, err := os.Stat(destFile)
		if err == nil && destInfo.IsDir() {
			destFile = filepath.Join(destFile, filepath.Base(sourceFile))
		}

		destDir := filepath.Dir(destFile)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Printf("Error creating destination directory: %v\n", err)
			return err
		}

		if _, err := os.Stat(destFile); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("Destination file '%s' already exists. Overwrite? (y/n): ", destFile)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Copy cancelled")
					return nil
				}
			}
		}

		source, err := os.Open(sourceFile)
		if err != nil {
			fmt.Printf("Error opening source file: %v\n", err)
			return err
		}
		defer source.Close()

		destination, err := os.Create(destFile)
		if err != nil {
			fmt.Printf("Error creating destination file: %v\n", err)
			return err
		}
		defer destination.Close()

		bytesWritten, err := io.Copy(destination, source)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			return err
		}

		fmt.Printf("File copied successfully from '%s' to '%s' (%d bytes)\n",
			args[0], args[1], bytesWritten)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.Flags().BoolP("force", "f", false, "Force overwrite without confirmation")
}
