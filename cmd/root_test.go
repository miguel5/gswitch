package cmd

import "testing"

func TestTicketFromBranch(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		want   string
	}{
		{"feature branch", "feature/ticket-42-login", "TICKET-42"},
		{"bugfix branch", "bugfix/ABC-1-fix-thing", "ABC-1"},
		{"mixed case", "feature/Deploy-7-stuff", "DEPLOY-7"},
		{"no ticket", "main", ""},
		{"detached HEAD", "HEAD", ""},
		{"no digit", "feature/foo-bar", ""},
		{"no prefix letter", "feature/42-foo", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ticketFromBranch(tt.branch); got != tt.want {
				t.Errorf("ticketFromBranch(%q) = %q, want %q", tt.branch, got, tt.want)
			}
		})
	}
}
