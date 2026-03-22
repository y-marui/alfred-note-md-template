"""Tests for command handlers."""

from __future__ import annotations

import json
from pathlib import Path
from unittest.mock import patch

from app.commands import config_cmd, help_cmd, note_cmd, open_cmd, search


class TestSearchCommand:
    def test_empty_query_returns_prompt(self, capsys):
        search.handle("")
        data = json.loads(capsys.readouterr().out)
        assert data["items"][0]["valid"] is False
        assert "Type to search" in data["items"][0]["title"]

    def test_whitespace_only_query_returns_prompt(self, capsys):
        """Whitespace-only query must be treated as empty."""
        search.handle("   ")
        data = json.loads(capsys.readouterr().out)
        assert "Type to search" in data["items"][0]["title"]

    def test_returns_results(self, capsys):
        with patch.object(
            search._service,
            "search",
            return_value=[
                {"id": "1", "title": "Result 1", "subtitle": "Sub", "url": "https://example.com"}
            ],
        ):
            search.handle("query")

        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == 1
        assert data["items"][0]["title"] == "Result 1"

    def test_no_results_returns_empty_message(self, capsys):
        with patch.object(search._service, "search", return_value=[]):
            search.handle("noresults")

        data = json.loads(capsys.readouterr().out)
        assert "No results" in data["items"][0]["title"]


class TestOpenCommand:
    def test_no_args_shows_all_shortcuts(self, capsys):
        open_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == len(open_cmd._SHORTCUTS)

    def test_filter_by_name(self, capsys):
        open_cmd.handle("repo")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert all("repo" in t for t in titles)

    def test_unknown_shortcut_shows_error(self, capsys):
        open_cmd.handle("nonexistent")
        data = json.loads(capsys.readouterr().out)
        assert "No shortcut" in data["items"][0]["title"]


class TestConfigCommand:
    def test_empty_config_shows_no_settings(self, capsys):
        config_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert any("No settings" in t for t in titles)

    def test_reset_clears_config(self, capsys):
        config_cmd._config.set("key", "value")
        config_cmd.handle("reset")
        data = json.loads(capsys.readouterr().out)
        assert "reset" in data["items"][0]["title"].lower()
        assert config_cmd._config.all() == {}

    def test_shows_existing_settings(self, capsys):
        config_cmd._config.set("api_key", "secret")
        config_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert any("api_key" in t for t in titles)

    def test_unknown_subcommand_shows_current_config(self, capsys):
        """Unrecognised sub-commands fall through to showing the current config."""
        config_cmd.handle("unknown-subcommand")
        data = json.loads(capsys.readouterr().out)
        # Should return config view, not crash
        assert len(data["items"]) > 0


class TestHelpCommand:
    def test_shows_all_commands(self, capsys):
        help_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == len(help_cmd._COMMANDS)

    def test_all_items_invalid(self, capsys):
        help_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert all(not it["valid"] for it in data["items"])


class TestNoteCommand:
    def test_missing_directory_shows_error(self, capsys, monkeypatch):
        monkeypatch.setenv("NOTE_TEMPLATES_DIR", "/nonexistent/path")
        note_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert "not found" in data["items"][0]["title"].lower()
        assert data["items"][0]["valid"] is False

    def test_empty_directory_shows_error(self, capsys, tmp_path: Path, monkeypatch):
        monkeypatch.setenv("NOTE_TEMPLATES_DIR", str(tmp_path))
        note_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert "No templates" in data["items"][0]["title"]

    def test_lists_md_files(self, capsys, tmp_path: Path, monkeypatch):
        monkeypatch.setenv("NOTE_TEMPLATES_DIR", str(tmp_path))
        (tmp_path / "article.md").touch()
        (tmp_path / "review.md").touch()
        (tmp_path / "notes.txt").touch()  # should be ignored
        note_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert "article" in titles
        assert "review" in titles
        assert len(titles) == 2

    def test_filter_by_query(self, capsys, tmp_path: Path, monkeypatch):
        monkeypatch.setenv("NOTE_TEMPLATES_DIR", str(tmp_path))
        (tmp_path / "book-review.md").touch()
        (tmp_path / "travel-log.md").touch()
        note_cmd.handle("book")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == 1
        assert data["items"][0]["title"] == "book-review"

    def test_no_match_shows_error(self, capsys, tmp_path: Path, monkeypatch):
        monkeypatch.setenv("NOTE_TEMPLATES_DIR", str(tmp_path))
        (tmp_path / "article.md").touch()
        note_cmd.handle("xyz")
        data = json.loads(capsys.readouterr().out)
        assert "No templates" in data["items"][0]["title"]

    def test_item_arg_is_full_path(self, capsys, tmp_path: Path, monkeypatch):
        monkeypatch.setenv("NOTE_TEMPLATES_DIR", str(tmp_path))
        md = tmp_path / "my-template.md"
        md.touch()
        note_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert data["items"][0]["arg"] == str(md)
