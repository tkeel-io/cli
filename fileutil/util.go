package fileutil

import (
	"os"
	path "path/filepath"

	"github.com/pkg/errors"
)

const (
	_tkeelRudderDir = ".tkeel/rudder"
	_tokenFile      = ".token"
)

func LocateAdminToken(flag int) (*os.File, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "get user home dir failed")
	}
	return LocateFile(flag, homedir, _tkeelRudderDir, _tokenFile)
}

func LocateFile(flag int, dir string, files ...string) (*os.File, error) {
	filepath := dir
	if len(files) != 0 {
		files = append([]string{dir}, files...)
		filepath = path.Join(files...)
	}
	_, err := os.Stat(filepath)
	if err != nil {
		err = os.MkdirAll(path.Dir(filepath), os.ModeDir|os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "create target file dir failed")
		}
	}

	f, err := os.OpenFile(filepath, flag, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open target file error")
	}

	return f, nil
}

func RWFlag() int {
	return os.O_RDWR | os.O_CREATE
}

func RewriteFlag() int {
	return os.O_RDWR | os.O_CREATE | os.O_TRUNC
}
