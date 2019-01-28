package pipeline

import (
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

var (
	ErrNeedPointer  = errors.New("expected pointer to struct")
	ErrInvalidValue = errors.New("expected string, bool or int")
)

// SetPluginOptions takes a pointer to a struct containing configuration
// values for a plugin, and populates it with values found in the config file
func (o OptionMap) SetPluginOptions(pluginOption interface{}) (err error) {

	plugOptTyp := reflect.TypeOf(pluginOption).Elem()
	plugOptVal := reflect.ValueOf(pluginOption).Elem()

	if plugOptVal.Kind() != reflect.Struct || !plugOptVal.CanSet() {
		return errors.Wrapf(ErrNeedPointer, "found %s", plugOptTyp.Kind().String())
	}

	for i := 0; i < plugOptVal.NumField(); i++ {

		// values in the configuration file are expected to be in snake case
		optionValue := o[strcase.ToSnake(plugOptTyp.Field(i).Name)]

		switch optionValue.(type) {
		case string:
			if plugOptVal.Field(i).Kind() == reflect.String {
				plugOptVal.Field(i).SetString(optionValue.(string))
			}
		case bool:
			if plugOptVal.Field(i).Kind() == reflect.Bool {
				plugOptVal.Field(i).SetBool(optionValue.(bool))
			}
		case int, int64:
			if plugOptVal.Field(i).Kind() == reflect.Int64 {
				plugOptVal.Field(i).SetInt(optionValue.(int64))
			}
		}
	}
	return
}
