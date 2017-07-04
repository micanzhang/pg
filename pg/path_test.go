package pg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	path, err := Path()
	if err != nil {
		t.Error(err)
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		t.Error(err)
	}
}
