# gobt - Go BitTorrent 4.2.0 实现

> 一个用Go从零实现的完整BitTorrent 4.2.0客户端，支持高效的P2P文件传输。

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Development-yellow)](docs/SUMMARY.md)

## 🎯 项目目标

实现一个**与标准BitTorrent 4.2.0完全兼容**的P2P文件共享系统，支持：
- ✅ 与任何标准BitTorrent客户端互操作
- ✅ 高效的多peer并发下载
- ✅ 完整的DHT网络支持
- ✅ 从零开发，无重依赖
- ✅ 生产级别的性能和稳定性

## 📋 快速导航

### 📚 文档
| 文档 | 用途 | 读者 |
|------|------|------|
| [REQUIREMENTS.md](docs/REQUIREMENTS.md) | 详细功能需求说明 | 产品经理、测试 |
| [DESIGN.md](docs/DESIGN.md) | 系统架构和设计决策 | 架构师、核心开发 |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | 代码结构和实现细节 | 开发工程师 |
| [SUMMARY.md](docs/SUMMARY.md) | 项目总体概览 | 所有人 |

### 🚀 快速开始

```bash
# 1. 克隆项目
git clone https://github.com/yourname/gobt.git
cd gobt

# 2. 初始化
go mod download
go mod tidy

# 3. 运行
go run ./cmd start ubuntu.iso.torrent

# 4. 查看帮助
go run ./cmd help
```

---

## 📊 项目进度概览

### MVP完成度: 0% → 25%

```
Phase 1 - 基础实现        ⏳ ░░░░░░░░░░░░░░░░░░░░ 0%
├─ Bencode编解码          ⏳
├─ Torrent解析             ⏳
├─ Tracker通信             ⏳
├─ 握手与消息              ⏳
└─ 基础下载               ⏳

Phase 2 - 完整功能        ⏳ ░░░░░░░░░░░░░░░░░░░░ 0%
├─ DHT实现               ⏳
├─ 下载优化              ⏳
├─ Upload/Choking        ⏳
└─ 长期稳定性             ⏳

Phase 3 - 生产就绪         ⏳ ░░░░░░░░░░░░░░░░░░░░ 0%
├─ 性能优化              ⏳
├─ 全面测试              ⏳
├─ 文档完善              ⏳
└─ 发布v1.0              ⏳
```

---

## 📦 项目结构

```
gobt/
├── docs/                          # 📚 文档
│   ├── REQUIREMENTS.md            # 需求文档
│   ├── DESIGN.md                  # 设计文档
│   ├── ARCHITECTURE.md            # 架构文档
│   └── SUMMARY.md                 # 项目总结
│
├── pkg/                           # 📦 核心库
│   ├── bencode/                   # Bencode编解码
│   ├── torrent/                   # Torrent文件处理
│   ├── protocol/                  # BitTorrent协议
│   ├── tracker/                   # Tracker通信
│   ├── dht/                       # DHT网络
│   ├── network/                   # 网络层
│   ├── storage/                   # 存储管理
│   ├── hash/                      # 数据验证
│   ├── core/                      # 核心业务逻辑
│   ├── config/                    # 配置管理
│   ├── logger/                    # 日志系统
│   └── utils/                     # 工具库
│
├── cmd/                           # ⚙️ 可执行程序
│   ├── main.go
│   ├── cli.go
│   └── flags.go
│
├── tests/                         # ✅ 测试
│   ├── integration/
│   ├── fixtures/
│   └── mocks/
│
├── config/                        # ⚙️ 配置示例
│   └── gobt.yaml
│
├── go.mod                         # 依赖定义
├── go.sum                         # 依赖校验
├── Makefile                       # 构建脚本
├── README.md                      # 本文件
└── LICENSE                        # MIT许可证
```

---

## 🔧 核心功能

### 已实现 ✅
- [ ] 待开始

### 开发中 🔨
- [ ] Bencode编解码
- [ ] Torrent文件解析
- [ ] Tracker通信

### 计划中 📅
- [ ] DHT网络
- [ ] 完整下载引擎
- [ ] 上传和做种
- [ ] Web UI

---

## 💾 关键数据结构

### Torrent元数据
```go
type Torrent struct {
    Announce     string              // tracker地址
    AnnounceList [][]string          // 备用tracker
    CreationDate int64
    Comment      string
    Info         *Info
}

type Info struct {
    Length       int64               // 文件大小
    Name         string              // 文件名
    PieceLength  int                 // piece大小（通常256KB）
    Pieces       string              // SHA1 hash列表
}
```

### Peer信息
```go
type Peer struct {
    ID           [20]byte            // Peer标识
    IP           string              // IP地址
    Port         uint16              // 端口号
    
    // 协议状态
    AmChoking    bool                // 我们是否choke
    PeerChoking  bool                // 对方是否choke
    
    // 进度
    Bitfield     *Bitfield          // 已有pieces
    Uploaded     int64               // 上传字节数
    Downloaded   int64               // 下载字节数
}
```

---

## 🔄 工作流程

### 下载流程示意

```
┌─────────────────────────────────┐
│   用户提供.torrent文件          │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│   解析torrent，计算info_hash   │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│   连接tracker获取peer列表       │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│   建立peer连接 (并发)           │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│   发送request下载piece          │
├─────────────────────────────────┤
│   ├─ 接收piece数据              │
│   ├─ 验证SHA1 hash             │
│   ├─ 保存到磁盘                 │
│   └─ 广播已有pieces             │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│   下载完成，继续做种             │
└─────────────────────────────────┘
```

---

## 🎓 学习资源

### BitTorrent规范
- [BEP 3: Protocol Specification](http://www.bittorrent.org/beps/bep_0003.html) - 核心协议
- [BEP 6: Fast Extension](http://www.bittorrent.org/beps/bep_0006.html) - Choking算法
- [BEP 14: DHT](http://www.bittorrent.org/beps/bep_0014.html) - DHT网络

### 算法资源
- [Kademlia算法](https://en.wikipedia.org/wiki/Kademlia) - P2P路由
- [Rarest First](https://en.wikipedia.org/wiki/BitTorrent#Strategy) - 下载策略

### Go相关
- [Go官方文档](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)

---

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./pkg/bencode/...

# 生成测试覆盖率报告
make coverage

# 运行性能测试
go test -bench=. -benchmem ./...

# 集成测试
go test ./tests/integration/...
```

---

## 🏗️ 构建与部署

### 本地构建
```bash
# 编译
make build

# 编译并运行
make run TORRENT=example.torrent

# 清理
make clean
```

### 跨平台编译
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o gobt-linux ./cmd

# macOS
GOOS=darwin GOARCH=amd64 go build -o gobt-darwin ./cmd

# Windows
GOOS=windows GOARCH=amd64 go build -o gobt.exe ./cmd
```

### Docker
```bash
# 构建镜像
docker build -t gobt:latest .

# 运行容器
docker run -v /data:/app/downloads gobt:latest start example.torrent
```

---

## 💡 设计原则

| 原则 | 描述 |
|------|------|
| **简洁性** | 最小化依赖，优先使用标准库 |
| **兼容性** | 严格遵守BEP规范，与标准客户端互操作 |
| **性能** | 高效的并发模型和内存使用 |
| **可维护性** | 清晰的代码结构和充分的文档 |
| **可测试性** | 完善的单元测试和集成测试 |

---

## 🔐 安全性考虑

### 已实现
- [ ] SHA1 hash验证确保数据完整性
- [ ] Peer黑名单防止恶意peer
- [ ] 连接超时保护

### 规划
- [ ] 连接加密（Protocol Encryption）
- [ ] 防止DDoS攻击
- [ ] 隐私保护选项

---

## 📈 性能目标

### 下载速度
| 文件大小 | 预期速度 | 目标时间 |
|---------|---------|---------|
| 10MB | 1MB/s | 10秒 |
| 100MB | 5MB/s | 20秒 |
| 1GB | 20MB/s | 50秒 |
| 10GB | 50MB/s | 3-5分钟 |

### 资源占用
| 指标 | 目标 | 备注 |
|------|------|------|
| 内存占用 | <500MB | 200个peer连接 |
| CPU使用率 | <20% | 单核使用率 |
| 磁盘I/O | <100MB/s | 受限于硬件 |

---

## 🔄 版本计划

### v1.0.0 - MVP (2026年Q3)
- ✅ 基础协议实现
- ✅ 完整下载功能
- ✅ Tracker和DHT支持
- ✅ 命令行界面

### v1.1.0 - 优化 (2026年Q4)
- Upload/Choking优化
- 性能调优
- 广泛的兼容性测试

### v2.0.0 - 扩展 (2027年H1)
- Web UI
- 加密支持
- 高级功能（magnet link等）

---

## 🤝 贡献指南

### 参与开发
1. Fork项目
2. 创建feature分支: `git checkout -b feature/your-feature`
3. 提交更改: `git commit -am 'Add feature'`
4. 推送到分支: `git push origin feature/your-feature`
5. 创建Pull Request

### 代码风格
- 遵循 `gofmt` 格式
- 使用 `golangci-lint` 检查
- 编写充分的单元测试
- 提供清晰的commit信息

### 报告Bug
使用GitHub Issues提交bug报告，请包含：
- 问题描述
- 复现步骤
- 环境信息（OS、Go版本等）
- 错误日志

---

## 📝 许可证

本项目采用MIT许可证。详见 [LICENSE](LICENSE)

```
MIT License

Copyright (c) 2026 gobt contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
...
```

---

## 👥 团队

- **项目发起**: 个人开发项目
- **贡献者**: 欢迎加入！

---

## 📞 支持

### 获取帮助
- 📖 查看 [文档](docs/)
- 🐛 提交 [Issues](https://github.com/yourname/gobt/issues)
- 💬 讨论 [Discussions](https://github.com/yourname/gobt/discussions)

### 联系方式
- Email: your-email@example.com
- Twitter: [@yourhandle](https://twitter.com/yourhandle)

---

## 🎉 致谢

感谢以下资源和项目的参考和启发：
- BitTorrent官方规范 (BEP)
- Transmission 客户端
- qBittorrent 项目
- Kademlia 论文

---

## 📊 项目统计

```
预期代码量:        ~10,000 LOC
测试覆盖率:        >80%
文档页数:          50+
支持的BEP:         10+
```

---

**最后更新**: 2026-07-01  
**状态**: 需求分析完成，准备开发  
**下一步**: 开始Phase 1实现

---

## 相关链接

- 🌐 [项目主页](https://github.com/yourname/gobt)
- 📚 [完整文档](docs/)
- 🎯 [实现路线图](docs/ROADMAP.md)
- ✅ [测试报告](docs/TEST_RESULTS.md)

