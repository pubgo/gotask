package gotask_test

import (
	"fmt"
	"github.com/pubgo/errors"
	"github.com/pubgo/gotask"
	"testing"
	"time"
)

func TestTasks(t *testing.T) {
	defer errors.Debug()

	_fn1 := gotask.TaskOf(func(i int) {
		//fmt.Println(i)
		errors.T(i == 10999, "90999 error")
	}, func(err error) {
		errors.Wrap(err, "wrap")
	})

	var task = gotask.NewTask(2000, time.Second+time.Millisecond*10)
	for i := 0; i < 100000; i++ {
		errors.Panic(task.Do(_fn1, i))
	}
	task.Wait()
}

func TestErrLog(t *testing.T) {
	defer errors.Debug()

	_fn := gotask.TaskOf(func(i int) {
		//fmt.Println(i)
		errors.T(i == 90999, "90999 error")
	})

	var task = gotask.NewTask(500, time.Second+time.Millisecond*10)
	for i := 0; i < 100000; i++ {
		errors.Panic(task.Do(_fn, i))
	}

	task.Wait()
}

func parserArticleWithReadability(i int) {
	defer errors.Handle(func() {})

	errChan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 4)
		errChan <- true
	}()

	for {
		select {
		case <-time.After(3 * time.Second):
			errors.Wrap("readbility timeout", "等待 %d", i)
		case <-errChan:
			return
		}
	}

}

func TestW(t *testing.T) {
	defer errors.Debug()

	var _fn = gotask.TaskOf(func(i int) {
		parserArticleWithReadability(i)
		fmt.Println("ok", i)
	}, func(err error) {
		errors.ErrHandle(err, func(err *errors.Err) {
			fmt.Println("tag: ", err.Tag())
			errors.Wrap(err, "testW")
		})
	})

	var sss = gotask.NewTask(10000, time.Second*2)
	for i := 0; i < 1000000; i++ {
		errors.Panic(sss.Do(_fn, i))
	}
	sss.Wait()
}
