package common

import (
	"time"
)

func ParseInt64From(data map[string]interface{}, name string) int64 {
	if v, ok := data[name]; ok {
		return ParseInt64(v)
	}
	return 0
}

func ParseIntFrom(data map[string]interface{}, name string) int {
	if v, ok := data[name]; ok {
		return ParseInt(v)
	}
	return -1
}

func ParseInt64(value any) int64 {
	switch value.(type) {
	case int:
		return int64(value.(int))
	case float64:
		return int64(value.(float64))
	default:
		return value.(int64)
	}
}

func ParseInt(value any) int {
	switch value.(type) {
	case int:
		return value.(int)
	case float64:
		return int(value.(float64))
	default:
		return value.(int)
	}
}

func CopyArray[T any | string](source []T) []T {
	if len(source) == 0 {
		return nil
	}
	dst := make([]T, 5)
	copy(source, dst)
	return dst
}

func NowTime() string {
	tm := time.Now()
	return tm.Format("2006-01-02 15:04:05")
}
