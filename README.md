# gswitch

Interactive git branch switcher with fuzzy search.

## Install

Requires [Go](https://go.dev/dl/). `git` must be on your `PATH` at runtime.

Install directly:

```bash
go install github.com/miguel5/gswitch@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`.

Or build from source:

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
gswitch -c new-feature
```

Create a new branch named `feature/new-feature` and switch to it.

```bash
gswitch -cu new-feature
```

Create and switch to `feature/new-feature`, then push it to the default remote
and set its upstream. `-u` must be combined with `-c`; the remote is `origin`
when present, otherwise the only configured remote.