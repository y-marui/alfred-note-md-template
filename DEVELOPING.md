# Developing

This document covers the development workflow, conventions, and guidelines for contributors to this project.

## Development Setup

```bash
git clone https://github.com/y-marui/alfred-note-md-template
cd alfred-note-md-template
go build ./...
```

**Prerequisites:**
- macOS (required for Alfred, the pasteboard, and `osascript` keystroke automation)
- Go (see `go.mod` for the toolchain version)
- Alfred 5 with Powerpack
- `gh` CLI (required for releases): `brew install gh`

## Development Workflow

### Daily commands

```bash
templates_dir=~/Documents/Note\ Templates go run ./cmd/note-md-template-alfred ""       # list all templates
templates_dir=~/Documents/Note\ Templates go run ./cmd/note-md-template-alfred "review"  # filter by query
go run ./cmd/note-md-template-paste-alfred /path/to/template.md                          # paste a template
make test                 # run test suite
make lint                 # gofmt -l + go vet
make fmt                  # gofmt -w (auto-fix)
make build                # go build ./...
make build-workflow       # build dist/*.alfredworkflow
```

Pipe the Script Filter binary through `jq` for pretty-printed JSON:

```bash
go run ./cmd/note-md-template-alfred "" | jq .
```

### Testing in Alfred

1. `make build-workflow` — generates `dist/*.alfredworkflow`
2. Double-click the `.alfredworkflow` file to install in Alfred
3. Open Alfred, type `note`, select a template, and confirm the paste

`cmd/note-md-template-paste-alfred` drives real clipboard writes and
simulated keystrokes — always verify it against the actual Alfred +
note.com editor rather than from `go run` output alone.

## Naming Conventions

| Scope | Convention | Example |
|---|---|---|
| Go packages | short, lowercase, no underscores | `mdtemplate`, `templatelist`, `scriptfilter` |
| Exported functions / types | `PascalCase` | `Parse`, `List`, `Response`, `Item` |
| Unexported functions / variables | `camelCase` | `resolvePath`, `templatesDir` |
| Alfred variable names | `lowercase_with_underscores` | `templates_dir` |
| Commit messages | Conventional Commits | `feat:`, `fix:`, `docs:`, `chore:` |
| Branch names | `feat/`, `fix/`, `docs/`, `chore/` | `feat/add-open-browser` |

## Code Style

- **Formatter:** `gofmt`. CI enforces this (`make lint`).
- **Linter:** `go vet`.
- **Comments:** Write *why*, not *what*. Do not comment self-evident code.
- **No debug prints:** Remove all stray `fmt.Print*` statements before committing;
  the only writer to stdout is `scriptfilter.Response.Write`.
- **No third-party dependencies** unless clearly justified — keep `go.mod` dependency-free.

## Commit Guidelines

- Commit per **feature unit**, after confirming it works.
- **No WIP commits** — do not commit code that does not run.
- **No `--no-verify`** — never skip pre-commit hooks.

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add clipboard copy action
fix: cache miss on special characters in query
chore: update Go toolchain to 1.28
docs: update README usage section
refactor: simplify paste sequencing logic
```

## Pull Request Checklist

- [ ] `make lint` passes
- [ ] `make test` passes
- [ ] `make build-workflow` succeeds
- [ ] New packages/behavior have tests
- [ ] `README.md` updated if user-facing changes
- [ ] `CHANGELOG.md` entry added under `[Unreleased]`

## Code Review Guidelines

**Reviewers check for:**
- Architectural constraints respected (no business logic in `cmd/note-md-template-alfred` or `cmd/note-md-template-paste-alfred`)
- No hardcoded absolute paths (use `$HOME` / env vars)
- No debug prints in production code
- No Unicode emoji in Alfred result item `title` / `subtitle`
- Tests cover the new or changed behavior
- Alfred env variables managed via Config Builder, not `variables` key

**Security-sensitive changes** (clipboard/keystroke automation, file access) require
explicit security review before merge.

**Self-review:** Individual contributors open a PR and self-review before merging
to `main`.
