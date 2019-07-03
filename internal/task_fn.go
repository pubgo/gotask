package internal

import "reflect"

type TaskFnDef struct {
	Fn           reflect.Value
	Args         []reflect.Value
	VariadicType reflect.Value
	IsVariadic   bool
}

type Stat struct {
	QL        int     `json:"q_l,omitempty"`
	CurDur    float64 `json:"cur_dur,omitempty"`
	MaxQ      int     `json:"max_q,omitempty"`
	MaxDur    float64 `json:"max_dur,omitempty"`
	ErrCount  int     `json:"err_count,omitempty"`
	TaskCount int     `json:"task_count,omitempty"`
}
