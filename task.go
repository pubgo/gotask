package gotask

import (
	"github.com/pubgo/errors"
	"runtime"
	"sync"
	"time"
)

func _NewTask(max int, maxDur time.Duration) *_Task {
	_t := &_Task{
		max:     max,
		maxDur:  maxDur,
		taskL:   make(chan bool, max),
		taskQ:   make(chan func(...interface{}) (err error), max),
		_curDur: make(chan time.Duration, max),
		mux:     &sync.Mutex{},
	}
	go _t._loop()
	return _t
}

type _Task struct {
	max    int
	maxDur time.Duration

	curDur  time.Duration
	_curDur chan time.Duration

	taskL chan bool
	taskQ chan func(...interface{}) (err error)

	mux *sync.Mutex
}

func (t *_Task) Size() int {
	return len(t.taskL)
}

// ATLen current active task size
func (t *_Task) CurSize() int {
	t.mux.Lock()
	defer t.mux.Unlock()

	return len(t.taskL) - len(t.taskQ)
}

func (t *_Task) Wait() {
	for len(t.taskL) > 0 {
		time.Sleep(time.Second)
	}
}

func (t *_Task) Stat() Stat {
	return Stat{
		QL:     t.Size(),
		CurDur: t.curDur.Seconds(),
		MaxQ:   t.max,
		MaxDur: t.maxDur.Seconds(),
	}
}

func (t *_Task) Do(name string, args ...interface{}) {
	t.mux.Lock()
	defer t.mux.Unlock()

	errors.T(!TaskMatch(name), "the task %s is not existed", name)
	for {
		if len(t.taskL) < t.max && t.curDur < t.maxDur {
			t.taskQ <- TaskGet(name)(args...)
			t.taskL <- true
			return
		}

		if len(t.taskL) < runtime.NumCPU()*2 {
			t.curDur = 0
		}

		if _l := logger.Info(); _l.Enabled() {
			_l.Int("q_l", len(t.taskQ)).
				Str("cur_dur", t.curDur.String()).
				Int("max_q", t.max).
				Str("max_dur", t.maxDur.String()).
				Msg("task info")
		}
		time.Sleep(time.Microsecond)
	}
}

func (t *_Task) _taskHandle(fn func(...interface{}) (err error)) {
	_t := time.Now()
	errors.ErrHandle(fn(), func(err *errors.Err) {
		if _l := logger.Warn(); _l.Enabled() {
			_l.Err(err).
				Int("taskQ_len", len(t.taskQ)).
				Int("max_taskQ_len", t.max).
				Str("cur_dur", t.curDur.String()).
				Str("max_dur", t.maxDur.String()).
				Str("method", "task").
				Msg("")
		}
	})
	<-t.taskL
	t._curDur <- time.Now().Sub(_t)
}

func (t *_Task) _loop() {
	for {
		select {
		case _fn := <-t.taskQ:
			go t._taskHandle(_fn)
		case _curDur := <-t._curDur:
			t.curDur = (t.curDur + _curDur) / 2
		}
	}
}
