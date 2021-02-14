package log

import (
	"context"
)

// Info is
func Info(ctx context.Context, message string, args ...interface{}) {
	getLogImpl().Info(ctx, message, args...)
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

func SetOperationIDFunc(f func() string) {
	operationIDFunc = f
}

func ContextWithOperationID(ctx context.Context) context.Context {
	opIDInterface := ctx.Value(operationIDField)
	if opIDInterface == nil {
		initOperationIDFunc()
		return context.WithValue(ctx, operationIDField, operationIDFunc())
	}
	return ctx
}
