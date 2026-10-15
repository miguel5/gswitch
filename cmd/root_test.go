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

func TestStashRef(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{"wip entry", "stash@{0}: WIP on main: some message", "stash@{0}"},
		{"on entry", "stash@{3}: On feature/foo: message", "stash@{3}"},
		{"message with colons", "stash@{1}: WIP on main: fix: a: b", "stash@{1}"},
		{"no colon", "stash@{2}", "stash@{2}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stashRef(tt.line); got != tt.want {
				t.Errorf("stashRef(%q) = %q, want %q", tt.line, got, tt.want)
			}
		})
	}
}

func TestMatchingBranches(t *testing.T) {
	branches := []string{"main", "feature/Login-Page", "feature/logout", "bugfix/login-crash"}

	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"empty query returns all", "", branches},
		{"case-insensitive substring", "LOGIN", []string{"feature/Login-Page", "bugfix/login-crash"}},
		{"partial match", "log", []string{"feature/Login-Page", "feature/logout", "bugfix/login-crash"}},
		{"exact match", "main", []string{"main"}},
		{"no match", "nope", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchingBranches(branches, tt.query)

			if len(got) != len(tt.want) {
				t.Fatalf("matchingBranches(%q) = %v, want %v", tt.query, got, tt.want)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("matchingBranches(%q)[%d] = %q, want %q", tt.query, i, got[i], tt.want[i])
				}
			}
		})
	}
}
