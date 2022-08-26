package fmtcore

import (
	"context"
	log "github.com/rctrj/glog"
	"testing"
)

func Test_fmtLog_Log(t *testing.T) {
	conf := NewConfig("debug", "minimal", true, true)
	logger := New(conf)
	ctx := context.Background()

	logger.Info(ctx, "Test Info", log.Entry("Test_Entry", 0), log.MinimalEntry("Min_entry", struct {
		Key   string `json:"key"`
		Param int    `json:"param"`
	}{
		Key:   "asdf",
		Param: 100,
	}))

	logger.Debug(ctx, "Test Debug", log.Entry("Test_Entry", 0), log.MinimalEntry("Min_entry", struct {
		Key   string `json:"key"`
		Param int    `json:"param"`
	}{
		Key:   "asdf",
		Param: 100,
	}))
}
