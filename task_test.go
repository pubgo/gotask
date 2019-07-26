package gotask_test

import (
	"fmt"
	"github.com/pubgo/errors"
	"github.com/pubgo/gotask"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"testing"
	"time"
)

func init() {
	gotask.InitDebugLog()
}

func TestTasks(t *testing.T) {
	defer errors.Assert()

	//zerolog.SetGlobalLevel(zerolog.WarnLevel)
	//fn := errors.Try(func(i int) {
	//	defer errors.Resp(func(err *errors.Err) {
	//
	//	})
	//
	//	errors.T(i == 29, "90999 error")
	//})

	gotask.TaskRegister("fn", func(i int) {
		//defer errors.Resp(func(err *errors.Err) {
		//})

		errors.T(i == 29, "90999 error")
	})

	var task = gotask.NewTask(10, time.Second+time.Millisecond*10)
	defer task.Stop()

	for i := 0; i < 100; i++ {
		task.Do("fn", i)
	}

	task.Wait()
	errors.P(task.Stat())
}

func TestErrLog(t *testing.T) {
	defer errors.Assert()

	//zerolog.SetGlobalLevel(zerolog.WarnLevel)

	gotask.TaskRegister("fn", func(i int) {
		//fmt.Println(i)
		errors.T(i == 90999, "90999 error")
	})

	var task = gotask.NewTask(500, time.Second+time.Millisecond*10)
	for i := 0; i < 100000; i++ {
		go task.Do("fn", i)
	}
	task.Wait()
	errors.P(task.Stat())
}

func parserArticleWithReadability(i int) {
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
	defer errors.Assert()

	gotask.TaskRegister("fn", func(i int) {
		errors.ErrHandle(errors.Try(func() {})(func() {
			parserArticleWithReadability(i)
			fmt.Println("ok", i)
		}), func(err *errors.Err) {
			errors.Wrap(err, "testW")
		})
	})

	var task = gotask.NewTask(10000, time.Second*2)
	for i := 0; i < 1000000; i++ {
		task.Do("fn", i)
	}
	task.Wait()
	errors.P(task.Stat())
}

func isEOF(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}

func TestUrl(t *testing.T) {
	defer errors.Assert()

	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	client := &http.Client{Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}}
	client.Timeout = 5 * time.Second

	gotask.TaskRegister("fn", func(c *http.Client, i int) {
		errors.Panic(errors.Retry(3, func() {
			fmt.Println("try: ", i)
			req, err := http.NewRequest(http.MethodGet, "https://www.yuanben.io", nil)
			errors.Panic(err)

			resp, err := c.Do(req)
			errors.Panic(err)
			errors.T(resp.StatusCode != http.StatusOK, "状态不正确%d", resp.StatusCode)
		}))
	})

	var task = gotask.NewTask(200, time.Second*2)
	for i := 0; i < 3000; i++ {
		task.Do("fn", client, i)
	}
	task.Wait()
	fmt.Println(task.Stat())
}

func TestTaskRegistry(t *testing.T) {
	errors.TestRun(gotask.TaskRegister, func(desc func(string) *errors.Test) {
		desc("int params").In("fn", func(i int) {
			errors.ErrHandle(errors.Try(func() {
				errors.T(i == 29, "90999 error")
			}), func(err *errors.Err) {
				errors.Wrap(err, "wrap")
			})
		}).IsNil()
	})

	errors.TestRun(gotask.NewTask, func(desc func(string) *errors.Test) {
		desc("ok").In(10, time.Second+time.Millisecond*10).IsNil(func(task *gotask.Task) {
			for i := 0; i < 100; i++ {
				task.Do("fn", i)
			}
			task.Wait()
			errors.P(task.Stat())
		})
	})

}
