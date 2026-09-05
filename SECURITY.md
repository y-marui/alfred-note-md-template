# Security Policy

## Supported Versions

Only the latest release is supported with security fixes.

## Reporting a Vulnerability

Please **do not** open a public GitHub issue for security vulnerabilities.

Instead, report them privately via
[GitHub Security Advisories](https://github.com/y-marui/alfred-note-md-template/security/advisories/new)
or email the maintainer directly.

We aim to acknowledge reports within 48 hours and provide a fix within 7 days
for confirmed vulnerabilities.

## Scope

This workflow reads local Markdown template files and writes them to the
macOS clipboard. Common areas of concern:

- **Credential handling** — never store secrets in `workflow/info.plist` or
  committed files; use Alfred's built-in encrypted keychain instead.
- **Input handling** — the Alfred query and the selected template path are
  passed to `cmd/note-md-template-alfred` and `cmd/note-md-template-paste-alfred`
  as plain arguments; they must not be interpolated into shell commands or
  AppleScript without proper quoting (see `internal/clipboard`'s AppleScript
  string escaping).
- **Clipboard/keystroke automation** — `internal/keystroke` simulates Cmd+V
  and Return via `osascript`/System Events; this requires the user to grant
  Accessibility permission to Alfred, and only ever acts on content the user
  selected.

## Automated security checks

| Hook | What it detects |
|---|---|
| `gitleaks` (`.gitleaks.toml`) | Hardcoded secrets, API keys, local absolute paths |
| `detect-private-key` | SSH/TLS private key headers |
| `no-commit-dotenv` | `.env` files accidentally staged |
| `check-added-large-files` | Files over 500 KB |

These hooks run on every commit (pre-commit) and in CI (`security` job).
