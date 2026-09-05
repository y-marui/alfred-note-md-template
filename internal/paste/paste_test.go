package paste

import (
	"reflect"
	"testing"
	"time"

	"github.com/y-marui/alfred-note-md-template/internal/mdtemplate"
)

type fakeClipboard struct {
	calls []string
}

func (f *fakeClipboard) WriteText(text string) error {
	f.calls = append(f.calls, "text:"+text)
	return nil
}

func (f *fakeClipboard) WriteImage(path string) error {
	f.calls = append(f.calls, "image:"+path)
	return nil
}

type fakeKeyboard struct {
	calls []string
}

func (f *fakeKeyboard) Paste() error {
	f.calls = append(f.calls, "paste")
	return nil
}

func (f *fakeKeyboard) PressReturn() error {
	f.calls = append(f.calls, "return")
	return nil
}

func (f *fakeKeyboard) Wait(d time.Duration) {
	f.calls = append(f.calls, "wait")
}

func TestRun_TextBlock_PastesOnce(t *testing.T) {
	clip, key := &fakeClipboard{}, &fakeKeyboard{}

	err := Run([]mdtemplate.Block{mdtemplate.TextBlock{Text: "hello"}}, clip, key)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertEqual(t, clip.calls, []string{"text:hello"})
	assertEqual(t, key.calls, []string{"paste"})
}

func TestRun_ImageBlock_PastesThenSettles(t *testing.T) {
	clip, key := &fakeClipboard{}, &fakeKeyboard{}

	err := Run([]mdtemplate.Block{mdtemplate.ImageBlock{Path: "/tmp/a.png"}}, clip, key)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertEqual(t, clip.calls, []string{"image:/tmp/a.png"})
	assertEqual(t, key.calls, []string{"paste", "wait"})
}

func TestRun_CaptionBlock_PastesThenPressesReturnTwice(t *testing.T) {
	clip, key := &fakeClipboard{}, &fakeKeyboard{}

	err := Run([]mdtemplate.Block{mdtemplate.CaptionBlock{Text: "caption"}}, clip, key)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertEqual(t, clip.calls, []string{"text:caption"})
	assertEqual(t, key.calls, []string{"paste", "return", "return"})
}

func TestRun_FullSequence(t *testing.T) {
	clip, key := &fakeClipboard{}, &fakeKeyboard{}

	blocks := []mdtemplate.Block{
		mdtemplate.TextBlock{Text: "intro"},
		mdtemplate.ImageBlock{Path: "/tmp/a.png"},
		mdtemplate.CaptionBlock{Text: "caption"},
		mdtemplate.TextBlock{Text: "outro"},
	}
	if err := Run(blocks, clip, key); err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertEqual(t, clip.calls, []string{"text:intro", "image:/tmp/a.png", "text:caption", "text:outro"})
	assertEqual(t, key.calls, []string{
		"paste",
		"paste", "wait",
		"paste", "return", "return",
		"paste",
	})
}

func assertEqual(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
