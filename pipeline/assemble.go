package pipeline

import (
	"fmt"
	"github.com/awillis/sluus/core"
	"os"
	"path/filepath"
)

func FindConfigTOML() {
	err := core.Logger.Sync()
	if err != nil {
		panic(err)
	}
	core.Logger.Info("about to do filepath walk")
	err = filepath.Walk(core.CONFDIR, func(path string, info os.FileInfo, err error) error {
		println("walking files")
		fmt.Println(path)
		return err
	})

	if err != nil {
		//core.Logger.Panic("error searching for config files")
	}
}
