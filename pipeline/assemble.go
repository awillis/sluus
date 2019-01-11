package pipeline

import (
	"fmt"
	"github.com/awillis/sluus/core"
	"os"
	"path/filepath"
)

func FindConfigTOML() {
	println("about to do filepath walk")
	err := filepath.Walk(core.CONFDIR, func(path string, info os.FileInfo, err error) error {
		println("walking files")
		fmt.Println(path)
		core.Logger.Infow("file listing",
			"path", path, "info", info,
		)
		return err
	})

	if err != nil {
		core.Logger.Panic("error searching for config files")
	}
}
