package lock

import (
	"github.com/DavidYou21st/go-tools/utils"
	"sync"
	"sync/atomic"
)

const none = -1

// 可重入锁
type locker struct {
	heldCount int64      // 锁被持有的次数
	heldBy    int64      // 锁被持有的协程Id
	cond      *sync.Cond // 释放条件，利用 sync.Cond 实现
}

func NewReEntrantLock() sync.Locker {
	return &locker{heldBy: none}
}

// Lock locks the locker with reentrant mode
func (l *locker) Lock() {
	// get current goroutine id
	gid := utils.GetCurrentGoRoutineId()
	// if current goroutine is the lock owner
	if l.tryOwnerAcquire(gid) {
		return
	}
	// if lock is not held by any goroutine
	if l.tryNoneAcquire(gid) {
		return
	}
	// if current goroutine is not the lock owner
	l.acquireSlow(gid)
}

// Unlock unlocks the locker using reentrant mode
func (l *locker) Unlock() {
	// only the owner goroutine can enter this method
	heldBy := atomic.LoadInt64(&l.heldBy)
	if heldBy == none {
		panic("unlock of unlocked reentrantLocker")
	}
	if heldBy == utils.GetCurrentGoRoutineId() {
		l.releaseOnce()
		return
	}
	panic("unlock by different goroutine")
}

// releaseOnce release once on locker
func (l *locker) releaseOnce() {
	if l.heldCount == 0 {
		panic("unlocks more than locks")
	}
	l.heldCount--
	if l.heldCount == 0 {
		l.heldBy = none
		// signal one and only one goroutine
		l.cond.Signal()
	}
	return
}

// tryOwnerAcquire acquires the lock if and only if lock is currently held by same goroutine
func (l *locker) tryOwnerAcquire(gid int64) bool {
	// CAS gid -> gid, 若成功则代表当前线程持有锁
	if atomic.CompareAndSwapInt64(&l.heldBy, gid, gid) {
		atomic.AddInt64(&l.heldCount, 1)
		return true
	}
	return false
}

// tryNoneAcquire uses CAS to try to swap l.heldBy and l.heldCount field
func (l *locker) tryNoneAcquire(gid int64) bool {
	// none = -1
	// CAS none -> gid, 若成功则代表原本没有线程持有锁，当前线程获取到了锁
	if atomic.CompareAndSwapInt64(&l.heldBy, none, gid) {
		atomic.AddInt64(&l.heldCount, 1)
		return true
	}
	return false
}

// acquireSlow waits for cond to signal
func (l *locker) acquireSlow(gid int64) {
	l.cond.L.Lock()
	// 若当前还是有其他线程持有锁
	for atomic.LoadInt64(&l.heldBy) != none {
		// 继续等待唤醒
		l.cond.Wait()
	}
	// 若当前没有线程持有锁
	l.acquireLocked(gid)
	l.cond.L.Unlock()
}

// acquireLocked does acquire by current goroutine while locked
func (l *locker) acquireLocked(gid int64) {
	// 由于sync.Cond的Signal方法每次只会唤醒一个线程，所以这里直接替换即可
	atomic.SwapInt64(&l.heldBy, gid)
	atomic.SwapInt64(&l.heldCount, 1)
}
