package gotask

import (
	"github.com/pubgo/errors"
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
