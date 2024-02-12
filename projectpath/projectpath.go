package projectpath

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b), "../")
)

func Absolute(path string) string {
	return filepath.Join(Root, path)
}
