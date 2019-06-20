package internal

import (
	"sync"
)

func NewWaitGroup(wg *sync.WaitGroup, done chan bool) *WaitGroup {
	return &WaitGroup{wg: wg, _done: done, mtx: &sync.Mutex{}}
}

type WaitGroup struct {
	mtx   *sync.Mutex
	wg    *sync.WaitGroup
	_done chan bool
}

func (t *WaitGroup) Add() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.wg.Add(1)
	t._done <- true
}

func (t *WaitGroup) Done() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.wg.Done()
	<-t._done
}

func (t *WaitGroup) Len() int {
	return len(t._done)
}

func (t *WaitGroup) Wait() {
	t.wg.Wait()
}
