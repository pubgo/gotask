package gotask

import (
	"github.com/pubgo/errors"
	"github.com/pubgo/gotask/internal"
	"reflect"
)

var _tasks = make(map[string]*internal.TaskFnDef)

func TaskRegistry(name string, fn interface{}) {
	defer errors.Assert()
	if _, ok := _tasks[name]; ok {
		errors.T(ok, "%s has existed", name)
	}

	_fn := reflect.ValueOf(fn)
	errors.T(errors.IsZero(_fn) ||
		_fn.Kind() != reflect.Func ||
		_fn.Type().NumOut() != 0, "fn error")

	var variadicType reflect.Value
	var isVariadic = _fn.Type().IsVariadic()
	if isVariadic {
		variadicType = reflect.New(_fn.Type().In(_fn.Type().NumIn() - 1).Elem()).Elem()
	}

	_tasks[name] = &internal.TaskFnDef{
		Fn:           _fn,
		VariadicType: variadicType,
		IsVariadic:   isVariadic,
	}
}

func GetTasks() map[string]*internal.TaskFnDef {
	return _tasks
}

func GetTask(name string) (tsk *internal.TaskFnDef) {
	if _dt, ok := _tasks[name]; ok {
		return _dt
	}
	return
}

func MatchTask(name string) (ok bool) {
	_, ok = _tasks[name]
	return
}
