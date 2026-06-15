package param

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetValue(t *testing.T) {
	{
		fmt.Println("----")
		var source = "Hello, world!"
		var target string
		err := SetValue(reflect.ValueOf(source), reflect.ValueOf(&target))
		if assert.NoError(t, err) {
			assert.Equal(t, "Hello, world!", target)
		}
	}
	{
		fmt.Println("----")
		type TargetData struct {
			Value string
		}
		var source = "Hello, world!"
		var targetData TargetData
		err := SetValue(reflect.ValueOf(source), reflect.ValueOf(&targetData.Value))
		if assert.NoError(t, err) {
			assert.Equal(t, "Hello, world!", targetData.Value)
		}
	}
	{
		fmt.Println("----")
		type TargetData struct {
			Value *string
		}
		var source = "Hello, world!"
		var targetData TargetData
		err := SetValue(reflect.ValueOf(source), reflect.ValueOf(&targetData.Value))
		if assert.NoError(t, err) {
			if assert.NotNil(t, targetData.Value) {
				assert.Equal(t, "Hello, world!", *targetData.Value)
			}
		}
	}
}
