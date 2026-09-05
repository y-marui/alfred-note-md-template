package clipboard

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestMain saves the real clipboard's plain-text content once before any
// test in this package runs, and restores it once after they all finish, so
// a local run doesn't clobber the developer's clipboard. Non-text clipboard
// content present before the suite runs cannot be restored this way — an
// accepted limitation of testing against the real pasteboard.
func TestMain(m *testing.M) {
	original, _ := exec.Command("pbpaste").Output()
	code := m.Run()
	_ = WriteText(string(original))
	os.Exit(code)
}

// requireMacOS skips the test unless running on macOS with the required
// command-line tools on PATH.
func requireMacOS(t *testing.T, bins ...string) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("clipboard package is macOS-only")
	}
	for _, bin := range bins {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not found", bin)
		}
	}
}

func TestWriteText_RoundTrip(t *testing.T) {
	requireMacOS(t, "pbcopy", "pbpaste")

	want := "line1\nline2​ZWSP"
	if err := WriteText(want); err != nil {
		t.Fatalf("WriteText: %v", err)
	}

	got, err := exec.Command("pbpaste").Output()
	if err != nil {
		t.Fatalf("pbpaste: %v", err)
	}
	if string(got) != want {
		t.Errorf("pbpaste = %q, want %q", got, want)
	}
}

func TestWriteImage_SetsImageClipboardType(t *testing.T) {
	requireMacOS(t, "osascript")

	path := writeTestPNG(t)
	if err := WriteImage(path); err != nil {
		t.Fatalf("WriteImage: %v", err)
	}

	out, err := exec.Command("osascript", "-e", "clipboard info").Output()
	if err != nil {
		t.Fatalf("clipboard info: %v", err)
	}
	if !strings.Contains(string(out), "TIFF picture") {
		t.Errorf("clipboard info = %q, want it to mention TIFF picture", out)
	}
}

// writeTestPNG writes a minimal valid 1x1 PNG file and returns its path.
func writeTestPNG(t *testing.T) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}

	path := filepath.Join(t.TempDir(), "test.png")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}
