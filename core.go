package log

import (
	"context"
	"fmt"
)

type core interface {
	Log(ctx context.Context, level Level, message string, entries []IEntry)
}

type IEntry interface {
	Key() string
	Value() interface{}
}

type entry struct {
	key           string
	value         interface{}
	isMinimalInfo bool
}

func (e *entry) Key() string {
	return e.key
}

func (e *entry) Value() interface{} {
	return e.value
}

func Entry(key string, value interface{}) *entry {
	return &entry{
		key:           key,
		value:         value,
		isMinimalInfo: false,
	}
}

func MinimalEntry(key string, value interface{}) *entry {
	return &entry{
		key:           key,
		value:         value,
		isMinimalInfo: true,
	}
}

func ErrorEntry(err error) *entry {
	return &entry{
		key:           "error",
		value:         err,
		isMinimalInfo: true,
	}
}

func (e *entry) convertToString(jsonMarshaller func(interface{}) ([]byte, error)) (string, error) {
	var printValue interface{}

	if errValue, ok := e.value.(error); ok {
		printValue = errValue
	} else {
		json, err := jsonMarshaller(e.value)
		if err != nil {
			return "", err
		}
		printValue = string(json)
	}

	return fmt.Sprintf("{\"%v\": %v}", e.key, printValue), nil
}
