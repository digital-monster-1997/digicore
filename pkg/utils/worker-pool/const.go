package worker_pool

import (
	"errors"
	"runtime"
	"time"
)

var (
	// errQueueIsFull 如果 Worker 的可存放序列，已經最大了，無法再放更多，返回這個錯誤
	errQueueIsFull = errors.New("the queue is full")
	// errQueueIsReleased 如果嘗試把一個新任務放盡已經要被 release的 queue 當中
	errQueueIsReleased = errors.New("the queue length is zero")
	// ErrPoolClosed 提交任務時，池子已是關閉狀態
	ErrPoolClosed = errors.New("this pool has been closed")
	// ErrPoolOverload 池子已滿，且沒有可以用的 worker 可以用了
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")
	ErrTimeout = errors.New("operation timed out")
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	ErrInvalidPreAllocSize = errors.New("can not set up a negative capacity under PreAlloc mode")
	workerChanCap = func() int {
		// Use blocking channel if GOMAXPROCS=1.
		// This switches context from sender to receiver immediately,
		// which results in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the sender might be dragged down if the receiver is CPU-bound.
		return 1
	}()
)

const (
	// OPENED 池子的狀態是打開的
	OPENED = iota

	// CLOSED 池子目前的狀態是關起來的
	CLOSED
)

const (
	DefaultCleanIntervalTime = time.Second
)