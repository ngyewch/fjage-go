package param

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type GenericValue struct {
	Value any
}

type internalGenericValue struct {
	Clazz string `json:"clazz,omitempty"`
	Data  []byte `json:"data,omitempty"`
}

type internalGenericValueStrings struct {
	Clazz string   `json:"clazz,omitempty"`
	Data  []string `json:"data,omitempty"`
}

func (v *GenericValue) UnmarshalJSON(data []byte) error {
	var igv internalGenericValue
	err := json.Unmarshal(data, &igv)
	if err == nil {
		r := bytes.NewReader(igv.Data)
		switch igv.Clazz {
		case "[B":
			v.Value = igv.Data
		case "[D":
			var values []float64
			var n float64
			for {
				err = binary.Read(r, binary.LittleEndian, &n)
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}
				}
				values = append(values, n)
			}
			v.Value = values
		case "[F":
			var values []float32
			var n float32
			for {
				err = binary.Read(r, binary.LittleEndian, &n)
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}
				}
				values = append(values, n)
			}
			v.Value = values
		case "[I":
			var values []int32
			var n int32
			for {
				err = binary.Read(r, binary.LittleEndian, &n)
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}
				}
				values = append(values, n)
			}
			v.Value = values
		case "[J":
			var values []int64
			var n int64
			for {
				err = binary.Read(r, binary.LittleEndian, &n)
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}
				}
				values = append(values, n)
			}
			v.Value = values
		default:
			return fmt.Errorf("unknown clazz %q", igv.Clazz)
		}
		return nil
	}

	var igvStrings internalGenericValueStrings
	err = json.Unmarshal(data, &igvStrings)
	if err == nil {
		switch igvStrings.Clazz {
		case "[Ljava.lang.String;":
			v.Value = igvStrings.Data
		default:
			return fmt.Errorf("unknown clazz %q", igv.Clazz)
		}
		return nil
	}

	var iv any
	err = json.Unmarshal(data, &iv)
	if err == nil {
		return err
	}
	v.Value = iv

	return nil
}

func (v *GenericValue) MarshalJSON() ([]byte, error) {
	if (v == nil) || (v.Value == nil) {
		return json.Marshal(nil)
	}
	t := reflect.TypeOf(v.Value)
	switch t.Kind() {
	case reflect.Slice:
		return v.marshalJSONArray()
	case reflect.Array:
		return v.marshalJSONArray()
	case reflect.Int32:
		return json.Marshal(v.Value.(int32))
	case reflect.Int64:
		return json.Marshal(v.Value.(int64))
	case reflect.Float32:
		return json.Marshal(v.Value.(float32))
	case reflect.Float64:
		return json.Marshal(v.Value.(float64))
	case reflect.String:
		return json.Marshal(v.Value.(string))
	case reflect.Bool:
		return json.Marshal(v.Value.(bool))
	default:
		return nil, fmt.Errorf("cannot marshal type: %T", v.Value)
	}
}

func (v *GenericValue) marshalJSONArray() ([]byte, error) {
	switch arrayElements := v.Value.(type) {
	case []byte:
		igv := internalGenericValue{
			Clazz: "[B",
			Data:  arrayElements,
		}
		return json.Marshal(igv)
	case []int32:
		var b []byte
		for _, arrayElement := range arrayElements {
			b = binary.LittleEndian.AppendUint32(b, uint32(arrayElement))
		}
		igv := internalGenericValue{
			Clazz: "[I",
			Data:  b,
		}
		return json.Marshal(igv)
	case []int64:
		var b []byte
		for _, arrayElement := range arrayElements {
			b = binary.LittleEndian.AppendUint64(b, uint64(arrayElement))
		}
		igv := internalGenericValue{
			Clazz: "[J",
			Data:  b,
		}
		return json.Marshal(igv)
	case []float32:
		b := bytes.NewBuffer(nil)
		for _, arrayElement := range arrayElements {
			err := binary.Write(b, binary.LittleEndian, arrayElement)
			if err != nil {
				return nil, err
			}
		}
		igv := internalGenericValue{
			Clazz: "[F",
			Data:  b.Bytes(),
		}
		return json.Marshal(igv)
	case []float64:
		b := bytes.NewBuffer(nil)
		for _, arrayElement := range arrayElements {
			err := binary.Write(b, binary.LittleEndian, arrayElement)
			if err != nil {
				return nil, err
			}
		}
		igv := internalGenericValue{
			Clazz: "[D",
			Data:  b.Bytes(),
		}
		return json.Marshal(igv)
	case []string:
		return json.Marshal(internalGenericValueStrings{
			Clazz: "[Ljava.lang.String;",
			Data:  arrayElements,
		})
	default:
		return nil, fmt.Errorf("cannot marshal type: %T", v.Value)
	}
}
