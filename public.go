package log

import (
	"context"
)

// Info is
func Info(ctx context.Context, message string, args ...interface{}) {
	getLogImpl().Info(ctx, message, args...)
}

// Warn is
func Warn(ctx context.Context, message string, args ...interface{}) {
	getLogImpl().Warn(ctx, message, args...)
}

// Error is
func Error(ctx context.Context, message string, args ...interface{}) {
	getLogImpl().Error(ctx, message, args...)
}

// UseRotateFile is
func UseRotateFile(path, name string, maxAgeInDays int) {
	setFile(path, name, maxAgeInDays)
}

func UseJSONFormat() {
	setJSONFormat()
}

func UseSimpleFormat() {
	setSimpleFormat()
}

func SetRpcIDFunc(f func() string) {
	rpcidFunc = f
}

func ContextWithRpcID(ctx context.Context) context.Context {
	rpcIDInterface := ctx.Value(rpcidField)
	if rpcIDInterface == nil {
		initRpcFunc()
		return context.WithValue(ctx, rpcidField, rpcidFunc())
	}
	return ctx
}
