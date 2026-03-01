# 重构：抽离默认配置到 runner/configs.go

## TL;DR
> 把 main.go 中的硬编码 defaultConfig 移到 runner/configs.go，简化 main.go

## Context

### 当前问题
- main.go 第 46-92 行有约 40 行硬编码的 defaultConfig
- 与 runner/configs.go 中的 Config 定义重复

### 解决方案
在 runner/configs.go 添加 DefaultConfig() 函数

---

## Work Objectives

### 需要修改的文件
1. `cmd/wutils/runner/configs.go` - 添加 DefaultConfig() 函数
2. `cmd/wutils/main.go` - 简化为调用 DefaultConfig()

---

## TODOs

- [x] 1. 在 runner/configs.go 添加 DefaultConfig() 函数

  **What to do**:
  - 添加 `func DefaultConfig() Config` 函数
  - 返回 main.go 第 46-92 行的默认配置

  **References**:
  - `cmd/wutils/main.go:46-92` - 当前的默认配置实现

  **Acceptance Criteria**:
  - [ ] 函数签名: `func DefaultConfig() Config`
  - [ ] 返回值与 main.go 中的 defaultConfig 一致
  - [ ] go build 通过

- [x] 2. 修改 main.go 使用 DefaultConfig()

  **What to do**:
  - 删除 main.go 第 46-92 行的 defaultConfig 变量
  - 改为调用 `runner.DefaultConfig()`

  **References**:
  - `cmd/wutils/main.go:46` - 当前 defaultConfig 定义
  - `cmd/wutils/runner/configs.go` - 新增的 DefaultConfig 函数

  **Acceptance Criteria**:
  - [ ] go build 通过
  - [ ] 首次运行行为不变（自动创建默认配置）

---

## Commit Strategy
- Message: `refactor(config): extract default config to runner package`
- Files: runner/configs.go, main.go
