package plugin

type (
	Option  interface{}
	Default func(Option)
)

func Validate(opt Option, def ...Default) {
	for _, o := range def {
		o(opt)
	}
}
