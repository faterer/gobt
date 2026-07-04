# Project Plan (Week 1-12)

gobt 项目执行总计划（精简版）。

## 终极目标

交付一个可用的 BitTorrent 客户端，具备：
- torrent 解析
- Tracker + DHT + Peer Discovery
- Peer 连接与消息交换
- 分片下载、校验、续传
- 可稳定运行与发布

## 执行规则

- 每周更新一次本文件
- 状态仅用：`Not Started` / `In Progress` / `Done`
- 每周备注只写：完成项、阻塞项、下周动作

## 12周路线图

| Week | 主题 | 核心交付 |
|---|---|---|
| 1 | 项目基础 | 工程结构、测试流程、CLI 入口 |
| 2 | Bencode | 编解码 + 测试 |
| 3 | Torrent | 模型、解析、校验、Info Hash |
| 4 | Tracker | Announce 请求/响应、Peers 解析 |
| 5 | DHT | 节点与基础查询（find_node/get_peers） |
| 6 | Peer Discovery | 聚合 Tracker + DHT，去重与评分 |
| 7 | Peer 连接 | TCP 生命周期 + Handshake |
| 8 | 消息协议 | 核心消息与状态机 |
| 9 | 下载流水线 | Piece/Block 请求、重试、组装 |
| 10 | 存储校验 | 落盘、SHA1 校验、坏块重下 |
| 11 | 续传稳定性 | 断点恢复、长稳运行 |
| 12 | 性能发布 | 调优、文档、版本发布 |

## 当前进度

| Week | Status | Notes | Last Update |
|---|---|---|---|
| 1 | Done | 项目基础完成 | 2026-07-04 |
| 2 | Done | Bencode 完成并有测试 | 2026-07-04 |
| 3 | Done | Torrent 解析与校验完成 | 2026-07-04 |
| 4 | Done | Tracker 通信完成并测试通过 | 2026-07-04 |
| 5 | Not Started | 下一步：DHT 基础 | 2026-07-04 |
| 6 | Not Started | 依赖 Week 5 | 2026-07-04 |
| 7 | Not Started | 依赖 Discovery 输出 | 2026-07-04 |
| 8 | Not Started | 协议与状态机 | 2026-07-04 |
| 9 | Not Started | 下载流水线 | 2026-07-04 |
| 10 | Not Started | 存储与校验 | 2026-07-04 |
| 11 | Not Started | 续传与稳定性 | 2026-07-04 |
| 12 | Not Started | 性能与发布 | 2026-07-04 |

## 本周行动（Week 5）

- Day 1: DHT 消息格式与 bootstrap 策略
- Day 2: 节点距离计算与节点表
- Day 3: find_node 流程跑通
- Day 4: get_peers 流程跑通
- Day 5: 单元测试 + 集成检查
