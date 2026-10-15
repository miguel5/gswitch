# gswitch

Interactive git branch switcher with fuzzy search.

## Install

`git` must be on your `PATH` at runtime.

### Homebrew (macOS/Linux)

```bash
brew install miguel5/tap/gswitch
```

### Prebuilt binaries

Download the archive for your OS/arch from the
[latest release](https://github.com/miguel5/gswitch/releases/latest), extract it,
and move `gswitch` somewhere on your `PATH`.

### With Go

Requires [Go](https://go.dev/dl/).

```bash
go install github.com/miguel5/gswitch@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`.

### Build from source

```bash
git clone https://github.com/miguel5/gswitch.git
cd gswitch
go build -o gswitch .
```

Move the `gswitch` binary somewhere on your `PATH`.

## Usage

```bash
gswitch
```

Open fuzzy finder with all branches.

```bash
gswitch auth
```

Filter branches matching "auth".

```bash
gswitch -l auth
```

List branches matching "auth", one per line, without switching. Run
`gswitch -l` with no pattern to list every local branch.

```bash
gswitch -c new-feature
```

Create a new branch named `feature/new-feature` and switch to it.

```bash
gswitch -cu new-feature
```

Create and switch to `feature/new-feature`, then push it to the default remote
and set its upstream. `-u` must be combined with `-c`; the remote is `origin`
when present, otherwise the only configured remote.

```bash
gswitch -m "fix login redirect"
```

Commit staged changes with the current branch's ticket prefixed automatically.
On a branch like `feature/ticket-42-login` this runs
`git commit -m "TICKET-42: fix login redirect"`. If no `PREFIX-123` ticket is
found in the branch name, the message is committed as-is.