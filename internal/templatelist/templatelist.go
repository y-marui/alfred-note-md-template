// Package templatelist lists and filters the user's note.com paste
// templates for Alfred's Script Filter.
//
// Templates are read from the directory configured via the templates_dir
// Alfred Config Builder variable. Default: ~/Documents/Note Templates.
package templatelist

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/y-marui/alfred-note-md-template/internal/scriptfilter"
)

const envTemplatesDir = "templates_dir"

// List returns the Script Filter response listing .md template files in the
// configured templates directory, optionally filtered by query.
func List(query string) scriptfilter.Response {
	dir := templatesDir()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return errorResponse("Templates directory not found", "Create "+dir+" and add .md files")
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(stem(names[i])) < strings.ToLower(stem(names[j]))
	})

	if len(names) == 0 {
		return errorResponse("No templates found", "Add .md files to "+dir)
	}

	trimmed := strings.ToLower(strings.TrimSpace(query))
	filtered := names
	if trimmed != "" {
		filtered = nil
		for _, name := range names {
			if strings.Contains(strings.ToLower(stem(name)), trimmed) {
				filtered = append(filtered, name)
			}
		}
	}

	if len(filtered) == 0 {
		return errorResponse(`No templates matching "`+query+`"`, "Try a different keyword")
	}

	items := make([]scriptfilter.Item, 0, len(filtered))
	for _, name := range filtered {
		path := filepath.Join(dir, name)
		items = append(items, scriptfilter.Item{
			Title:    stem(name),
			Subtitle: path,
			Arg:      path,
			UID:      path,
			Valid:    scriptfilter.BoolPtr(true),
		})
	}
	return scriptfilter.Response{Items: items}
}

func templatesDir() string {
	if dir := os.Getenv(envTemplatesDir); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("Documents", "Note Templates")
	}
	return filepath.Join(home, "Documents", "Note Templates")
}

func stem(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func errorResponse(title, subtitle string) scriptfilter.Response {
	return scriptfilter.Response{
		Items: []scriptfilter.Item{
			{Title: title, Subtitle: subtitle, Valid: scriptfilter.BoolPtr(false)},
		},
	}
}
