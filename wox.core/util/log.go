package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logInstance *Log
var logOnce sync.Once

type Log struct {
	logger *zap.Logger
	writer io.Writer
}

func GetLogger() *Log {
	logOnce.Do(func() {
		logFolder := GetLocation().GetLogDirectory()
		logInstance = CreateLogger(logFolder)
		setCrashOutput(logInstance)
	})
	return logInstance
}

func setCrashOutput(logInstance *Log) {
	logFile := path.Join(GetLocation().GetLogDirectory(), "crash.log")
	// Open file in append mode instead of create
	crashFile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logInstance.Error(context.Background(), fmt.Sprintf("failed to open crash log file: %s", err.Error()))
		return
	}
	defer crashFile.Close()

	// Verify write permission
	if err := crashFile.Chmod(0644); err != nil {
		logInstance.Error(context.Background(), fmt.Sprintf("failed to set crash log file permission: %s", err.Error()))
		return
	}

	setCrashOutputErr := debug.SetCrashOutput(crashFile, debug.CrashOptions{})
	if setCrashOutputErr != nil {
		logInstance.Error(context.Background(), fmt.Sprintf("failed to set crash output: %s", setCrashOutputErr.Error()))
		return
	}
}

func CreateLogger(logFolder string) *Log {
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		os.MkdirAll(logFolder, os.ModePerm)
	}

	logImpl := &Log{}
	logImpl.logger, logImpl.writer = createLogger(logFolder)
	log.SetFlags(0) // remove default timestamp
	log.SetOutput(logImpl.writer)
	return logImpl
}

func (l *Log) GetWriter() io.Writer {
	return logInstance.writer
}

func GetContextTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value(ContextKeyTraceId).(string); ok {
		return traceId
	}

	return ""
}

func GetContextComponentName(ctx context.Context) string {
	if componentName, ok := ctx.Value(ContextKeyComponentName).(string); ok {
		return componentName
	}

	return "Wox"
}

func formatMsg(context context.Context, msg string, level string) string {
	var builder strings.Builder
	builder.Grow(256)
	builder.WriteString(FormatTimestampWithMs(GetSystemTimestamp()))
	builder.WriteString(" G")
	builder.WriteString(LeftPad(strconv.FormatInt(GetGID(), 10), 7, '0'))
	builder.WriteString(" ")
	if traceId := GetContextTraceId(context); traceId != "" {
		builder.WriteString(traceId)
		builder.WriteString(" ")
	}
	builder.WriteString(fmt.Sprintf("[%s] ", level))
	if componentName := GetContextComponentName(context); componentName != "" {
		builder.WriteString(fmt.Sprintf("[%s] ", componentName))
	}
	builder.WriteString(msg)
	return builder.String()
}

func (l *Log) Debug(context context.Context, msg string) {
	l.logger.Debug(formatMsg(context, msg, "DBG"))
}

func (l *Log) Warn(context context.Context, msg string) {
	l.logger.Warn(formatMsg(context, msg, "WRN"))
}

func (l *Log) Info(context context.Context, msg string) {
	l.logger.Info(formatMsg(context, msg, "INF"))
}

func (l *Log) Error(context context.Context, msg string) {
	l.logger.Error(formatMsg(context, msg, "ERR"))
}

func createLogger(logFolder string) (*zap.Logger, io.Writer) {
	writeSyncer := zapcore.AddSync(&Lumberjack{
		Filename:  path.Join(logFolder, "log"),
		LocalTime: true,
		MaxSize:   500, // megabytes
		MaxAge:    3,   // days
	})

	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = nil
	cfg.EncodeLevel = nil

	zapLogger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		writeSyncer,
		zap.DebugLevel,
	))

	reader, writer := io.Pipe()
	Go(NewTraceContext(), "log reader", func() {
		defer reader.Close()
		defer writer.Close()
		buf := make([]byte, 2048)
		for {
			n, err := reader.Read(buf)
			if err != nil {
				break
			}
			msg := string(buf[:n])
			//remove newline in msg
			msg = strings.TrimRight(msg, "\n")
			zapLogger.Info(formatMsg(NewTraceContext(), fmt.Sprintf("[SYS LOG] %s", msg), "INF"))
		}
	})

	return zapLogger, writer
}
