package param

import (
	"context"
	"errors"
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
	err = SetValue(source, target)
	if err != nil {
		return err
	}
	return nil
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
			p := strings.LastIndex(name, ".")
			if p >= 0 {
				shortName := name[p+1:]
				field, ok = fieldMap[shortName]
				if !ok {
					return nil
				}
			} else {
				return nil
			}
		}
		sourceValue := reflect.ValueOf(value)
		if sourceValue.IsValid() {
			targetValue := structValue.FieldByName(field.Name)
			err := SetValue(sourceValue, targetValue)
			if err != nil {
				if errors.Is(err, ErrCannotConvert) {
					return fmt.Errorf("cannot set field %s (type=%s, param=%s) with value %v (%T): %w", field.Name, field.Type, name, value, value, err)
				}
				return err
			}
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
			if value != nil {
				err = consumer(name, value.Value)
				if err != nil {
					return err
				}
			} else {
				err = consumer(name, nil)
				if err != nil {
					return err
				}
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
