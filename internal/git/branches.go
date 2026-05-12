package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func LocalBranches() ([]string, error) {
	cmd := exec.Command(
		"git",
		"for-each-ref",
		"--sort=-committerdate",
		"--format=%(refname:short)",
		"refs/heads",
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	lines := strings.Split(string(out), "\n")

	branches := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line != "" {
			branches = append(branches, line)
		}
	}

	return branches, nil
}

func CreateAndSwitchBranch(branch string) error {
	cmd := exec.Command("git", "switch", "-c", branch)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
