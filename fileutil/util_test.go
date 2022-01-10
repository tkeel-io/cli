package fileutil

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocateFile(t *testing.T) {
	wd, _ := os.Getwd()
	tempDir, _ := os.MkdirTemp(wd, "")
	file := "test.file"

	tests := []struct {
		params []string
		want   struct {
			NotNil    bool
			ErrNotNil bool
		}
	}{
		{[]string{tempDir, file}, struct {
			NotNil    bool
			ErrNotNil bool
		}{true, false}},
		{[]string{path.Join(tempDir, file)}, struct {
			NotNil    bool
			ErrNotNil bool
		}{true, false}},
	}

	var f *os.File
	var err error
	for _, test := range tests {
		err = checkFile(path.Join(tempDir, file))
		assert.NotNil(t, err)

		if len(test.params) == 1 {
			f, err = LocateFile(RewriteFlag(), test.params[0])
		} else {
			f, err = LocateFile(RewriteFlag(), test.params[0], test.params[1:]...)
		}

		if test.want.ErrNotNil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		if test.want.NotNil {
			assert.NotNil(t, f)
		} else {
			assert.Nil(t, f)
		}

		if len(test.params) == 1 {
			os.RemoveAll(test.params[0])
		} else {
			os.RemoveAll(path.Join(test.params...))
		}
	}
	os.RemoveAll(tempDir)
}

func checkFile(file string) error {
	_, err := os.Stat(file)
	return err //nolint
}
