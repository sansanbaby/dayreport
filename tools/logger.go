package tools

import (
	"fmt"
	"runtime"
)

// LogError 统一的错误日志输出函数
// 输出格式：文件名 + 方法名/函数名 + 行号
func LogError(err error) error {
	if err == nil {
		return nil
	}

	// 获取调用者信息（跳过 runtime.Callers 和 LogError 本身）
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("ERROR: %v", err)
	}

	// 获取函数名
	funcName := runtime.FuncForPC(pc).Name()
	if funcName == "" {
		return fmt.Errorf("[%s:%d] ERROR: %v", file, line, err)
	}

	// 提取包名和函数名
	return fmt.Errorf("[%s:%d] %s ERROR: %v", file, line, funcName, err)
}

// LogErrorf 带格式化的错误日志输出函数
func LogErrorf(format string, args ...interface{}) error {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf(format, args...)
	}

	funcName := runtime.FuncForPC(pc).Name()
	if funcName == "" {
		return fmt.Errorf("[%s:%d] "+format, append([]interface{}{file, line}, args...)...)
	}

	return fmt.Errorf("[%s:%d] %s "+format, append([]interface{}{file, line, funcName}, args...)...)
}
