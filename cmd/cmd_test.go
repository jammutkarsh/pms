package cmd

import (
	"path/filepath"
	"testing"
)

func TestPathFinder(t *testing.T) {
	testCases := []struct {
		path     string
		expected string
	}{
		{filepath.Join("user", "home", "username", "project", "..", "project2"), filepath.Join("user", "home", "username", "project2")},
		{filepath.Join("user", "..", "home", "username", "project"), filepath.Join("home", "username", "project")},
		{filepath.Join("user", "home", "username", ".", "project"), filepath.Join("user", "home", "username", "project")},
		{filepath.Join("user", "home", "username", "project", "..", "..", "project2"), filepath.Join("user", "home", "project2")},
		{filepath.Join("user", "home", "..", "username", "project", "..", "..", "project2"), filepath.Join("user", "project2")},
		{filepath.Join("user", "..", "..", "..", "..", "..", "project"), filepath.Join("..", "..", "..", "..", "project")},
		{filepath.Join("..", "..", ".."), filepath.Join("..", "..", "..")},
	}

	for _, tc := range testCases {
		result := pathFinder(tc.path)
		if result != tc.expected {
			t.Errorf("pathFinder(%q) = %q; want %q", tc.path, result, tc.expected)
		}
	}
}
