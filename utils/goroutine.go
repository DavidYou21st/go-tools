package utils

import (
	"runtime"
	"strconv"
	"strings"
)

// GetCurrentGoRoutineId 从当前线程的栈信息中抽取协程id
func GetCurrentGoRoutineId() int64 {
	var buf [64]byte
	// 获取栈信息
	n := runtime.Stack(buf[:], false)
	// 抽取id
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine"))[0]
	// 转为64位整数
	id, _ := strconv.Atoi(idField)
	return int64(id)
}
