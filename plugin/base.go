package plugin

type Plugin interface {
	PluginInit() bool
	PluginLoad(name string) bool
}

type PlugBase struct {
}

func (p *PlugBase) PluginLoad(name string) bool {

	//plug, err := plugin.Open(filename)
	//
	//if err != nil {
	//	fmt.Errorf("error loading plugin: %v", err)
	//	return false
	//}
	//
	//symPlug, err := plug.Lookup("KapilaryPlugin")
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return false
	//}
	//
	//var stogo symPlug.KapilaryPlugin
	//stogo, ok := symPlug.(KapilaryPlugin)
	//
	//if !ok {
	//	fmt.Println("unexpected type from module symbol")
	//	return false
	//}
	//_ = stogo
	//// load plugin by filename using plugin.Open
	//// using filename to derive plugin name, attempt to load
	//// a symbol called 'NamedPlugin'
	//// Assert that the loaded symbol meets the KapilaryPlugin
	//// interface
	//// Call PluginInit method

	return true
}
