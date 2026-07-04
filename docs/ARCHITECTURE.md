# BitTorrent 4.2.0 架构详解文档

## 1. 项目初始化

### 1.1 目录结构创建

```bash
gobt/
├── cmd/                          # 可执行程序入口
│   ├── main.go                   # 程序主入口
│   ├── cli.go                    # 命令行解析
│   └── flags.go                  # 命令行标志定义
│
├── pkg/                          # 可复用的库代码
│   ├── bencode/                  # Bencode编码/解码
│   │   ├── bencode.go           # 主模块
│   │   ├── encoder.go           # 编码器实现
│   │   ├── decoder.go           # 解码器实现
│   │   ├── types.go             # 类型定义
│   │   └── bencode_test.go
│   │
│   ├── torrent/                 # Torrent文件处理
│   │   ├── parser.go            # 文件解析
│   │   ├── metadata.go          # 元数据结构
│   │   ├── validator.go         # 验证逻辑
│   │   └── torrent_test.go
│   │
│   ├── protocol/                # BitTorrent协议
│   │   ├── handshake.go         # 握手协议
│   │   ├── messages.go          # 消息定义
│   │   ├── parser.go            # 消息解析
│   │   ├── constants.go         # 常量定义
│   │   └── protocol_test.go
│   │
│   ├── tracker/                 # Tracker通信
│   │   ├── manager.go           # 管理器
│   │   ├── http.go              # HTTP tracker
│   │   ├── udp.go               # UDP tracker
│   │   ├── announce.go          # Announce请求
│   │   └── tracker_test.go
│   │
│   ├── dht/                     # DHT网络
│   │   ├── node.go              # DHT节点
│   │   ├── kademlia.go          # Kademlia算法
│   │   ├── routing_table.go     # 路由表
│   │   ├── message.go           # DHT消息
│   │   └── dht_test.go
│   │
│   ├── network/                 # 网络层
│   │   ├── peer.go              # Peer定义
│   │   ├── peer_manager.go      # Peer管理
│   │   ├── connection.go        # 连接处理
│   │   ├── connection_pool.go   # 连接池
│   │   ├── message_handler.go   # 消息处理
│   │   └── network_test.go
│   │
│   ├── storage/                 # 存储管理
│   │   ├── file_manager.go      # 文件管理
│   │   ├── piece_store.go       # Piece存储
│   │   ├── metadata_store.go    # 元数据存储
│   │   ├── cache.go             # 缓存层
│   │   └── storage_test.go
│   │
│   ├── hash/                    # 哈希验证
│   │   ├── verifier.go          # 验证器
│   │   └── hash_test.go
│   │
│   ├── core/                    # 核心业务逻辑
│   │   ├── session.go           # 会话管理
│   │   ├── download.go          # 下载管理
│   │   ├── upload.go            # 上传管理
│   │   ├── scheduler.go         # 调度器
│   │   ├── bitfield.go          # Bitfield
│   │   ├── state.go             # 状态管理
│   │   └── core_test.go
│   │
│   ├── config/                  # 配置管理
│   │   ├── config.go            # 配置结构
│   │   ├── defaults.go          # 默认配置
│   │   └── loader.go            # 配置加载
│   │
│   ├── logger/                  # 日志系统
│   │   ├── logger.go            # 日志记录
│   │   └── metrics.go           # 性能指标
│   │
│   └── utils/                   # 工具库
│       ├── peer_id.go           # PeerID生成
│       ├── info_hash.go         # InfoHash计算
│       ├── net_utils.go         # 网络工具
│       ├── file_utils.go        # 文件工具
│       └── misc.go              # 杂项工具
│
├── internal/                    # 内部代码（不对外暴露）
│   └── version.go               # 版本信息
│
├── docs/                        # 文档
│   ├── REQUIREMENTS.md          # 需求文档
│   ├── DESIGN.md                # 设计文档
│   ├── ARCHITECTURE.md          # 架构文档（本文件）
│   ├── API.md                   # API文档
│   ├── PROTOCOL.md              # 协议细节
│   └── CONTRIBUTING.md          # 贡献指南
│
├── config/                      # 配置文件示例
│   ├── gobt.yaml              # YAML配置
│   └── gobt.toml              # TOML配置
│
├── tests/                       # 测试文件
│   ├── integration/
│   │   ├── download_test.go
│   │   ├── upload_test.go
│   │   └── e2e_test.go
│   ├── fixtures/
│   │   ├── sample.torrent      # 测试用torrent文件
│   │   └── testdata/           # 测试数据
│   └── mocks/
│       └── mock_tracker.go     # Mock Tracker
│
├── scripts/                     # 脚本
│   ├── build.sh                # 构建脚本
│   ├── test.sh                 # 测试脚本
│   └── deploy.sh               # 部署脚本
│
├── .github/                     # GitHub配置
│   ├── workflows/
│   │   ├── ci.yml              # CI/CD配置
│   │   └── release.yml         # 发布配置
│   └── ISSUE_TEMPLATE/
│
├── go.mod                       # Go模块文件
├── go.sum                       # 依赖checksum
├── Makefile                     # Makefile
├── README.md                    # 项目说明
├── LICENSE                      # 许可证
└── .gitignore                   # Git忽略配置
```

---

## 2. 核心数据结构

### 2.1 Bencode类型系统

```go
// bencode/types.go

type Value interface {
    MarshalBencode() ([]byte, error)
    UnmarshalBencode([]byte) error
}

type Integer int64

type String []byte

type List []Value

type Dict map[string]Value
```

### 2.2 Torrent元数据

```go
// torrent/metadata.go

type Torrent struct {
    // 可选：是否为单文件模式
    Announce     string              // Primary tracker
    AnnounceList [][]string          // Backup trackers
    CreationDate int64               // Unix timestamp
    Comment      string
    CreatedBy    string
    Encoding     string              // 默认 UTF-8
    
    // Info部分（必需）
    Info         *Info
}

type Info struct {
    // 单文件
    Length       int64               // 文件大小
    Name         string              // 文件名
    
    // 多文件
    Files        []FileInfo          // 文件列表
    
    // 共同
    PieceLength  int                 // 通常262144 (256KB)
    Pieces       string              // SHA1 hashes (raw binary)
    Private      int                 // 0 or 1
}

type FileInfo struct {
    Length       int64
    Path         []string            // 目录路径
}
```

### 2.3 Peer结构

```go
// network/peer.go

type Peer struct {
    // 标识
    ID           [20]byte
    IP           net.IP
    Port         uint16
    Address      string              // "ip:port"
    Source       PeerSource          // tracker, dht, pex, lsd
    
    // 连接
    Conn         net.Conn
    Connected    bool
    HandshakeDone bool
    
    // 协议状态 (BEP 3)
    AmChoking    bool                // 我们是否choke该peer
    AmInterested bool                // 我们是否interested
    PeerChoking  bool                // 对方是否choke我们
    PeerInterested bool              // 对方是否interested
    
    // Bitfield (对方拥有哪些pieces)
    Bitfield     *Bitfield
    
    // 统计
    Uploaded     int64
    Downloaded   int64
    UploadRate   float64             // bytes/sec
    DownloadRate float64             // bytes/sec
    
    // 时间
    ConnectTime  time.Time
    LastActivity time.Time
    LastChoke    time.Time
    
    // 性能
    Score        float64             // 用于排序
}

type PeerSource string

const (
    SourceTracker PeerSource = "tracker"
    SourceDHT     PeerSource = "dht"
    SourcePEX     PeerSource = "pex"
    SourceLSD     PeerSource = "lsd"
)
```

### 2.4 Piece结构

```go
// core/piece.go

type Piece struct {
    Index        int
    Size         int
    Hash         [20]byte
    
    // 下载进度
    Downloaded   int
    Data         []byte
    
    // 状态
    State        PieceState
    
    // 优先级
    Rarity       int                 // 有多少个peer拥有
    Priority     float64             // 调度优先级
    
    // 时间戳
    LastActivity time.Time
}

type PieceState int

const (
    PieceMissing PieceState = iota
    PieceDownloading
    PieceDone
    PieceVerified
)
```

### 2.5 Session结构

```go
// core/session.go

type Session struct {
    // 标识
    InfoHash     [20]byte
    PeerID       [20]byte
    
    // 元数据
    Torrent      *Torrent
    Metadata     *Metadata
    
    // 状态
    State        SessionState
    Progress     float64
    
    // 核心组件
    PeerMgr      *PeerManager
    DownloadMgr  *DownloadManager
    UploadMgr    *UploadManager
    Scheduler    *PieceScheduler
    StorageMgr   *StorageManager
    
    // 发现
    TrackerMgr   *TrackerManager
    DHT          *DHTNode
    
    // 统计
    Stats        SessionStats
    
    // 控制
    ctx          context.Context
    cancel       context.CancelFunc
    wg           sync.WaitGroup
}

type SessionState int

const (
    StateIdle SessionState = iota
    StateInitializing
    StateDownloading
    StatePaused
    StateSeeding
    StateStopping
    StateStopped
)
```

---

## 3. 模块间通信

### 3.1 事件系统

```go
// core/events.go

type Event interface {
    Type() EventType
}

type EventType int

const (
    EventPeerConnected EventType = iota
    EventPeerDisconnected
    EventPieceDownloaded
    EventPieceVerified
    EventPieceFailed
    EventDownloadComplete
    EventError
    EventStatusUpdate
)

type EventBus struct {
    subscribers map[EventType][]Handler
    mu          sync.RWMutex
}

type Handler func(Event) error

func (eb *EventBus) Subscribe(et EventType, h Handler)
func (eb *EventBus) Publish(e Event) error
```

### 3.2 消息队列

```go
// network/message_queue.go

type MessageQueue struct {
    ch          chan Message
    maxSize     int
    mu          sync.Mutex
}

type Message struct {
    ID        MessageID
    Payload   interface{}
    Timestamp time.Time
    Retry     int
}

type MessageID int

const (
    MsgHandshake MessageID = iota
    MsgKeepAlive
    MsgChoke
    MsgUnchoke
    MsgInterested
    MsgNotInterested
    MsgHave
    MsgBitfield
    MsgRequest
    MsgPiece
    MsgCancel
    MsgPort
)
```

### 3.3 状态同步

```go
// core/state_sync.go

type StateSync struct {
    bitfield    *Bitfield              // 已下载pieces
    uploads     map[string]int64       // peer_addr -> bytes
    downloads   map[string]int64       // peer_addr -> bytes
    peers       map[string]*Peer       // 活跃peer
    mu          sync.RWMutex
}

func (ss *StateSync) UpdatePiece(idx int, verified bool)
func (ss *StateSync) GetProgress() float64
func (ss *StateSync) GetStats() SessionStats
```

---

## 4. 关键算法实现

### 4.1 Kademlia算法 (DHT)

```go
// dht/kademlia.go

// 距离度量：XOR距离
func XORDistance(a, b [20]byte) [20]byte {
    var result [20]byte
    for i := 0; i < 20; i++ {
        result[i] = a[i] ^ b[i]
    }
    return result
}

// K-bucket选择
func (rt *RoutingTable) FindClosest(target [20]byte, k int) []*NodeInfo {
    bucketIndex := rt.getBucketIndex(target)
    nodes := []*NodeInfo{}
    
    // 从目标bucket开始搜索
    for i := bucketIndex; i < len(rt.buckets) && len(nodes) < k; i++ {
        nodes = append(nodes, rt.buckets[i].nodes...)
    }
    
    // 搜索较早的buckets
    for i := bucketIndex - 1; i >= 0 && len(nodes) < k; i-- {
        nodes = append(nodes, rt.buckets[i].nodes...)
    }
    
    return nodes[:min(len(nodes), k)]
}
```

### 4.2 Rarest First算法

```go
// core/scheduler.go

type PieceScheduler struct {
    pieces      []*Piece
    queue       *PriorityQueue
}

func (ps *PieceScheduler) SelectNextPiece() *Piece {
    // 收集所有peer的bitfield
    rarityMap := make(map[int]int)
    
    for peer := range activePeers {
        for i, have := range peer.Bitfield.Bits {
            if have && !downloaded[i] {
                rarityMap[i]++
            }
        }
    }
    
    // 选择最稀有的piece
    minRarity := len(activePeers)
    var selected *Piece
    
    for i, rarity := range rarityMap {
        if rarity < minRarity {
            minRarity = rarity
            selected = ps.pieces[i]
        }
    }
    
    return selected
}
```

### 4.3 Choking算法 (BEP 6)

```go
// core/upload.go

const (
    UnchokePeriod       = 10 * time.Second
    OptimisticTickPeriod = 30 * time.Second
    MaxUnchokeCount     = 4
)

func (um *UploadManager) UpdateChoking() {
    // 1. 根据上传速度排序interested的peers
    interestedPeers := um.getInterestedPeers()
    sort.Slice(interestedPeers, func(i, j int) bool {
        return interestedPeers[i].UploadRate > interestedPeers[j].UploadRate
    })
    
    // 2. Unchoke前N个
    for i := 0; i < len(interestedPeers) && i < MaxUnchokePeerCount; i++ {
        um.unchokedPeers[interestedPeers[i].Address] = true
        um.sendUnchoke(interestedPeers[i])
    }
    
    // 3. Choke其他的
    for i := MaxUnchokePeerCount; i < len(interestedPeers); i++ {
        delete(um.unchokedPeers, interestedPeers[i].Address)
        um.sendChoke(interestedPeers[i])
    }
    
    // 4. Optimistic unchoke（每30秒轮换一个）
    if time.Since(um.lastOptimisticTick) > OptimisticTickPeriod {
        um.rotateOptimisticUnchoke()
        um.lastOptimisticTick = time.Now()
    }
}
```

### 4.4 Endgame模式

```go
// core/scheduler.go

func (ps *PieceScheduler) EnterEndgame() {
    // 当99%已下载时启动
    if ps.getProgress() > 0.99 {
        ps.endgameMode = true
    }
}

func (ps *PieceScheduler) EndgameTick() {
    if !ps.endgameMode {
        return
    }
    
    // 对于所有未完成的pieces
    for _, piece := range ps.incompletePieces() {
        // 向所有拥有该piece的peer发送request
        for peer := range peersMissingThis(piece) {
            if !hasActivePendingRequest(peer, piece) {
                ps.sendRequest(peer, piece)
            }
        }
    }
}
```

---

## 5. 网络I/O设计

### 5.1 TCP连接管理

```go
// network/connection.go

type Connection struct {
    peer       *Peer
    conn       net.Conn
    reader     *bufio.Reader
    writer     *bufio.Writer
    
    sendCh     chan Message
    recvCh     chan Message
    errCh      chan error
    
    ctx        context.Context
    cancel     context.CancelFunc
}

func (c *Connection) Run() {
    go c.readLoop()
    go c.writeLoop()
    
    select {
    case err := <-c.errCh:
        c.close(err)
    case <-c.ctx.Done():
        c.close(nil)
    }
}

func (c *Connection) readLoop() {
    for {
        msg, err := c.readMessage()
        if err != nil {
            c.errCh <- err
            return
        }
        c.recvCh <- msg
    }
}

func (c *Connection) writeLoop() {
    for {
        select {
        case msg := <-c.sendCh:
            err := c.writeMessage(msg)
            if err != nil {
                c.errCh <- err
                return
            }
        case <-c.ctx.Done():
            return
        }
    }
}
```

### 5.2 连接池

```go
// network/connection_pool.go

type ConnectionPool struct {
    connections map[string]*Connection
    maxSize     int
    semaphore   chan struct{}
    mu          sync.RWMutex
}

func (cp *ConnectionPool) Dial(addr string) (*Connection, error) {
    select {
    case cp.semaphore <- struct{}{}:
        // 有可用的连接额度
    case <-time.After(5 * time.Second):
        return nil, ErrConnectionPoolFull
    }
    
    conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
    if err != nil {
        <-cp.semaphore
        return nil, err
    }
    
    c := &Connection{conn: conn}
    cp.mu.Lock()
    cp.connections[addr] = c
    cp.mu.Unlock()
    
    return c, nil
}
```

### 5.3 消息编解码

```go
// protocol/parser.go

type MessageParser struct {
    lengthBuf [4]byte
}

func (mp *MessageParser) ReadMessage(r io.Reader) (interface{}, error) {
    // 读取4字节长度
    _, err := io.ReadFull(r, mp.lengthBuf[:])
    if err != nil {
        return nil, err
    }
    
    length := binary.BigEndian.Uint32(mp.lengthBuf[:])
    
    // Keep-alive消息
    if length == 0 {
        return &KeepAliveMsg{}, nil
    }
    
    // 读取消息ID
    msgID := make([]byte, 1)
    _, err = io.ReadFull(r, msgID)
    if err != nil {
        return nil, err
    }
    
    // 读取payload
    payloadLen := length - 1
    payload := make([]byte, payloadLen)
    _, err = io.ReadFull(r, payload)
    if err != nil {
        return nil, err
    }
    
    // 根据msgID解析
    return mp.parseMessage(msgID[0], payload)
}
```

---

## 6. 配置与参数

### 6.1 可调参数

```go
// config/defaults.go

const (
    // Network
    DefaultMaxConnections       = 200
    DefaultMaxPeersRequest      = 50
    DefaultDialTimeout          = 10 * time.Second
    DefaultKeepaliveInterval    = 2 * time.Minute
    
    // Download
    DefaultMaxConcurrentRequests = 16
    DefaultRequestTimeout        = 15 * time.Second
    DefaultPieceBatchSize        = 4
    
    // Upload
    DefaultMaxUnchokeCount       = 4
    DefaultChokeInterval         = 10 * time.Second
    DefaultOptimisticTick        = 30 * time.Second
    
    // Storage
    DefaultCacheSize             = 100 * 1024 * 1024  // 100MB
    
    // Tracker
    DefaultTrackerRetryInterval  = 30 * time.Second
    DefaultTrackerMaxRetries     = 3
    DefaultTrackerMinInterval    = 60 * time.Second
    
    // DHT
    DefaultDHTBootstrapTimeout   = 30 * time.Second
    DefaultDHTNodeCacheSize      = 1000
    DefaultDHTBucketSize         = 20  // K in Kademlia
)
```

### 6.2 性能调优指南

```
高速下载 (>100MB/s):
- MaxConnections: 500
- MaxConcurrentRequests: 32
- CacheSize: 500MB
- UploadRateLimit: 10MB/s

标准配置 (10-100MB/s):
- MaxConnections: 200
- MaxConcurrentRequests: 16
- CacheSize: 100MB
- UploadRateLimit: 0 (无限)

低带宽 (<1MB/s):
- MaxConnections: 50
- MaxConcurrentRequests: 8
- CacheSize: 50MB
- UploadRateLimit: 256KB/s
```

---

## 7. 错误处理框架

```go
// core/errors.go

type ErrorCode int

const (
    ErrOK ErrorCode = iota
    ErrTrackerUnavailable
    ErrInvalidTorrent
    ErrHashMismatch
    ErrDiskFull
    ErrNetworkError
    ErrPeerUnresponsive
    ErrDHTBootstrapFailed
)

type ErrorWithContext struct {
    Code      ErrorCode
    Message   string
    Component string
    Timestamp time.Time
    Err       error
}

func (e *ErrorWithContext) Error() string {
    return fmt.Sprintf("[%s] %s: %v", e.Component, e.Message, e.Err)
}

func (e *ErrorWithContext) IsRetryable() bool {
    switch e.Code {
    case ErrNetworkError, ErrPeerUnresponsive:
        return true
    default:
        return false
    }
}
```

---

## 8. 测试策略详解

### 8.1 单元测试示例

```go
// pkg/bencode/bencode_test.go

func TestBencodeInteger(t *testing.T) {
    tests := []struct {
        input    int64
        expected string
    }{
        {42, "i42e"},
        {-3, "i-3e"},
        {0, "i0e"},
    }
    
    for _, tt := range tests {
        result, _ := bencode.Encode(tt.input)
        if result != tt.expected {
            t.Errorf("Encode(%d) = %s, want %s", tt.input, result, tt.expected)
        }
    }
}

func TestBencodeDecoding(t *testing.T) {
    input := "i42e"
    result, _ := bencode.Decode(input)
    if result != int64(42) {
        t.Errorf("Decode(%s) = %v, want 42", input, result)
    }
}
```

### 8.2 集成测试示例

```go
// tests/integration/session_test.go

func TestSessionLifecycle(t *testing.T) {
    // 1. 创建Session
    sess, _ := core.NewSession("test.torrent")
    
    // 2. 启动
    sess.Start()
    
    // 3. 等待下载完成
    <-time.After(10 * time.Second)
    
    // 4. 验证进度
    if sess.Progress < 0.1 {
        t.Fatal("Progress too slow")
    }
    
    // 5. 停止
    sess.Stop()
}
```

---

## 9. 部署与运行

### 9.1 构建

```bash
# 编译
make build

# 编译特定平台
make build-linux
make build-windows
make build-darwin

# 编译并运行测试
make test

# 编译覆盖率测试
make coverage
```

### 9.2 运行

```bash
# 基本用法
./gobt start ubuntu.iso.torrent

# 指定输出目录
./gobt start ubuntu.iso.torrent --output-dir /tmp/downloads

# 使用配置文件
./gobt start ubuntu.iso.torrent --config gobt.yaml

# 显示帮助
./gobt help
```

### 9.3 Docker部署

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gobt ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/gobt /usr/local/bin/
ENTRYPOINT ["gobt"]
```

---

## 10. 开发流程

### 10.1 设置开发环境

```bash
# 克隆仓库
git clone https://github.com/yourname/gobt.git
cd gobt

# 安装依赖
go mod download
go mod tidy

# 运行linter
golangci-lint run

# 运行测试
go test ./...

# 运行benchmarks
go test -bench=. -benchmem ./...
```

### 10.2 贡献流程

1. Fork项目
2. 创建feature分支 (`git checkout -b feature/awesome-feature`)
3. 提交更改 (`git commit -m 'Add awesome feature'`)
4. 推送到分支 (`git push origin feature/awesome-feature`)
5. 创建Pull Request

---

## 11. 参考文档

- [BEP 3: The BitTorrent Protocol Specification](http://www.bittorrent.org/beps/bep_0003.html)
- [BEP 6: Fast Extension](http://www.bittorrent.org/beps/bep_0006.html)
- [BEP 10: Extension Protocol](http://www.bittorrent.org/beps/bep_0010.html)
- [BEP 11: Peer Exchange (PEX)](http://www.bittorrent.org/beps/bep_0011.html)
- [BEP 14: Local Service Discovery](http://www.bittorrent.org/beps/bep_0014.html)
- [BEP 20: Peer ID Specification](http://www.bittorrent.org/beps/bep_0020.html)
- [Kademlia: A Peer-to-peer Information System](https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf)

