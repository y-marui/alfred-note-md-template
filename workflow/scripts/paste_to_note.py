#!/usr/bin/env python3
"""Paste a markdown template into note.com editor.

Called as a Run Script action in Alfred after the user selects a template.
argv[1] is the absolute path to the selected .md template file.

Requires pyobjc-framework-Cocoa (available via system Python on macOS,
or invoke with: uv run --with pyobjc-framework-Cocoa paste_to_note.py <path>).

Paste sequence per block:
  - TextBlock   : set clipboard to plain text -> Cmd+V
  - ImageBlock  : set clipboard to TIFF data  -> Cmd+V -> sleep
  - CaptionBlock: set clipboard to plain text -> Cmd+V -> Return -> Return
"""

from __future__ import annotations

import subprocess
import sys
import time
from pathlib import Path

# ---------------------------------------------------------------------------
# Path setup — same pattern as entry.py
# ---------------------------------------------------------------------------
_workflow_root = Path(__file__).resolve().parent.parent
for _p in (str(_workflow_root / "vendor"), str(_workflow_root / "src")):
    if _p not in sys.path:
        sys.path.insert(0, _p)

from app.services.template_parser import CaptionBlock, ImageBlock, TextBlock, parse  # noqa: E402

# ---------------------------------------------------------------------------
# Clipboard helpers
# ---------------------------------------------------------------------------


def _set_text(text: str) -> None:
    from AppKit import NSPasteboard  # type: ignore[import]

    pb = NSPasteboard.generalPasteboard()
    pb.clearContents()
    pb.declareTypes_owner_(["public.utf8-plain-text"], None)
    pb.setString_forType_(text, "public.utf8-plain-text")


def _set_image(path: Path) -> None:
    from AppKit import NSImage, NSPasteboard  # type: ignore[import]

    img = NSImage.alloc().initWithContentsOfFile_(str(path))
    if img is None:
        raise FileNotFoundError(f"Cannot load image: {path}")
    tiff_data = img.TIFFRepresentation()
    pb = NSPasteboard.generalPasteboard()
    pb.clearContents()
    pb.declareTypes_owner_(["public.tiff"], None)
    pb.setData_forType_(tiff_data, "public.tiff")


def _paste() -> None:
    subprocess.run(
        ["osascript", "-e", 'tell application "System Events" to keystroke "v" using command down'],
        check=True,
    )
    time.sleep(0.8)


def _press_enter() -> None:
    subprocess.run(
        ["osascript", "-e", 'tell application "System Events" to keystroke return'],
        check=True,
    )
    time.sleep(0.3)


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------


def main() -> None:
    if len(sys.argv) < 2:
        print("Usage: paste_to_note.py <template_path>", file=sys.stderr)
        sys.exit(1)

    template_path = Path(sys.argv[1])
    if not template_path.exists():
        print(f"Template not found: {template_path}", file=sys.stderr)
        sys.exit(1)

    # Allow Alfred window to close and note.com editor to regain focus
    time.sleep(0.5)

    blocks = parse(template_path)

    for block in blocks:
        if isinstance(block, TextBlock):
            _set_text(block.text)
            _paste()
        elif isinstance(block, ImageBlock):
            _set_image(block.path)
            _paste()
            time.sleep(0.5)
        elif isinstance(block, CaptionBlock):
            _set_text(block.text)
            _paste()
            _press_enter()
            _press_enter()


if __name__ == "__main__":
    main()
