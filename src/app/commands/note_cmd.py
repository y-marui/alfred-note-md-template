"""note command - list markdown templates for note.com paste.

Usage in Alfred:  note
                  note <query>

Templates are read from the directory configured via the
``templates_dir`` Alfred workflow variable (config builder).
Default: ~/Documents/Note Templates
"""

from __future__ import annotations

import os
from pathlib import Path

from alfred.response import item, output

_ENV_KEY = "templates_dir"
_DEFAULT_DIR = Path.home() / "Documents" / "Note Templates"


def handle(args: str) -> None:
    """List .md template files, optionally filtered by query."""
    templates_dir = Path(os.environ.get(_ENV_KEY) or _DEFAULT_DIR)

    if not templates_dir.exists():
        output(
            [
                item(
                    title="Templates directory not found",
                    subtitle=f"Create {templates_dir} and add .md files",
                    valid=False,
                )
            ]
        )
        return

    templates = sorted(templates_dir.glob("*.md"), key=lambda p: p.stem.lower())

    if not templates:
        output(
            [
                item(
                    title="No templates found",
                    subtitle=f"Add .md files to {templates_dir}",
                    valid=False,
                )
            ]
        )
        return

    query = args.strip().lower()
    filtered = [t for t in templates if query in t.stem.lower()] if query else templates

    if not filtered:
        output(
            [
                item(
                    title=f'No templates matching "{args}"',
                    subtitle="Try a different keyword",
                    valid=False,
                )
            ]
        )
        return

    output(
        [
            item(
                title=t.stem,
                subtitle=str(t),
                arg=str(t),
                uid=str(t),
            )
            for t in filtered
        ]
    )
