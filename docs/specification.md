# Specification

> Functional specification, behavior definition, and data flow for
> alfred-note-md-template. For the end-user guide, see
> [docs/usage.md](usage.md).

## Overview

This workflow is an Alfred 5 Script Filter that lists/filters `.md`
template files, then a Run Script action that plays the selected
template's content (text, images, captions) into the frontmost app via
simulated clipboard writes and keystrokes.

## Commands

### `note` (the only command)

**Trigger:** `note` or `note <query>`

**Behavior:**
1. `internal/templatelist.List` lists `.md` files under the `templates_dir`
   Config Builder variable (default `~/Documents/Note Templates`).
2. Filters by `<query>` matched against the filename, case-insensitively.
   An empty query returns every template.
3. If no templates match → display an informative item with `valid: false`.
4. Otherwise → one result item per matching file; pressing Enter passes its
   absolute path as `arg` to the Run Script action.

**Result item fields:**

| Field | Source | Notes |
|---|---|---|
| `title` | filename (without extension) | Primary display text |
| `subtitle` | absolute path of the `.md` file | Secondary display text |
| `arg` | absolute path of the `.md` file | Passed to the Run Script action on Enter |
| `uid` | absolute path of the `.md` file | Used by Alfred for learned ordering |

### Paste action (Run Script, not directly invoked by a keyword)

**Trigger:** pressing Enter on a `note` result.

**Behavior:**
1. `internal/mdtemplate.Parse` reads the selected template file and splits
   it into an ordered list of `TextBlock`/`ImageBlock`/`CaptionBlock` values.
2. `internal/paste.Play` sequences each block into the frontmost app:
   - `TextBlock` — pasted as-is.
   - `ImageBlock` — pasted, then given extra settle time to render.
   - `CaptionBlock` — pasted, then confirmed with two Returns (note.com's
     editor requires this to attach a caption to the preceding image).
3. Clipboard writes and keystrokes are the real `internal/clipboard`/
   `internal/keystroke` implementations — `pbcopy`/`osascript` for text and
   image clipboard writes, `osascript`/System Events for Cmd+V/Return.

## Data Flow

```
Alfred input (keyword "note" + query string)
  │
  ▼
cmd/note-md-template-alfred/main.go        reads os.Args[1]
  │
  ▼
internal/templatelist.List(query)          filters .md files in templates_dir
  │
  ▼
internal/scriptfilter.Response.Write()     prints JSON to stdout → Alfred renders result list
  │
  ▼ (user presses Enter — arg is the template's path)
cmd/note-md-template-paste-alfred/main.go  reads os.Args[1]
  │
  ▼
internal/mdtemplate.Parse(path)            splits the file into blocks
  │
  ▼
internal/paste.Play(blocks, ...)           writes clipboard + sends keystrokes per block
```

## Error Handling

- `templates_dir` missing or unreadable → the Script Filter shows an
  informative "Templates directory not found" item, not an empty result.
- Any panic during dispatch is recovered into a visible Script Filter error
  item; the paste binary has no JSON contract (a Run Script action), so its
  failures are reported to stderr for Alfred's debugger (⌘D) instead.

## Configuration Variables

Managed via Alfred Configuration Builder (see `docs/configuration-builder.md`).

| Variable | Type | Default | Effect |
|---|---|---|---|
| `templates_dir` | file (folder picker) | `~/Documents/Note Templates` | Directory scanned for `.md` template files |

## Constraints

- Script Filter response time target: **< 100 ms**.
- All Script Filter output must go through `scriptfilter.Response.Write()` —
  never `fmt.Print*` directly.
- Neither binary contains business logic beyond argument handling and
  wiring — see `AI_CONTEXT.md`'s Architecture Constraints.
- Pasting requires the target app to be frontmost and macOS Accessibility
  permission granted to Alfred (for simulated keystrokes) — see
  [docs/usage.md](usage.md)'s Troubleshooting section.
