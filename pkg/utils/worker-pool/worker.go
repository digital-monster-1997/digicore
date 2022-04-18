package worker_pool

import (
	"log"
	"time"
)

type Worker struct{
	// pool 這個 worker 所屬的池子
	pool *Pool
	// task 欲執行的任務
	task chan func()
	// recycleTime 將worker 放回池子時，會更新，
	// 表示最後執行完成的時間，
	// 如果現在時間 - recycleTime > ExpiryDuration 的話，這個就要回
	recycleTime time.Time
}

// 創造新的 worker 時，要一起執行 run ，一個 worker 理應只執行一次 run
func(w *Worker)run(){
	// 將 running（有任務正在執行的計算） 數量加一
	w.pool.incRunning()
	// 啟動這個 worker 然後執行任務
	go func(){
		// 結束前執行
		defer func(){
			// 將執行任務的人的數量減少一
			w.pool.decRunning()
			// 把worker 放入快取持，一段時間後再返回原本的
			// 如果有事都會先去快取池，取資料速度較快
			w.pool.workerCache.Put(w)
			// 如果有壞掉的話
			if p:= recover(); p != nil{
				// 如果有登記錯誤控制的 handler 給他用專屬的，沒有就直接印出來
				if ph := w.pool.options.PanicHandler; ph != nil{
					ph(p)
				} else {
					log.Printf("worker exits from a panic: %v\n", p)
				}
			}
			// 隨機從阻塞的隊列當中喚醒一個
			w.pool.cond.Signal()
		}()
		// 執行任務，for range 一個 chan ，在還沒讀取到下一個資料前，就會阻塞，如果沒人放資料就會造成 deadlock!
		for f := range w.task{
			// 如果沒任務了，就把自己放回池子當中
			if f == nil{
				return
			}
			f()
			// 執行完任務，主動把自己放回池子中，如果沒有放回去，就走 defer 放回去的流程
			if ok := w.pool.revertWorker(w); !ok{
				return
			}
		}
	}()
}
// -------------------

type QueueType int

const (
	stackType QueueType = 1 << iota
	loopQueueType
)

type workerQueue interface {
	// len worker 長度
	len() int
	// isEmpty worker 是不是空的
	isEmpty() bool
	// insert 放回 queue當中
	insert(worker *Worker) error
	// detach 從 queue 當中取出一個worker
	detach() *Worker
	// retrieveExpiry 檢查到期的 worker
	retrieveExpiry(d time.Duration)[]*Worker
	reset()
}

func newWorkerQueue(qType QueueType, size int)workerQueue{
	switch qType {
	case stackType:
		return newWorkerStack(size)
	//case loopQueueType:
	//	return newWorkerLoopQueue(size)
	default:
		return newWorkerStack(size)
	}
}