package mdtemplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeTemplate writes content to a template.md file under dir and returns
// its path.
func makeTemplate(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "template.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func TestParse_TextOnly(t *testing.T) {
	dir := t.TempDir()
	md := makeTemplate(t, dir, "Hello\n\nWorld")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("len(blocks) = %d, want 1", len(blocks))
	}
	text, ok := blocks[0].(TextBlock)
	if !ok {
		t.Fatalf("blocks[0] = %T, want TextBlock", blocks[0])
	}
	if !strings.Contains(text.Text, "Hello") {
		t.Errorf("text.Text = %q, want it to contain %q", text.Text, "Hello")
	}
}

func TestParse_SingleImage(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "photo.png")
	touch(t, img)
	md := makeTemplate(t, dir, "Before\n\n![alt]("+img+")\n\nAfter")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(blocks) != 3 {
		t.Fatalf("len(blocks) = %d, want 3", len(blocks))
	}
	if _, ok := blocks[0].(TextBlock); !ok {
		t.Errorf("blocks[0] = %T, want TextBlock", blocks[0])
	}
	image, ok := blocks[1].(ImageBlock)
	if !ok {
		t.Fatalf("blocks[1] = %T, want ImageBlock", blocks[1])
	}
	if _, ok := blocks[2].(TextBlock); !ok {
		t.Errorf("blocks[2] = %T, want TextBlock", blocks[2])
	}
	if image.Alt != "alt" {
		t.Errorf("image.Alt = %q, want %q", image.Alt, "alt")
	}
	if image.Path != img {
		t.Errorf("image.Path = %q, want %q", image.Path, img)
	}
}

func TestParse_RelativeImagePath(t *testing.T) {
	dir := t.TempDir()
	imgDir := filepath.Join(dir, "images")
	if err := os.Mkdir(imgDir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	img := filepath.Join(imgDir, "photo.png")
	touch(t, img)
	md := makeTemplate(t, dir, "Text\n\n![cap](./images/photo.png)")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	image, ok := blocks[1].(ImageBlock)
	if !ok {
		t.Fatalf("blocks[1] = %T, want ImageBlock", blocks[1])
	}
	if image.Path != img {
		t.Errorf("image.Path = %q, want %q", image.Path, img)
	}
}

func TestParse_ImageAltTextPreserved(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "img.png")
	touch(t, img)
	md := makeTemplate(t, dir, "![My Caption]("+img+")")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	image, ok := blocks[0].(ImageBlock)
	if !ok {
		t.Fatalf("blocks[0] = %T, want ImageBlock", blocks[0])
	}
	if image.Alt != "My Caption" {
		t.Errorf("image.Alt = %q, want %q", image.Alt, "My Caption")
	}
}

func TestParse_MultipleImages(t *testing.T) {
	dir := t.TempDir()
	img1 := filepath.Join(dir, "a.png")
	img2 := filepath.Join(dir, "b.png")
	touch(t, img1)
	touch(t, img2)
	content := "Intro\n\n![A](" + img1 + ")\n\nMiddle\n\n![B](" + img2 + ")\n\nEnd"
	md := makeTemplate(t, dir, content)

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var images []ImageBlock
	for _, b := range blocks {
		if img, ok := b.(ImageBlock); ok {
			images = append(images, img)
		}
	}
	if len(images) != 2 {
		t.Fatalf("len(images) = %d, want 2", len(images))
	}
	if images[0].Alt != "A" {
		t.Errorf("images[0].Alt = %q, want %q", images[0].Alt, "A")
	}
	if images[1].Alt != "B" {
		t.Errorf("images[1].Alt = %q, want %q", images[1].Alt, "B")
	}
}

func TestParse_NoTextAroundImagesIgnored(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "img.png")
	touch(t, img)
	md := makeTemplate(t, dir, "![cap]("+img+")")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("len(blocks) = %d, want 1", len(blocks))
	}
	if _, ok := blocks[0].(ImageBlock); !ok {
		t.Errorf("blocks[0] = %T, want ImageBlock", blocks[0])
	}
}

func TestParse_ImageWithItalicCaption(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "img.png")
	touch(t, img)
	md := makeTemplate(t, dir, "![alt]("+img+")\n*My caption*\n\nNext text")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(blocks) != 3 {
		t.Fatalf("len(blocks) = %d, want 3", len(blocks))
	}
	if _, ok := blocks[0].(ImageBlock); !ok {
		t.Errorf("blocks[0] = %T, want ImageBlock", blocks[0])
	}
	caption, ok := blocks[1].(CaptionBlock)
	if !ok {
		t.Fatalf("blocks[1] = %T, want CaptionBlock", blocks[1])
	}
	if caption.Text != "My caption" {
		t.Errorf("caption.Text = %q, want %q", caption.Text, "My caption")
	}
	if _, ok := blocks[2].(TextBlock); !ok {
		t.Errorf("blocks[2] = %T, want TextBlock", blocks[2])
	}
}

func TestParse_ImageWithUnderscoreCaption(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "img.png")
	touch(t, img)
	md := makeTemplate(t, dir, "![alt]("+img+")\n_My caption_")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	caption, ok := blocks[1].(CaptionBlock)
	if !ok {
		t.Fatalf("blocks[1] = %T, want CaptionBlock", blocks[1])
	}
	if caption.Text != "My caption" {
		t.Errorf("caption.Text = %q, want %q", caption.Text, "My caption")
	}
}

func TestParse_ImageWithoutCaption(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "img.png")
	touch(t, img)
	md := makeTemplate(t, dir, "![alt]("+img+")\n\nNext text")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("len(blocks) = %d, want 2", len(blocks))
	}
	if _, ok := blocks[0].(ImageBlock); !ok {
		t.Errorf("blocks[0] = %T, want ImageBlock", blocks[0])
	}
	if _, ok := blocks[1].(TextBlock); !ok {
		t.Errorf("blocks[1] = %T, want TextBlock", blocks[1])
	}
}

func TestParse_ImageAtEndWithCaption(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "img.png")
	touch(t, img)
	md := makeTemplate(t, dir, "![alt]("+img+")\n*caption only*")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("len(blocks) = %d, want 2", len(blocks))
	}
	caption, ok := blocks[1].(CaptionBlock)
	if !ok {
		t.Fatalf("blocks[1] = %T, want CaptionBlock", blocks[1])
	}
	if caption.Text != "caption only" {
		t.Errorf("caption.Text = %q, want %q", caption.Text, "caption only")
	}
}

func TestParse_ImageAtEndNoCaption(t *testing.T) {
	dir := t.TempDir()
	img := filepath.Join(dir, "img.png")
	touch(t, img)
	md := makeTemplate(t, dir, "Text\n\n![alt]("+img+")")

	blocks, err := Parse(md)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("len(blocks) = %d, want 2", len(blocks))
	}
	if _, ok := blocks[0].(TextBlock); !ok {
		t.Errorf("blocks[0] = %T, want TextBlock", blocks[0])
	}
	if _, ok := blocks[1].(ImageBlock); !ok {
		t.Errorf("blocks[1] = %T, want ImageBlock", blocks[1])
	}
}
