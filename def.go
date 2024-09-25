package GTI_OBD_Scanner

// Level 日志级别
type Level int

// Format 日志输出格式
type Format int

// CutRule 日志切割规则
type CutRule int

const (
	NoneLevel  Level = iota // NONE
	DebugLevel              // DEBUG
	InfoLevel               // INFO
	WarnLevel               // WARN
	ErrorLevel              // ERROR
	FatalLevel              // FATAL
	PanicLevel              // PANIC
)

const (
	TextFormat Format = iota // 文本格式
	JsonFormat               // JSON格式
)

const (
	CutByYear   CutRule = iota + 1 // 按照年切割
	CutByMonth                     // 按照月切割
	CutByDay                       // 按照日切割
	CutByHour                      // 按照时切割
	CutByMinute                    // 按照分切割
	CutBySecond                    // 按照秒切割
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	}
	return "NONE"
}
