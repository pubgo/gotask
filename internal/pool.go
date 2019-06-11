package internal

import (
	"sync"
)

func NewWaitGroup(wg *sync.WaitGroup, done chan bool) *WaitGroup {
	return &WaitGroup{wg: wg, _done: done}
}

type WaitGroup struct {
	wg    *sync.WaitGroup
	_done chan bool
}

func (t *WaitGroup) Add() {
	t.wg.Add(1)
}

func (t *WaitGroup) Done() {
	t.wg.Done()
}

func (t *WaitGroup) Len() int {
	return len(t._done)
}

func (t *WaitGroup) Wait() {
	t.wg.Wait()
}
