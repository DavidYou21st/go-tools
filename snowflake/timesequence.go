package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	startTimestamp = int64(1517937700000) // 默认起始的时间戳 1517937700000， 计算时，减去这个值
	districtIdBits = uint(2)              //区域 所占位数
	nodeIdBits     = uint(7)              //节点 所占位数
	sequenceBits   = uint(10)             //自增ID 所占位数

	/*
	 * 1位 符号位  |  41位时间戳                                   | 2位区域位  |  3位时钟序列    |  7位节点       | 10 （毫秒内）自增ID
	 * 0          | 0 00000000 00000000 00000000 00000000 00000000 | 00  	| 	000   	   |  0000000 		| 0000000000
	 *
	 */
	maxNodeId     = -1 ^ (-1 << nodeIdBits)     //节点 ID 最大范围
	maxDistrictId = -1 ^ (-1 << districtIdBits) //最大区域范围

	nodeIdShift        = sequenceBits //左移次数
	districtIdShift    = sequenceBits + nodeIdBits
	timestampLeftShift = sequenceBits + nodeIdBits + districtIdBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
	maxNextIdsNum      = 100 //单次获取ID的最大数量
)

type IDGenerator struct {
	sequence       int64 //序号
	lastTimestamp  int64 //最后一次请求的时间戳
	nodeId         int64 //节点ID
	startTimestamp int64
	districtId     int64
	mutex          sync.Mutex
}

// NewIDGenerator return a snowflake id generator
func NewIDGenerator(NodeId int64) (*IDGenerator, error) {
	var districtId int64
	districtId = 1 //暂时默认给1 ，方便以后扩展
	generator := &IDGenerator{}
	if NodeId > maxNodeId || NodeId < 0 {
		return nil, errors.New(fmt.Sprintf("NodeId Id: %d error", NodeId))
	}
	if districtId > maxDistrictId || districtId < 0 {
		return nil, errors.New(fmt.Sprintf("District Id: %d error", districtId))
	}
	generator.nodeId = NodeId
	generator.districtId = districtId
	generator.lastTimestamp = -1
	generator.sequence = 0
	generator.startTimestamp = startTimestamp
	generator.mutex = sync.Mutex{}

	return generator, nil
}

// timeGen generate a unix millisecond.
func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// tilNextMillis spin wait till next millisecond.
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

// NextID get a snowflake id.
func (id *IDGenerator) NextID() (int64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	return id.nextID()
}

func (id *IDGenerator) nextID() (int64, error) {
	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		return 0, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp))
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp
	return ((timestamp - id.startTimestamp) << timestampLeftShift) | (id.districtId << districtIdShift) | (id.nodeId << nodeIdShift) | id.sequence, nil
}
