package types

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteArray(t *testing.T) {
	data := []byte("abcd")
	base64EncodedData := base64.StdEncoding.EncodeToString(data)
	expectedEncodedData := fmt.Sprintf("{\"clazz\":\"[B\",\"data\":\"%s\"}", base64EncodedData)
	actualEncodedData, err := ByteArray(data).MarshalJSON()
	if assert.NoError(t, err) {
		assert.Equal(t, []byte(expectedEncodedData), actualEncodedData)

		expectedByteArray := ByteArray(data)
		{
			var actualByteArray ByteArray
			err = actualByteArray.UnmarshalJSON(actualEncodedData)
			if assert.NoError(t, err) {
				assert.Equal(t, expectedByteArray, actualByteArray)
			}
		}
		{
			var actualByteArray ByteArray
			err = actualByteArray.UnmarshalJSON([]byte("[97,98,99,100]"))
			if assert.NoError(t, err) {
				assert.Equal(t, expectedByteArray, actualByteArray)
			}
		}
	}
}
