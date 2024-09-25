package GTI_OBD_Scanner

import (
	"io"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

const (
	defaultFileExt  = "log"
	defaultFileName = "GIT-OBD-Scanner"
)

type WriterOptions struct {
	Path    string
	Level   Level
	MaxAge  time.Duration
	MaxSize int64
	CutRule CutRule
}

func NewWriter(opts WriterOptions) (io.Writer, error) {
	var (
		fileExt      string
		fileName     string
		rotationTime time.Duration
		srcFileParts = make([]string, 0, 3)
		newFileParts = make([]string, 0, 4)
	)

	// 构建源文件名
	path, file := filepath.Split(opts.Path)
	list := strings.Split(file, ".")
	switch c := len(list); c {
	case 1:
		if list[0] == "" {
			fileName, fileExt = defaultFileName, defaultFileExt
		} else {
			fileName, fileExt = list[0], defaultFileExt
		}
	case 2:
		fileName, fileExt = list[0], list[1]
	default:
		fileName, fileExt = strings.Join(list[:c-1], "."), list[c-1]
	}

	srcFileParts = append(srcFileParts, fileName)
	newFileParts = append(newFileParts, fileName)

	// 如果指定日志级别，将日志级别添加到文件名中
	if opts.Level != 0 {
		srcFileParts = append(srcFileParts, strings.ToLower(opts.Level.String()))
		newFileParts = append(newFileParts, strings.ToLower(opts.Level.String()))
	}

	// 根据 CutRule 切割日志
	switch opts.CutRule {
	case CutByYear:
		newFileParts = append(newFileParts, "%Y")
		rotationTime = 365 * 24 * time.Hour
	case CutByMonth:
		newFileParts = append(newFileParts, "%Y%m")
		rotationTime = 31 * 24 * time.Hour
	case CutByDay:
		newFileParts = append(newFileParts, "%Y%m%d")
		rotationTime = 24 * time.Hour
	case CutByHour:
		newFileParts = append(newFileParts, "%Y%m%d%H")
		rotationTime = time.Hour
	case CutByMinute:
		newFileParts = append(newFileParts, "%Y%m%d%H%M")
		rotationTime = time.Minute
	case CutBySecond:
		newFileParts = append(newFileParts, "%Y%m%d%H%M%S")
		rotationTime = time.Second
	}

	srcFileParts = append(srcFileParts, fileExt)
	newFileParts = append(newFileParts, fileExt)

	srcFileName := filepath.Join(path, strings.Join(srcFileParts, "."))
	newFileName := filepath.Join(path, strings.Join(newFileParts, "."))

	// 日志轮转
	options := make([]rotatelogs.Option, 0, 4)
	options = append(options, rotatelogs.WithLinkName(srcFileName))
	if opts.MaxAge > 0 {
		options = append(options, rotatelogs.WithMaxAge(opts.MaxAge))
	}
	if opts.MaxSize > 0 {
		options = append(options, rotatelogs.WithRotationSize(opts.MaxSize))
	}
	if rotationTime > 0 {
		options = append(options, rotatelogs.WithRotationTime(rotationTime))
	}

	return rotatelogs.New(newFileName, options...)
}
