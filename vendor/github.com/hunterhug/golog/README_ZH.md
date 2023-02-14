# Golang日志库，非常简单的记录一切

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/golog.svg?style=social&label=Forks)](https://github.com/hunterhug/golog/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/golog.svg?style=social&label=Stars)](https://github.com/hunterhug/golog/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/golog.svg)](https://github.com/hunterhug/golog)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/golog.svg)](https://github.com/hunterhug/golog/issues)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

感谢 Uber 的开源项目 ZapLog，它的速度很快，非常快，且是工业级别应用于很多大型的生产环境。我从没见过这么快的日志库，于是我封装了一层，奥利给。

1. 支持日志输出到控制台或文件，支持很多配置，如打印函数调用行，配置日志级别等。
2. 支持日志层级输出到不同文件，且支持文件切割，清除过期日志文件等。
3. 使用接口契约，你可以自行再次封装该库，我留有 `SetCallerSkip` 给你继续封装。

使用起来非常简单，你值得拥有，快点用起来吧，替代你那又慢又难用的日志库。

[English README](README.md)

## 如何使用

非常简单，您只需要常规执行:

```
go get -v -u github.com/hunterhug/golog
```

## 例子

默认什么都不配置的话，日志是打印到终端控制台的，默认日志级别是 InfoLevel，会打印出函数调用者的路径，默认长路径，打印出来的日志是文本格式，高亮模式。

您可以修改某些配置，来进行定制化，比如修改输出调用者函数路径为短路径，修改打印出 JSON 格式等等。

### 例子1：默认用法

```go
package main

import . "github.com/hunterhug/golog"

func main() {
	// use default log
	Info("now is Info", 2, " good")
	Debug("now is Debug", 2, " good")
	Warn("now is Warn", 2, " good")
	Error("now is Error", 2, " good")
	Infof("now is Infof: %d,%s", 2, "good")
	Debugf("now is Debugf: %d,%s", 2, "good")
	Warnf("now is Warnf: %d,%s", 2, "good")
	Errorf("now is Errorf: %d,%s", 2, "good")
	Sync()

	// config log
	SetLevel(DebugLevel).SetCallerShort(true).SetOutputJson(true).InitLogger()

	Info("now is Info", 2, " good")
	Debug("now is Debug", 2, " good")
	Warn("now is Warn", 2, " good")
	Error("now is Error", 2, " good")
	Infof("now is Infof: %d,%s", 2, "good")
	Debugf("now is Debugf: %d,%s", 2, "good")
	Warnf("now is Warnf: %d,%s", 2, "good")
	Errorf("now is Errorf: %d,%s", 2, "good")
	Sync()

}
```

输出是：

```
2021-08-27T11:16:10.455+0800    INFO    /Users/pika/Documents/code/github/golog/demo/demo1/main.go:7    main.main       now is Info2 good
2021-08-27T11:16:10.455+0800    WARN    /Users/pika/Documents/code/github/golog/demo/demo1/main.go:9    main.main       now is Warn2 good
2021-08-27T11:16:10.455+0800    ERROR   /Users/pika/Documents/code/github/golog/demo/demo1/main.go:10   main.main       now is Error2 good
2021-08-27T11:16:10.455+0800    INFO    /Users/pika/Documents/code/github/golog/demo/demo1/main.go:11   main.main       now is Infof: 2,good
2021-08-27T11:16:10.455+0800    WARN    /Users/pika/Documents/code/github/golog/demo/demo1/main.go:13   main.main       now is Warnf: 2,good
2021-08-27T11:16:10.455+0800    ERROR   /Users/pika/Documents/code/github/golog/demo/demo1/main.go:14   main.main       now is Errorf: 2,good
{"l":"info","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:19","func":"main.main","msg":"now is Info2 good"}
{"l":"debug","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:20","func":"main.main","msg":"now is Debug2 good"}
{"l":"warn","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:21","func":"main.main","msg":"now is Warn2 good"}
{"l":"error","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:22","func":"main.main","msg":"now is Error2 good"}
{"l":"info","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:23","func":"main.main","msg":"now is Infof: 2,good"}
{"l":"debug","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:24","func":"main.main","msg":"now is Debugf: 2,good"}
{"l":"warn","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:25","func":"main.main","msg":"now is Warnf: 2,good"}
{"l":"error","t":"2021-08-27T11:16:10.455+0800","caller":"demo1/main.go:26","func":"main.main","msg":"now is Errorf: 2,good"}
```

### 例子2：输出到文件

最重要的是，您可以输出日志到文件中，文件还支持按日志级别输出到不同的文件中，且支持文件切割，文件自动清理等。

```go
package main

import (
	"context"
	"fmt"
	. "github.com/hunterhug/golog"
	"time"
)

func main() {
	SetName("log_demo")
	SetLevel(InfoLevel)

	SetCallerShort(true).SetOutputJson(true)

	AddFieldFunc(func(ctx context.Context, m map[string]interface{}) {
		m["diy_filed"] = ctx.Value("diy")
	})

	SetOutputFile("./log", "demo").SetFileRotate(30*24*time.Hour, 24*time.Hour)
	SetIsOutputStdout(true)
	InitLogger()

	Info("now is Info", 2, " good")
	Debug("now is Debug", 2, " good")
	Warn("now is Warn", 2, " good")
	Error("now is Error", 2, " good")
	Infof("now is Infof: %d,%s", 2, "good")
	Debugf("now is Debugf: %d,%s", 2, "good")
	Warnf("now is Warnf: %d,%s", 2, "good")
	Errorf("now is Errorf: %d,%s", 2, "good")

	ctx := context.WithValue(context.Background(), "diy", []interface{}{"ahhahahahahh"})
	InfoContext(ctx, "InfoContext")
	InfoContext(ctx, "InfoContext, %s:InfoContext, %d", "ss", 333)
	InfoWithFields(map[string]interface{}{"k1": "sss"}, "InfoWithFields:%s，%d", "sss", 33333)
	InfoWithFields(map[string]interface{}{"k1": "sss"}, "InfoWithFields")

	err := Sync()
	if err != nil {
		fmt.Println(err.Error())
	}
}
```

您可以看到目录下有个文件夹 `·/log`，里面就是日志了，可以打开看看。 你会很好奇，控制台也打印出了日志，因为我们同时配置了日志输出到终端，使用 `SetIsOutputStdout` 函数即可。

## 用法一览

学习这个库非常简单，看看下面的接口方法。

```go
type LoggerInterface interface {
	SetOutputFile(logPath, fileName string)
	SetFileRotate(fileMaxAge, fileRotation time.Duration)
	SetLevel(level Level)
	SetCallerShort(short bool)
	SetName(name string)
	SetIsOutputStdout(isOutputStdout bool)
	SetCallerSkip(skip int)
	SetOutputJson(json bool)

	GetOutputFile() (logPath, fileName string)
	GetFileRotate() (fileMaxAge, fileRotation time.Duration)
	GetLevel() (level Level)
	GetCallerShort() (short bool)
	GetName() (name string)
	GetIsOutputStdout() (isOutputStdout bool)
	GetCallerSkip() (skip int)
	GetOutputJson() bool

	InitLogger()
	Sync() error

	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})

	Panic(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})

	PanicWithFields(fields map[string]interface{}, template string, args ...interface{})
	FatalWithFields(fields map[string]interface{}, template string, args ...interface{})
	ErrorWithFields(fields map[string]interface{}, template string, args ...interface{})
	WarnWithFields(fields map[string]interface{}, template string, args ...interface{})
	InfoWithFields(fields map[string]interface{}, template string, args ...interface{})
	DebugWithFields(fields map[string]interface{}, template string, args ...interface{})

	PanicContext(ctx context.Context, template string, args ...interface{})
	FatalContext(ctx context.Context, template string, args ...interface{})
	ErrorContext(ctx context.Context, template string, args ...interface{})
	WarnContext(ctx context.Context, template string, args ...interface{})
	InfoContext(ctx context.Context, template string, args ...interface{})
	DebugContext(ctx context.Context, template string, args ...interface{})

	PanicContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	FatalContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	ErrorContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	WarnContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	InfoContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	DebugContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})

	AddFieldFunc(func(context.Context, map[string]interface{}))
}
```

有些看不懂没关系，因为你可能用不上。

# License

```
Copyright [2021-2021] [github.com/hunterhug]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```