# BitTorrent 4.2.0 需求文档

## 1. 项目概述

本项目使用Go语言从零实现一个与BitTorrent 4.2.0兼容的完整分布式文件传输系统。该系统遵循BitTorrent协议规范，支持高效的点对点（P2P）文件共享。

**版本**: 4.2.0  
**开发语言**: Go 1.25+  
**目标平台**: Windows, Linux, macOS

---

## 2. 核心功能需求

### 2.1 Torrent文件解析与处理
- **要求**: 能解析标准.torrent文件（Bencode编码格式）
- **功能**:
  - 解析torrent元数据（名称、大小、文件列表）
  - 提取piece hash值用于数据验证
  - 支持单文件和多文件torrent
  - 支持announce和announce-list（多tracker）
  - 解析info_hash和piece长度

### 2.2 Tracker通信
- **DHT (Distributed Hash Table)** 支持
  - 节点ID生成和管理
  - Kademlia算法实现
  - 支持DHT网络引导
  
- **Tracker HTTP/UDP通讯**
  - 发送announce请求获取peer列表
  - 支持event参数：started, stopped, completed
  - 解析tracker响应获取peer信息
  - 心跳机制，定期向tracker报告状态

### 2.3 Peer发现与连接
- **Peer discovery**:
  - 从tracker获取peer列表
  - DHT网络peer发现
  - PEX (Peer Exchange) 协议
  - 本地peer发现（LSD - Local Service Discovery）

- **Peer连接管理**:
  - 并发连接多个peers
  - TCP连接建立和管理
  - 连接池维护（默认50-200个连接）
  - 自适应连接数调整

### 2.4 BitTorrent协议实现
- **握手阶段**:
  - 19字节协议头："BitTorrent protocol"
  - 8字节reserved标志位
  - 20字节info_hash
  - 20字节peer_id

- **消息类型支持**:
  - keep-alive
  - choke/unchoke
  - interested/not interested
  - have/bitfield
  - request/piece/cancel
  - port
  - extended messages (BEP 10)

### 2.5 下载引擎
- **数据下载**:
  - Piece级别下载管理
  - 自适应piece选择算法
  - 稀有piece优先（Rarest First）
  - End-game模式（下载完成阶段加速）

- **并发下载优化**:
  - 单个piece多source并行下载
  - Pipeline式请求管理
  - 智能去重和流控制

- **数据校验**:
  - SHA1 hash校验（BEP 3）
  - 错误数据丢弃和重新请求
  - 坏peer检测和黑名单

### 2.6 上传引擎
- **上传管理**:
  - 选择unchoke策略（BEP 6 - Choking algorithm）
  - 配合下载进度的unchoke决策
  - 防止leeching攻击
  - 带宽限制和QoS

### 2.7 配置与性能
- **可配置参数**:
  - 上下行带宽限制
  - 最大连接数
  - Piece大小配置
  - 请求超时时间
  - Tracker更新间隔

- **性能指标**:
  - 实时下载/上传速度
  - 连接状态监控
  - 内存使用优化
  - CPU使用率控制

---

## 3. 非功能性需求

### 3.1 可靠性
- 异常断线自动重连
- 数据损坏自动恢复
- Graceful shutdown处理
- 完整性日志记录

### 3.2 性能
- 支持1000+个并发peers
- 毫秒级响应时间
- 内存占用<500MB（默认配置）
- 支持大文件（TB级）

### 3.3 兼容性
- 兼容标准BitTorrent客户端
- 支持IPv4和IPv6
- 跨平台运行
- 支持各种编码格式的路径

### 3.4 安全性
- 防止恶意peer攻击
- 验证所有接收数据
- 隐私保护（可选）
- PEX和DHT安全防护

---

## 4. 详细功能规范

### 4.1 Bencode编码支持
- 整数：i{number}e (例: i42e)
- 字符串：{length}:{string} (例: 4:spam)
- 列表：l{elements}e
- 字典：d{key}{value}e (按key字典序排序)

### 4.2 Info Hash计算
- 对torrent文件的info部分进行SHA1 hash
- 结果为20字节的二进制数据
- 用于tracker和peer通信

### 4.3 Peer ID生成
- 20字节标识
- 格式: -GO4200-{random12bytes}
- 其中GO = Go实现标志，4200 = 版本4.2.0

### 4.4 协议扩展（BEP支持）
- BEP 3: Protocol specification
- BEP 6: Choking algorithm
- BEP 10: Extension Protocol
- BEP 11: Peer Exchange (PEX)
- BEP 14: DHT
- BEP 20: Peer ID convention

### 4.5 下载流程
```
1. 用户提供.torrent文件
2. 解析torrent元数据
3. 启动Tracker和DHT，获取peer列表
4. 连接多个peers
5. 请求piece数据
6. 验证hash，保存到磁盘
7. 协调上传已有pieces
8. 下载完成后做种
```

### 4.6 上传流程
```
1. 维护bitfield表示已有pieces
2. 响应peer的interested请求
3. 根据choking算法决定unchoke
4. 发送piece数据给请求的peer
5. 带宽限制和流控制
```

---

## 5. 接口需求

### 5.1 命令行接口
```bash
gop2p start <torrent-file>           # 开始下载/上传
gop2p info <torrent-file>            # 显示torrent信息
gop2p status                         # 显示当前状态
gop2p stop                           # 停止任务
gop2p config [key] [value]           # 配置参数
```

### 5.2 RPC/API接口
- RESTful API查询下载状态
- 实时事件流推送
- 带宽控制API
- Peer管理API

### 5.3 Web UI（可选）
- 显示下载进度
- 连接peer列表
- 网络统计
- 性能图表

---

## 6. 存储需求

### 6.1 文件存储
- 完整下载的文件
- 部分下载的临时数据
- Resume支持（记录下载状态）

### 6.2 元数据存储
- Torrent信息缓存
- 已知peers缓存
- DHT节点缓存
- Session状态

---

## 7. 限制与约束

### 7.1 协议限制
- Piece大小: 16KB - 16MB（通常256KB-1MB）
- 单个请求大小: 最大16KB
- 最大peer连接数: 500（可配置）
- Tracker重试间隔: 30秒 - 1小时

### 7.2 资源限制
- 内存: <1GB
- 磁盘I/O优化
- 网络连接优化
- CPU核心利用

---

## 8. 扩展功能（Phase 2+）

### 8.1 加密支持
- Protocol Encryption (BEP 20)
- MSE (Message Stream Encryption)

### 8.2 高级功能
- Magnet link支持
- WebTorrent兼容性
- Streaming模式
- Bandwidth estimation

### 8.3 生态集成
- VPN支持
- 代理支持
- 分析和日志系统
- 插件系统

---

## 9. 成功指标

- [ ] 能成功解析标准.torrent文件
- [ ] 能与标准BitTorrent tracker通信
- [ ] 能发现并连接真实peers
- [ ] 能完整下载小型torrent（<100MB）
- [ ] 能完整下载大型torrent（>1GB）
- [ ] 上传速度>10MB/s（测试环境）
- [ ] 内存占用<500MB
- [ ] 兼容性测试通过（与Transmission, qBittorrent等）
