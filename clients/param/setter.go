package param

import (
	"fmt"
	"reflect"
)

var (
	ErrCannotSet     = fmt.Errorf("cannot set")
	ErrCannotConvert = fmt.Errorf("cannot convert")
)

func SetValue(source reflect.Value, target reflect.Value) error {
	if source.Kind() == reflect.Pointer {
		return fmt.Errorf("source cannot be a pointer")
	}
	if target.Kind() == reflect.Pointer {
		if target.Type().Elem().Kind() == reflect.Pointer {
			return SetValue(source, target.Elem())
		}
		if target.CanSet() {
			newSourcePointer := reflect.New(target.Type().Elem())
			convertedSource, err := Convert(source, target.Type().Elem())
			if err != nil {
				return err
			}
			newSourcePointer.Elem().Set(convertedSource)
			target.Set(newSourcePointer)
			return nil
		}
		return SetValue(source, target.Elem())
	}
	if !target.CanSet() {
		return ErrCannotSet
	}
	convertedSource, err := Convert(source, target.Type())
	if err != nil {
		return err
	}
	target.Set(convertedSource)
	return nil
}
