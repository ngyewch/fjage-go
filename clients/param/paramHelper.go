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

const (
	parameterTagName = "fjage.parameter"
)

type Helper struct {
	gw gateway.Gateway
}

func NewParamHelper(gw gateway.Gateway) *Helper {
	return &Helper{
		gw: gw,
	}
}

func (helper *Helper) newMessage(performative fjage.Performative, recipient string) (*fjage.Message, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &fjage.Message{
		MsgID:        id.String(),
		Performative: performative,
		Recipient:    recipient,
		Sender:       helper.gw.AgentID(),
		SentAt:       time.Now().UnixMilli(),
	}, nil
}

func (helper *Helper) request(ctx context.Context, agentID string, initializer func(req *param.ParameterReq)) (*param.ParameterRsp, error) {
	header, err := helper.newMessage(fjage.PerformativeRequest, agentID)
	if err != nil {
		return nil, err
	}
	req := &param.ParameterReq{
		Message: header,
	}
	initializer(req)
	sendResponse, err := helper.gw.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	if sendResponse.Message.Header().Performative != fjage.PerformativeInform {
		return nil, fjage.NewPerformativeError(sendResponse.Message.Header().Performative)
	}
	rsp, ok := sendResponse.Message.(*param.ParameterRsp)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", sendResponse.Message)
	}
	return rsp, nil
}

func (helper *Helper) GetParam(ctx context.Context, agentID string, name string, v any) error {
	target := reflect.ValueOf(v)
	rsp, err := helper.request(ctx, agentID, func(req *param.ParameterReq) {
		req.Param = name
	})
	if err != nil {
		return err
	}
	if rsp.Value == nil {
		return fmt.Errorf("missing param value")
	}
	source := reflect.ValueOf(rsp.Value.Value)
	err = setValue(source, target)
	if err != nil {
		return err
	}
	return nil
}

func setValue(source reflect.Value, target reflect.Value) error {
	if source.Kind() == reflect.Pointer {
		return fmt.Errorf("source cannot be a pointer")
	}
	if source.Kind() == target.Kind() {
		target.Set(source)
		return nil
	}
	if target.Kind() == reflect.Pointer {
		if target.CanSet() {
			if source.Type().AssignableTo(target.Type().Elem()) {
				sourcePointer := reflect.New(target.Type().Elem())
				sourcePointer.Elem().Set(source)
				target.Set(sourcePointer)
				return nil
			} else if source.Type().ConvertibleTo(target.Type().Elem()) {
				sourcePointer := reflect.New(target.Type().Elem())
				sourcePointer.Elem().Set(source.Convert(target.Type().Elem()))
				target.Set(sourcePointer)
				return nil
			}
		} else if target.Elem().CanSet() {
			if source.Type().AssignableTo(target.Type().Elem()) {
				target.Elem().Set(source)
				return nil
			} else if source.Type().ConvertibleTo(target.Type().Elem()) {
				target.Elem().Set(source.Convert(target.Type().Elem()))
				return nil
			}

		}
	} else {
		if source.Type().AssignableTo(target.Type()) {
			target.Set(source)
		} else if source.Type().ConvertibleTo(target.Type()) {
			target.Set(source.Convert(target.Type()))
		}
	}
	return fmt.Errorf("param value (%v) not assignable to target (%v)", source.Type(), target.Type())
}

func (helper *Helper) GetParamsAndHandle(ctx context.Context, agentID string, names []string, paramHandler func(name string, value any) error) error {
	rsp, err := helper.request(ctx, agentID, func(req *param.ParameterReq) {
		for _, name := range names {
			req.Requests = append(req.Requests, param.ParameterReqEntry{
				Param: name,
			})
		}
	})
	if err != nil {
		return err
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
		parameterTagValue := field.Tag.Get(parameterTagName)
		parts := strings.Split(parameterTagValue, ",")
		parameterName := parts[0]
		if parameterName != "" {
			names = append(names, parameterName)
			fieldMap[parameterName] = field
		}
	}
	err := helper.GetParamsAndHandle(ctx, agentID, names, func(name string, value any) error {
		field, ok := fieldMap[name]
		if !ok {
			return nil
		}
		sourceValue := reflect.ValueOf(value)
		targetValue := structValue.FieldByName(field.Name)
		err := setValue(sourceValue, targetValue)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (helper *Helper) GetParams(ctx context.Context, agentID string) (map[string]any, error) {
	paramMap := make(map[string]any)
	err := helper.GetParamsAndHandle(ctx, agentID, []string{}, func(name string, value any) error {
		paramMap[name] = value
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paramMap, nil
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
	_, err := helper.request(ctx, agentID, func(req *param.ParameterReq) {
		req.Param = name
		req.Value = &param.GenericValue{
			Value: v,
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (helper *Helper) SetParamsFromStruct(ctx context.Context, agentID string, v any) error {
	var inputValue reflect.Value
	inputType := reflect.TypeOf(v)
	if (inputType.Kind() == reflect.Pointer) && (inputType.Elem().Kind() == reflect.Struct) {
		inputValue = reflect.ValueOf(v).Elem()
		inputType = inputValue.Type()
	} else if inputType.Kind() == reflect.Struct {
		inputValue = reflect.ValueOf(v)
	} else {
		return fmt.Errorf("input value must be a struct (or a pointer to a struct)")
	}
	_, err := helper.request(ctx, agentID, func(req *param.ParameterReq) {
		for i := 0; i < inputType.NumField(); i++ {
			field := inputType.Field(i)
			parameterTagValue := field.Tag.Get(parameterTagName)
			parts := strings.Split(parameterTagValue, ",")
			parameterName := parts[0]
			if parameterName != "" {
				fieldValue := inputValue.Field(i)
				if fieldValue.Kind() == reflect.Pointer {
					if fieldValue.IsZero() {
						continue
					}
					fieldValue = fieldValue.Elem()
				}
				req.Requests = append(req.Requests, param.ParameterReqEntry{
					Param: parameterName,
					Value: &param.GenericValue{
						Value: fieldValue.Interface(),
					},
				})
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}
