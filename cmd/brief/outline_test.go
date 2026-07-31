package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/git-pkgs/outline"
)

var benchmarkOutlineResult *outline.Result

func TestRunOutlineSkipsDetectedBinaryContent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")
	writeFile(t, dir, "document.pdf", "%PDF-1.7\nbody\n")

	var output bytes.Buffer
	if code := runOutline(dir, outline.Options{}, false, &output); code != 0 {
		t.Fatalf("runOutline() = %d, want 0", code)
	}

	got := output.String()
	if !strings.Contains(got, "### document.pdf\n\n_(skipped: binary,") {
		t.Errorf("output does not mark PDF as binary:\n%s", got)
	}
	if strings.Contains(got, "%PDF-1.7") {
		t.Errorf("output includes PDF content:\n%s", got)
	}
	if !strings.Contains(got, "package main") {
		t.Errorf("output omits text content:\n%s", got)
	}
}

func BenchmarkOutlinePack(b *testing.B) {
	dir := b.TempDir()
	for i := range 100 {
		name := "source-" + strconv.Itoa(i) + ".go"
		content := []byte("package source\n\nfunc Value() int { return 1 }\n")
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			b.Fatal(err)
		}
	}
	for i := range 10 {
		name := "document-" + strconv.Itoa(i) + ".pdf"
		if err := os.WriteFile(filepath.Join(dir, name), []byte("%PDF-1.7\nbody\n"), 0o644); err != nil {
			b.Fatal(err)
		}
	}

	opts := outline.Options{Concurrency: 1}
	b.ResetTimer()
	for b.Loop() {
		result, err := outline.Pack(dir, opts)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkOutlineResult = result
	}
}
