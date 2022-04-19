package ddefer

import (
	"testing"
)

func TestRegister(t *testing.T){
	nd := NewStack()
	state := "mr. tiger,"
	task1 := func()error{
		state+= " go in la?"
		return nil
	}
	task2 := func()error{
		state+= " are you"
		return nil
	}
	task3 := func()error{
		state+= " where"
		return nil
	}

	nd.Push(task1, task2)
	nd.Push(task3)
	nd.DoAllTask()
	want := "mr. tiger, where are you go in la?"
	if state != want{
		t.Fatalf("Stack has error,want:%v ret:%v", want, state)
	}
}