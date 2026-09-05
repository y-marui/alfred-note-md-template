# UI Design

Alfred Script Filter workflows present results as a list of items in the Alfred
launcher. This document defines the UI conventions for result items in this
workflow.

## Result Item Structure

Alfred result items are JSON objects with the following fields used in this workflow:

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string | yes | Primary text (large, always visible) — the template's filename |
| `subtitle` | string | no | Secondary text (small, below title) — the template's absolute path |
| `arg` | string | no | The template's absolute path; passed to the Run Script paste action on Enter |
| `uid` | string | no | The template's absolute path; used by Alfred for learned ordering |
| `valid` | bool | yes | If false, Enter does not trigger an action |

## Text Guidelines

### No Unicode Emoji in `title` / `subtitle`

- **Prohibited:** `🔍 Search`, `✅ Done`, `📄 Document`
- **Allowed:** ASCII symbols — `>`, `*`, `[x]`, `(!)`, `--`
- **Reason:** Emoji rendering is inconsistent across Alfred versions and macOS
  updates. ASCII symbols are universally stable.

### Empty / Error States

- No templates match the query → an informative item (e.g. `"No matching
  templates"`) with `valid: false`.
- `templates_dir` missing or unreadable → `"Templates directory not found"`
  with `valid: false` — see [docs/usage.md](usage.md) Troubleshooting.
- Error in the Script Filter step → panic recovery automatically shows a
  `"Workflow Error"` item; do not hide errors silently. The paste action has
  no JSON contract (a Run Script action), so its failures surface via
  Alfred's debugger (⌘D) instead.

## Icon

- Workflow icon: `workflow/icon.png` (PNG, any size — Alfred scales it).
- Alfred controls light/dark mode; do not ship separate light/dark icons.
- No per-item icons are used in this workflow.

## Keyboard Shortcuts

These are standard Alfred behaviors — do not override them in the workflow:

| Key | Action |
|---|---|
| ↩ Enter | Run the paste action with the selected template's path |
| ⌘C | Copy `arg` (the template's path) to clipboard |
| ⌘L | Show `title` in Large Type |

## Layout Conventions

### `note` template list

```
title:    <filename without extension>
subtitle: <absolute path>
arg:      <absolute path>
uid:      <absolute path>
valid:    true
```

### No-match / error rows

```
title:    "No matching templates" | "Templates directory not found" | "Workflow Error"
subtitle: <detail, if any>
valid:    false
```
