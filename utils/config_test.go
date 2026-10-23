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
	}, now)
	want := []string{"few-today", "often-old", "once-now", "never", "never2"}
	for i, p := range got {
		if p.ProjectName != want[i] {
			t.Fatalf("rank order = %v, want %v", got, want)
		}
	}
}
