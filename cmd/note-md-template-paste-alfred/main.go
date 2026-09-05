// Command note-md-template-paste-alfred is the Run Script action binary
// Alfred invokes after the user selects a template from the note-md-template-alfred
// Script Filter (see workflow/info.plist). It receives the selected
// template's absolute path as $1, parses it, and pastes each block into the
// frontmost app (expected to be note.com's editor).
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/y-marui/alfred-note-md-template/internal/clipboard"
	"github.com/y-marui/alfred-note-md-template/internal/keystroke"
	"github.com/y-marui/alfred-note-md-template/internal/mdtemplate"
	"github.com/y-marui/alfred-note-md-template/internal/paste"
)

// alfredWindowCloseDelay gives Alfred's window time to close and the
// previously-frontmost app (note.com's editor) time to regain focus before
// the first paste.
const alfredWindowCloseDelay = 500 * time.Millisecond

type systemClipboard struct{}

func (systemClipboard) WriteText(text string) error  { return clipboard.WriteText(text) }
func (systemClipboard) WriteImage(path string) error { return clipboard.WriteImage(path) }

type systemKeyboard struct{}

func (systemKeyboard) Paste() error         { return keystroke.Paste() }
func (systemKeyboard) PressReturn() error   { return keystroke.PressReturn() }
func (systemKeyboard) Wait(d time.Duration) { keystroke.Wait(d) }

func main() {
	if len(os.Args) < 2 {
		fail("Usage: note-md-template-paste-alfred <template_path>")
	}
	templatePath := os.Args[1]
	if _, err := os.Stat(templatePath); err != nil {
		fail(fmt.Sprintf("Template not found: %s", templatePath))
	}

	time.Sleep(alfredWindowCloseDelay)

	blocks, err := mdtemplate.Parse(templatePath)
	if err != nil {
		fail(fmt.Sprintf("Parsing template failed: %v", err))
	}

	if err := paste.Run(blocks, systemClipboard{}, systemKeyboard{}); err != nil {
		fail(fmt.Sprintf("Pasting template failed: %v", err))
	}
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
