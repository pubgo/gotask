package gotask

import "github.com/pubgo/assert"

var Debug = true

func errorLog(err error) {
	if Debug {
		assert.P(err)
	}
}
