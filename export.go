package gotask

// handle
var NewTask = _NewTask

type Task = _Task
type Stat = _Stat
type TaskFn func(...interface{}) func(...interface{}) (err error)

var InitDebugLog = _InitDebugLog
