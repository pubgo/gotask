package internal

type TaskFnDef struct {
	Fn   interface{}
	Args []interface{}
}

func NewTaskFn(fn interface{}, args []interface{}) TaskFnDef {
	return TaskFnDef{
		Fn:   fn,
		Args: args,
	}
}

type TaskFn func(args ...interface{}) TaskFnDef

type Stat struct {
	QL        int     `json:"q_l,omitempty"`
	CurDur    float64 `json:"cur_dur,omitempty"`
	MaxQ      int     `json:"max_q,omitempty"`
	MaxDur    float64 `json:"max_dur,omitempty"`
	ErrCount  int     `json:"err_count,omitempty"`
	TaskCount int     `json:"task_count,omitempty"`
}
