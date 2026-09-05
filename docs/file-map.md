# File Map

> File-level dependency map for alfred-note-md-template.
> Add to this as you explore the codebase during development.

## Entry Points

| File | Role |
|---|---|
| `cmd/note-md-template-alfred/main.go` | Alfred's Script Filter node executes this — lists/filters templates |
| `cmd/note-md-template-paste-alfred/main.go` | Alfred's Run Script action executes this — pastes the selected template |

## Call Flow

```
cmd/note-md-template-alfred/main.go
  └─ internal/templatelist.List(query)        ← lists/filters .md files in templates_dir
       └─ internal/scriptfilter.Response.Write()

cmd/note-md-template-paste-alfred/main.go
  └─ internal/mdtemplate.Parse(path)           ← parses the template into blocks
       └─ internal/paste.Play(blocks, clipboard, keyboard)
            ├─ internal/clipboard (real impl: pbcopy / osascript)
            └─ internal/keystroke (real impl: osascript / System Events)
```

## Module Dependency Table

### `internal/`

| File | Imports from | Notes |
|---|---|---|
| `mdtemplate/mdtemplate.go` | stdlib only | Parses a `.md` file into `TextBlock`/`ImageBlock`/`CaptionBlock` values |
| `templatelist/templatelist.go` | stdlib only | Lists/filters `.md` files under `templates_dir` |
| `paste/paste.go` | none (depends on `Clipboard`/`Keyboard` interfaces only) | Sequences clipboard writes + keystrokes per block; tested with fakes |
| `clipboard/clipboard.go` | stdlib (`os/exec`) | Real `Clipboard` implementation: `pbcopy` (text), `osascript` (image) |
| `keystroke/keystroke.go` | stdlib (`os/exec`) | Real `Keyboard` implementation: `osascript`/System Events |
| `scriptfilter/scriptfilter.go` | stdlib only | Script Filter JSON types + writer |

### `cmd/`

| File | Imports from | Notes |
|---|---|---|
| `note-md-template-alfred/main.go` | `internal/templatelist`, `internal/scriptfilter` | Alfred boundary for the Script Filter node |
| `note-md-template-paste-alfred/main.go` | `internal/mdtemplate`, `internal/paste`, `internal/clipboard`, `internal/keystroke` | Alfred boundary for the Run Script action |

## Key Files for Customization

| File | What to change |
|---|---|
| `workflow/info.plist` | `bundleid`, keyword, UIDs, category, description |
| `internal/mdtemplate/mdtemplate.go` | Template block-parsing rules |
| `internal/paste/paste.go` | Paste sequencing/timing logic |
| `go.mod` | Module path |
| `workflow/icon.png` | Workflow icon |
