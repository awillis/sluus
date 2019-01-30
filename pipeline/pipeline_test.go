package pipeline

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	pipe := New("test")
	assert.IsType(t, &Pipe{}, pipe)
	assert.NotNil(t, pipe.ID())
}

//func TestPipe_SetSource(t *testing.T) {
//	pipe := New()
//	proc := processor.New("grpc", plugin.SOURCE)
//	err := pipe.SetSource(proc)
//	assert.NotNil(t, err)
//}
//
//func TestPipe_SetSink(t *testing.T) {
//	pipe := New()
//	proc := processor.New("grpc", plugin.SINK)
//	err := pipe.SetSink(proc)
//	assert.NotNil(t, err)
//}
//
//func TestPipe(t *testing.T) {
//	pipe := New()
//	src := processor.New("grpc", plugin.SOURCE)
//	sink := processor.New("grpc", plugin.SINK)
//	err := pipe.SetSource(src)
//	assert.NotNil(t, err)
//	err = pipe.SetSink(sink)
//	assert.NotNil(t, err)
//	sluus := NewSluus(src, sink)
//	err = pipe.AddConduit(sluus)
//	assert.NotNil(t, err)
//}
