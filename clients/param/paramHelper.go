package param

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/gateway"
	"github.com/ngyewch/fjage-go/param"
)

type Helper struct {
	gw gateway.Gateway
}

func NewParamHelper(gw gateway.Gateway) *Helper {
	return &Helper{
		gw: gw,
	}
}

func (helper *Helper) GetParam(ctx context.Context, agentID string, name string, v any) error {
	target := reflect.ValueOf(v)
	if target.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot set target value")
	}
	if !target.Elem().CanSet() {
		return fmt.Errorf("cannot set target value")
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	req := &param.ParameterReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    agentID,
			Sender:       helper.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Param: name,
	}
	sendResponse, err := helper.gw.Send(ctx, req)
	if err != nil {
		return err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeInform {
		return fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	rsp, ok := sendResponse.Message.(*param.ParameterRsp)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", sendResponse.Message)
	}
	if rsp.Value == nil {
		return fmt.Errorf("missing param value")
	}
	source := reflect.ValueOf(rsp.Value.Value)
	if !source.Type().AssignableTo(target.Type().Elem()) {
		return fmt.Errorf("param value (%v) not assignable to target (%v)", source.Type(), target.Type())
	}
	target.Elem().Set(source)
	return nil
}

func (helper *Helper) GetParamsAndHandle(ctx context.Context, agentID string, names []string, paramHandler func(name string, value any) error) error {
	var requests []param.ParameterReqEntry
	for _, name := range names {
		requests = append(requests, param.ParameterReqEntry{
			Param: name,
		})
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	req := &param.ParameterReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    agentID,
			Sender:       helper.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Requests: requests,
	}
	sendResponse, err := helper.gw.Send(ctx, req)
	if err != nil {
		return err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeInform {
		return fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	rsp, ok := sendResponse.Message.(*param.ParameterRsp)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", sendResponse.Message)
	}
	err = iterateParameters(rsp, paramHandler)
	if err != nil {
		return err
	}
	return nil
}

func (helper *Helper) GetParamsAndPopulate(ctx context.Context, agentID string, v any) error {
	var names []string
	inputType := reflect.TypeOf(v)
	if inputType.Kind() != reflect.Pointer {
		return fmt.Errorf("value must be a pointer")
	}
	structType := inputType.Elem()
	structValue := reflect.ValueOf(v).Elem()
	fieldMap := make(map[string]reflect.StructField)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonTag := field.Tag.Get("json")
		parts := strings.Split(jsonTag, ",")
		jsonFieldName := parts[0]
		if jsonFieldName != "" {
			names = append(names, jsonFieldName)
			fieldMap[jsonFieldName] = field
		}
	}
	err := helper.GetParamsAndHandle(ctx, agentID, names, func(name string, value any) error {
		field, ok := fieldMap[name]
		if !ok {
			return nil
		}
		sourceValue := reflect.ValueOf(value)
		targetValue := structValue.FieldByName(field.Name)
		if field.Type.Kind() == reflect.Pointer {
			if !sourceValue.Type().AssignableTo(field.Type.Elem()) {
				return fmt.Errorf("param value (%v) not assignable to target (%v)", sourceValue.Type(), targetValue.Type())
			}
			sourcePointerValue := reflect.New(sourceValue.Type())
			sourcePointerValue.Elem().Set(sourceValue)
			targetValue.Set(sourcePointerValue)
		} else {
			if !sourceValue.Type().AssignableTo(field.Type) {
				return fmt.Errorf("param value (%v) not assignable to target (%v)", sourceValue.Type(), targetValue.Type())
			}
			targetValue.Set(sourceValue)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (helper *Helper) GetParams(ctx context.Context, agentID string, nameValueMap map[string]any) error {
	var requests []param.ParameterReqEntry
	for name, value := range nameValueMap {
		target := reflect.ValueOf(value)
		if target.Kind() != reflect.Ptr {
			return fmt.Errorf("%s: cannot set target value", name)
		}
		if !target.Elem().CanSet() {
			return fmt.Errorf("%s: cannot set target value", name)
		}
		requests = append(requests, param.ParameterReqEntry{
			Param: name,
		})
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	req := &param.ParameterReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    agentID,
			Sender:       helper.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Requests: requests,
	}
	sendResponse, err := helper.gw.Send(ctx, req)
	if err != nil {
		return err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeInform {
		return fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	nameValueMapIsEmpty := len(nameValueMap) == 0
	rsp, ok := sendResponse.Message.(*param.ParameterRsp)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", sendResponse.Message)
	}
	handleParameter := func(name string, value any) error {
		if nameValueMapIsEmpty {
			nameValueMap[name] = value
		} else {
			target, ok := nameValueMap[name]
			if ok {
				targetValue := reflect.ValueOf(target)
				sourceValue := reflect.ValueOf(value)
				if !sourceValue.Type().AssignableTo(targetValue.Type().Elem()) {
					return fmt.Errorf("param value (%v) not assignable to target (%v)", sourceValue.Type(), targetValue.Type())
				}
				targetValue.Elem().Set(sourceValue)
			}
		}
		return nil
	}
	err = iterateParameters(rsp, handleParameter)
	if err != nil {
		return err
	}
	return nil
}

func iterateParameters(rsp *param.ParameterRsp, consumer func(name string, value any) error) error {
	err := consumer(rsp.Param, rsp.Value.Value)
	if err != nil {
		return err
	}
	if rsp.Values != nil {
		for name, value := range rsp.Values {
			err = consumer(name, value.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (helper *Helper) SetParam(ctx context.Context, agentID string, name string, v any) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	req := &param.ParameterReq{
		Message: &fjage.Message{
			MsgID:        id.String(),
			Performative: fjage.PerformativeRequest,
			Recipient:    agentID,
			Sender:       helper.gw.AgentID(),
			SentAt:       time.Now().UnixMilli(),
		},
		Param: name,
		Value: &param.GenericValue{
			Value: v,
		},
	}
	sendResponse, err := helper.gw.Send(ctx, req)
	if err != nil {
		return err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeInform {
		return fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	_, ok := sendResponse.Message.(*param.ParameterRsp)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", sendResponse.Message)
	}
	return nil
}
