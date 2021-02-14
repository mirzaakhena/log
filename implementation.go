package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/ksuid"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type contextType string

const funcConst = "func"
const operationIDConst = "opID"
const operationIDField = contextType(operationIDConst)

func setJSONFormat() {
	format = &logrus.JSONFormatter{
		TimestampFormat: "0102 150405.000",
	}
}

func setSimpleFormat() {
	format = &nested.Formatter{
		NoColors:        true,
		HideKeys:        true,
		TimestampFormat: "0102 150405.000",
		FieldsOrder:     []string{"operationID", "func"},
	}
}

var (
	singleton       sync.Once
	defaultLogger   logrusImpl
	useFile         bool
	path            string
	name            string
	maxAge          int
	format          logrus.Formatter
	operationIDFunc func() string
)

type logrusImpl struct {
	theLogger *logrus.Logger
}

func setFile(pathFile, nameFile string, maxAgeInDays int) {
	path = pathFile
	name = nameFile
	maxAge = maxAgeInDays
	useFile = true
}

func initFormat() {
	if format == nil {
		setJSONFormat()
	}
}

func initOperationIDFunc() {
	if operationIDFunc == nil {
		operationIDFunc = func() string {
			return ksuid.New().String()
		}
	}
}

func getLogImpl() Logger {
	singleton.Do(func() {

		initFormat()

		initOperationIDFunc()

		defaultLogger = logrusImpl{theLogger: logrus.New()}
		defaultLogger.theLogger.SetFormatter(format)

		if !useFile {
			return
		}

		writer, _ := rotatelogs.New(
			fmt.Sprintf("%s/logs/%s.log.%s", path, name, "%Y%m%d"),
			rotatelogs.WithLinkName(fmt.Sprintf("%s/%s.log", path, name)),
			rotatelogs.WithMaxAge(time.Duration(maxAge*24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(1*24)*time.Hour),
		)

		defaultLogger.theLogger.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.ErrorLevel: writer,
			},
			defaultLogger.theLogger.Formatter,
		))

	})

	return &defaultLogger
}

func getFunctionCall(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	x := strings.LastIndex(funcName, "/")
	return fmt.Sprintf("%s:%d", funcName[x+1:], line)
}

func fetchOperationID(ctx context.Context) string {

	operationIDInterface := ctx.Value(operationIDField)
	if operationIDInterface == nil {
		return "-"
	}

	operationID, ok := operationIDInterface.(string)
	if !ok {
		return "-"
	}

	return operationID
}

func (x *logrusImpl) additionalField(ctx context.Context) *logrus.Entry {

	var theLogger *logrus.Entry

	theLogger = x.theLogger.WithContext(ctx)

	operationID := fetchOperationID(ctx)
	if operationID != "" {
		theLogger = theLogger.WithField(operationIDConst, operationID)
	}

	funcCall := getFunctionCall(4)
	if funcCall != "" {
		theLogger = theLogger.WithField(funcConst, funcCall)
	}

	return theLogger
}

func (x *logrusImpl) Info(ctx context.Context, message string, args ...interface{}) {
	logMessage := fmt.Sprintf(message, args...)
	x.additionalField(ctx).Info(logMessage)
}

func (x *logrusImpl) Error(ctx context.Context, message string, args ...interface{}) {
	logMessage := fmt.Sprintf(message, args...)
	x.additionalField(ctx).Error(logMessage)
}
