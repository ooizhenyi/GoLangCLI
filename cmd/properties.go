package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var propertiesCmd = &cobra.Command{
	Use:   "ppt [folder name]",
	Short: "View folder properties",
	Long:  `Display properties of a folder such as size, creation time, and permissions.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		folderName := args[0]
		fullPath := filepath.Join(dir, folderName)

		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			fmt.Printf("Folder '%s' does not exist\n", folderName)
			return
		}

		if !info.IsDir() {
			fmt.Printf("'%s' is not a folder\n", folderName)
			return
		}

		var size int64
		filepath.Walk(fullPath, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				size += info.Size()
			}
			return nil
		})

		var fileCount, folderCount int
		filepath.Walk(fullPath, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				folderCount++
			} else {
				fileCount++
			}
			return nil
		})

		fmt.Printf("Properties for folder '%s':\n", folderName)
		fmt.Printf("  Path: %s\n", fullPath)
		fmt.Printf("  Size: %d bytes\n", size)
		fmt.Printf("  Created: %s\n", formatTime(info.ModTime()))
		fmt.Printf("  Modified: %s\n", formatTime(info.ModTime()))
		fmt.Printf("  Permissions: %s\n", info.Mode().String())
		fmt.Printf("  Files: %d\n", fileCount-1)
		fmt.Printf("  Subfolders: %d\n", folderCount-1)
	},
}

func formatTime(t time.Time) string {
	return t.Format("Jan 02, 2006 15:04:05")
}

func init() {
	rootCmd.AddCommand(propertiesCmd)
}
