package utils

import (
	"path"
	"path/filepath"
	"runtime"
)

// RootDir returns root directory of project.
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d) + "/"
}
