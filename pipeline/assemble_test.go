package pipeline

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindConfigTOML(t *testing.T) {
	_, err := FindConfigTOML()
	assert.NoError(t, err)
}
