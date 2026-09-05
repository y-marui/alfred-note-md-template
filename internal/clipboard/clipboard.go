// Package clipboard writes text and image content to the macOS general
// pasteboard. It never logs the content it handles.
package clipboard

import (
	"fmt"
	"os/exec"
	"strings"
)

// WriteText replaces the general pasteboard's content with text, discarding
// any other representation the pasteboard held.
func WriteText(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard: writing text: %w", err)
	}
	return nil
}

// WriteImage replaces the general pasteboard's content with the image file
// at path, coerced to a TIFF picture. pbcopy has no equivalent for image
// data, so this shells out to osascript, which can read a file and coerce
// it to an image class before setting the clipboard.
func WriteImage(path string) error {
	script := `set the clipboard to (read (POSIX file ` + appleScriptString(path) + `) as TIFF picture)`
	out, err := exec.Command("osascript", "-e", script).CombinedOutput()
	if err != nil {
		return fmt.Errorf("clipboard: writing image %s: %w: %s", path, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// appleScriptString quotes s as an AppleScript string literal.
func appleScriptString(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
	return `"` + replacer.Replace(s) + `"`
}
