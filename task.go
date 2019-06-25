package gotask

import (
	"github.com/pubgo/errors"
	"github.com/pubgo/gotask/internal"
	"github.com/rs/zerolog/log"
	"reflect"
	"runtime"
	"sync"
	"time"
)

func TaskOf(fn interface{}) internal.TaskFn {
	defer errors.Handle(func() {})

	errors.T(errors.IsZero(fn) ||
		reflect.TypeOf(fn).Kind() != reflect.Func ||
		reflect.TypeOf(fn).NumOut() != 0, "fn error")

	return func(args ...interface{}) *internal.TaskFnDef {
		defer errors.Handle(func() {})
		return internal.NewTaskFn(fn, args)
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

	_stopQ    chan error
	_stop     error
	errCount  int
	taskCount int

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

func (t *Task) Stat() internal.Stat {
	return internal.Stat{
		QL:        len(t.q),
		CurDur:    t.curDur.Seconds(),
		MaxQ:      t.max,
		MaxDur:    t.maxDur.Seconds(),
		ErrCount:  t.errCount,
		TaskCount: t.taskCount,
	}
}

func (t *Task) Err() error {
	return t._stop
}

func (t *Task) Do(f internal.TaskFn, args ...interface{}) {
	defer errors.Handle(func() {})

	for {

		if t.Len() < t.max && t.curDur < t.maxDur {
			t.wg.Add()
			t.q <- f(args...)
			return
		}

		if t.Len() < runtime.NumCPU()*2 {
			t.curDur = 0
		}

		log.Info().
			Int("q_l", len(t.q)).
			Str("cur_dur", t.curDur.String()).
			Int("max_q", t.max).
			Str("max_dur", t.maxDur.String()).
			Msg("task info")

		time.Sleep(time.Millisecond)
	}
}

func (t *Task) _loop() {
	defer errors.Handle(func() {})

	for {
		select {
		case _fn := <-t.q:
			t.taskCount++

			go func() {
				t._curDur <- errors.FnCost(func() {
					errors.ErrHandle(errors.Try(_fn.Fn, _fn.Args...), func(err *errors.Err) {
						t._stopQ <- err
					})
					t.done()
				})
			}()
		case _curDur := <-t._curDur:
			t.curDur = (t.curDur + _curDur) / 2
		case t._stop = <-t._stopQ:
			t.errCount++
		}
	}
}
