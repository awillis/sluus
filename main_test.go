package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Main(t *testing.T) {
	assert.NotPanics(t, main)
}
