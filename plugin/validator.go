package plugin

type (
	Option  interface{}
	Default func(Option)
)

// Validate() is used by plugins to check configured values passed in from
// the configuration file, and set reasonable default values if necessary
func Validate(opt Option, def ...Default) {
	for _, o := range def {
		o(opt)
	}
}
