# Dokan IO 缓冲虚拟文件系统

## TL;DR
> 基于 Dokan (Windows FUSE) 实现虚拟文件系统，对指定盘符进行透明读写缓冲
> 
> **核心功能**：
> - 挂载为盘符 (如 X:)，拦截所有文件操作
> - 写缓冲：内存缓冲 + 定期批量刷盘 (write-behind)
> - 读缓冲：热点数据缓存 (read-ahead)
> 
> **预计工作量**：Large (6-8 周)
> **并行执行**：YES - 多模块并行

---

## Context

### 原始需求
- HDD 在频繁小文件随机读写时 IOPS 达到瓶颈
- 系统未提供好的 IO 合并和缓冲
- 目标：通过透明拦截实现 IO 缓冲，减少 IOPS

### 技术选型
- **方案**：Dokan (Windows FUSE) 用户态虚拟文件系统
- **Go 库**：github.com/binzume/dkango (Pure Go, 支持 io/fs.FS)
- **驱动**：需要安装 Dokany (>= 1.0)

### 风险提示
- 需要安装驱动（管理员权限）
- 用户态文件系统有性能开销
- 需要处理复杂的 Dokan 接口

---

## Work Objectives

### 核心交付
1. Dokan 虚拟文件系统基础框架
2. 写缓冲实现 (write-behind)
3. 读缓冲实现 (read-ahead)
4. 配置系统（内存限制、缓冲策略）
5. 性能测试与基准测试

### 定义完成
- [ ] `wutils buffer mount X:` 命令可挂载缓冲盘
- [ ] 挂载后对 X: 的读写自动经过缓冲层
- [ ] 配置文件可设置内存上限、刷新策略
- [ ] 基准测试显示 IOPS 改善

---

## Verification Strategy

### 测试策略
- **基础设施**：Go testing.B 基准测试
- **测试场景**：
  - 小文件随机写入 (4KB, 1000次)
  - 小文件随机读取 (4KB, 1000次)
  - 混合读写工作负载

### QA 场景
- 挂载/卸载稳定性
- 大文件写入
- 并发读写
- 错误处理（磁盘满、权限问题）

---

## Execution Strategy

### 模块划分

```
Wave 1 (基础 - 2 周):
├── Task 1: 项目初始化 + 依赖配置
├── Task 2: Dokan 基础框架 (挂载/卸载)
├── Task 3: 基础文件系统实现 (读文件)
└── Task 4: 写文件接口实现

Wave 2 (缓冲核心 - 3 周):
├── Task 5: 写缓冲层 - 内存缓冲
├── Task 6: 写缓冲层 - 异步刷盘 (write-behind)
├── Task 7: 读缓冲层 - LRU 缓存
├── Task 8: 读缓冲层 - 预读策略 (read-ahead)
└── Task 9: 缓冲策略配置系统

Wave 3 (完善 + 测试 - 2 周):
├── Task 10: 命令行接口 (mount/unmount)
├── Task 11: 性能基准测试
├── Task 12: 错误处理完善
└── Task 13: 文档与示例
```

---

## TODOs

### Wave 1: 基础框架

- [x] 1. 项目初始化 + 依赖配置

  **What to do**:
  - 创建 `cmd/wutils/buffer` 目录
  - 添加依赖：github.com/binzume/dkango
  - 创建 go.mod 更新（如需要）
  - 验证 Dokany 驱动可用

  **References**:
  - dkango: https://github.com/binzume/dkango
  - Dokan: https://github.com/dokan-dev/dokany

  **Acceptance Criteria**:
  - [ ] `go build` 通过
  - [ ] Dokany 驱动已安装验证

- [x] 2. Dokan 基础框架 (挂载/卸载)

  **What to do**:
  - 实现 `BufferFS` 结构体
  - 实现 `Mount(drive string, config *BufferConfig) error`
  - 实现 `Unmount(drive string) error`
  - 实现基本的 FileSystem 接口

  **References**:
  - dkango examples: https://github.com/binzume/dkango/tree/master/examples

  **Acceptance Criteria**:
  - [ ] `buffer mount X:` 可挂载虚拟盘 X:
  - [ ] `buffer unmount X:` 可卸载

- [x] 3. 基础文件系统实现 (读文件)

  **What to do**:
  - 实现 File.ReadFile 读取后端存储
  - 实现目录遍历
  - 实现文件属性查询

  **Acceptance Criteria**:
  - [x] 可读取虚拟盘中的文件
  - [x] 目录列表正常工作

- [x] 4. 写文件接口实现

  **What to do**:
  - 实现 File.WriteFile 写入后端存储
  - 实现文件创建、删除

  **Acceptance Criteria**:
  - [x] 可写入文件到虚拟盘
  - [x] 文件创建/删除正常

### Wave 2: 缓冲核心

- [x] 5. 写缓冲层 - 内存缓冲
- [x] 6. 写缓冲层 - 异步刷盘 (write-behind)
- [x] 7. 读缓冲层 - LRU 缓存
- [x] 8. 读缓冲层 - 预读策略 (read-ahead)

  **What to do**:
  - 实现 LRU 缓存结构
  - 实现缓存淘汰策略
  - 配置缓存大小上限

  **Acceptance Criteria**:
  - [ ] 热点数据缓存在内存
  - [ ] 缓存大小可配置

- [ ] 8. 读缓冲层 - 预读策略 (read-ahead)

  **What to do**:
  - 实现顺序访问预读
  - 实现热点块预取
  - 配置预读块大小

  **Acceptance Criteria**:
  - [ ] 顺序读时预加载后续数据
  - [ ] 预读可配置开关

- [ ] 9. 缓冲策略配置系统

  **What to do**:
  - 定义 BufferConfig 结构
  - 支持 YAML 配置
  - 与现有 config 系统集成

  **Acceptance Criteria**:
  - [ ] 配置可控制内存上限、刷新间隔等

### Wave 3: 完善 + 测试

- [ ] 10. 命令行接口 (mount/unmount)

  **What to do**:
  - 添加 `wutils buffer mount` 命令
  - 添加 `wutils buffer unmount` 命令
  - 添加 `wutils buffer status` 命令

  **Acceptance Criteria**:
  - [ ] CLI 可用

- [ ] 11. 性能基准测试

  **What to do**:
  - 编写 benchmark 测试
  - 对比缓冲 vs 无缓冲性能
  - 测试不同工作负载

  **Acceptance Criteria**:
  - [ ] 有量化性能数据
  - [ ] IOPS 改善显著

- [ ] 12. 错误处理完善

  **What to do**:
  - 处理磁盘满情况
  - 处理权限错误
  - 处理后端存储失败

  **Acceptance Criteria**:
  - [ ] 错误有清晰提示

- [ ] 13. 文档与示例

  **What to do**:
  - 添加命令使用文档
  - 添加配置示例

  **Acceptance Criteria**:
  - [ ] README 更新

---

## Commit Strategy
- Message: `feat(buffer): add Dokan-based IO buffer filesystem`
- Files: cmd/wutils/buffer/

---

## Success Criteria

### 验证命令
```bash
# 挂载
wutils buffer mount X: --source D: --memory-limit 512MB --flush-interval 5s

# 测试写入
for i in {1..1000}; do echo "test $i" > "X:/file_$i.txt"; done

# 卸载
wutils buffer unmount X:
```

### 性能目标
- 小文件随机写 IOPS 提升 >= 10x
- 内存缓冲延迟 < 10ms
- 定期刷盘不阻塞主操作
