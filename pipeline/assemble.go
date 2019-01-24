package pipeline

import (
	"github.com/awillis/sluus/core"
	"os"
	"path/filepath"
)

func FindConfigTOML() (filelist []string, err error) {

	err = filepath.Walk(core.CONFDIR, func(path string, info os.FileInfo, err error) (rerr error) {
		if info.IsDir() {
			return
		}

		filelist = append(filelist, path)
		return
	})

	return
}
