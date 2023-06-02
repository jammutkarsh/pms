/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project to the list",
	Long: `add subcommand takes one argument which is the path to the project.
If no argument is provided, the current directory will be added to the list.
`,
	Example: `  pms add
  pms add .
  pms add ~/path/to/project
  pms add '$GOPATH/src/github.com/username/project' **`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Return the current directory
		} else if len(args) == 1 {
			// User provides a path to the project; Could be relative path or absolute path
		} else {
			// Invalid number of arguments

		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
