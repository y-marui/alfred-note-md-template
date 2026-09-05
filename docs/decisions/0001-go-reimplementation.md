# ADR 0001: Reimplement in Go

## Status

Accepted on 2026-09-05.

## Decision

Rewrite this workflow from the Python `alfred-workflow-template` scaffold to
Go, mirroring the `cmd/`+`internal/` layout already used by
`alfred-clean-invisible-text`, `alfred-markdown-ref`, and
`alfred-password-generator`. Drop the unused generic scaffold commands
(`search`, `open`, `config`, `help`, `example_service`, `api_client`) rather
than port them — the packaged workflow only ever wired up the `note`
Script Filter node, so those commands were unreachable dead code, not part
of the shipped product.

Keep the implementation self-contained in this repository rather than
splitting a general-purpose CLI into a separate repo (contrast
`alfred-clean-invisible-text`'s split with `go-clean-invisible-text`):
nothing about template parsing or paste automation is useful outside this
Alfred workflow.

## Rationale

- Consistency: half of the sibling `alfred-*` workflows are already Go-based.
- Distribution: a single static, universal (amd64+arm64) binary avoids
  bundling a Python runtime and vendored wheels into the `.alfredworkflow`
  package, and starts faster on every keystroke in the Script Filter.
- Low migration risk: the Python side had no third-party runtime
  dependencies beyond the in-repo `alfred` SDK, so there was no library to
  find a Go equivalent for.

## Consequences

- The workflow is packaged as two universal binaries
  (`note-md-template-alfred`, `note-md-template-paste-alfred`) instead of a
  Python `src/` tree plus vendored packages — see `docs/architecture.md`.
- `templates_dir` remains the only Config Builder variable; the Python
  version's `use_uv` toggle no longer applies (no interpreter to select).
- The existing pytest suite for `note_cmd` and `template_parser` was the
  behavior specification for the Go rewrite's tests
  (`internal/templatelist`, `internal/mdtemplate`); it was retired along
  with the Python source. The paste automation (`internal/paste`,
  `internal/clipboard`, `internal/keystroke`) is new: the Python version's
  `paste_to_note.py` had no test coverage of its own.

## Alternatives considered (post-migration)

Whether `internal/clipboard`/`internal/keystroke`'s `osascript` calls could
be replaced by Alfred's own native output nodes
(`alfred.workflow.output.clipboard`'s `autopaste`,
`alfred.workflow.output.dispatchkeycombo`) was investigated separately and
rejected:

- `alfred.workflow.output.clipboard`'s `clipboardtext` is text-only — Alfred
  has no native output node for writing image data to the pasteboard, so
  `internal/clipboard.WriteImage` (image blocks) has no native equivalent
  regardless of the text-block decision.
- A template's block count is variable (0..N images), and Alfred's
  workflow object schema (`docs/alfred-workflow-notes/workflow-object-schema.md`)
  has no native loop construct — `alfred.workflow.utility.junction` is a
  branch/merge point, not iteration. The graph can't express "paste N
  blocks in sequence" statically.
- A Go binary can't synchronously drive a native output node mid-run and
  wait for it to finish before continuing — `alfred.workflow.output.callexternaltrigger`
  is fire-and-forget, with no completion signal back to the caller.
- Accessibility permission is required by Alfred either way (native
  Dispatch Key Combo or `osascript` System Events), so there's no
  permission-scope benefit to splitting keystrokes out natively.

Given the paste automation's core requirement — sequencing an
Alfred-unaware number of text/image/caption blocks — stays in
`internal/paste` either way, keeping clipboard and keystroke actions there
too was simpler than a partial split with no clear benefit.
