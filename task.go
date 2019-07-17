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

func NewTask(max int, maxDur time.Duration) *Task {
	_t := &Task{
		max:     max,
		maxDur:  maxDur,
		q:       make(chan *_TaskFn, max),
		_curDur: make(chan time.Duration, max),
		_stopQ:  make(chan error, max),
		wg:      internal.NewWaitGroup(&sync.WaitGroup{}, make(chan bool, max)),
	}
	go _t._loop()
	return _t
}

type Task struct {
	max    int
	maxDur time.Duration

	curDur  time.Duration
	_curDur chan time.Duration

	q chan *_TaskFn

	_stopQ    chan error
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

var _TaskFnPool = &sync.Pool{
	New: func() interface{} {
		return &_TaskFn{
			Fn: reflect.Value{},
		}
	},
}

func getTaskFn() *_TaskFn {
	return _TaskFnPool.Get().(*_TaskFn)
}

type _TaskFn struct {
	Fn   reflect.Value
	Args []reflect.Value
}

func (t *_TaskFn) reset() {
	t.Args = t.Args[:0]
	t.Fn = reflect.Value{}
	_TaskFnPool.Put(t)
}

func (t *Task) Do(fName string, args ...interface{}) {
	f, ok := _tasks[fName]
	errors.T(!ok, "the task %s is not existed", fName)

	var _args = make([]reflect.Value, len(args))
	for i, k := range args {
		_v := reflect.ValueOf(k)
		if k != nil && !errors.IsZero(_v) {
			_args[i] = _v
			continue
		}

		if f.IsVariadic {
			_args[i] = f.VariadicType
			continue
		}

		_args[i] = reflect.New(f.Fn.Type().In(i)).Elem()
	}

	tsk := getTaskFn()
	tsk.Fn = f.Fn
	tsk.Args = _args

	for {
		if t.Len() < t.max && t.curDur < t.maxDur {
			t.wg.Add()
			t.q <- tsk
			return
		}

		if t.Len() < runtime.NumCPU()*2 {
			t.curDur = 0
		}

		if _l := log.Info(); _l.Enabled() {
			_l.Int("q_l", len(t.q)).
				Str("cur_dur", t.curDur.String()).
				Int("max_q", t.max).
				Str("max_dur", t.maxDur.String()).
				Msg("task info")
		}
		time.Sleep(time.Microsecond)
	}
}

func (t *Task) _loop() {
	for {
		select {
		case _fn := <-t.q:
			t.taskCount++

			go func() {
				_t := time.Now()
				errors.ErrHandle(errors.Try(func() {
					_fn.Fn.Call(_fn.Args)
				}), func(err *errors.Err) {
					t._stopQ <- err
				})
				t.done()
				_fn.reset()
				t._curDur <- time.Now().Sub(_t)
			}()
		case _curDur := <-t._curDur:
			t.curDur = (t.curDur + _curDur) / 2

		case _err := <-t._stopQ:
			t.errCount++
			if _l := log.Warn(); _l.Enabled() {
				_l.Err(_err).
					Int("q_l", len(t.q)).
					Int("err_count", t.errCount).
					Int("task_count", t.taskCount).
					Str("cur_dur", t.curDur.String()).
					Int("max_q", t.max).
					Str("max_dur", t.maxDur.String()).
					Str("method", "task").
					Msg("")
			}
		}
	}
}
