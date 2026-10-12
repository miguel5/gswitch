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

var version = "dev"

var createBranch bool
var createUpstream bool

var rootCmd = &cobra.Command{
	Use:     "gswitch [pattern]",
	Short:   "Interactive git branch switcher",
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if createUpstream && !createBranch {
			return fmt.Errorf("--upstream requires --create")
		}

		if createBranch {
			if len(args) == 0 {
				return fmt.Errorf("branch name is required")
			}

			rawName := args[0]

			normalized := normalizeBranchName(rawName)

			fullBranchName := cfg.BranchPrefix + normalized

			fmt.Printf("Creating branch %s\n", fullBranchName)

			if err := git.CreateAndSwitchBranch(fullBranchName); err != nil {
				return err
			}

			if createUpstream {
				return git.PushAndSetUpstream(fullBranchName)
			}

			return nil
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

	rootCmd.Flags().BoolVarP(
		&createUpstream,
		"upstream",
		"u",
		false,
		"Push the new branch and set its upstream (requires --create)",
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
