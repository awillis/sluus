package pipeline

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfigurationFile(t *testing.T) {

	filelist, err := FindConfigTOML()
	assert.NoError(t, err)
	for _, file := range filelist {
		config, err := ReadConfigurationFile(file)
		assert.NoError(t, err)
		fmt.Printf("final config: %+v\n", config)
	}
}
