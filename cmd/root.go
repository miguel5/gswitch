package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"

	"github.com/miguel5/gswitch/internal/config"
	"github.com/miguel5/gswitch/internal/git"
)

var createBranch bool

var rootCmd = &cobra.Command{
	Use:   "gswitch [pattern]",
	Short: "Interactive git branch switcher",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if createBranch {
			if len(args) == 0 {
				return fmt.Errorf("branch name is required")
			}

			rawName := args[0]

			normalized := normalizeBranchName(rawName)

			fullBranchName := cfg.BranchPrefix + normalized

			fmt.Printf("Creating branch %s\n", fullBranchName)

			return git.CreateAndSwitchBranch(fullBranchName)
		}

		query := ""

		if len(args) > 0 {
			query = strings.ToLower(args[0])
		}

		branches, err := git.LocalBranches()
		if err != nil {
			return err
		}

		filtered := make([]string, 0)

		for _, branch := range branches {
			if query == "" || strings.Contains(strings.ToLower(branch), query) {
				filtered = append(filtered, branch)
			}
		}

		if len(filtered) == 0 {
			return fmt.Errorf("no matching branches found")
		}

		if len(filtered) == 1 {
			return switchBranch(filtered[0])
		}

		idx, err := fuzzyfinder.Find(
			filtered,
			func(i int) string {
				return filtered[i]
			},
		)

		if err != nil {
			return err
		}

		return switchBranch(filtered[idx])
	},
}

func init() {
	rootCmd.Flags().BoolVarP(
		&createBranch,
		"create",
		"c",
		false,
		"Create a new branch",
	)
}

func switchBranch(branch string) error {
	fmt.Printf("Switching to %s\n", branch)

	c := exec.Command("git", "switch", branch)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	return c.Run()
}

func normalizeBranchName(name string) string {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)

	re := regexp.MustCompile(`[^a-z0-9]+`)

	name = re.ReplaceAllString(name, "-")

	name = strings.Trim(name, "-")

	return name
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
