// Package keystroke simulates the keypresses needed to paste content into
// the frontmost app via macOS's Accessibility APIs (System Events).
package keystroke

import (
	"fmt"
	"os/exec"
	"time"
)

// Delays after each keypress, giving the target app time to process the
// paste or line break before the next clipboard write.
const (
	pasteSettleDelay  = 800 * time.Millisecond
	returnSettleDelay = 300 * time.Millisecond
)

// Paste sends Cmd+V to the frontmost app and waits for it to settle.
func Paste() error {
	if err := run(pasteScript); err != nil {
		return err
	}
	time.Sleep(pasteSettleDelay)
	return nil
}

// PressReturn sends Return to the frontmost app and waits for it to settle.
func PressReturn() error {
	if err := run(returnScript); err != nil {
		return err
	}
	time.Sleep(returnSettleDelay)
	return nil
}

// Wait pauses for d, e.g. to give a pasted image time to render before the
// next clipboard write.
func Wait(d time.Duration) {
	time.Sleep(d)
}

const (
	pasteScript  = `tell application "System Events" to keystroke "v" using command down`
	returnScript = `tell application "System Events" to keystroke return`
)

func run(script string) error {
	if out, err := exec.Command("osascript", "-e", script).CombinedOutput(); err != nil {
		return fmt.Errorf("keystroke: %w: %s", err, out)
	}
	return nil
}
