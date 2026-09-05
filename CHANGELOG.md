# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Rewrote the implementation from Python to Go (#30). The workflow now
  packages two universal (amd64+arm64) binaries
  (`note-md-template-alfred`, `note-md-template-paste-alfred`) instead of a
  vendored Python runtime; no user-facing behavior change. See
  `docs/decisions/0001-go-reimplementation.md`.
- Dropped the unused generic scaffold commands inherited from
  `alfred-workflow-template` (`search`, `open`, `config`, `help`) — they
  were never wired into this workflow's Alfred package.

## [0.1.0] - 2024-01-01

### Added

- Initial release, based on the `alfred-workflow-template` Python scaffold
- `note` command: list and filter `.md` templates for note.com paste
- GitHub Actions CI (lint, test, build)
- GitHub Actions Release (tag → `.alfredworkflow` → GitHub Release)

[Unreleased]: https://github.com/y-marui/alfred-note-md-template/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/y-marui/alfred-note-md-template/releases/tag/v0.1.0
