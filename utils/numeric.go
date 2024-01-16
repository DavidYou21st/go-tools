package utils

import (
	"net"
	"strconv"
)

// IsNumeric 判断一个元素是否为数值类型
func IsNumeric(s any) bool {
	switch v := s.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(v, 64)
		return err == nil
	default:
		return false
	}
}

func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
