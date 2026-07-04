# gobt

gobt 是一个用 Go 逐步实现 BitTorrent 客户端的学习型项目。

## 当前进度

已完成模块:
- pkg/bencode: bencode 编码与解码
- pkg/torrent: torrent 元数据模型、解析/编码、校验、info hash
- pkg/tracker: HTTP announce 请求构建与响应解析
- pkg/utils: 版本与通用工具

当前可用状态:
- 全仓测试可通过
- 主程序可运行
- examples 目录包含可直接运行的示例和样例 torrent 文件

下一阶段:
- DHT 基础
- Peer discovery 聚合 (Tracker + DHT)
- Peer 握手与消息协议
- 分片下载、校验与续传

## 快速开始

```bash
# 1) 克隆项目
git clone https://github.com/faterer/gobt.git
cd gobt

# 2) 运行测试
go test ./...

# 3) 运行主程序
go run ./cmd

# 4) 运行示例
cd examples
go run parse_torrent.go init.go

# 可选: 带构建标签的独立示例
go run -tags bencode_example bencode_simple.go
go run -tags tracker_example tracker_announce.go <torrent-file> <tracker-url>
```

## 示例文件

样例 .torrent 文件已统一放在 examples 目录:
- examples/example-demo.torrent
- examples/example-file.torrent
- examples/demo-multifile.torrent

## 目录结构

```text
gobt/
├── cmd/                # 主程序入口
├── pkg/                # 核心模块
│   ├── bencode/
│   ├── torrent/
│   ├── tracker/
│   └── utils/
├── examples/           # 示例程序与样例数据
├── docs/               # 当前维护文档
│   └── archive/        # 历史归档文档
├── project_plan.md     # Week 1-12 执行计划
└── go.mod
```

## 文档入口

建议阅读顺序:
1. project_plan.md
2. docs/GETTING_STARTED.md
3. docs/INDEX.md
4. docs/ROADMAP.md
5. docs/ARCHITECTURE.md
6. docs/QUICKREF.md

历史阶段文档已归档在 docs/archive/。

## 常用验证命令

```bash
# 全量测试
go test ./...

# tracker 包测试
go test ./pkg/tracker -v

# 覆盖率 (可选)
go test -coverprofile=coverage ./...
```

## License

MIT
