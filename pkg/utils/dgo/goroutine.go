package dgo

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sync"
)

func RecoverGo(task func(), recoverFn func())(ret error){
	// 如果有傳，revover 的 function ，退出前執行
	if recoverFn != nil {
		defer recoverFn()
	}
	// 退出前執行
	defer func(){
		if err := recover(); err != nil{
			_,file,line,_ := runtime.Caller(2)
			// TODO 換成自己的格式
			log.Fatalf("%s:%d", file, line)
			if _,ok:= err.(error); ok{
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
			ret = errors.New(fmt.Sprintf("%s, %s:%d",ret,runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name(), line))
		}
	}()
	// 執行 func
	task()
	return nil
}

// Serial 序列化執行，一個做完換一個
func Serial(tasks ...func()) func(){
	return func(){
		for _, task := range tasks{
			task()
		}
	}
}

// Parallel  併發執行
func Parallel(tasks ...func())func(){
	var wg sync.WaitGroup
	return func(){
		wg.Add(len(tasks))
		for _, task := range tasks{
			go RecoverGo(task,wg.Done)
		}
	}
}