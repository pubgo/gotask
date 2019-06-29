package gotask_test

import (
	"fmt"
	"github.com/pubgo/errors"
	"github.com/pubgo/gotask"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestTasks(t *testing.T) {
	defer errors.Debug()

	_fn1 := gotask.TaskOf(func(i int) {
		errors.ErrHandle(errors.Try(func() {})(func() {
			//fmt.Println(i)
			errors.T(i == 29, "90999 error")
		}), func(err *errors.Err) {
			errors.Wrap(err, "wrap")
		})
	})

	var task = gotask.NewTask(10, time.Second+time.Millisecond*10)
	for i := 0; i < 100; i++ {
		task.Do(_fn1, i)
	}
	task.Wait()
	errors.P(task.Stat())
	fmt.Println(task.Err())
}

func TestErrLog(t *testing.T) {
	defer errors.Debug()

	_fn := gotask.TaskOf(func(i int) {
		//fmt.Println(i)
		errors.T(i == 90999, "90999 error")
	})

	var task = gotask.NewTask(500, time.Second+time.Millisecond*10)
	for i := 0; i < 100000; i++ {
		task.Do(_fn, i)
	}

	task.Wait()
	errors.P(task.Stat())
	fmt.Println(task.Err())
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
		errors.ErrHandle(errors.Try(func() {})(func() {
			parserArticleWithReadability(i)
			fmt.Println("ok", i)
		}), func(err *errors.Err) {
			errors.Wrap(err, "testW")
		})
	})

	var task = gotask.NewTask(10000, time.Second*2)
	for i := 0; i < 1000000; i++ {
		task.Do(_fn, i)
	}
	task.Wait()
	errors.P(task.Stat())
	fmt.Println(task.Err())
}

func isEOF(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}

var _fn = gotask.TaskOf(func(c *http.Client, i int) {
	errors.Retry(3, func() {
		fmt.Println("try: ", i)
		req, err := http.NewRequest(http.MethodGet, "http://baidu.com", nil)
		errors.Panic(err)
		req.Close = true

		resp, err := c.Do(req)
		errors.Panic(err)
		errors.T(resp.StatusCode != http.StatusOK, "状态不正确%d", resp.StatusCode)
	})

	//dt, err := ioutil.ReadAll(resp.Body)
	//errors.Panic(err)
	//fmt.Println(string(dt))

})

func TestUrl(t *testing.T) {
	gotask.Cfg.Debug = false

	client := &http.Client{Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}}
	client.Timeout = 5 * time.Second

	var task = gotask.NewTask(50, time.Second*2)
	for i := 0; i < 300; i++ {
		fmt.Println(i)
		task.Do(_fn, client, i)
	}
	task.Wait()
	errors.P(task.Stat())
	fmt.Println(task.Err())

}
