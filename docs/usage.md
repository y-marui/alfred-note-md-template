# Usage

## Quick Start

1. Create `.md` files in your templates folder (default `~/Documents/Note Templates`)
   using standard Markdown image syntax:

   ```markdown
   Intro paragraph.

   ![caption](./images/photo.png)

   *This is the image's caption.*

   More text.
   ```

2. Open note.com's editor in your browser.
3. In Alfred, type `note` followed by a space.

## Commands

```
note
note <query>
```

Lists `.md` files in your templates folder, filtered by `<query>` (matched
against the filename, case-insensitively). Press Enter on a result to paste
that template — text, images, and captions — into the frontmost app (your
note.com editor).

## Template format

- Plain text between images is pasted as-is.
- `![alt](path)` is pasted as an image. Relative paths resolve against the
  template file's own directory.
- An italic line (`*text*` or `_text_`) immediately following an image
  becomes that image's caption in note.com — pasted, then confirmed with
  Return.

## Configuration

Alfred Preferences → Workflows → this workflow → Configure Workflow:

| Variable | Default | Description |
|---|---|---|
| `templates_dir` | `~/Documents/Note Templates` | Directory containing your `.md` template files. |

## Tips

- Accessibility permission: the first time you use `note`, macOS may prompt
  you to grant Alfred Accessibility permission — required for the simulated
  Cmd+V / Return keystrokes that drive the paste.
- Keep the note.com editor focused before pressing Enter on a result; the
  paste plays into whichever app is frontmost.

## Troubleshooting

**No results appear / "Templates directory not found"**
- Confirm `templates_dir` (or the default `~/Documents/Note Templates`) exists and contains `.md` files.
- Check Alfred's debugger: open Alfred → ⌘D

**Paste doesn't happen / nothing types**
- Check that Alfred has Accessibility permission: System Settings → Privacy & Security → Accessibility.
- Make sure the target app (note.com editor) is frontmost before the paste starts.
