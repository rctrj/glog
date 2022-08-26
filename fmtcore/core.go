package fmtcore

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/rctrj/glog"
)

type (
	marshalFunc = func(interface{}) ([]byte, error)

	fmtCore struct {
		jsonMarshaller marshalFunc
	}
)

func newCore(isPrettyPrint bool) *fmtCore {
	jsonMarshaller := json.Marshal
	if isPrettyPrint {
		jsonMarshaller = func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "    ")
		}
	}

	return &fmtCore{jsonMarshaller: jsonMarshaller}
}

func (f fmtCore) Log(_ context.Context, level log.Level, message string, entries []log.IEntry) {
	interfaces := f.toInterfaceList(level.String(), message, entries)
	fmt.Println(interfaces...)
}

func (f fmtCore) toInterfaceList(level, message string, entries []log.IEntry) []interface{} {
	var interfaces []interface{}
	interfaces = append(interfaces, "["+level+"]")
	interfaces = append(interfaces, message)
	for _, entry := range entries {
		strVal, err := f.convertToString(entry)
		if err == nil {
			interfaces = append(interfaces, strVal)
		}
	}
	return interfaces
}

func (f fmtCore) convertToString(e log.IEntry) (string, error) {
	var printValue interface{}

	v := e.Value()
	if errValue, ok := v.(error); ok {
		printValue = errValue
	} else {
		j, err := f.jsonMarshaller(v)
		if err != nil {
			return "", err
		}
		printValue = string(j)
	}

	return fmt.Sprintf("{\"%v\": %v}", e.Key(), printValue), nil
}
