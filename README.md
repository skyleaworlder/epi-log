# :memo: epi-log

琢磨出来这个名字的时候，我觉得简直太妙了。但转念一想，一定早有歪果仁想到了，事实也是如此……

## :wrench: Feature

没什么 `Feature`，只实现了日志该有的基本功能。

* 异步刷新，多用户线程 & 单消费进程，先存 `buffer`，统一刷到各个对端；
* 故障处理，通过捕捉信号，尽量避免意外对日志完整性的影响；
* 外部接口，模仿 `logrus` 使用 `export.go` 提供使用默认配置的导出函数。

## :hammer: Usage

### :one: 使用导出函数

由于一开始没设计好，所以使用起来很丑陋（

```go
func main() {
    // call Use to use epilog.
    // Use generates a signal-process goroutine,
    // and begin to Monitor buffer in epilog.
    // Defer End function is followed closely.
    // End processes items still in buffer,
    // when main-goroutine exits or terminates.
    go epilog.Use()
    defer epilog.End()
}
```

和著名日志库 `logrus` 不同，`epilog` 内使用缓冲，并批量写入。所以在 `main` 函数启用 `epilog` 后，需要通过运行在 `main` 函数最后一条语句后的 `End` 函数，来处理缓冲区中剩余未写入的日志项。

### :two: 使用构造函数

```go
func main() {
    // use default constructor to generate a epilog.Logger.
    logger := epilog.New()

    // Change output level.
    logger.ChangeLevel(epilog.WARNING)

    // register self-defined appender for epilog.Logger.
    fa := epilog.NewFileAppender("fa_1", "./1.txt")
    logger.RegisterAppender("fa_1", fa)

    // go Use & defer End
    go logger.Use()
    defer logger.End()

    // for test
    logger.Print("test print")
    logger.Debugln("test debugln")
    logger.Infoln("test infoln")
    logger.Warningln("test warningln")
}
```

* 通过 `New` 创建 `struct`；
* 此后通过 `RegisterAppender` 注册 `Appender`；
* 启用 `epilog` 的 `Monitor`，`Signal-Processor` 和 `End`。
