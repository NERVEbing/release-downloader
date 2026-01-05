package main

import "testing"

func TestMatchPattern(t *testing.T) {	tests := []struct {
		name        string
		str         string
		pattern     string
		wantMatch   bool
		wantErr     bool
	}{
		{
			name:      "exact match",
			str:       "linux-amd64",
			pattern:   "linux-amd64",
			wantMatch: true,
		},
		{
			name:      "regex match",
			str:       "linux-amd64",
			pattern:   "linux.*",
			wantMatch: true,
		},
		{
			name:    "invalid regex",
			str:     "linux-amd64",
			pattern: "[",
			wantErr: true,
		},
		{
			name:      "substring but not regex",
			str:       "linux-amd64",
			pattern:   "linux",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			match, err := matchPattern(tt.str, tt.pattern)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if match != tt.wantMatch {
				t.Fatalf("match = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}
