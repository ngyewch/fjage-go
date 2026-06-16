package param

import (
	"reflect"
)

func Convert(sourceValue reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	if sourceValue.Type().AssignableTo(targetType) {
		return sourceValue, nil
	} else if sourceValue.Type().ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType), nil
	}
	v := sourceValue.Interface()
	sourceValue = reflect.ValueOf(v)
	if sourceValue.Type().AssignableTo(targetType) {
		return sourceValue, nil
	} else if sourceValue.Type().ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType), nil
	}
	if (sourceValue.Kind() == reflect.Slice) && (targetType.Kind() == reflect.Slice) {
		newSlice := reflect.MakeSlice(targetType, sourceValue.Len(), sourceValue.Cap())
		for i := 0; i < sourceValue.Len(); i++ {
			convertedValue, err := Convert(sourceValue.Index(i), targetType.Elem())
			if err != nil {
				return reflect.Value{}, err
			}
			newSlice.Index(i).Set(convertedValue)
		}
		return newSlice, nil
	}
	//fmt.Printf("sourceValue=%v, sourceType=%v, targetType=%v\n", sourceValue, sourceValue.Type(), targetType)
	return reflect.Value{}, ErrCannotConvert
}
