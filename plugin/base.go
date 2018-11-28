package plugin

import (
	"fmt"
	"github.com/golang-collections/go-datastructures/queue"
	"kapilary/core"
	"os"
	"path/filepath"
	"plugin"
)

type Base struct {
	input    chan core.Event
	output   chan core.Event
	queue    queue.PriorityQueue
	category string
}

func (b Base) Input() chan core.Event {
	return b.input
}

func (b Base) Output() chan core.Event {
	return b.output
}

func Load(filename string) bool {

	plug, err := plugin.Open(filename)

	if err != nil {
		fmt.Errorf("error loading plugin: %v", err)
		return false
	}

	symPlug, err := plug.Lookup("KapilaryPlugin")

	if err != nil {
		fmt.Println(err)
		return false
	}

	var stogo splug.KapilaryPlugin
	stogo, ok := symPlug.(KapilaryPlugin)

	if !ok {
		fmt.Println("unexpected type from module symbol")
		return false
	}
	_ = stogo
	// load plugin by filename using plugin.Open
	// using filename to derive plugin name, attempt to load
	// a symbol called 'NamedPlugin'
	// Assert that the loaded symbol meets the KapilaryPlugin
	// interface
	// Call PluginInit method

	return true
}

func FindPlugins(config config.StogoConfig) []string {

	list := make([]string, 10)

	filepath.Walk(config.PlugDir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("Unable to walk path")
		}

		if info.IsDir() == false {
			list = append(list, path)
		}

		return nil
	})

	// using config object, find plugins in path
	return list
}

func (b Base) Queue() *queue.PriorityQueue {
	return &b.queue
}

func (b Base) Category() string {
	return b.category
}
