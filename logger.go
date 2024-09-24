package GTI_OBD_Scanner

import (
	"github.com/gti-obd-scanner/log/internal/encoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.SugaredLogger
	opts   *options
}

func NewLogger(opts ...Option) *Logger {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	// Zap 将日志数据转换为最终输出格式的关键组件
	var (
		fileEncoder     zapcore.Encoder // 输出到文件
		terminalEncoder zapcore.Encoder // 输出到终端
	)

	switch o.format {
	case JsonFormat:
		fileEncoder = encoder.NewJsonEncoder(o.timeFormat, o.callerFullPath)
		terminalEncoder = fileEncoder
	default:
		fileEncoder = encoder.NewTextEncoder(o.timeFormat, o.callerFullPath, false)
		terminalEncoder = encoder.NewTextEncoder(o.timeFormat, o.callerFullPath, true)
	}

	return nil
}
