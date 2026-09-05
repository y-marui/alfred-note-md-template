package templatelist

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func TestList_MissingDirectoryShowsError(t *testing.T) {
	t.Setenv(envTemplatesDir, "/nonexistent/path")

	resp := List("")

	item := resp.Items[0]
	if !strings.Contains(strings.ToLower(item.Title), "not found") {
		t.Errorf("Title = %q, want it to contain %q", item.Title, "not found")
	}
	if item.Valid == nil || *item.Valid {
		t.Errorf("Valid = %v, want false", item.Valid)
	}
}

func TestList_EmptyDirectoryShowsError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envTemplatesDir, dir)

	resp := List("")

	if !strings.Contains(resp.Items[0].Title, "No templates") {
		t.Errorf("Title = %q, want it to contain %q", resp.Items[0].Title, "No templates")
	}
}

func TestList_ListsMdFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envTemplatesDir, dir)
	touch(t, filepath.Join(dir, "article.md"))
	touch(t, filepath.Join(dir, "review.md"))
	touch(t, filepath.Join(dir, "notes.txt")) // should be ignored

	resp := List("")

	if len(resp.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(resp.Items))
	}
	titles := map[string]bool{}
	for _, it := range resp.Items {
		titles[it.Title] = true
	}
	if !titles["article"] || !titles["review"] {
		t.Errorf("titles = %v, want article and review", titles)
	}
}

func TestList_FilterByQuery(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envTemplatesDir, dir)
	touch(t, filepath.Join(dir, "book-review.md"))
	touch(t, filepath.Join(dir, "travel-log.md"))

	resp := List("book")

	if len(resp.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].Title != "book-review" {
		t.Errorf("Title = %q, want %q", resp.Items[0].Title, "book-review")
	}
}

func TestList_NoMatchShowsError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envTemplatesDir, dir)
	touch(t, filepath.Join(dir, "article.md"))

	resp := List("xyz")

	if !strings.Contains(resp.Items[0].Title, "No templates") {
		t.Errorf("Title = %q, want it to contain %q", resp.Items[0].Title, "No templates")
	}
}

func TestList_ItemArgIsFullPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envTemplatesDir, dir)
	md := filepath.Join(dir, "my-template.md")
	touch(t, md)

	resp := List("")

	if resp.Items[0].Arg != md {
		t.Errorf("Arg = %q, want %q", resp.Items[0].Arg, md)
	}
}
