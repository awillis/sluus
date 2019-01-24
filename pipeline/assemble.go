package pipeline

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrNoConfigFound = errors.New("no configuration files found")
)

func Assemble() (err error) {

	confFiles, err := FindConfigurationFiles()

	if err != nil {
		return
	}

	if len(confFiles) == 0 {
		return ErrNoConfigFound
	}

	for _, file := range confFiles {
		config, err := ReadConfigurationFile(file)

		if err != nil {
			return err
		}

		if err = assembleConfig(config); err != nil {
			return err
		}
	}
	return
}

func assembleConfig(config Config) (err error) {
	fmt.Println(config)
	return
}
