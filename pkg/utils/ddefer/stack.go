package ddefer

import "sync"

type DeferTask func()error

// DeferStack 放到這個裡面的最後需要被一一執行
type DeferStack struct {
	// fns 要放入 stack 的 function
	fns 	[]DeferTask
	sync.RWMutex
}

// Push 推入 Stack 當中
func(d *DeferStack)Push(tasks ...DeferTask){
	d.Lock()
	defer d.Unlock()
	d.fns = append(d.fns, tasks...)
}

// DoAllTask 執行所有 stack 當中的資料
func(d *DeferStack)DoAllTask(){
	d.Lock()
	defer d.Unlock()
	// 倒序輸出
	for item := len(d.fns) -1; item>=0; item --{
		_ = d.fns[item]()
	}
}

func NewStack() * DeferStack {
	return &DeferStack{
		fns: make([]DeferTask, 0),
	}
}