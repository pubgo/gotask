package gotask

type _TaskDef interface {
	Size() int
	CurSize() int
	Wait()
	Stat() Stat
	Do(name TaskFn, args ...interface{})
	Stop()
	_taskHandle(fn func(...interface{}) (err error))
	_loop()
}
