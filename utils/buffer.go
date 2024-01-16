package utils

import (
	"bytes"
	"sync"
)

// 定义一个bytes.Buffer的对象缓冲池
var buffers sync.Pool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func GetBufFromPool() *bytes.Buffer {
	buf := buffers.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func PubBufToPool(buf *bytes.Buffer) {
	buffers.Put(buf)
}
