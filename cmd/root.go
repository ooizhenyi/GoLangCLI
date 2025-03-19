/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "filemanager",
	Short: "A file management CLI application",
	Long: `A file management CLI application that allows you to:
- Create folders
- Move folders
- Delete folders
- Edit folder names
- View folder properties
- View files inside folders`,
}


func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("dir", "d", ".", "Directory to work with")
}


