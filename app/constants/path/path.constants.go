package path

import (
	"path"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Basepath   = filepath.Dir(b)
)

func GetBasepath() string {
	//../Basepath
	return path.Join(Basepath, "../../")
}
