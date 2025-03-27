package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls [folderName]",
	Short: "List contents of a folder",
	Long:  `Display the files and subfolders inside a folder.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}
		targetFolder := dir

		if len(args) > 0 {
			targetFolder = filepath.Join(dir, args[0])
		}

		info, err := os.Stat(targetFolder)
		if os.IsNotExist(err) {
			fmt.Printf("Folder '%s' does not exist\n", targetFolder)
			return err
		}

		if !info.IsDir() {
			fmt.Printf("'%s' is not a folder\n", targetFolder)
			return errors.New("folder name provided is not a folder")
		}

		listType, _ := cmd.Flags().GetString("type")

		entries, err := os.ReadDir(targetFolder)
		if err != nil {
			fmt.Printf("Error reading folder: %v\n", err)
			return err
		}

		fmt.Printf("Contents of '%s':\n", targetFolder)

		var filteredEntries []fs.DirEntry
		for _, entry := range entries {
			switch listType {
			case "all":
				filteredEntries = append(filteredEntries, entry)
			case "files":
				if !entry.IsDir() {
					filteredEntries = append(filteredEntries, entry)
				}
			case "folders":
				if entry.IsDir() {
					filteredEntries = append(filteredEntries, entry)
				}
			}
		}
		if len(filteredEntries) == 0 {
			fmt.Println("  (empty)")
			return errors.New("folder is empty")
		}
		detailed, _ := cmd.Flags().GetBool("detailed")
		for _, entry := range filteredEntries {
			if detailed {
				info, _ := entry.Info()
				modTime := info.ModTime().Format("Jan 02 15:04")
				size := info.Size()
				mode := info.Mode().String()

				entryType := "F"
				if entry.IsDir() {
					entryType = "D"
					size = 0
				}

				fmt.Printf("  %s %s %8d %s %s\n", entryType, mode, size, modTime, entry.Name())
			} else {
				if entry.IsDir() {
					fmt.Printf("  [DIR] %s\n", entry.Name())
				} else {
					fmt.Printf("  [FILE] %s\n", entry.Name())
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("type", "t", "all", "Type of entries to list (all, files, folders)")
	listCmd.Flags().BoolP("detailed", "l", false, "Show detailed information")
}
