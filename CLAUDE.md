# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 常用命令

### 测试
```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestLoggerConcurrent

# 运行带竞态检测的测试
go test -race -v
```

### 构建
```bash
# 构建模块（这是一个库，没有 main 包）
go build ./...

# 构建示例
go build ./examples/...

# 运行示例
go run examples/main.go
```

## 项目概述

**cuslog** 是一个轻量级的 Go 日志库（v0.1.0），采用模块化架构，支持多日志级别、多格式和多输出目标。

## 核心架构

### 处理流程

日志处理遵循以下管道流程：

```
用户调用 (Info/Error/Warn 等)
    ↓
级别检查（低于最小级别的日志会被忽略，不进行格式化）
    ↓
创建 Entry（收集时间戳、调用者信息、消息）
    ↓
格式化（TxtFormatter 或 JsonFormatter）
    ↓
输出写入（ConsoleWriter、FileWriter 或 MultiWriter）
    ↓
特殊处理（Fatal 调用 os.Exit(1)，Panic 触发 panic）
```

### 核心组件

#### Logger 结构（logger.go）
```go
type Logger struct {
    config    Config      // 配置
    formatter Formatter   // 格式化器
    writer    Writer      // 输出目标
    mu        sync.Mutex  // 并发安全
}
```

#### 关键接口

**Formatter 接口** (formatter.go)
```go
type Formatter interface {
    Format(entry *Entry) ([]byte, error)
}
```
- `TxtFormatter`: 人类可读的文本格式，支持可选的 ANSI 颜色
- `JsonFormatter`: 机器可读的 JSON 结构化格式

**Writer 接口** (writer.go)
```go
type Writer interface {
    Write(p []byte) (n int, err error)
    Close() error
}
```
- `ConsoleWriter`: 输出到 stdout（自动检测是否启用颜色）
- `FileWriter`: 输出到文件（支持自动重新打开）
- `MultiWriter`: 同时输出到多个目标

#### 数据结构

**Entry** (entry.go)
```go
type Entry struct {
    Time    time.Time              // 时间戳
    Level   Level                  // 日志级别
    Message string                 // 日志消息
    File    string                 // 调用者文件名
    Line    int                    // 调用者行号
    Data    map[string]interface{} // 自定义结构化字段
}
```

**Level** (level.go)
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

### 设计模式

- **策略模式**: Formatter 和 Writer 接口允许运行时交换实现
- **工厂模式**: `GetFormatter()` 和 `createWriter()` 根据配置创建实例
- **适配器模式**: `StdLogger` 将 cuslog 适配到 Go 标准库 log 接口
- **观察者模式**: MultiWriter 使多个输出目标能接收相同的日志条目

## 关键实现细节

### 并发安全
- 所有写操作都通过 `sync.Mutex` 保护
- 每个 Writer 实现也有自己的互斥锁
- 测试验证了 100 个 goroutine × 100 条消息 = 10,000 次并发写入的安全性

### 调用者信息
- 使用 `runtime.Caller()` 获取文件路径和行号
- Skip 参数正确导航调用栈，获取实际调用位置
- 只提取文件名（而非完整路径）以保持输出简洁

### 颜色输出
- 仅对控制台输出启用颜色（自动检测）
- 使用 ANSI 颜色代码（在终端中有效，文件中无效）
- 颜色映射到日志级别（Debug=灰色, Info=蓝色, Warn=黄色, Error=红色等）

### 动态配置
- `SetLevel(level)`: 运行时更改最小日志级别
- `SetFormat(format)`: 在 txt 和 json 格式之间切换
- 颜色会在格式更改时自动重新计算

### 错误处理
- 格式化错误输出到 stderr（不会导致应用崩溃）
- 写入错误被报告但不会阻止其他写入
- FileWriter 支持自动重新打开文件

## 版本规划

项目遵循清晰的版本策略（详见 `docs/版本开发计划.md`）：

- **v0.1.0** (当前): 基础日志功能，基于互斥锁的并发控制
- **v0.2.0**: 文件轮转、配置文件、context 支持
- **v0.3.0**: 无锁性能优化、异步写入
- **v0.4.0**: 生产级特性（Hook、采样、网络输出）

## 扩展点

添加新功能的方法：

1. **新输出目标**: 实现 `Writer` 接口（如 TCP、syslog）
2. **新格式**: 实现 `Formatter` 接口（如 XML、自定义格式）
3. **自定义字段**: 使用 `Entry.Data` map 添加结构化元数据
4. **Hook**: v0.4.0 计划支持 - 为特定日志级别添加拦截器

## 测试指南

- 始终使用多个 goroutine 测试并发场景
- 使用 `bufferWriter` 测试辅助工具捕获输出
- 验证级别过滤（低于最小级别的日志不应被格式化）
- 同时测试 TXT 和 JSON 格式
- 验证动态配置更改（`SetLevel`、`SetFormat`）

## 标准库兼容性

**stdlib.go** 提供与 Go 标准 `log` 包的兼容性：

```go
// 创建兼容的 logger
stdLogger := NewStdLogger(logger)

// 或使用全局函数
cuslog.Print("消息")
cuslog.Printf("格式化: %s", "消息")
```

这允许存量代码平滑迁移到 cuslog。

## 相关文档

- **[README.md](README.md)**: 项目介绍、快速开始指南
- **[docs/设计文档.md](docs/设计文档.md)**: 详细的技术设计和架构说明
- **[docs/版本开发计划.md](docs/版本开发计划.md)**: 版本迭代规划（v0.1.0 → v0.4.0）
- **[DEVELOPMENT.md](DEVELOPMENT.md)**: 开发相关文档
