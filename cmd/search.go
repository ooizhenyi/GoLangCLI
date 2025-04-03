package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [item]",
	Short: "Search for files and folders",
	Long:  `Search for files and folders matching a specific term.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		dir, _ := cmd.Flags().GetString("dir")
		searchTerm := args[0]
		recursive, _ := cmd.Flags().GetBool("recursive")
		fileType, _ := cmd.Flags().GetString("type")
		caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")

		if !caseSensitive {
			searchTerm = strings.ToLower(searchTerm)
		}
		matches := 0
		searchFunc := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if fileType == "files" && info.IsDir() {
				if path == dir {
					return nil
				}
				return filepath.SkipDir
			}
			if fileType == "folders" && !info.IsDir() {
				return nil
			}

			name := info.Name()
			if !caseSensitive {
				name = strings.ToLower(name)
			}

			if strings.Contains(name, searchTerm) {
				relativePath, _ := filepath.Rel(dir, path)
				if info.IsDir() {
					fmt.Printf("  [DIR] %s\n", relativePath)
				} else {
					fmt.Printf("  [FILE] %s (%d bytes)\n", relativePath, info.Size())
				}
				matches++
			}

			if !recursive && path != dir && info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		err := filepath.Walk(dir, searchFunc)
		if err != nil {
			fmt.Printf("Error during search: %v\n", err)
			return err
		}

		if matches == 0 {
			fmt.Println("No matches found")
		} else {
			fmt.Printf("Found %d matches\n", matches)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolP("recursive", "r", false, "Search recursively")
	searchCmd.Flags().StringP("type", "t", "all", "Type to search (all, files, folders)")
	searchCmd.Flags().BoolP("case-sensitive", "c", false, "Use case-sensitive search")
}
