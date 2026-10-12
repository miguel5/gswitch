# AGENTS.md

## Project
`gswitch` — small Go CLI: interactive git branch switcher with fuzzy search
(`go-fuzzyfinder`). Module `github.com/miguel5/gswitch`, Go 1.26.3.

## Commands
- Build: `go build -o gswitch .`
- Run: `go run . [pattern]` or `go run . -c <name>`
- Test all: `go test ./...`; single test: `go test ./internal/git -run TestName`
- Format: `gofmt -w .`
No Makefile, linter config, CI, or tests exist in this repo.

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
- Viper read errors are ignored, so a missing/broken config silently uses defaults.
- `Execute()` prints the error to stdout and calls `os.Exit(1)`.

## Gotchas
- Every dependency in `go.mod` is marked `// indirect` despite being direct
  imports; `go mod tidy` will (correctly) change this.
