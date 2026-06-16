package param

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValueSetter func(source reflect.Value, target reflect.Value) error

func TestSetValue(t *testing.T) {
	doTestSetValueSuite(t, SetValue)
}

func doTestSetValueSuite(t *testing.T, setValue ValueSetter) {
	doTestSetValue[string](t, setValue, "Hello, world!", "Hello, world!")
}

func doTestSetValue[T any](t *testing.T, setValue ValueSetter, source any, expected any) {
	{
		var target T
		err := setValue(reflect.ValueOf(source), reflect.ValueOf(&target))
		if assert.NoError(t, err) {
			assert.Equal(t, expected, target)
		}
	}
	{
		type Holder struct {
			Value T
		}
		var targetData Holder
		err := setValue(reflect.ValueOf(source), reflect.ValueOf(&targetData.Value))
		if assert.NoError(t, err) {
			assert.Equal(t, expected, targetData.Value)
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
				assert.Equal(t, expected, *targetData.Value)
			}
		}
	}
}
