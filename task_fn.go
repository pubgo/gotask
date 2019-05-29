package gotask

type _taskFn struct {
	fn   interface{}
	args []interface{}
	efn  []func(err error)
}

type TaskFn func(args ...interface{}) *_taskFn

func TaskOf(fn interface{}, efn ...func(err error)) TaskFn {
	_AssertFn(fn)
	_T(len(efn) != 0 && efn[0] == nil, "efn is nil")

	return func(args ...interface{}) *_taskFn {
		return &_taskFn{
			fn:   fn,
			args: args,
			efn:  efn,
		}
	}
}
