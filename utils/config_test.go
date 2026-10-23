package utils

import (
	"testing"
	"time"
)

func TestRank(t *testing.T) {
	now := time.Now()
	ago := func(d time.Duration) int64 { return now.Add(-d).Unix() }
	got := rank([]Project{
		{ProjectName: "never"},
		{ProjectName: "often-old", Opens: 20, LastOpened: ago(30 * 24 * time.Hour)}, // 5
		{ProjectName: "once-now", Opens: 1, LastOpened: ago(time.Minute)},           // 4
		{ProjectName: "few-today", Opens: 3, LastOpened: ago(2 * time.Hour)},        // 6
		{ProjectName: "never2"},
		{ProjectName: "pinned", Pinned: true},
	}, now)
	want := []string{"pinned", "few-today", "often-old", "once-now", "never", "never2"}
	for i, p := range got {
		if p.ProjectName != want[i] {
			t.Fatalf("rank order = %v, want %v", got, want)
		}
	}
}

func TestParseGitStatus(t *testing.T) {
	out := `# branch.oid 6988e3c
# branch.head main
# branch.upstream origin/main
# branch.ab +0 -0
1 .M N... 100644 100644 100644 a b cmd/root.go
2 R. N... 100644 100644 100644 a b R100 new.go	old.go
? scratch.txt
`
	if got := parseGitStatus(out); got != (gitStatus{branch: "main", changed: 3}) {
		t.Errorf("parseGitStatus = %+v", got)
	}
}
