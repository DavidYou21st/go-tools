package utils

import (
	"strconv"
	"strings"
)

// GetMajorVer 获取主版本号
func GetMajorVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	if first == 0 {
		return 0, nil
	}
	return strconv.ParseUint(v[:first], 10, 64)
}

// GetMinorVer 获取次版本号
func GetMinorVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

// GetRevisionVer 获取修订版本号
func GetRevisionVer(v string) (uint64, error) {
	last := strings.LastIndexByte(v, '.')
	if last == 0 {
		return 0, nil
	}
	return strconv.ParseUint(v[last+1:], 10, 64)
}
