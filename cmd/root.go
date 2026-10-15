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
var message string
var listBranches bool
var stash bool

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

		if stash {
			if createBranch || createUpstream || listBranches || cmd.Flags().Changed("message") {
				return fmt.Errorf("--stash cannot be combined with --create/--upstream/--message/--list")
			}

			if len(args) > 0 {
				return fmt.Errorf("--stash cannot be combined with a pattern")
			}

			return browseStashes()
		}

		if listBranches {
			if createBranch || createUpstream || cmd.Flags().Changed("message") {
				return fmt.Errorf("--list cannot be combined with --create/--upstream/--message")
			}

			branches, err := git.LocalBranches()
			if err != nil {
				return err
			}

			query := ""

			if len(args) > 0 {
				query = args[0]
			}

			for _, branch := range matchingBranches(branches, query) {
				fmt.Println(branch)
			}

			return nil
		}

		if cmd.Flags().Changed("message") {
			if createBranch || createUpstream {
				return fmt.Errorf("--message cannot be combined with --create/--upstream")
			}

			if len(args) > 0 {
				return fmt.Errorf("--message cannot be combined with a pattern")
			}

			if message == "" {
				return fmt.Errorf("message cannot be empty")
			}

			return commitWithTicket(message)
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
			query = args[0]
		}

		branches, err := git.LocalBranches()
		if err != nil {
			return err
		}

		filtered := matchingBranches(branches, query)

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

	rootCmd.Flags().StringVarP(
		&message,
		"message",
		"m",
		"",
		"Commit staged changes with the branch's ticket prefix",
	)

	rootCmd.Flags().BoolVarP(
		&listBranches,
		"list",
		"l",
		false,
		"List branches matching the pattern without switching",
	)

	rootCmd.Flags().BoolVarP(
		&stash,
		"stash",
		"s",
		false,
		"Browse stashes and apply or pop the selected one",
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

var ticketRe = regexp.MustCompile(`[A-Za-z]+-[0-9]+`)

func commitWithTicket(msg string) error {
	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	full := msg

	if ticket := ticketFromBranch(branch); ticket != "" {
		full = ticket + ": " + msg
	} else {
		fmt.Fprintf(os.Stderr, "note: no ticket found in branch %q; committing without a prefix\n", branch)
	}

	return git.Commit(full)
}

func ticketFromBranch(branch string) string {
	if branch == "HEAD" {
		return ""
	}

	match := ticketRe.FindString(branch)
	if match == "" {
		return ""
	}

	return strings.ToUpper(match)
}

func browseStashes() error {
	stashes, err := git.StashList()
	if err != nil {
		return err
	}

	if len(stashes) == 0 {
		return fmt.Errorf("no stashes found")
	}

	idx, err := fuzzyfinder.Find(
		stashes,
		func(i int) string {
			return stashes[i]
		},
	)
	if err != nil {
		return err
	}

	ref := stashRef(stashes[idx])

	actions := []string{"apply", "pop"}

	actionIdx, err := fuzzyfinder.Find(
		actions,
		func(i int) string {
			return actions[i]
		},
	)
	if err != nil {
		return err
	}

	switch actions[actionIdx] {
	case "apply":
		fmt.Printf("Applying %s\n", ref)

		return git.StashApply(ref)
	default:
		fmt.Printf("Popping %s\n", ref)

		return git.StashPop(ref)
	}
}

func stashRef(line string) string {
	if i := strings.Index(line, ":"); i != -1 {
		return strings.TrimSpace(line[:i])
	}

	return strings.TrimSpace(line)
}

func matchingBranches(branches []string, query string) []string {
	query = strings.ToLower(query)

	matched := make([]string, 0)

	for _, branch := range branches {
		if query == "" || strings.Contains(strings.ToLower(branch), query) {
			matched = append(matched, branch)
		}
	}

	return matched
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
