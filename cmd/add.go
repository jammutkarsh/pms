/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/JammUtkarsh/pms/utils"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [path...]",
	Short: "Add new projects to the list",
	Long: `add subcommand takes one or more paths to projects.
If no argument is provided, the current directory will be added to the list.
`,
	Example: `  pms add
  pms add .
  pms add ~/path/to/project
  pms add dir1 dir2 dir3
  pms add '$GOPATH/src/github.com/username/project' **`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = []string{"."}
		}
		// keep going past a bad path so the rest still get added
		var errs []error
		for _, arg := range args {
			path, err := filepath.Abs(arg)
			if err == nil {
				err = utils.AddProject(path)
			}
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", arg, err))
			}
		}
		cobra.CheckErr(errors.Join(errs...))
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
