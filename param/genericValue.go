package param

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
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
					break
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
					break
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
					break
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
					break
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
