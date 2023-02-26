package tempdir

import (
	"log"
	"os"
	"path/filepath"
)

type Dir struct {
	Path string
}

func New(name string) *Dir {
	dir, err := os.MkdirTemp("", name)
	if err != nil {
		log.Fatalf("Failed to create a temporary directory for %s: %v", name, err)
	}
	return &Dir{
		Path: dir,
	}
}

func (d *Dir) Join(path string) string {
	return filepath.Join(d.Path, path)
}

func (d *Dir) Close() {
	os.RemoveAll(d.Path)
}
