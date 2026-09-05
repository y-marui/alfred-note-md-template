// Package mdtemplate parses a note.com paste template into an ordered list
// of text, image, and caption blocks.
//
// Images are identified by standard markdown image syntax: ![alt](path)
// A CaptionBlock is produced when an italic line (*text* or _text_)
// immediately follows an image — note.com treats the first text after a
// pasted image as the image caption.
//
// Images with relative paths are resolved relative to the template file's
// directory.
package mdtemplate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Block is one piece of parsed template content: TextBlock, ImageBlock, or
// CaptionBlock.
type Block interface {
	isBlock()
}

// TextBlock is a run of plain text pasted as-is.
type TextBlock struct {
	Text string
}

func (TextBlock) isBlock() {}

// ImageBlock is an image to be pasted from disk, resolved to an absolute
// path.
type ImageBlock struct {
	Path string
	Alt  string
}

func (ImageBlock) isBlock() {}

// CaptionBlock is the italic line immediately following an image, pasted
// and confirmed with Return.
type CaptionBlock struct {
	Text string
}

func (CaptionBlock) isBlock() {}

var (
	imagePattern = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	// italicLine matches a single line wrapped in * or _ (italic syntax).
	italicLine = regexp.MustCompile(`^\*(.+)\*$|^_(.+)_$`)
)

// Parse reads templatePath and splits it into an ordered list of content
// blocks.
func Parse(templatePath string) ([]Block, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	return splitBlocks(string(content), filepath.Dir(templatePath)), nil
}

func splitBlocks(content, baseDir string) []Block {
	var blocks []Block
	lastEnd := 0
	lastWasImage := false

	for _, m := range imagePattern.FindAllStringSubmatchIndex(content, -1) {
		start, end := m[0], m[1]
		textBefore := content[lastEnd:start]

		if lastWasImage {
			blocks = append(blocks, parsePostImageText(textBefore)...)
		} else if strings.TrimSpace(textBefore) != "" {
			blocks = append(blocks, TextBlock{Text: strings.Trim(textBefore, "\n")})
		}

		alt := content[m[2]:m[3]]
		rawPath := strings.TrimSpace(content[m[4]:m[5]])
		blocks = append(blocks, ImageBlock{Path: resolvePath(rawPath, baseDir), Alt: alt})

		lastEnd = end
		lastWasImage = true
	}

	remaining := content[lastEnd:]
	if lastWasImage {
		blocks = append(blocks, parsePostImageText(remaining)...)
	} else if strings.TrimSpace(remaining) != "" {
		blocks = append(blocks, TextBlock{Text: strings.Trim(remaining, "\n")})
	}

	return blocks
}

// parsePostImageText parses text that immediately follows an image.
//
// If the first non-empty line is italic (*text* or _text_), it becomes a
// CaptionBlock. Any remaining text becomes a TextBlock.
func parsePostImageText(text string) []Block {
	stripped := strings.TrimLeft(text, "\n")
	if strings.TrimSpace(stripped) == "" {
		return nil
	}

	firstLine, rest, _ := strings.Cut(stripped, "\n")
	if m := italicLine.FindStringSubmatch(strings.TrimSpace(firstLine)); m != nil {
		caption := m[1]
		if caption == "" {
			caption = m[2]
		}
		blocks := []Block{CaptionBlock{Text: strings.TrimSpace(caption)}}
		if strings.TrimSpace(rest) != "" {
			blocks = append(blocks, TextBlock{Text: strings.Trim(rest, "\n")})
		}
		return blocks
	}

	return []Block{TextBlock{Text: strings.Trim(stripped, "\n")}}
}

func resolvePath(raw, baseDir string) string {
	p := expandUser(raw)
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Clean(filepath.Join(baseDir, p))
}

func expandUser(p string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	if p == "~" {
		return home
	}
	if rest, ok := strings.CutPrefix(p, "~/"); ok {
		return filepath.Join(home, rest)
	}
	return p
}
