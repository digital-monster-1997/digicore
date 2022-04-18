package worker_pool

import "time"

type workerStack struct {
	workers []*Worker
	expiry 	[]*Worker
	size 	int
}

// len worker 長度
func(w *workerStack) len()int{
	return len(w.workers)
}

// isEmpty worker 是不是空的
func(w *workerStack) isEmpty()bool{
	return len(w.workers) == 0
}

// insert(worker *Worker) error
func(w *workerStack) insert(worker *Worker) error{
	w.workers = append(w.workers, worker)
	return nil
}

//detach 從 queue 當中取出一個worker
func(w *workerStack) detach() *Worker{
	// 如果不為零，就從最後面拿出一個，如果為零，就說拿不到
	workerCount := w.len()
	if workerCount ==0{
		return nil
	}
	// 取出最後一個 worker
	worker := w.workers[workerCount-1]
	w.workers[workerCount-1] = nil // 避免內存外洩
	w.workers = w.workers[:workerCount -1]
	return worker
}

// retrieveExpiry 把過期的 worker 都清除，並返回沒有過期的
func(w *workerStack) retrieveExpiry(d time.Duration)[]*Worker{
	// 取得目前有多少 worker，如果為零，就不用清理了
	workerCount := w.len()
	if workerCount ==0{
		return nil
	}
	// 過期的時間點，目前往前多久是過期的時間
	expiryTime := time.Now().Add(-d)
	// 依照過期時間早晚排序完成的，找到第一個過期的 index 把往後的都刪除
	index := w.binarySearch(0, workerCount-1, expiryTime)
	// 先清除過期的
	w.expiry = w.expiry[:0]
	if index != -1{
		// 把過期的放到過期的陣列中
		w.expiry = append(w.expiry, w.workers[:index+1]...)
		m := copy(w.workers, w.workers[index+1:])
		for item := m; item< workerCount; item++{
			w.workers[item] = nil
		}
		w.workers = w.workers[:m]
	}
	return w.expiry
}

// binarySearch 找尋這個過期時間
func(w *workerStack) binarySearch(leftPointer, rightPointer int, expiryTime time.Time) int {
	var mid int
	for leftPointer <= rightPointer{
		mid = (leftPointer + rightPointer) /2
		if expiryTime.Before(w.workers[mid].recycleTime){
			rightPointer-=mid
		} else{
			leftPointer+=mid
		}
	}
	return rightPointer
}

func(w *workerStack) reset()  {
	for item := 0; item <w.len(); item++{
		w.workers[item].task <- nil
		w.workers[item] = nil
	}
	w.workers = w.workers[:0]
}

func newWorkerStack(size int) *workerStack {
	return &workerStack{
		workers: make([]*Worker, 0, size),
		size:  size,
	}
}