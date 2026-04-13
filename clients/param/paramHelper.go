package param

import (
	"context"
	"fmt"
	"reflect"
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
	rsp, ok := sendResponse.Message.(*param.ParameterRsp)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", sendResponse.Message)
	}
	if rsp.Values == nil {
		return fmt.Errorf("missing param values")
	}
	{
		target, ok := nameValueMap[rsp.Param]
		if ok {
			targetValue := reflect.ValueOf(target)
			sourceValue := reflect.ValueOf(rsp.Value.Value)
			if !sourceValue.Type().AssignableTo(targetValue.Type().Elem()) {
				return fmt.Errorf("param value (%v) not assignable to target (%v)", sourceValue.Type(), targetValue.Type())
			}
			targetValue.Elem().Set(sourceValue)
		}
	}
	for name, value := range rsp.Values {
		target, ok := nameValueMap[name]
		if !ok {
			continue
		}
		targetValue := reflect.ValueOf(target)
		sourceValue := reflect.ValueOf(value.Value)
		if !sourceValue.Type().AssignableTo(targetValue.Type().Elem()) {
			return fmt.Errorf("param value (%v) not assignable to target (%v)", sourceValue.Type(), targetValue.Type())
		}
		targetValue.Elem().Set(sourceValue)
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
