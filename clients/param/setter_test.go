package param

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValueSetter func(source reflect.Value, target reflect.Value) error

func TestSetValue2(t *testing.T) {
	doTestSetValueSuite(t, SetValue2)
}

func TestSetValue(t *testing.T) {
	doTestSetValueSuite(t, SetValue)
}

func doTestSetValueSuite(t *testing.T, setValue ValueSetter) {
	doTestSetValue[string](t, setValue, "Hello, world!", "Hello, world!")
	doTestSetValue[[]float64](t, setValue, []any{float64(12.34), float32(-56.78)}, []float64{12.34, -56.78})
}

func doTestSetValue[T any](t *testing.T, setValue ValueSetter, source any, expected any) {
	{
		var target T
		err := setValue(reflect.ValueOf(source), reflect.ValueOf(&target))
		if assert.NoError(t, err) {
			doAssertEqual(t, expected, target)
		}
	}
	{
		type Holder struct {
			Value T
		}
		var targetData Holder
		err := setValue(reflect.ValueOf(source), reflect.ValueOf(&targetData.Value))
		if assert.NoError(t, err) {
			doAssertEqual(t, expected, targetData.Value)
		}
	}
	{
		type Holder struct {
			Value *T
		}
		var targetData Holder
		err := setValue(reflect.ValueOf(source), reflect.ValueOf(&targetData.Value))
		if assert.NoError(t, err) {
			if assert.NotNil(t, targetData.Value) {
				doAssertEqual(t, expected, *targetData.Value)
			}
		}
	}
}

func doAssertEqual(t *testing.T, expected any, actual any) {
	expectedValue := reflect.ValueOf(expected)
	expectedType := expectedValue.Type()
	if (expectedType.Kind() == reflect.Slice) && (expectedType.Elem().Kind() == reflect.Float64) {
		assert.InEpsilonSlice(t, expected, actual, epsilon)
	} else if expectedType.Kind() == reflect.Float64 {
		assert.InEpsilon(t, expected, actual, epsilon)
	} else {
		assert.Equal(t, expected, actual)
	}
}
