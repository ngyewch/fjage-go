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
	fmt.Printf("target kind=%v canSet=%v\n", target.Kind(), target.CanSet())
	if (target.Kind() == reflect.Pointer) && !target.CanSet() {
		return SetValue(source, target.Elem())
	}
	if !target.CanSet() {
		return ErrCannotSet
	}
	//fmt.Printf("source=%v (%s) / target=%v (%s)\n", source, source.Type().String(), target, target.Type().String())
	if source.Kind() == reflect.Pointer {
		return fmt.Errorf("source cannot be a pointer")
	}
	if (source.Type() == target.Type()) && target.CanSet() {
		target.Set(source)
		return nil
	}
	if (target.Kind() == reflect.Slice) && (source.Kind() == reflect.Slice) {
		newSlice := reflect.MakeSlice(target.Type(), source.Len(), source.Cap())
		for i := 0; i < source.Len(); i++ {
			err := SetValue(source.Index(i), newSlice.Index(i))
			if err != nil {
				return err
			}
		}
		target.Set(newSlice)
		return nil
	} else if target.Kind() == reflect.Pointer {
		if target.CanSet() {
			if source.Type().AssignableTo(target.Type().Elem()) {
				sourcePointer := reflect.New(target.Type().Elem())
				sourcePointer.Elem().Set(source)
				target.Set(sourcePointer)
				return nil
			} else if source.Type().ConvertibleTo(target.Type().Elem()) {
				sourcePointer := reflect.New(target.Type().Elem())
				sourcePointer.Elem().Set(source.Convert(target.Type().Elem()))
				target.Set(sourcePointer)
				return nil
			}
		} else if target.Elem().CanSet() {
			if source.Type().AssignableTo(target.Type().Elem()) {
				target.Elem().Set(source)
				return nil
			} else if source.Type().ConvertibleTo(target.Type().Elem()) {
				target.Elem().Set(source.Convert(target.Type().Elem()))
				return nil
			}
		}
	} else {
		if source.Type().AssignableTo(target.Type()) {
			target.Set(source)
			return nil
		} else if source.Type().ConvertibleTo(target.Type()) {
			target.Set(source.Convert(target.Type()))
			return nil
		}
		// TODO properly handle any
		sourceType := reflect.TypeOf(source.Interface())
		sourceValue := reflect.ValueOf(source.Interface())
		if sourceType.AssignableTo(target.Type()) {
			target.Set(sourceValue)
			return nil
		} else if sourceType.ConvertibleTo(target.Type()) {
			target.Set(sourceValue.Convert(target.Type()))
			return nil
		}
	}
	return fmt.Errorf("param value (%v) not assignable to target (%v)", source.Type(), target.Type())
}
