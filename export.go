package gotask

// handle
var NewAsyncTask = _NewAsyncTask

type AsyncTask = _AsyncTask
type Stat = _Stat
type TaskDef = _TaskDef
type TaskFn func(...interface{}) func(...interface{}) (err error)

var InitDebugLog = _InitDebugLog
