# cuslog - 轻量级 Go 日志库

一个轻量、高性能、可扩展的 Golang 日志库，支持多级别、多格式、多输出目标的日志记录能力。

## v0.1.0 版本特性

### 核心功能

1. **6个日志级别**
   - Debug / Info / Warn / Error / Fatal / Panic
   - 支持动态级别配置
   - 级别过滤机制

2. **多种输出格式**
   - TXT 文本格式（支持彩色输出）
   - JSON 结构化格式
   - 格式可动态切换

3. **多输出目标**
   - 控制台输出
   - 文件输出
   - 支持同时输出到多个目标

4. **丰富的元信息**
   - 时间戳
   - 日志级别
   - 调用文件路径和行号
   - 自定义结构化字段

5. **并发安全**
   - 基于 `sync.Mutex` 实现互斥锁
   - 支持多 goroutine 并发写入

6. **标准库兼容**
   - 完全兼容 Go 标准库 `log` 包接口
   - 支持平滑迁移

7. **异常处理**
   - Fatal 级别输出后调用 `os.Exit(1)`
   - Panic 级别输出并触发 panic（包含堆栈信息）

## 安装

```bash
go get cuslog
```

## 快速开始

### 基础使用

```go
package main

import (
    cuslog "cuslog"
)

func main() {
    // 使用默认日志器
    cuslog.Debug("调试信息")
    cuslog.Info("普通信息")
    cuslog.Warn("警告信息")
    cuslog.Error("错误信息")
}
```

### 自定义配置

```go
package main

import (
    cuslog "cuslog"
)

func main() {
    // 创建自定义日志器
    config := cuslog.Config{
        Level:        cuslog.LevelDebug,
        Format:       "txt",
        EnableCaller: true,
        EnableColor:  true,
        Outputs: []cuslog.OutputConfig{
            {Type: "console"},
        },
    }

    logger, err := cuslog.New(config)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    logger.Debug("这是调试信息")
    logger.Infof("这是 %s 信息", "格式化的")
}
```

### JSON 格式输出

```go
config := cuslog.Config{
    Level:        cuslog.LevelInfo,
    Format:       "json",
    EnableCaller: true,
    Outputs: []cuslog.OutputConfig{
        {Type: "console"},
    },
}

logger, _ := cuslog.New(config)
logger.Info("JSON 格式的日志")
```

输出示例：
```json
{"time":"2023-10-01T12:00:00Z","level":"INFO","msg":"JSON 格式的日志","file":"main.go","line":15}
```

### 文件输出

```go
config := cuslog.Config{
    Level:  cuslog.LevelDebug,
    Format: "txt",
    Outputs: []cuslog.OutputConfig{
        {Type: "console"},              // 同时输出到控制台
        {Type: "file", FilePath: "app.log"},  // 和文件
    },
    EnableCaller: true,
    EnableColor:  true,
}

logger, _ := cuslog.New(config)
defer logger.Close()

logger.Info("这条日志会同时输出到控制台和文件")
```

### 动态修改配置

```go
logger, _ := cuslog.New(cuslog.DefaultConfig())

// 动态修改日志级别
logger.SetLevel(cuslog.LevelDebug)

// 动态修改输出格式
logger.SetFormat("json")
```

### 标准库兼容

```go
// 方式1：使用全局函数
cuslog.Print("使用标准库 Print 方法")
cuslog.Printf("使用标准库 %s 方法", "Printf")
cuslog.Println("使用标准库 Println 方法")

// 方式2：创建兼容的 logger
config := cuslog.DefaultConfig()
logger, _ := cuslog.New(config)
stdLogger := cuslog.NewStdLogger(logger)

stdLogger.Print("通过 stdLogger 打印")
stdLogger.Printf("格式化: %s", "输出")
```

## 配置选项

### Config 结构体

| 字段 | 类型 | 说明 |
|------|------|------|
| Level | Level | 日志级别（Debug/Info/Warn/Error/Fatal/Panic） |
| Format | string | 输出格式（"json" 或 "txt"） |
| Outputs | []OutputConfig | 输出目标配置列表 |
| EnableCaller | bool | 是否显示调用文件路径和行号 |
| EnableColor | bool | 是否启用彩色输出（仅 txt 格式） |

### OutputConfig 结构体

| 字段 | 类型 | 说明 |
|------|------|------|
| Type | string | 输出类型（"console" 或 "file"） |
| FilePath | string | 文件路径（file 类型时必需） |
| Addr | string | 网络地址（v0.4.0+ 支持） |

## 日志级别

```go
const (
    LevelDebug Level = iota  // 0
    LevelInfo                // 1
    LevelWarn                // 2
    LevelError               // 3
    LevelFatal               // 4
    LevelPanic               // 5
)
```

- 只有大于等于设定级别的日志才会被输出
- 例如：设置级别为 `LevelInfo`，则 `Debug` 日志不会输出

## 输出示例

### TXT 格式（带颜色）

```
2023-10-01 12:00:00 [DEBUG] main.go:10: 这是调试信息
2023-10-01 12:00:01 [INFO] main.go:11: 这是普通信息
2023-10-01 12:00:02 [WARN] main.go:12: 这是警告信息
2023-10-01 12:00:03 [ERROR] main.go:13: 这是错误信息
```

### JSON 格式

```json
{"time":"2023-10-01T12:00:00Z","level":"DEBUG","msg":"这是调试信息","file":"main.go","line":10}
{"time":"2023-10-01T12:00:01Z","level":"INFO","msg":"这是普通信息","file":"main.go","line":11}
{"time":"2023-10-01T12:00:02Z","level":"WARN","msg":"这是警告信息","file":"main.go","line":12}
{"time":"2023-10-01T12:00:03Z","level":"ERROR","msg":"这是错误信息","file":"main.go","line":13}
```

## 测试

运行单元测试：

```bash
go test -v
```

运行并发测试：

```bash
go test -v -run TestLoggerConcurrent
```

## 性能特点

- **并发安全**：使用互斥锁保证多 goroutine 写入安全
- **零依赖**：仅使用 Go 标准库
- **轻量级**：核心代码简洁，易于理解和维护
- **高效格式化**：优化的格式化逻辑，减少内存分配

## 版本规划

- **v0.1.0** (当前)：基础可用版，核心日志输出能力
- **v0.2.0**：功能增强版，文件输出、日志轮转
- **v0.3.0**：性能优化版，异步写入、无锁机制
- **v0.4.0**：生产级完善版，Hook、采样、网络输出

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
