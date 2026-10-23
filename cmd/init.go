package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// shellInit wraps pms in a function so the cd it prints runs in the current shell.
const shellInit = `pms() {
  if [ $# -eq 0 ]; then
    local cmd
    cmd="$(command pms --shell)" && [ -n "$cmd" ] && eval "$cmd"
  else
    command pms "$@"
  fi
}
`

var initCmd = &cobra.Command{
	Use:       "init [zsh|bash]",
	Short:     "Print shell integration so pms can cd in your current shell",
	Example:   `  echo 'eval "$(pms init zsh)"' >> ~/.zshrc`,
	ValidArgs: []string{"zsh", "bash"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(shellInit)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
