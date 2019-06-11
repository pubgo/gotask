package gotask

import (
	"log"
	"runtime"
	"sync"
	"time"
)

func NewTask(max int, maxDur time.Duration) *Task {
	_t := &Task{
		max:     max,
		maxDur:  maxDur,
		q:       make(chan *_taskFn, max),
		_curDur: make(chan time.Duration, max),
		_stopQ:  make(chan error),
		wg: &_WaitGroup{
			_done: make(chan bool, max),
			wg:    &sync.WaitGroup{},
		},
	}
	go _t._handle()
	return _t
}

type Task struct {
	maxDur time.Duration

	curDur  time.Duration
	_curDur chan time.Duration

	max int

	q chan *_taskFn

	_stopQ chan error
	_stop  error

	wg *_WaitGroup
}

func (t *Task) Wait() {
	t.wg.Wait()
}

func (t *Task) Do(f TaskFn, args ...interface{}) error {
	for {
		if t._stop != nil {
			return t._stop
		}

		if t.wg.Len() < t.max && t.curDur < t.maxDur {
			t.q <- f(args...)
			t.wg.Add()
			return nil
		}

		if t.wg.Len() < runtime.NumCPU()*2 {
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
		case __fn := <-t.q:
			go func(_fn *_taskFn) {
				t._curDur <- _FnCost(func() {
					err := _KTry(_fn.fn, _fn.args...)
					if err == nil {
						return
					}

					if len(_fn.efn) == 0 || _fn.efn[0] == nil {
						return
					}

					if _err := _KTry(_fn.efn[0], err); _err != nil {
						t._stopQ <- _KTry(_fn.efn[0], err)
					}
				})

				t.wg.Done()
			}(__fn)
		case _c := <-t._curDur:
			t.curDur = t.curDur/2 + _c/2
		case t._stop = <-t._stopQ:
		}
	}
}
