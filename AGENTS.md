# AGENTS.md

## Project
`gswitch` — small Go CLI: interactive git branch switcher with fuzzy search
(`go-fuzzyfinder`). Module `github.com/miguel5/gswitch`, Go 1.26.3.

## Commands
- Build: `go build -o gswitch .`
- Run: `go run . [pattern]` or `go run . -c <name>`
- Test all: `go test ./...`; single test: `go test ./internal/git -run TestName`
- Format: `gofmt -w .`
- Release check: `goreleaser check`; local dry run: `goreleaser release --snapshot --clean`
No Makefile or linter config; `.github/workflows/ci.yml` runs
gofmt/vet/test/build on pushes and PRs, and `release.yml` (GoReleaser) runs on
`v*` tags.

## Layout
- `main.go` → `cmd.Execute()` only.
- `cmd/root.go` — all CLI logic: arg parsing, branch filtering, fuzzy-find
  prompt, switch/create, branch-name normalization, `-c/--create` and
  `-u/--upstream` flags.
- `internal/git/branches.go` — shells out to the real `git` binary via
  `os/exec` (`git for-each-ref`, `git switch`, `git push`, `git remote`); no
  git library, `git` must be on PATH.
- `internal/config/config.go` — Viper config at
  `$HOME/.config/gswitch/config.yaml`, key `branch_prefix` (default `feature/`).

## Behavior worth knowing before editing
- Branches are pre-filtered case-insensitively by substring; if exactly one
  matches, it switches immediately without showing the fuzzy finder.
- `-c/--create <name>` lowercases, replaces non-alphanumeric runs with `-`,
  trims, then prepends `branch_prefix`.
- `-u/--upstream` (written `gswitch -cu <name>`) requires `-c`; after creating
  it runs `git push --set-upstream <remote> <branch>`. The remote is `origin`
  if present, else the sole configured remote, else it errors.
- `-m/--message <text>` is a standalone mode: it extracts the first
  `PREFIX-123` token from the current branch name (uppercased) and runs
  `git commit -m "<TICKET>: <text>"`. If no ticket is found (or on detached
  `HEAD`) it commits the message as-is with a note to stderr. It cannot be
  combined with `-c`/`-u` or a pattern, and empty messages error.
- `-l/--list [pattern]` is a standalone mode: it prints local branches whose
  name contains `pattern` (case-insensitive, same `matchingBranches` helper as
  the interactive path), one per line, and exits 0 with no output when nothing
  matches. With no pattern it lists every local branch. It cannot be combined
  with `-c`/`-u`/`-m`.
- Viper read errors are ignored, so a missing/broken config silently uses defaults.
- `Execute()` prints the error to stdout and calls `os.Exit(1)`.

## Releasing
- Pushing a `v*` tag triggers `.github/workflows/release.yml` (GoReleaser v2):
  builds darwin/linux/windows × amd64/arm64, publishes a GitHub Release with
  `checksums.txt`, and updates the `miguel5/homebrew-tap` cask.
- Config: `.goreleaser.yaml`. Version comes from `cmd/root.go`'s `version` var,
  injected via ldflags; local builds show `dev`.

## Gotchas
- Homebrew tap publishing needs the `HOMEBREW_TAP_GITHUB_TOKEN` repo secret;
  without it the release job fails at the cask step (binaries still publish).
- Release version ldflags target `github.com/miguel5/gswitch/cmd.version`; renaming
  that package/var breaks `--version` in released binaries.
