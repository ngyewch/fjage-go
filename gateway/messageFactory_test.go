package gateway

import (
	"encoding/json"
	"testing"

	"github.com/goforj/godump"
	"github.com/stretchr/testify/assert"
)

const testJSON = "{\"clazz\":\"org.arl.fjage.param.ParameterRsp\",\"data\":{\"index\":-1,\"values\":{\"org.arl.unet.device.DeviceParam.storage\":[470695096320,502116614144],\"org.arl.unet.device.DeviceParam.thermal\":{}},\"param\":\"org.arl.unet.device.DeviceParam.health\",\"value\":\"OK\",\"readonly\":[\"org.arl.unet.device.DeviceParam.storage\"],\"msgID\":\"019d8595-09cf-7f9f-9715-66cda836c9c0\",\"perf\":\"INFORM\",\"recipient\":\"gateway-1adfee69-dbe1-4544-95c3-ac3b9fc12533\",\"sender\":\"device\",\"inReplyTo\":\"c851cabf-9c22-4834-82a1-6740add5008f\",\"sentAt\":1776062630352}}"

func Test1(t *testing.T) {
	messageFactory := NewMessageFactory()
	var messageEnvelope MessageEnvelope
	err := json.Unmarshal([]byte(testJSON), &messageEnvelope)
	if assert.NoError(t, err) {
		m, err := messageFactory.UnmarshalMessage(&messageEnvelope)
		if assert.NoError(t, err) {
			godump.Dump(m)
		}
	}
}
