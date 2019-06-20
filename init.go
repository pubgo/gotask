package gotask

import "github.com/pubgo/assert"

var Cfg = struct {
	Debug bool
}{
	Debug: true,
}

func errorLog(err error) {
	if Cfg.Debug {
		assert.P(err)
	}
}
