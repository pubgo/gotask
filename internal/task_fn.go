package internal

type TaskFnDef struct {
	Fn   interface{}
	Args []interface{}
	Efn  func(err error)
}

func NewTaskFn(fn interface{}, args []interface{}, efn func(err error)) *TaskFnDef {
	return &TaskFnDef{
		Fn:   fn,
		Args: args,
		Efn:  efn,
	}
}

type TaskFn func(args ...interface{}) *TaskFnDef
