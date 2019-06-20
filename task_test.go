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
	_fn1 := gotask.TaskOf(func(i int) {
		//fmt.Println(i)
		assert.T(i == 90999, "90999 error")
	}, func(err error) {
		assert.Throw(err)
	})

	var task = gotask.NewTask(500, time.Second+time.Millisecond*10)

	fmt.Println("time cost: ", assert.FnCost(func() {
		for i := 0; i < 100000; i++ {
			if err := task.Do(_fn1, i); err != nil {
				assert.P(err)
				break
			}
		}
	}))
	task.Wait()
}

func TestErrLog(t *testing.T) {
	_fn := gotask.TaskOf(func(i int) {
		//fmt.Println(i)
		assert.T(i == 90999, "90999 error")
	})

	var task = gotask.NewTask(500, time.Second+time.Millisecond*10)

	fmt.Println("time cost: ", assert.FnCost(func() {
		for i := 0; i < 100000; i++ {
			if err := task.Do(_fn, i); err != nil {
				assert.P(err)
				break
			}
		}
	}))
	task.Wait()
}

func parserArticleWithReadability(i int) error {
	errChan := make(chan bool)
	go func() {
		time.Sleep(time.Second*4)
		errChan <- true
	}()

	for {
		select {
		case <-time.After(3 * time.Second):
			return assert.Wrap(errors.New("readbility timeout"), "等待 %d",i)
		case <-errChan:
			return nil
		}
	}

}

func TestW(t *testing.T) {
	var _fn = gotask.TaskOf(func(i int) {
		assert.ErrWrap(parserArticleWithReadability(i),"yyyy")
		fmt.Println("ok", i)
	}, func(err error) {
		assert.ErrHandle(err, func(err *assert.KErr) {
			err.P()
		})
	})

	var sss = gotask.NewTask(10000, time.Second*2)
	for i := 0; i < 1000000; i++ {
		if err := sss.Do(_fn, i); err != nil {
			panic(err)
		}
		fmt.Println(i)
		fmt.Println(sss.Len())
	}
	sss.Wait()
}