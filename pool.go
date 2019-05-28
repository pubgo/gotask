package gotask

import "sync"

type _WaitGroup struct {
	wg    *sync.WaitGroup
	_done chan bool
}

func (t *_WaitGroup) Add() {
	t.wg.Add(1)
}

func (t *_WaitGroup) Done() {
	t.wg.Done()
}

func (t *_WaitGroup) Len() int {
	return len(t._done)
}

func (t *_WaitGroup) Wait() {
	t.wg.Wait()
}
