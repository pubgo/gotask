package gotask

import (
	"github.com/pubgo/errors"
	"github.com/pubgo/gotask/internal"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"
)

func TaskOf(fn interface{}, efn ...func(err error)) internal.TaskFn {
	defer errors.Handle(func() {})

	errors.T(errors.IsZero(fn) || reflect.TypeOf(fn).Kind() != reflect.Func, "fn error")
	errors.T(len(efn) != 0 && errors.IsZero(efn[0]), "efn is nil")

	return func(args ...interface{}) *internal.TaskFnDef {
		defer errors.Handle(func() {})

		var _log = errorLog
		if len(efn) != 0 {
			_log = efn[0]
		}

		return internal.NewTaskFn(fn, args, _log)
	}
}

func NewTask(max int, maxDur time.Duration) *Task {
	_t := &Task{
		max:     max,
		maxDur:  maxDur,
		q:       make(chan *internal.TaskFnDef, max),
		_curDur: make(chan time.Duration, max),
		_stopQ:  make(chan error, max),
		wg:      internal.NewWaitGroup(&sync.WaitGroup{}, make(chan bool, max)),
	}
	go _t._loop()
	return _t
}

type Task struct {
	maxDur time.Duration

	curDur  time.Duration
	_curDur chan time.Duration

	max int

	q chan *internal.TaskFnDef

	_stopQ chan error
	_stop  error

	wg *internal.WaitGroup
}

func (t *Task) Len() int {
	return t.wg.Len()
}

func (t *Task) Wait() {
	t.wg.Wait()
}

func (t *Task) done() {
	t.wg.Done()
}

func (t *Task) Do(f internal.TaskFn, args ...interface{}) error {
	for {
		if t._stop != nil {
			return t._stop
		}

		if t.Len() < t.max && t.curDur < t.maxDur {
			t.wg.Add()
			t.q <- f(args...)
			return nil
		}

		if t.Len() < runtime.NumCPU()*2 {
			t.curDur = 0
		}

		if Cfg.Debug {
			log.Printf("q_l:%d cur_dur:%s max_q:%d max_dur:%s", len(t.q), t.curDur.String(), t.max, t.maxDur.String())
		}

		time.Sleep(time.Millisecond)
	}
}

func (t *Task) _loop() {
	for {
		select {
		case _fn := <-t.q:
			go func() {
				t._curDur <- errors.FnCost(func() {
					errors.ErrHandle(errors.Try(_fn.Fn, _fn.Args...), func(err *errors.Err) {
						errors.ErrHandle(errors.Try(_fn.Efn, err), func(err *errors.Err) {
							t._stopQ <- err
						})
					})
				})

				t.done()
			}()
		case t.curDur = <-t._curDur:
		case t._stop = <-t._stopQ:
		}
	}
}
