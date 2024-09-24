package GTI_OBD_Scanner

import (
	"time"
)

const (
	defaultFile              = "./log/gti.log"
	defaultLevel             = InfoLevel
	defaultFormat            = TextFormat
	defaultStdout            = true
	defaultFileMaxAge        = 7 * 24 * time.Hour
	defaultFileMaxSize       = 100
	defaultFileCutRule       = CutByDay
	defaultTimeFormat        = "2006/01/02 15:04:05.000000"
	defaultCallerFullPath    = false
	defaultClassifiedStorage = false
)

type Option func(o *options)

type options struct {
	file              string        // 输出的文件路径，有文件路径才会输出到文件，否则只会输出到终端
	level             Level         // 输出的最低日志级别，默认Info
	format            Format        // 输出的日志格式，Text或者Json，默认Text
	stdout            bool          // 是否输出到终端，debug模式下默认输出到终端
	timeFormat        string        // 时间格式，标准库时间格式，默认2006/01/02 15:04:05.000000
	stackLevel        Level         // 堆栈的最低输出级别，默认不输出堆栈
	fileMaxAge        time.Duration // 文件最大留存时间，默认7天
	fileMaxSize       int64         // 文件最大尺寸限制，单位（MB），默认100MB
	fileCutRule       CutRule       // 文件切割规则，默认按照天
	callerSkip        int           // 调用者跳过的层级深度
	callerFullPath    bool          // 是否启用调用文件全路径，默认短路径
	classifiedStorage bool          // 是否启用分级存储，默认不分级
}

func defaultOptions() *options {
	opts := &options{
		file:              defaultFile,
		level:             defaultLevel,
		format:            defaultFormat,
		stdout:            defaultStdout,
		timeFormat:        defaultTimeFormat,
		fileMaxAge:        defaultFileMaxAge,
		fileMaxSize:       defaultFileMaxSize,
		fileCutRule:       defaultFileCutRule,
		callerFullPath:    defaultCallerFullPath,
		classifiedStorage: defaultClassifiedStorage,
	}

	// 读取项目配置
	return opts
}
