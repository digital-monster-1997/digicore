package worker_pool

import (
	"errors"
	"github.com/digital-monster-1997/digicore/pkg/utils/spin_lock"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct{
	// capacity 池子大小
	capacity int32
	// running 正在執行中的 worker 數量
	running int32
	// lock 鎖的介面，可以自己實現，也可以交由別人實現
	lock sync.Locker
	// goroutine 的 workers
	workers workerQueue
	// status 用於通知池子關閉
	status int32
	// cond 用來取得一個空閑的 goroutine
	cond *sync.Cond
	// workerCache 放入 Sync Pool 當中當作快取
	workerCache sync.Pool
	// blockingNum 在 pool.submit 阻塞的上限值，受到 pool.lock 保護
	blockingNum int

	stopHeartBeat chan struct{}
	options Options
}

// ---------------------------------------------- 公有函數區 ----------------------------------------------

// IsClosed 判斷池子是否已關閉
func (p *Pool)IsClosed()bool{
	return atomic.LoadInt32(&p.status) == CLOSED
}

// Running 看看目前有多少正在運作的 goroutine
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

// Cap 返回目前池子的容量大小
func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// Free 取得目前閒置的 goroutine 數量
func(p *Pool)Free() int{
	c := p.Cap()
	if c<0{
		return -1
	}
	return c - p.Running()
}

// Submit 提交任務
func (p *Pool)Submit(task func())error{
	if p.IsClosed(){
		// 池子已關閉，不接受提交
		return ErrPoolClosed
	}
	var w *Worker
	// 取出 worker
	if w = p.retrieveWorker(); w==nil{
		return ErrPoolOverload
	}
	w.task <- task
	return nil
}

// ChangeSize 調整池子大小，對於 pre alloc 的池子無效
func(p *Pool)ChangeSize(size int){
	// 取得池子大小
	cap := p.Cap()
	// 如果池子的大小為 -1 或者要調整的數量是負數，或者調整的數量跟池子大小一樣，不用改變
	// 預先建立的池子大小也不可改變
	if cap == -1 || size <=0 || size == cap || p.options.PreAlloc{
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
	// 如果條小，正在執行的執行完成後，沒有人在使用它，時間到了會自行淘汰
	// 如果條大，要通知正在阻塞的，有新的 worker 可以用了，隨機挑選一個來使用
	if size > cap {
		if size-cap == 1{
			p.cond.Signal()
			return
		}
		p.cond.Broadcast()
	}
}

// Release 關閉這個池子，並釋放工作隊列
func(p *Pool)Release(){
	// 如果交換失敗了，就直接回，可能已經關閉了，不用再關一次
	if !atomic.CompareAndSwapInt32(&p.status, OPENED,CLOSED){
		return
	}
	// 有一些人可能在 retrieveWorker() 等待喚醒，以防止那些呼叫的人無限阻塞
	p.lock.Lock()
	p.workers.reset()
	p.lock.Unlock()
	p.cond.Broadcast()
}

// ReleaseTimeout is like Release but with a timeout, it waits all workers to exit before timing out.
func (p *Pool) ReleaseTimeout(timeout time.Duration) error {
	if p.IsClosed() {
		return errors.New("pool is already closed")
	}
	select {
	case p.stopHeartBeat <- struct{}{}:
		<-p.stopHeartBeat
	default:
	}
	p.Release()
	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		if p.Running() == 0 {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return ErrTimeout
}

// Reboot reboots a closed pool.
func (p *Pool) Reboot() {
	if atomic.CompareAndSwapInt32(&p.status, CLOSED, OPENED) {
		go p.purgePeriodically()
	}
}

// ---------------------------------------------- 私有函數區 ----------------------------------------------

// purgePeriodically 定期清除過期的 worker
func(p *Pool)purgePeriodically(){
	heartBeat := time.NewTicker(p.options.ExpiryDuration)
	defer heartBeat.Stop()
	for{
		select{
		case <-heartBeat.C:
		case <-p.stopHeartBeat:
			p.stopHeartBeat <- struct{}{}
			return
		}
		// 關閉了就不用清了，直接停止
		if p.IsClosed(){
			break
		}
		// 把過期與沒過期切分開ㄌㄞ˙
		p.lock.Lock()
		expiredWorkers := p.workers.retrieveExpiry(p.options.ExpiryDuration)
		p.lock.Unlock()
		// 逐個清理過期的
		for i := range expiredWorkers{
			expiredWorkers[i].task <- nil
			expiredWorkers[i] = nil
		}
		if p.Running() ==0 {
			p.cond.Broadcast()
		}
	}


}
// revertWorker 放回一個 worker to pool
func (p *Pool) revertWorker(w *Worker) bool{
	if capacity := p.Cap();(capacity>0 && p.Running() > capacity) || p.IsClosed(){
		p.cond.Broadcast()
		return false
	}
	w.recycleTime = time.Now()
	p.lock.Lock()
	if p.IsClosed(){
		p.lock.Unlock()
		return false
	}
	err := p.workers.insert(w)
	if err != nil {
		p.lock.Unlock()
		return false
	}
	p.cond.Signal()
	p.lock.Unlock()
	return true
}

// retrieveWorker returns an available worker to run the tasks.
func (p *Pool) retrieveWorker() (w *Worker) {
	spawnWorker := func() {
		w = p.workerCache.Get().(*Worker)
		w.run()
	}

	p.lock.Lock()

	w = p.workers.detach()
	if w != nil { // first try to fetch the worker from the queue
		p.lock.Unlock()
	} else if capacity := p.Cap(); capacity == -1 || capacity > p.Running() {
		// if the worker queue is empty and we don't run out of the pool capacity,
		// then just spawn a new worker goroutine.
		p.lock.Unlock()
		spawnWorker()
	} else { // otherwise, we'll have to keep them blocked and wait for at least one worker to be put back into pool.
		if p.options.NonBlocking {
			p.lock.Unlock()
			return
		}
	retry:
		if p.options.MaxBlockingTasks != 0 && p.blockingNum >= p.options.MaxBlockingTasks {
			p.lock.Unlock()
			return
		}
		p.blockingNum++
		p.cond.Wait() // block and wait for an available worker
		p.blockingNum--
		var nw int
		if nw = p.Running(); nw == 0 { // awakened by the scavenger
			p.lock.Unlock()
			if !p.IsClosed() {
				spawnWorker()
			}
			return
		}
		if w = p.workers.detach(); w == nil {
			if nw < capacity {
				p.lock.Unlock()
				spawnWorker()
				return
			}
			goto retry
		}

		p.lock.Unlock()
	}
	return
}

func (p *Pool)incRunning(){
	atomic.AddInt32(&p.running,1)
}

func (p *Pool)decRunning(){
	atomic.AddInt32(&p.running,-1)
}

// ---------------------------------------------- 新建一個 pool ----------------------------------------------


// NewPool ...
func NewPool(size int, options ...Option) (*Pool, error) {
	opts := loadOptions(options...)

	if size <= 0 {
		size = -1
	}

	if expiry := opts.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.ExpiryDuration = DefaultCleanIntervalTime
	}


	p := &Pool{
		capacity:      int32(size),
		lock:         	spin_lock.NewSpinLock(),
		stopHeartBeat: make(chan struct{}, 1),
		options:       *opts,
	}
	p.workerCache.New = func() interface{} {
		return &Worker{
			pool: p,
			task: make(chan func(), workerChanCap),
		}
	}
	if p.options.PreAlloc {
		if size == -1 {
			return nil, ErrInvalidPreAllocSize
		}
		p.workers = newWorkerQueue(loopQueueType, size)
	} else {
		p.workers = newWorkerQueue(stackType, 0)
	}

	p.cond = sync.NewCond(p.lock)

	// Start a goroutine to clean up expired workers periodically.
	go p.purgePeriodically()

	return p, nil
}