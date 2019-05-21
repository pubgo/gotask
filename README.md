# gotask
task for go 

```go
fmt.Println(_FnCost(func() {
    Debug = false

    _fn := TaskOf(func(i int) {
        //fmt.Println(i)
        _T(i == 90999, "90999 error")
    }, func(err error) {
        _Throw(err)
    })

    var task = NewTask(500, time.Second+time.Millisecond*10)
    for i := 0; i < 100000; i++ {
        if err := task.Do(_fn, i); err != nil {
            fmt.Println(err)
            break
        }
    }

    task.Wait()
}))
```