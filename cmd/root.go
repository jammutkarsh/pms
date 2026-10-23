/*
Copyright © 2023 Utkarsh Chourasia 5.utkarshc@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/JammUtkarsh/pms/utils"
	"github.com/spf13/cobra"
)

// shellMode prints a shell command for the `pms init` wrapper to eval, so cd sticks in the caller's shell.
var shellMode bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pms",
	Short: "A project launcher",
	Long: `PMS - Project Management S̶y̶s̶t̶e̶m̶ Sucks

Search your projects and open them in your editor, IDE, AI agents or just cd into them.
Keys and actions are configured in ~/.pms.json.`,
	Example: `  pms
  pms https://github.com/user/repo                  # clone into current dir and add
  pms https://github.com/user/repo ~/working/path   # clone into an existing dir
  pms https://github.com/user/repo myRepo           # clone as myRepo
  eval "$(pms init zsh)"   # add to ~/.zshrc so cd works in your shell
`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if !isGitURL(args[0]) {
				cobra.CheckErr(fmt.Errorf("unknown command or git URL %q", args[0]))
			}
			cobra.CheckErr(clone(args))
			return
		}
		m := utils.NewModel()
		// The UI goes to stderr so stdout stays clean for the shell wrapper.
		if _, err := tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run(); err != nil {
			cobra.CheckErr(err)
		}
		if m.Action == nil {
			return
		}
		if _, err := os.Stat(m.Chosen.ProjectPath); err != nil {
			cobra.CheckErr(fmt.Errorf("%s no longer exists, remove it from the list with ctrl+d", m.Chosen.ProjectPath))
		}
		cobra.CheckErr(utils.UpdateProjectList(m.Chosen.ProjectPath))

		if shellMode {
			line := "cd -- " + shellQuote(m.Chosen.ProjectPath)
			if m.Action.Run != "" {
				line += " && " + m.Action.Run
			}
			fmt.Println(line)
			return
		}

		// Without the wrapper a child process can't change our parent's directory,
		// so run the action (or a subshell for cd) inside the project instead.
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "sh"
		}
		var shellArgs []string
		if m.Action.Run != "" {
			shellArgs = []string{"-c", m.Action.Run}
		} else {
			fmt.Fprintln(os.Stderr, `pms: opening a subshell; add eval "$(pms init zsh)" to your shell rc to cd in place`)
		}
		c := exec.Command(shell, shellArgs...)
		c.Dir, c.Stdin, c.Stdout, c.Stderr = m.Chosen.ProjectPath, os.Stdin, os.Stdout, os.Stderr
		if err := c.Run(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				os.Exit(exitErr.ExitCode())
			}
			cobra.CheckErr(err)
		}
	},
	Version: "0.5.0",
}

func isGitURL(s string) bool {
	for _, p := range []string{"https://", "http://", "ssh://", "git://", "git@"} {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// cloneDest picks where `git clone` puts the repo: inside dest when it's an
// existing directory, at dest otherwise (so it doubles as a rename), or in the
// current directory when no dest is given.
func cloneDest(url, dest string) string {
	name := strings.TrimSuffix(strings.TrimRight(url, "/"), ".git")
	name = name[strings.LastIndexAny(name, "/:")+1:]
	if dest == "" {
		return name
	}
	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		return filepath.Join(dest, name)
	}
	return dest
}

// clone runs `git clone url [dest]` and adds the result to the project list.
func clone(args []string) error {
	var dest string
	if len(args) > 1 {
		dest = args[1]
	}
	dest = cloneDest(args[0], dest)
	c := exec.Command("git", "clone", "--", args[0], dest)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stderr, os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	path, err := filepath.Abs(dest)
	if err != nil {
		return err
	}
	return utils.AddProject(path)
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pms.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVar(&shellMode, "shell", false, "print a shell command instead of running it (used by `pms init`)")
	rootCmd.Flags().MarkHidden("shell")
}
