package plugin

import (
	"context"
	"fmt"
	"go/token"
	"golang.org/x/tools/go/packages"

	"github.com/awillis/sluus/core"
)

func Load(name string, ptype core.PluginType) (Processor, error) {

	var err error
	pkglist, err := packages.Load(&packages.Config{
		Mode:       packages.LoadSyntax,
		Context:    context.Background(),
		Dir:        "",
		Env:        nil,
		BuildFlags: nil,
		Fset:       &token.FileSet{},
		ParseFile:  nil,
		Tests:      false,
		Overlay:    nil,
	}, name)

	if err != nil {
		return nil, fmt.Errorf("error loading plugin: %s", err)
	}

	for _, pkg := range pkglist {
		_ = pkg
	}

	p := new(Plugin)

	return p, err
}
