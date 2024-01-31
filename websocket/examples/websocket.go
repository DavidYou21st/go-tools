package examples

import (
	"bytes"
	"github.com/DavidYou21st/go-tools/websocket/impl"
	"github.com/gorilla/websocket"
	"net/http"
)

//
//func main() {
//	http.HandleFunc("/ws", wsHandle)
//
//	http.ListenAndServe("0.0.0.0:8866", nil)
//}

// 定义转换器
var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
		conn   *impl.Connection
	)

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	//初始化连接
	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}
	//发心跳包
	go conn.Heartbeat()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}
ERR:
	//关闭连接操作
	conn.Close()
}

func BytesCombine1(pBytes ...[]byte) []byte {
	length := len(pBytes)
	s := make([][]byte, length)
	for index := 0; index < length; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
