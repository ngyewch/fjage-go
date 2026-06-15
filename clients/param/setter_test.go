package param

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetValue(t *testing.T) {
	doTestSetValue[string](t, "Hello, world!", "Hello, world!")
}

func doTestSetValue[T any](t *testing.T, source any, expected any) {
	{
		var target T
		err := SetValue(reflect.ValueOf(source), reflect.ValueOf(&target))
		if assert.NoError(t, err) {
			assert.Equal(t, expected, target)
		}
	}
	{
		type Holder struct {
			Value T
		}
		var targetData Holder
		err := SetValue(reflect.ValueOf(source), reflect.ValueOf(&targetData.Value))
		if assert.NoError(t, err) {
			assert.Equal(t, expected, targetData.Value)
		}
	}
	{
		type Holder struct {
			Value *T
		}
		var targetData Holder
		err := SetValue(reflect.ValueOf(source), reflect.ValueOf(&targetData.Value))
		if assert.NoError(t, err) {
			if assert.NotNil(t, targetData.Value) {
				assert.Equal(t, expected, *targetData.Value)
			}
		}
	}

}
