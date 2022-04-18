package spin_lock

import (
	"runtime"
	"sync"
	"sync/atomic"
)

const maxBackoff = 16

type spinLock uint32

func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

func (sl *spinLock)Lock(){
	backoff := 1
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		// 退避指數，最多 2^4 方次，就是16 就每次都等這個退避指數
		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
}


func NewSpinLock()sync.Locker{
	return new(spinLock)
}