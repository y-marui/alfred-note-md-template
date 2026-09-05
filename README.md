# alfred-note-md-template

> **This is the English (reference) version.**
> For the Japanese canonical version, see [README-jp.md](README-jp.md).

> Paste Markdown templates — including images and captions — into note.com's editor from Alfred.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-note-md-template/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-note-md-template/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-note-md-template/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-note-md-template/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

| Field | Value |
|---|---|
| Target | Alfred 5 Script Filter workflow |
| Team size | Individual to small team (1–3 people) |
| Language | Go, no third-party dependencies |
| License | MIT |
| AI tools | Claude Code / GitHub Copilot / Gemini CLI |

## What it does

Keep your note.com article templates as local `.md` files — with images and
italic captions using standard Markdown syntax. Trigger Alfred with `note`,
pick a template, and it's pasted into note.com's editor block by block:
text, then each image, then its caption.

```markdown
Intro paragraph.

![caption](./images/photo.png)

*This is the image's caption.*

More text.
```

## Setup

1. Put your `.md` templates in a folder (default: `~/Documents/Note Templates`).
2. In Alfred Preferences, open this workflow and click **Configure Workflow**
   to set **Templates Directory** if you're not using the default.
3. Open note.com's editor in your browser.
4. In Alfred, type `note`, pick a template, press Enter.

No interpreter or third-party runtime to install — the workflow ships as a
single compiled binary.

The first paste may prompt macOS to grant Alfred Accessibility permission
(required for the simulated Cmd+V / Return keystrokes that drive the
paste); see [docs/usage.md](docs/usage.md) for details and troubleshooting.

## Features

- Single static Go binary — no vendored runtime, no interpreter startup cost
- Universal (amd64+arm64) build — native on both Intel and Apple Silicon
- Full test suite — `go test`, no Alfred required to run tests
- CI/CD — lint, test, build, and release via GitHub Actions
- AI-ready — `AI_CONTEXT.md` + `CLAUDE.md` for AI assistant context

## Requirements

- Alfred 5 (Powerpack required for Script Filter)
- macOS (for the pasteboard and `osascript` keystroke automation)

## Quick Start (developers)

```bash
git clone https://github.com/y-marui/alfred-note-md-template
cd alfred-note-md-template

# Simulate the Script Filter locally
templates_dir=~/Documents/Note\ Templates go run ./cmd/note-md-template-alfred ""

# Run tests
make test

# Build workflow package
make build-workflow
# → dist/note.com-template-paste-<version>.alfredworkflow
```

Double-click `dist/*.alfredworkflow` to install in Alfred.

## Usage

```
note              list all templates
note <query>      filter templates by name
```

See [docs/usage.md](docs/usage.md) for the full guide.

## Project Structure

```
alfred-note-md-template/
├── cmd/
│   ├── note-md-template-alfred/       # Script Filter binary
│   └── note-md-template-paste-alfred/ # Run Script action binary
├── internal/
│   ├── mdtemplate/     # Markdown template → block parser
│   ├── templatelist/   # list/filter templates
│   ├── paste/          # per-block paste sequencing
│   ├── clipboard/      # pasteboard writes (pbcopy / osascript)
│   ├── keystroke/       # simulated Cmd+V / Return
│   └── scriptfilter/   # Alfred Script Filter JSON types
├── workflow/           # Alfred package (info.plist, icon.png)
├── scripts/            # build-workflow.sh, extract-changelog.sh
└── docs/               # Architecture and usage documentation
```

## Documentation

| Document | Description |
|---|---|
| [DEVELOPING.md](DEVELOPING.md) | Development workflow, naming conventions, code review |
| [docs/architecture.md](docs/architecture.md) | Full architecture and package design |
| [docs/file-map.md](docs/file-map.md) | File-level dependency map |
| [docs/specification.md](docs/specification.md) | Functional specification and data flow |
| [docs/ui-design.md](docs/ui-design.md) | Alfred result item UI conventions |
| [docs/usage.md](docs/usage.md) | End-user usage guide |
| [docs/decisions/](docs/decisions/) | Architecture decision records |

## AI-Assisted Development

This project is configured for AI-assisted development.

| Tool | Role |
|---|---|
| Claude Code | Architecture, large-scale changes, refactoring |
| GitHub Copilot | Bug fixes, small implementation, unit tests |
| Gemini CLI | Documentation management |

See [`AI_CONTEXT.md`](AI_CONTEXT.md) and [`CLAUDE.md`](CLAUDE.md) for session context.

## Release

```bash
# 1. Bump version in workflow/info.plist
# 2. Tag and push
git tag v1.2.3
git push --tags
# GitHub Actions builds .alfredworkflow and creates a GitHub Release
```

## License

MIT — see [LICENSE](LICENSE)

---

*This document has a Japanese canonical version [README-jp.md](README-jp.md). Update both in the same commit when editing.*
