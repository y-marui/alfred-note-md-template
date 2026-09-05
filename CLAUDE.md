@AI_CONTEXT.md

## Adding a new internal package

1. Add pure logic (no Alfred/OS dependency) under `internal/<name>/`, tested with `go test`.
2. If it needs to build a Script Filter response, depend on `internal/scriptfilter`.
3. Wire it up from `cmd/note-md-template-alfred/main.go` or `cmd/note-md-template-paste-alfred/main.go` — those files hold no business logic themselves.

## Architecture rules

- `cmd/note-md-template-alfred/main.go` is the only file Alfred's Script Filter node executes. No business logic here — it dispatches to `internal/templatelist` and writes the response.
- `cmd/note-md-template-paste-alfred/main.go` is the only file Alfred's Run Script action node executes. It parses the selected template and plays it via `internal/paste`.
- `internal/` packages never depend on Alfred environment variables directly except where documented (`internal/templatelist` reads `templates_dir`).
- All Script Filter output goes through `internal/scriptfilter.Response.Write()`.
- `main()` in both binaries wraps its logic so an unhandled panic still produces visible output, not a blank Alfred result (see `note-md-template-alfred`'s `dispatch`).

## Testing rules

- Test `internal/mdtemplate`, `internal/templatelist`, and `internal/paste` — pure logic, no real Alfred environment needed.
- `internal/clipboard` tests exercise the real macOS pasteboard (`pbcopy`/`pbpaste`, `osascript`); they skip themselves outside macOS.
- `internal/paste` is tested via injected `Clipboard`/`Keyboard` interfaces — never invoke `osascript` from a test.

## Performance target

Script Filter response < 100ms. A compiled binary meets this with room to spare; avoid adding I/O to the hot path (`internal/templatelist.List`) beyond the one `os.ReadDir` call.

## Dependency management

Keep `go.mod` dependency-free unless a third-party package is clearly justified. Every dependency adds to workflow size and startup time.

## Pre-coding checklist

Before starting work, confirm if any of these are unclear:

- Goal / completion criteria
- Language / framework / version constraints
- New code vs existing code modification
- Whether tests are required
- Scope of impact

Do **not** ask about code style, file placement, or minor implementation
details — follow existing patterns in the codebase.

## Error handling stance

When an error occurs: **diagnose root cause → explain fix plan → implement**.
Never retry the same failing command. Never skip hooks (`--no-verify`).
