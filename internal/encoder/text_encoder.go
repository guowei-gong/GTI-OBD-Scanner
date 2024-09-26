package encoder

import (
	"fmt"
	"github.com/gti-obd-scanner/log/internal/utils"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type TextEncoder struct {
	zapcore.ObjectEncoder
	bufferPool     buffer.Pool // 缓冲池
	timeFormat     string      // 时间戳格式
	callerFullPath bool        // 是否全路径
	isTerminal     bool        // 是否输出到终端
}

// NewTextEncoder 终端环境中支持彩色输出
func NewTextEncoder(timeFormat string, callerFullPath, isTerminal bool) zapcore.Encoder {
	return &TextEncoder{
		bufferPool:     buffer.NewPool(),
		timeFormat:     timeFormat,
		callerFullPath: callerFullPath,
		isTerminal:     isTerminal,
	}
}

// Clone 创建并返回当前 Encoder 的一个副本, 保证线程安全, 允许多个 goroutine 并发使用同一个 Logger 而不会相互干扰
func (e *TextEncoder) Clone() zapcore.Encoder {
	return nil
}

// EncodeEntry 日志编码, 将一个 ent 日志条目和 fields 相关字段编码成 buffer 日志输出格式
func (e *TextEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := e.bufferPool.Get()
	stack := false

	// 检查是否需要包含堆栈信息
	if len(fields) > 0 && fields[0].Key == StackFlag && fields[0].Type == zapcore.BoolType {
		stack = utils.Bool(fields[0].Integer)
	}

	// 截取日志级别, 设置颜色
	levelText := ent.Level.CapitalString()[0:4]
	if e.isTerminal {
		var levelColor int
		switch ent.Level {
		case zapcore.DebugLevel:
			levelColor = gray
		case zapcore.WarnLevel:
			levelColor = yellow
		case zapcore.ErrorLevel, zapcore.FatalLevel, zapcore.PanicLevel:
			levelColor = red
		case zapcore.DPanicLevel:
			levelColor = red
			levelText = zapcore.PanicLevel.CapitalString()[0:4]
		case zapcore.InfoLevel:
			levelColor = blue
		default:
			levelColor = blue
		}
		line.AppendString(fmt.Sprintf("\x1b[%dm%s ", levelColor, levelText))
		// 添加格式化的时间戳
		line.AppendString(fmt.Sprintf("\x1b[0m[%s]", ent.Time.Format(e.timeFormat)))
	} else {
		line.AppendString(levelText)
		line.AppendString(fmt.Sprintf("[%s]", ent.Time.Format(e.timeFormat)))
	}

	// 添加调用者信息、文件名和行号
	if ent.Caller.Defined {
		if e.callerFullPath {
			// 是否使用完整路径
			line.AppendString(fmt.Sprintf(" %s:%d ", ent.Caller.File, ent.Caller.Line))
		} else {
			_, file := filepath.Split(ent.Caller.File)
			line.AppendString(fmt.Sprintf(" %s:%d ", file, ent.Caller.Line))
		}
	}

	// 日志正文
	line.AppendString(strings.TrimSuffix(ent.Message, "\n"))

	// 添加堆栈信息
	if stack && ent.Stack != "" {
		line.AppendByte('\n')
		line.AppendString("Stack:\n")

		stacks := strings.Split(ent.Stack, "\n")
		for i := range stacks {
			if i%2 == 0 {
				stacks[i] = strconv.Itoa(i/2+1) + ". " + stacks[i]
			}
		}
		line.AppendString(strings.Join(stacks, "\n"))
	}

	line.AppendString("\n")

	return line, nil
}
