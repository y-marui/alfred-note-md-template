// Package paste plays a parsed note.com template into the frontmost app by
// writing each block to the clipboard and simulating the keypresses that
// paste it.
package paste

import (
	"time"

	"github.com/y-marui/alfred-note-md-template/internal/mdtemplate"
)

// imageSettleDelay is extra time given to note.com after pasting an image,
// beyond the delay Keyboard.Paste already waits — larger images take
// longer to upload and render than a paste of plain text.
const imageSettleDelay = 500 * time.Millisecond

// Clipboard writes content to the system pasteboard.
type Clipboard interface {
	WriteText(text string) error
	WriteImage(path string) error
}

// Keyboard simulates the keypresses that act on pasteboard content once
// it's been written.
type Keyboard interface {
	Paste() error
	PressReturn() error
	Wait(d time.Duration)
}

// Run pastes each block of a parsed template in sequence:
//
//   - TextBlock:    set clipboard to plain text -> paste
//   - ImageBlock:   set clipboard to image data -> paste -> settle
//   - CaptionBlock: set clipboard to plain text -> paste -> Return -> Return
//
// The double Return after a caption matches note.com's editor: the first
// confirms the caption text, the second starts a new paragraph beneath the
// image so the next block doesn't append to the caption.
func Run(blocks []mdtemplate.Block, clip Clipboard, key Keyboard) error {
	for _, b := range blocks {
		switch v := b.(type) {
		case mdtemplate.TextBlock:
			if err := clip.WriteText(v.Text); err != nil {
				return err
			}
			if err := key.Paste(); err != nil {
				return err
			}
		case mdtemplate.ImageBlock:
			if err := clip.WriteImage(v.Path); err != nil {
				return err
			}
			if err := key.Paste(); err != nil {
				return err
			}
			key.Wait(imageSettleDelay)
		case mdtemplate.CaptionBlock:
			if err := clip.WriteText(v.Text); err != nil {
				return err
			}
			if err := key.Paste(); err != nil {
				return err
			}
			if err := key.PressReturn(); err != nil {
				return err
			}
			if err := key.PressReturn(); err != nil {
				return err
			}
		}
	}
	return nil
}
