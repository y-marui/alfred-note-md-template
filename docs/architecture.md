# Architecture

## Overview

This workflow is split into two independent binaries — one per Alfred node
— backed by shared, Alfred-independent internal packages.

```
Alfred
  │  keyword "note" + query
  ▼
cmd/note-md-template-alfred/main.go     ← Script Filter binary
  │
  ▼
internal/templatelist                   ← list + filter .md templates
  │
  ▼
internal/scriptfilter                   ← Script Filter JSON response

Alfred (user selects a result)
  │  selected template's absolute path
  ▼
cmd/note-md-template-paste-alfred/main.go  ← Run Script action binary
  │
  ▼
internal/mdtemplate                     ← parse template into blocks
  │
  ▼
internal/paste                          ← sequence clipboard + keystrokes
  │                                        per block
  ├─ internal/clipboard                 ← write text/image to pasteboard
  └─ internal/keystroke                 ← simulate Cmd+V / Return
```

## Packages

| Package | Purpose |
|---|---|
| `internal/mdtemplate` | Parses a `.md` template file into an ordered list of `TextBlock`/`ImageBlock`/`CaptionBlock` values. Pure logic, no I/O beyond reading the file. |
| `internal/templatelist` | Lists `.md` files under the `templates_dir` Config Builder variable (default `~/Documents/Note Templates`), filtered by the Alfred query. |
| `internal/paste` | Given parsed blocks, sequences calls to a `Clipboard` and `Keyboard` interface — text is pasted as-is, images are pasted then given extra time to render, captions are pasted and confirmed with two Returns. Tested with fakes; never shells out itself. |
| `internal/clipboard` | Writes plain text (`pbcopy`) or an image file (`osascript` coercing to `TIFF picture`) to the general pasteboard. |
| `internal/keystroke` | Simulates Cmd+V and Return via `osascript`/System Events, with the settle delays note.com's editor needs between actions. |
| `internal/scriptfilter` | Alfred Script Filter JSON types (`Item`, `Response`) and the writer that encodes them to stdout. |

## Binaries

| Binary | Alfred node | Role |
|---|---|---|
| `cmd/note-md-template-alfred` | Script Filter (keyword `note`) | Lists/filters templates, writes a Script Filter response. Wraps `internal/templatelist.List` in a `recover()` so an unhandled panic still shows a visible error item instead of blank output. |
| `cmd/note-md-template-paste-alfred` | Run Script action | Parses the selected template and plays it into the frontmost app via `internal/paste`, using the real `internal/clipboard` and `internal/keystroke` implementations. |

Neither binary contains business logic beyond argument handling and wiring
— see `CLAUDE.md`'s Architecture rules.

## Packaging

At build time (`make build-workflow` / `scripts/build-workflow.sh`):

```
.build/                                   ← temporary build directory
├── info.plist                            ← copied from workflow/
├── icon.png
├── note-md-template-alfred               ← universal (amd64+arm64) binary
└── note-md-template-paste-alfred         ← universal (amd64+arm64) binary
```

Both binaries are built for `darwin/amd64` and `darwin/arm64` and merged
into a single universal binary with `lipo`, so the packaged workflow runs
natively on both Intel and Apple Silicon Macs without needing a runtime
interpreter or vendored dependencies. The entire `.build/` directory is
zipped to `dist/<name>-<version>.alfredworkflow`.
