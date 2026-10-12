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

func PushAndSetUpstream(branch string) error {
	remote, err := defaultRemote()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "push", "--set-upstream", remote, branch)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func defaultRemote() (string, error) {
	cmd := exec.Command("git", "remote")

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list remotes: %w", err)
	}

	remotes := make([]string, 0)

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)

		if line != "" {
			remotes = append(remotes, line)
		}
	}

	for _, remote := range remotes {
		if remote == "origin" {
			return "origin", nil
		}
	}

	if len(remotes) == 1 {
		return remotes[0], nil
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("no git remotes configured")
	}

	return "", fmt.Errorf("multiple git remotes found; cannot determine default")
}
