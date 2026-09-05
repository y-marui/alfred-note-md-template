# Contributing

Thank you for contributing!

## Before you start

- Check existing issues and PRs to avoid duplicate work.
- For large changes, open an issue first to discuss the approach.

## Development setup

```bash
git clone https://github.com/y-marui/alfred-note-md-template
cd alfred-note-md-template
go build ./...
```

See [DEVELOPING.md](DEVELOPING.md) for the full development workflow, naming
conventions, and code review guidelines.

## Making changes

1. Create a branch: `git checkout -b feat/my-feature`
2. Make your changes
3. Run checks:

```bash
make lint
make test
make build
```

4. Test in Alfred: `make build-workflow` → double-click the `.alfredworkflow`
5. Open a PR using the template

## Code style

- `gofmt` + `go vet` enforced by CI
- Keep `go.mod` dependency-free unless clearly justified

## Commit guidelines

- Commit per **feature unit**, after confirming it works.
- **No WIP commits** — do not commit code that does not run.

### Commit message format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add clipboard copy action
fix: cache miss on special characters in query
chore: update Go toolchain to 1.28
docs: add examples to usage.md
refactor: simplify paste sequencing logic
```

## Pull Request checklist

See [.github/PULL_REQUEST_TEMPLATE.md](.github/PULL_REQUEST_TEMPLATE.md) for the current checklist.
