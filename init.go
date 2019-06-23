package gotask

import (
	"github.com/pubgo/errors"
	"reflect"
)

var Cfg = struct {
	Debug bool
}{
	Debug: true,
}

func errorLog(err error) {
	if Cfg.Debug {
		errors.ErrHandle(err, func(err *errors.Err) {
			err.P()
		})
	}
}

func assertFn(fn interface{}) {
	errors.T(errors.IsZero(fn), "the func is nil")

	_v := reflect.TypeOf(fn)
	errors.T(_v.Kind() != reflect.Func, "func type error(%s)", _v.String())
}
