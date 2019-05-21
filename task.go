package gotask

import (
	"log"
	"runtime"
	"sync"
	"time"
)

func NewTask(max int, maxDur time.Duration) *Task {
	_t := &Task{
		max: max, maxDur:
		maxDur, q: make(chan *_taskFn, max),
		_curDur:   make(chan time.Duration, max),
		_done:     make(chan bool, max),
		_stopQ:    make(chan error),
		lock:      &sync.Mutex{},
	}
	go _t._handle()
	return _t
}

type Task struct {
	maxDur time.Duration

	curDur  time.Duration
	_curDur chan time.Duration

	_done chan bool
	max   int

	q chan *_taskFn

	_stopQ chan error
	_stop  error

	lock *sync.Mutex
}

func (t *Task) Wait() {
	for len(t._done) > 0 {
		if t._stop != nil {
			return
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (t *Task) Do(f TaskFn, args ...interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	for {
		if t._stop != nil {
			return t._stop
		}

		if len(t._done) < t.max && t.curDur < t.maxDur {
			t.q <- f(args...)
			t._done <- true
			return nil
		}

		if len(t._done) < runtime.NumCPU()*2 {
			t.curDur = 0
		}

		if Debug {
			log.Printf("q_l:%d cur_dur:%s max_q:%d max_dur:%s", len(t.q), t.curDur.String(), t.max, t.maxDur.String())
		}

		time.Sleep(time.Millisecond * 200)
	}
}

func (t *Task) _handle() {
	for t._stop == nil {
		select {
		case _fn := <-t.q:
			go func() {
				t._curDur <- _FnCost(func() {
					if err := _KTry(_fn.fn, _fn.args...); err != nil {
						if len(_fn.efn) != 0 && _fn.efn[0] != nil {
							if _err := _KTry(_fn.efn[0], err); _err != nil {
								t._stopQ <- _err
							}
						}
					}
				})
				<-t._done
			}()
		case _c := <-t._curDur:
			t.curDur = t.curDur/2 + _c/2
		case _e := <-t._stopQ:
			t._stop = _e
		}
	}
}
