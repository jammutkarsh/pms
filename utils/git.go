package utils

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// gitStatus is what the list shows under each project.
type gitStatus struct {
	branch  string
	changed int // changed + untracked files
	missing bool
}

type gitStatusMsg struct {
	path   string
	status gitStatus
}

// gitSlots caps how many `git status` run at once so a long list doesn't fork-bomb.
var gitSlots = make(chan struct{}, 8)

// fetchGitStatus loads a project's git status in the background. Non-repos
// (and repos too slow to answer) just get no git details.
func fetchGitStatus(path string) tea.Cmd {
	return func() tea.Msg {
		gitSlots <- struct{}{}
		defer func() { <-gitSlots }()

		if _, err := os.Stat(path); err != nil {
			return gitStatusMsg{path, gitStatus{missing: true}}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		c := exec.CommandContext(ctx, "git", "-C", path, "status", "--porcelain=v2", "--branch")
		// don't take the index lock, so we never get in the way of the user's own git
		c.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
		out, err := c.Output()
		if err != nil {
			return nil
		}
		return gitStatusMsg{path, parseGitStatus(string(out))}
	}
}

// parseGitStatus reads `git status --porcelain=v2 --branch` output.
func parseGitStatus(out string) (s gitStatus) {
	for _, line := range strings.Split(out, "\n") {
		switch {
		case line == "":
		case strings.HasPrefix(line, "# branch.head "):
			s.branch = strings.TrimPrefix(line, "# branch.head ")
		case !strings.HasPrefix(line, "#"):
			s.changed++
		}
	}
	return s
}
