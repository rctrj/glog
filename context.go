package log

import (
	"context"
	"fmt"
)

type ctxExtractor func(ctx context.Context) interface{}

var ctxExtractorMap = make(map[fmt.Stringer]ctxExtractor)

func RegisterExtractor(key fmt.Stringer, value ctxExtractor) {
	ctxExtractorMap[key] = value
}

func extractEntriesFromCtx(ctx context.Context) []*entry {
	var arr []*entry
	for k, extractor := range ctxExtractorMap {
		val := extractor(ctx)
		if val == nil {
			continue
		}
		arr = append(arr, Entry(k.String(), val))
	}
	return arr
}
