"""Markdown template parser for note.com paste workflow.

Parses a markdown file into an ordered list of TextBlock, ImageBlock, and
CaptionBlock objects.

Images are identified by standard markdown image syntax: ![alt](path)
A CaptionBlock is produced when an italic line (*text* or _text_) immediately
follows an image — note.com treats the first text after a pasted image as the
image caption.

Images with relative paths are resolved relative to the template file's directory.
"""

from __future__ import annotations

import re
from dataclasses import dataclass
from pathlib import Path
from typing import Union


@dataclass
class TextBlock:
    text: str


@dataclass
class ImageBlock:
    path: Path
    alt: str


@dataclass
class CaptionBlock:
    text: str


Block = Union[TextBlock, ImageBlock, CaptionBlock]

_IMAGE_PATTERN = re.compile(r"!\[([^\]]*)\]\(([^)]+)\)")
# Matches a single line wrapped in * or _ (italic syntax)
_ITALIC_LINE = re.compile(r"^\*(.+)\*$|^_(.+)_$")


def parse(template_path: Path) -> list[Block]:
    """Parse a markdown template file into a list of content blocks.

    Args:
        template_path: Path to the .md template file.

    Returns:
        Ordered list of TextBlock, ImageBlock, and CaptionBlock objects.
    """
    content = template_path.read_text(encoding="utf-8")
    return _split_blocks(content, template_path.parent)


def _split_blocks(content: str, base_dir: Path) -> list[Block]:
    blocks: list[Block] = []
    last_end = 0
    last_was_image = False

    for match in _IMAGE_PATTERN.finditer(content):
        text_before = content[last_end : match.start()]

        if last_was_image:
            blocks.extend(_parse_post_image_text(text_before))
        elif text_before.strip():
            blocks.append(TextBlock(text=text_before.strip("\n")))

        alt = match.group(1)
        raw_path = match.group(2).strip()
        img_path = _resolve_path(raw_path, base_dir)
        blocks.append(ImageBlock(path=img_path, alt=alt))
        last_end = match.end()
        last_was_image = True

    remaining = content[last_end:]
    if last_was_image:
        blocks.extend(_parse_post_image_text(remaining))
    elif remaining.strip():
        blocks.append(TextBlock(text=remaining.strip("\n")))

    return blocks


def _parse_post_image_text(text: str) -> list[Block]:
    """Parse text that immediately follows an image.

    If the first non-empty line is italic (*text* or _text_), it becomes a
    CaptionBlock.  Any remaining text becomes a TextBlock.
    """
    stripped = text.lstrip("\n")
    if not stripped.strip():
        return []

    first_line, _, rest = stripped.partition("\n")
    m = _ITALIC_LINE.match(first_line.strip())
    if m:
        caption = (m.group(1) or m.group(2)).strip()
        result: list[Block] = [CaptionBlock(text=caption)]
        if rest.strip():
            result.append(TextBlock(text=rest.strip("\n")))
        return result

    return [TextBlock(text=stripped.strip("\n"))]


def _resolve_path(raw: str, base_dir: Path) -> Path:
    p = Path(raw).expanduser()
    if p.is_absolute():
        return p
    return (base_dir / p).resolve()
