package queue

import "errors"

// 环状队列

type LoopQueue struct {
	member []any

	// length & cap
	cap int
	// index
	front, rear int
}

func NewLoopQueue(size int) *LoopQueue {
	return &LoopQueue{
		member: make([]any, size),
		cap:    size,
		front:  0,
		rear:   0,
	}
}

func (q *LoopQueue) IsEmpty() bool {
	return q.front == q.rear
}

func (q *LoopQueue) IsFull() bool {
	return (q.rear+1)%q.cap == q.front
}

func (q *LoopQueue) Push(val any) error {
	if q.IsFull() {
		return errors.New("queue is full")
	}

	q.member[q.rear] = val
	// 当队尾达到最大index就不能简单自增而是要循环
	q.rear = (q.rear + 1) % q.cap

	return nil
}

func (q *LoopQueue) Pop() (any, error) {
	if q.IsEmpty() {
		return nil, errors.New("empty queue")
	}

	pop := q.member[q.front]
	// 当队头达到最大index就不能简单自增而是要循环
	q.front = (q.front + 1) % q.cap

	return pop, nil
}
