package tools

import (
	"fmt"
	"time"
)

// DateToMillisecondTimestamp 将日期格式（如 2026-3-17）转换为毫秒时间戳
// 参数 dateStr: 日期字符串，格式为 "2006-1-2"（例如："2026-3-17"）
// 返回：毫秒时间戳和错误信息
func DateToMillisecondTimestamp(dateStr string) (int64, error) {
	parsedTime, err := time.Parse("2006-1-2", dateStr)
	if err != nil {
		return 0, fmt.Errorf("解析日期失败：%v", err)
	}

	// 转换为毫秒时间戳
	return parsedTime.UnixMilli(), nil
}

// DatesToMillisecondTimestamps 将多个日期格式（如 2026-3-17）转换为毫秒时间戳切片
// 参数 dateStrs: 日期字符串切片，格式为 "2006-1-2"（例如：[]string{"2026-3-17", "2026-3-18"}）
// 返回：毫秒时间戳切片和错误信息
func DatesToMillisecondTimestamps(dateStrs []string) ([]int64, error) {
	timestamps := make([]int64, 0, len(dateStrs))

	for _, dateStr := range dateStrs {
		parsedTime, err := time.Parse("2006-1-2", dateStr)
		if err != nil {
			return nil, fmt.Errorf("解析日期 %s 失败：%v", dateStr, err)
		}
		timestamps = append(timestamps, parsedTime.UnixMilli())
	}

	return timestamps, nil
}
