package param

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	epsilon = 0.0000001
)

func TestConvert(t *testing.T) {
	doTestConvert(t, reflect.ValueOf(bool(true)), reflect.TypeFor[bool](), bool(true))

	doTestConvert(t, reflect.ValueOf(string("abc")), reflect.TypeFor[string](), string("abc"))

	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(int8(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(int16(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(int32(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(int64(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(uint8(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(uint16(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(uint32(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[float64](), float64(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[float32](), float32(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(uint64(123)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[float64](), float64(123.456))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[float32](), float32(123.456))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(float32(123.456)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[float64](), float64(123.456))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[float32](), float32(123.456))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[int64](), int64(123))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[int32](), int32(123))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[int16](), int16(123))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[int8](), int8(123))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[uint64](), uint64(123))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[uint32](), uint32(123))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[uint16](), uint16(123))
	doTestConvert(t, reflect.ValueOf(float64(123.456)), reflect.TypeFor[uint8](), uint8(123))

	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int8{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int16{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int32{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]int64{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint8{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint16{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint32{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]float64](), []float64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]float32](), []float32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]int32](), []int32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]int16](), []int16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]int8](), []int8{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]uint64](), []uint64{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]uint32](), []uint32{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]uint16](), []uint16{10, 20})
	doTestConvert(t, reflect.ValueOf([]uint64{10, 20}), reflect.TypeFor[[]uint8](), []uint8{10, 20})

	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]float64](), []float64{12.34, 56.78})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]float32](), []float32{12.34, 56.78})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]int64](), []int64{12, 56})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]int32](), []int32{12, 56})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]int16](), []int16{12, 56})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]int8](), []int8{12, 56})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]uint64](), []uint64{12, 56})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]uint32](), []uint32{12, 56})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]uint16](), []uint16{12, 56})
	doTestConvert(t, reflect.ValueOf([]float32{12.34, 56.78}), reflect.TypeFor[[]uint8](), []uint8{12, 56})

	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]float64](), []float64{12.34, 56.78})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]float32](), []float32{12.34, 56.78})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]int64](), []int64{12, 56})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]int32](), []int32{12, 56})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]int16](), []int16{12, 56})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]int8](), []int8{12, 56})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]uint64](), []uint64{12, 56})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]uint32](), []uint32{12, 56})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]uint16](), []uint16{12, 56})
	doTestConvert(t, reflect.ValueOf([]float64{12.34, 56.78}), reflect.TypeFor[[]uint8](), []uint8{12, 56})

	doTestConvert(t, reflect.ValueOf([]any{int32(10), int64(20)}), reflect.TypeFor[[]int64](), []int64{10, 20})
	doTestConvert(t, reflect.ValueOf([]any{float32(123.456), float64(-456.123)}), reflect.TypeFor[[]float64](), []float64{float64(123.456), float64(-456.123)})

	doTestConvert(t,
		reflect.ValueOf(map[string]any{
			"foo": "bar",
			"zoo": 123.456,
		}),
		reflect.TypeFor[map[string]any](),
		map[string]any{
			"foo": "bar",
			"zoo": 123.456,
		},
	)

	{
		_, err := Convert(reflect.ValueOf(nil), reflect.TypeFor[*int]())
		assert.ErrorIs(t, err, ErrCannotConvert)
	}
}

func doTestConvert(t *testing.T, sourceValue reflect.Value, targetType reflect.Type, expectedValue any) {
	v, err := Convert(sourceValue, targetType)
	if assert.NoError(t, err) {
		if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Float64 {
			assert.InEpsilonSlice(t, expectedValue, v.Interface(), epsilon)
		} else if v.Kind() == reflect.Float64 {
			assert.InEpsilon(t, expectedValue, v.Interface(), epsilon)
		} else {
			assert.Equal(t, expectedValue, v.Interface())
		}
	}
}
