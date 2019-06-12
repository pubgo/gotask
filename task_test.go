package gotask_test

import (
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
