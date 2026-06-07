package frontend

import (
	"embed"
	"io/fs"
)

//go:embed "build/*"
var buildFiles embed.FS

var Files fs.FS

func init() {
	var err error
	Files, err = fs.Sub(buildFiles, "build")
	if err != nil {
		panic(err)
	}
}
