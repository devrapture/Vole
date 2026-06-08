package scanner

import (
	"testing"
)

func TestUnusedList(t *testing.T) {
	tests := []struct {
		name   string
		assets []*ImageAsset
		want   int
	}{
		{
			name:   "no assets",
			assets: nil,
			want:   0,
		},
		{
			name: "all used",
			assets: []*ImageAsset{
				{Basename: "a.png", Used: true},
				{Basename: "b.png", Used: true},
			},
			want: 0,
		},
		{
			name: "all unused",
			assets: []*ImageAsset{
				{Basename: "a.png", Used: false},
				{Basename: "b.png", Used: false},
			},
			want: 2,
		},
		{
			name: "mixed",
			assets: []*ImageAsset{
				{Basename: "used.png", Used: true},
				{Basename: "unused.png", Used: false},
				{Basename: "also_used.png", Used: true},
				{Basename: "also_unused.png", Used: false},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ScanResult{Assets: tt.assets}
			got := r.UnusedList()
			if len(got) != tt.want {
				t.Errorf("UnusedList() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestUnusedListPreservesOrder(t *testing.T) {
	assets := []*ImageAsset{
		{Basename: "first.png", Used: false},
		{Basename: "second.png", Used: false},
		{Basename: "third.png", Used: true},
		{Basename: "fourth.png", Used: false},
	}
	r := &ScanResult{Assets: assets}
	got := r.UnusedList()
	if len(got) != 3 {
		t.Fatalf("expected 3 unused, got %d", len(got))
	}
	expected := []string{"first.png", "second.png", "fourth.png"}
	for i, a := range got {
		if a.Basename != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], a.Basename)
		}
	}
}
