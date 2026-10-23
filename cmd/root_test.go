package cmd

import (
	"path/filepath"
	"testing"
)

func TestCloneDest(t *testing.T) {
	dir := t.TempDir()
	for _, tc := range []struct{ url, dest, want string }{
		{"https://github.com/user/repo", "", "repo"},
		{"https://github.com/user/repo.git/", "", "repo"},
		{"git@github.com:user/repo.git", "", "repo"},
		{"git@host:repo.git", "", "repo"},
		{"https://github.com/user/repo", dir, filepath.Join(dir, "repo")},
		{"https://github.com/user/repo", "myRepo", "myRepo"},
	} {
		if got := cloneDest(tc.url, tc.dest); got != tc.want {
			t.Errorf("cloneDest(%q, %q) = %q, want %q", tc.url, tc.dest, got, tc.want)
		}
	}
	if isGitURL("ad") || !isGitURL("git@github.com:user/repo") {
		t.Error("isGitURL misclassified")
	}
}
