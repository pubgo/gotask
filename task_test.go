package gotask_test

import (
	"errors"
	"fmt"
	"github.com/pubgo/assert"
	"github.com/pubgo/gotask"
	"testing"
	"time"
)

func TestTasks(t *testing.T) {
	defer assert.Debug()

	_fn1 := gotask.TaskOf(func(i int) {
		//fmt.Println(i)
		assert.T(i == 10999, "90999 error")
	}, func(err error) {
		_e := assert.Wrap(err, "wrap")
		assert.Throw(_e)
	})

	var task = gotask.NewTask(2000, time.Second+time.Millisecond*10)

	for i := 0; i < 100000; i++ {
		assert.Throw(task.Do(_fn1, i))
	}
	task.Wait()
}

func TestErrLog(t *testing.T) {
	defer assert.Debug()

	_fn := gotask.TaskOf(func(i int) {
		//fmt.Println(i)
		assert.T(i == 90999, "90999 error")
	})

	var task = gotask.NewTask(500, time.Second+time.Millisecond*10)

	for i := 0; i < 100000; i++ {
		assert.Throw(task.Do(_fn, i))
	}

	task.Wait()
}

func parserArticleWithReadability(i int) {
	defer assert.Panic(func(m *assert.M) {
		m.Msg("parserArticleWithReadability")
	})

	errChan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 4)
		errChan <- true
	}()

	for {
		select {
		case <-time.After(3 * time.Second):
			assert.ErrWrap(errors.New("readbility timeout"), "等待 %d", i)
		case <-errChan:
			return
		}
	}

}

func TestW(t *testing.T) {
	defer assert.Debug()

	var _fn = gotask.TaskOf(func(i int) {
		parserArticleWithReadability(i)
		fmt.Println("ok", i)
	}, func(err error) {
		assert.ErrHandle(err, func(err *assert.KErr) {
			fmt.Println("tag: ", err.Tag())
			assert.ErrWrap(err, "testW")
		})
	})

	var sss = gotask.NewTask(10000, time.Second*2)
	for i := 0; i < 1000000; i++ {
		assert.Throw(sss.Do(_fn, i))
	}
	sss.Wait()
}
