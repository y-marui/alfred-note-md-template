"""Tests for service layer."""

from __future__ import annotations

from pathlib import Path
from textwrap import dedent
from unittest.mock import patch

from app.services.example_service import ExampleService
from app.services.template_parser import CaptionBlock, ImageBlock, TextBlock, parse

_STUB = [{"id": "1", "title": "Result", "subtitle": "Sub", "url": "https://example.com"}]


class TestExampleService:
    def test_search_returns_results(self):
        service = ExampleService()
        with patch.object(service._client, "search", return_value=_STUB):
            results = service.search("test")
        assert isinstance(results, list)
        assert len(results) > 0

    def test_search_caches_results(self):
        service = ExampleService()
        stub = [{"id": "1", "title": "T", "subtitle": "", "url": ""}]
        with patch.object(service._client, "search", return_value=stub) as mock_search:
            service.search("cached_query")
            service.search("cached_query")
            assert mock_search.call_count == 1  # second call uses cache

    def test_search_different_queries_are_cached_separately(self):
        service = ExampleService()
        with patch.object(service._client, "search", return_value=[]) as mock_search:
            service.search("query_a")
            service.search("query_b")
            assert mock_search.call_count == 2


class TestTemplateParser:
    def _make_template(self, tmp_path: Path, content: str) -> Path:
        md = tmp_path / "template.md"
        md.write_text(dedent(content), encoding="utf-8")
        return md

    def test_text_only(self, tmp_path: Path):
        md = self._make_template(tmp_path, "Hello\n\nWorld")
        blocks = parse(md)
        assert len(blocks) == 1
        assert isinstance(blocks[0], TextBlock)
        assert "Hello" in blocks[0].text

    def test_single_image(self, tmp_path: Path):
        img = tmp_path / "photo.png"
        img.touch()
        md = self._make_template(tmp_path, f"Before\n\n![alt]({img})\n\nAfter")
        blocks = parse(md)
        assert len(blocks) == 3
        assert isinstance(blocks[0], TextBlock)
        assert isinstance(blocks[1], ImageBlock)
        assert isinstance(blocks[2], TextBlock)
        assert blocks[1].alt == "alt"
        assert blocks[1].path == img.resolve()

    def test_relative_image_path(self, tmp_path: Path):
        img_dir = tmp_path / "images"
        img_dir.mkdir()
        img = img_dir / "photo.png"
        img.touch()
        md = self._make_template(tmp_path, "Text\n\n![cap](./images/photo.png)")
        blocks = parse(md)
        assert isinstance(blocks[1], ImageBlock)
        assert blocks[1].path == img.resolve()

    def test_image_alt_text_preserved(self, tmp_path: Path):
        img = tmp_path / "img.png"
        img.touch()
        md = self._make_template(tmp_path, f"![My Caption]({img})")
        blocks = parse(md)
        assert blocks[0].alt == "My Caption"  # type: ignore[union-attr]

    def test_multiple_images(self, tmp_path: Path):
        img1 = tmp_path / "a.png"
        img2 = tmp_path / "b.png"
        img1.touch()
        img2.touch()
        content = f"Intro\n\n![A]({img1})\n\nMiddle\n\n![B]({img2})\n\nEnd"
        md = self._make_template(tmp_path, content)
        blocks = parse(md)
        image_blocks = [b for b in blocks if isinstance(b, ImageBlock)]
        assert len(image_blocks) == 2
        assert image_blocks[0].alt == "A"
        assert image_blocks[1].alt == "B"

    def test_no_text_around_images_ignored(self, tmp_path: Path):
        img = tmp_path / "img.png"
        img.touch()
        md = self._make_template(tmp_path, f"![cap]({img})")
        blocks = parse(md)
        assert len(blocks) == 1
        assert isinstance(blocks[0], ImageBlock)

    def test_image_with_italic_caption(self, tmp_path: Path):
        img = tmp_path / "img.png"
        img.touch()
        md = self._make_template(tmp_path, f"![alt]({img})\n*My caption*\n\nNext text")
        blocks = parse(md)
        assert len(blocks) == 3
        assert isinstance(blocks[0], ImageBlock)
        assert isinstance(blocks[1], CaptionBlock)
        assert blocks[1].text == "My caption"
        assert isinstance(blocks[2], TextBlock)

    def test_image_with_underscore_caption(self, tmp_path: Path):
        img = tmp_path / "img.png"
        img.touch()
        md = self._make_template(tmp_path, f"![alt]({img})\n_My caption_")
        blocks = parse(md)
        assert isinstance(blocks[1], CaptionBlock)
        assert blocks[1].text == "My caption"

    def test_image_without_caption(self, tmp_path: Path):
        img = tmp_path / "img.png"
        img.touch()
        md = self._make_template(tmp_path, f"![alt]({img})\n\nNext text")
        blocks = parse(md)
        assert len(blocks) == 2
        assert isinstance(blocks[0], ImageBlock)
        assert isinstance(blocks[1], TextBlock)

    def test_image_at_end_with_caption(self, tmp_path: Path):
        img = tmp_path / "img.png"
        img.touch()
        md = self._make_template(tmp_path, f"![alt]({img})\n*caption only*")
        blocks = parse(md)
        assert len(blocks) == 2
        assert isinstance(blocks[1], CaptionBlock)
        assert blocks[1].text == "caption only"

    def test_image_at_end_no_caption(self, tmp_path: Path):
        img = tmp_path / "img.png"
        img.touch()
        md = self._make_template(tmp_path, f"Text\n\n![alt]({img})")
        blocks = parse(md)
        assert len(blocks) == 2
        assert isinstance(blocks[0], TextBlock)
        assert isinstance(blocks[1], ImageBlock)
