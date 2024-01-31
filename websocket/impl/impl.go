package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Connection struct {
	wsConn       *websocket.Conn
	inChannel    chan []byte
	outChannel   chan []byte
	closeChannel chan byte
	isClose      bool
	mutex        sync.Mutex
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:       wsConn,
		inChannel:    make(chan []byte, 1000),
		outChannel:   make(chan []byte, 1000),
		closeChannel: make(chan byte, 1),
	}
	//读协程
	go conn.readLoop()
	//写协程
	go conn.writeLoop()
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChannel:
	case <-conn.closeChannel:
		err = errors.New("connection 已关闭")
	}
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChannel <- data:
	case <-conn.closeChannel:
		err = errors.New("connection 已关闭")
	}
	return
}

func (conn *Connection) Close() {
	conn.wsConn.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClose {
		close(conn.closeChannel)
		conn.isClose = true
	}
	conn.mutex.Unlock()
}

// 将请求数据读入缓存
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}

		//阻塞在这里，等待inChan有空闲位置
		select {
		case conn.inChannel <- data:
		case <-conn.closeChannel:
			goto ERR
		}
	}
ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		data = <-conn.outChannel
		if conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}

func (conn *Connection) Heartbeat() {
	var (
		err error
	)
	for {
		if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
			return
		}
		time.Sleep(5 * time.Second)
	}
}
