package encoder

import (
	"fmt"
	"github.com/gti-obd-scanner/log/internal/utils"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"strings"
)

const StackFlag = "_stack"

const (
	fieldKeyLevel = "level" // 日志级别
	fieldKeyTime  = "time"  // 时间戳
	fieldKeyFile  = "file"  // 文件名和行号
	fieldKeyMsg   = "msg"   // 消息内容
	// 堆栈信息
	fieldKeyStack     = "stack"
	fieldKeyStackFunc = "func"
	fieldKeyStackFile = "file"
)

type JsonEncoder struct {
	zapcore.ObjectEncoder
	bufferPool     buffer.Pool // 缓冲池
	timeFormat     string      // 时间戳格式
	callerFullPath bool        // 是否全路径
}

func NewJsonEncoder(timeFormat string, callerFullPath bool) zapcore.Encoder {
	return &JsonEncoder{
		bufferPool:     buffer.NewPool(),
		timeFormat:     timeFormat,
		callerFullPath: callerFullPath,
	}
}

// Clone 创建并返回当前 Encoder 的一个副本, 保证线程安全, 允许多个 goroutine 并发使用同一个 Logger 而不会相互干扰
func (e *JsonEncoder) Clone() zapcore.Encoder {
	return nil
}

// EncodeEntry 日志编码, 将一个 ent 日志条目和 fields 相关字段编码成 buffer 日志输出格式
func (e *JsonEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := e.bufferPool.Get()
	stack := false

	// 检查是否需要包含堆栈信息
	if len(fields) > 0 && fields[0].Key == StackFlag && fields[0].Type == zapcore.BoolType {
		stack = utils.Bool(fields[0].Integer)
	}

	// 将日志级别转换为大写字符串
	var levelText string
	switch ent.Level {
	case zapcore.DPanicLevel:
		levelText = zapcore.PanicLevel.CapitalString()
	default:
		levelText = ent.Level.CapitalString()
	}

	// 构建 JSON
	line.AppendByte('{')
	line.AppendString(fmt.Sprintf(`"%s":"%s"`, fieldKeyLevel, levelText))
	line.AppendString(fmt.Sprintf(`,"%s":"%s"`, fieldKeyTime, ent.Time.Format(e.timeFormat)))

	if ent.Caller.Defined {
		var file string
		if e.callerFullPath {
			file = ent.Caller.File
		} else {
			_, file = filepath.Split(ent.Caller.File)
		}
		line.AppendString(fmt.Sprintf(`,"%s":"%s"`, fieldKeyFile, fmt.Sprintf("%s:%d", file, ent.Caller.Line)))
	}

	line.AppendString(fmt.Sprintf(`,"%s":"%s"`, fieldKeyMsg, utils.AddSlashes(strings.TrimSuffix(ent.Message, "\n"))))

	if stack && ent.Stack != "" {
		line.AppendString(fmt.Sprintf(`,"%s":[`, fieldKeyStack))

		stacks := strings.Split(ent.Stack, "\n")
		for i := range stacks {
			if i%2 == 0 {
				if i/2 == 0 {
					line.AppendString(fmt.Sprintf(`{"%s":"%s"`, fieldKeyStackFunc, stacks[i]))
				} else {
					line.AppendString(fmt.Sprintf(`,{"%s":"%s"`, fieldKeyStackFunc, stacks[i]))
				}
			} else {
				line.AppendString(fmt.Sprintf(`,"%s":"%s"}`, fieldKeyStackFile, strings.TrimPrefix(stacks[i], "\t")))
			}
		}
		line.AppendByte(']')
	}

	line.AppendByte('}')
	line.AppendString("\n")

	return line, nil
}
