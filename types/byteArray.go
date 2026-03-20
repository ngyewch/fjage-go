package types

import (
	"encoding/json"
	"fmt"
)

type ByteArray []byte

type byteArray struct {
	Clazz string `json:"clazz"`
	Data  []byte `json:"data"`
}

func (ba ByteArray) MarshalJSON() ([]byte, error) {
	if ba == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(byteArray{
		Clazz: "[B",
		Data:  ba,
	})
}

func (ba *ByteArray) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case 'n':
		if data[1] != 'u' || data[2] != 'l' || data[3] != 'l' {
			return fmt.Errorf("could not unmarshal byte array")
		}
		*ba = nil
		return nil

	case '{':
		var iba byteArray
		err := json.Unmarshal(data, &iba)
		if err != nil {
			return err
		}
		if iba.Clazz != "[B" {
			return fmt.Errorf("unexpected clazz value for byte array")
		}
		*ba = iba.Data
		return nil

	case '[':
		var b []byte
		err := json.Unmarshal(data, &b)
		if err != nil {
			return err
		}
		*ba = b
		return nil

	default:
		return fmt.Errorf("could not unmarshal byte array")
	}
}
