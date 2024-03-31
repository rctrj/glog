package main

import (
	"context"
	"github.com/rctrj/glog"
	"github.com/rctrj/glog/drivers/fmt"
)

func main() {
	driver := fmt.New(
		fmt.WithPrettyPrint(true),
		fmt.WithMode(fmt.ModeJson),
	)
	logger := log.New(
		driver,
		log.WithLevel(log.LevelTrace),
		log.WithFieldStackTrace(true),
	)
	ctx := context.Background()

	logger.Trace(ctx, "Trace Message", log.KV("some_key", "some_value"), log.KVContextual("contextual_key", "contextual_value"))
	logger.Debug(ctx, "Debug Message", log.KV("some_key", "some_value"), log.KVContextual("contextual_key", "contextual_value"))
	logger.Info(ctx, "Info Message", log.KV("some_key", "some_value"), log.KVContextual("contextual_key", "contextual_value"))
	logger.Warn(ctx, "Warn Message", log.KV("some_key", "some_value"), log.KVContextual("contextual_key", "contextual_value"))
	logger.Error(ctx, "Error Message", log.KV("some_key", "some_value"), log.KVContextual("contextual_key", "contextual_value"))
}
