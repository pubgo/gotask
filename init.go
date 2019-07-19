package gotask

import (
	"github.com/pubgo/errors"
	"reflect"
)

var _tasks = make(map[string]func(...interface{}) func(...interface{}) (err error))

func _TaskRegister(name string, fn interface{}) {
	defer errors.Assert()

	if _, ok := _tasks[name]; ok {
		errors.T(ok, "%s existed", name)
	}

	_fn := reflect.ValueOf(fn)
	errors.T(errors.IsZero(_fn) || _fn.Kind() != reflect.Func, "the func is nil(%#v) or type error(%s)", fn, _fn.Kind().String())
	_tasks[name] = errors.Try(fn)
}

func _TaskEach(fn func(name string, fn func(...interface{}) func(...interface{}) (err error))) {
	for k, v := range _tasks {
		fn(k, v)
	}
}

func _TaskGet(name string) func(...interface{}) func(...interface{}) (err error) {
	if _dt, ok := _tasks[name]; ok {
		return _dt
	}
	return nil
}

func _TaskMatch(name string) bool {
	_, ok := _tasks[name]
	return ok
}
