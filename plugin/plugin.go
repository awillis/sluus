package plugin

type Processor interface {
	Name() string
	Version() string
	Load(name string) bool
	Initialize() error
	Execute() error
	Shutdown() error
	Run()
}

type Plugin struct {
	Processor
	name       string
	version    string
	initialize func() error
	execute    func() error
	shutdown   func() error
}

func (p Plugin) Name() string {
	return p.name
}

func (p Plugin) Version() string {
	return p.version
}

func (p *Plugin) Initialize() error {
	return p.initialize()
}

func (p *Plugin) Execute() error {
	return p.execute()
}

func (p *Plugin) Shutdown() error {
	return p.shutdown()
}
