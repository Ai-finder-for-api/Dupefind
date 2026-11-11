package dupe

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, p, s string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFindDuplicates(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.txt"), "same content")
	write(t, filepath.Join(dir, "sub", "b.txt"), "same content")
	write(t, filepath.Join(dir, "c.txt"), "same length!")  // same size, different content
	write(t, filepath.Join(dir, "d.txt"), "unique")
	write(t, filepath.Join(dir, "tiny1"), "x")
	write(t, filepath.Join(dir, "tiny2"), "x")

	groups, errs := Find([]string{dir}, 2)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(groups) != 1 {
		t.Fatalf("want 1 group, got %d: %+v", len(groups), groups)
	}
	if len(groups[0].Paths) != 2 || groups[0].Size != int64(len("same content")) {
		t.Fatalf("bad group: %+v", groups[0])
	}
}

func TestHashFileKnownValue(t *testing.T) {
	p := filepath.Join(t.TempDir(), "f")
	write(t, p, "abc")
	h, err := HashFile(p)
	if err != nil {
		t.Fatal(err)
	}
	want := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if h != want {
		t.Fatalf("got %s want %s", h, want)
	}
}

func TestNoDuplicates(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "x"), "1")
	write(t, filepath.Join(dir, "y"), "22")
	groups, _ := Find([]string{dir}, 1)
	if len(groups) != 0 {
		t.Fatalf("expected none, got %+v", groups)
	}
}
