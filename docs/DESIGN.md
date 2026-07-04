# BitTorrent 4.2.0 设计文档

## 1. 系统架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                       CLI / Web UI / API                     │
└────────────┬────────────────────────────────────┬────────────┘
             │                                    │
┌────────────▼─────────────┐          ┌──────────▼──────────────┐
│   Client Core Module     │          │   Network Layer         │
│  - Session Management    │          │  - TCP Connection Pool  │
│  - Download Manager      │          │  - Message Framing      │
│  - Upload Manager        │          │  - Peer Handler         │
│  - Piece Scheduler       │          │  - Protocol Parser      │
└────────────┬─────────────┘          └──────────┬──────────────┘
             │                                   │
┌────────────▼──────────────────────────────────▼─────────────┐
│                  Discovery & Coordination                    │
│  ┌─────────────────┐  ┌──────────────┐  ┌────────────────┐ │
│  │  Tracker Comm   │  │ DHT Network  │  │ PEX & LSD      │ │
│  │  - Announce     │  │ - Kademlia   │  │ - Peer Share   │ │
│  │  - Tracker List │  │ - Bootstrap  │  │ - mDNS         │ │
│  │  - Retry Logic  │  │ - Node Table │  │ - Local Peer   │ │
│  └─────────────────┘  └──────────────┘  └────────────────┘ │
└──────────────────────────────────────────────────────────────┘
             │
┌────────────▼─────────────────────────────────────────────────┐
│               Storage & Persistence Layer                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │ File Manager │  │ Hash Verifier│  │ Resume State     │   │
│  │ - I/O Ops    │  │ - SHA1       │  │ - Bitfield       │   │
│  │ - Buffering  │  │ - Validation │  │ - Session Info   │   │
│  └──────────────┘  └──────────────┘  └──────────────────┘   │
└───────────────────────────────────────────────────────────────┘
             │
┌────────────▼─────────────────────────────────────────────────┐
│                    Storage (Disk & Memory)                    │
│  - Downloaded Pieces    - Torrent Metadata                   │
│  - Incomplete Pieces    - Peer Information                   │
│  - Cache Layer          - Configuration                       │
└───────────────────────────────────────────────────────────────┘
```

---

## 2. 核心模块设计

### 2.1 模块划分

```
gobt/
├── cmd/                          # 命令行入口
│   ├── main.go
│   └── cli.go
├── core/                         # 核心业务逻辑
│   ├── session.go               # 下载会话管理
│   ├── download_manager.go      # 下载管理
│   ├── upload_manager.go        # 上传管理
│   ├── piece_scheduler.go       # Piece调度
│   └── state_machine.go         # 状态机
├── protocol/                     # 协议实现
│   ├── bencode/
│   │   ├── encoder.go
│   │   ├── decoder.go
│   │   └── types.go
│   ├── torrent/
│   │   ├── parser.go            # Torrent文件解析
│   │   ├── metadata.go
│   │   └── validator.go
│   ├── messages/
│   │   ├── handshake.go
│   │   ├── wire.go              # 协议消息定义
│   │   └── parser.go
│   └── extensions/
│       ├── extended.go
│       ├── pex.go
│       └── dht_extension.go
├── network/                      # 网络层
│   ├── peer_manager.go          # Peer连接管理
│   ├── connection.go            # 单个连接
│   ├── connection_pool.go       # 连接池
│   ├── message_handler.go       # 消息处理
│   └── protocol_handler.go      # 协议处理
├── discovery/                    # 节点发现
│   ├── tracker/
│   │   ├── http_tracker.go
│   │   └── udp_tracker.go
│   ├── dht/
│   │   ├── node.go              # DHT节点
│   │   ├── routing_table.go    # Kademlia路由表
│   │   ├── kademlia.go          # Kademlia算法
│   │   └── message.go           # DHT消息
│   └── local/
│       ├── pex.go
│       └── lsd.go
├── storage/                      # 存储管理
│   ├── file_manager.go          # 文件I/O
│   ├── piece_store.go           # Piece存储
│   ├── metadata_store.go        # 元数据存储
│   └── cache.go                 # 缓存层
├── hash/                         # 数据验证
│   ├── sha1_verifier.go
│   └── checksum.go
├── config/                       # 配置管理
│   ├── config.go
│   ├── defaults.go
│   └── validator.go
├── logging/                      # 日志系统
│   ├── logger.go
│   └── metrics.go
├── utils/                        # 工具库
│   ├── bitfield.go
│   ├── peer_id.go
│   ├── info_hash.go
│   └── misc.go
└── tests/                        # 测试
    ├── integration_test.go
    └── unit_tests/
```

---

## 3. 数据流设计

### 3.1 启动流程

```
用户启动
  │
  ▼
加载配置 (config/)
  │
  ▼
解析Torrent文件 (protocol/torrent/)
  │
  ▼
计算Info Hash和Peer ID (utils/)
  │
  ▼
初始化Session (core/session.go)
  │
  ├─▶ 初始化Tracker管理 (discovery/tracker/)
  │
  ├─▶ 初始化DHT (discovery/dht/)
  │
  ├─▶ 初始化文件存储 (storage/)
  │
  ├─▶ 初始化连接池 (network/connection_pool.go)
  │
  └─▶ 启动主事件循环
```

### 3.2 下载流程

```
主事件循环
  │
  ├─▶ 获取Peer列表
  │   ├─ 从Tracker获取
  │   ├─ DHT网络查询
  │   ├─ PEX交换
  │   └─ LSD发现
  │
  ├─▶ 建立Peer连接 (network/peer_manager.go)
  │   ├─ TCP连接
  │   ├─ 握手验证 (protocol/messages/handshake.go)
  │   └─ 初始化消息交换
  │
  ├─▶ Piece调度 (core/piece_scheduler.go)
  │   ├─ 分析Bitfield
  │   ├─ 选择稀有pieces
  │   └─ 生成下载请求
  │
  ├─▶ 并发下载 (core/download_manager.go)
  │   ├─ 发送request消息
  │   ├─ 接收piece数据
  │   ├─ 缓冲数据 (storage/piece_store.go)
  │   └─ 处理超时/重试
  │
  ├─▶ 数据验证 (hash/sha1_verifier.go)
  │   ├─ 计算SHA1 hash
  │   ├─ 与预期hash比对
  │   ├─ 验证成功→保存 (storage/file_manager.go)
  │   └─ 验证失败→丢弃+重新请求
  │
  ├─▶ 广播已有pieces
  │   ├─ 发送bitfield消息
  │   ├─ 发送have消息
  │   └─ 更新PEX状态
  │
  └─▶ 下载完成
      ├─ 验证所有pieces
      ├─ 继续上传 (做种)
      └─ 可选：停止或移除任务
```

### 3.3 上传流程

```
Peer发送interested消息
  │
  ▼
记录该Peer为interested (core/upload_manager.go)
  │
  ▼
Choking算法决策 (core/upload_manager.go - BEP 6)
  │
  ├─▶ 标准策略：unchoke前4个速度最快的peers
  ├─▶ 特殊处理：optimistic unchoke (探测新peers)
  └─▶ 基于上传速度排序
  │
  ▼
发送unchoke消息
  │
  ▼
Peer发送request消息
  │
  ▼
读取piece数据 (storage/file_manager.go)
  │
  ▼
发送piece消息
  │
  ▼
更新上传统计
  │
  ▼
重复Choking决策 (每10秒一次)
```

---

## 4. 关键组件详设

### 4.1 会话管理 (Session)

```go
type Session struct {
    // 基本信息
    torrentFile      string
    infoHash         [20]byte
    peerID           [20]byte
    
    // 元数据
    metadata         *TorrentMetadata
    pieces           []PieceInfo
    pieceCount       int
    totalSize        int64
    
    // 状态
    state            SessionState    // idle, downloading, seeding, stopped
    bitfield         *Bitfield
    progress         float64
    
    // 核心管理器
    downloadMgr      *DownloadManager
    uploadMgr        *UploadManager
    peerMgr          *PeerManager
    storageMgr       *StorageManager
    
    // 发现
    trackerMgr       *TrackerManager
    dhtNode          *DHTNode
    
    // 统计
    stats            SessionStats
}

type SessionStats struct {
    Downloaded       int64
    Uploaded         int64
    DownloadRate     int64         // bytes/sec
    UploadRate       int64         // bytes/sec
    ConnectedPeers   int
    TotalPeers       int
    StartTime        time.Time
    ElapsedTime      time.Duration
}
```

### 4.2 Piece调度器 (PieceScheduler)

```go
type PieceScheduler struct {
    // 配置
    strategy         ScheduleStrategy  // rarest_first, sequential, etc.
    endgameThreshold float64           // 进入endgame的完成度 (>95%)
    
    // 状态
    downloadQueue    *PriorityQueue
    activeRequests   map[int][]*Request  // piece_id -> requests
    pendingPeers     map[int][]Peer     // piece_id -> peers_having_it
    
    // 统计
    metrics          SchedulerMetrics
}

type PieceInfo struct {
    Index            int
    Size             int
    Hash             [20]byte
    Downloaded       int
    State            PieceState  // missing, downloading, done
    Rarity           int         // how many peers have it
    Peers            []Peer
    LastActivity     time.Time
}
```

**调度算法**:
- **Rarest First**: 优先下载最稀有的piece
- **Sequential**: 顺序下载（内存受限时）
- **Random**: 随机选择（快速测试）
- **Endgame**: 最后阶段冗余请求加速

### 4.3 Peer管理器 (PeerManager)

```go
type PeerManager struct {
    peers            map[string]*Peer           // addr -> peer
    connections      *ConnectionPool
    
    // 统计
    totalSeen        int
    totalConnected   int
    maxConnections   int
    
    // 配置
    dialTimeout      time.Duration
    keepaliveInterval time.Duration
}

type Peer struct {
    ID               [20]byte
    Address          string                     // ip:port
    Port             uint16
    
    // 连接状态
    conn             *net.TCPConn
    connected        bool
    handshakeDone    bool
    
    // 协议状态 (BEP 3)
    amChoking        bool
    amInterested     bool
    peerChoking      bool
    peerInterested   bool
    
    // Bitfield
    bitfield         *Bitfield
    
    // 统计
    uploaded         int64
    downloaded       int64
    uploadRate       float64
    downloadRate     float64
    lastActivity     time.Time
}
```

### 4.4 下载管理器 (DownloadManager)

```go
type DownloadManager struct {
    session          *Session
    peerMgr          *PeerManager
    scheduler        *PieceScheduler
    fileStore        *FileStore
    
    // 请求管理
    pendingRequests  map[string]*Request    // key: peer_addr
    maxRequestsPerPeer int                  // 通常16
    requestTimeout   time.Duration
    
    // 统计
    totalDownloaded  int64
    currentRate      int64
    avgRate          float64
}

type Request struct {
    PieceIndex       int
    Offset           int
    Length           int
    Peer             *Peer
    Timestamp        time.Time
    Retries          int
    MaxRetries       int
}
```

**下载策略**:
1. 为每个active的peer维护请求队列
2. Pipeline式请求：一旦收到piece数据，立即发送下一个request
3. 超时处理：15秒无响应则重新请求
4. 智能分配：根据peer速度动态调整请求数

### 4.5 上传管理器 (UploadManager)

```go
type UploadManager struct {
    session          *Session
    peerMgr          *PeerManager
    fileStore        *FileStore
    
    // 上传控制
    maxUnchoked      int                    // 通常4
    unchokePeers     []*Peer
    optimistic       *Peer
    
    // 策略
    chokeInterval    time.Duration         // 通常10秒
    optimisticTick   time.Duration         // 通常30秒
    
    // 统计
    totalUploaded    int64
    currentRate      int64
    avgRate          float64
}

// Choking算法 (BEP 6)
// 每chokeInterval:
// 1. 选择upload_rate最高的maxUnchoked个interested peers
// 2. 其余peers置为choke
// 3. 定期轮换一个optimistic unchoke (探测新peers)
```

### 4.6 存储管理 (StorageManager)

```go
type StorageManager struct {
    baseDir          string
    torrentDir       string
    metadataStore    *MetadataStore
    fileStore        *FileStore
    pieceStore       *PieceStore
    cacheLayer       *Cache
}

type FileStore struct {
    files            []*File
    totalSize        int64
    pieceSize        int
}

type File struct {
    Path             string
    Size             int64
    Offset           int64              // in torrent
    Handle           *os.File
    WriteBuffer      *bufio.Writer
}

type PieceStore struct {
    pieces           map[int]*PieceMeta
    pieceSize        int
}

type PieceMeta struct {
    Index            int
    Hash             [20]byte
    Downloaded       int
    Written          int
    Buffer           []byte
    State            PieceState
}
```

### 4.7 Tracker管理 (TrackerManager)

```go
type TrackerManager struct {
    announces        []string               // 主tracker
    announceLists    [][]string            // backup trackers
    
    currentTracker   string
    session          *Session
    
    // 定时器
    updateInterval   time.Duration         // 初始60秒，根据tracker调整
    nextUpdate       time.Time
    
    // 状态
    lastEvent        AnnounceEvent
    peers            []PeerInfo
}

type AnnounceRequest struct {
    InfoHash         [20]byte
    PeerID           [20]byte
    Port             uint16
    Uploaded         int64
    Downloaded       int64
    Left             int64
    Event            string              // started, stopped, completed
    NumWant          int                  // 请求peer数量
    IP               net.IP
}

type AnnounceResponse struct {
    FailureReason    string
    WarningMessage   string
    Interval         int                 // 秒
    MinInterval      int                 // 最小更新间隔
    TrackerID        string
    Complete         int                 // seeders
    Incomplete       int                 // leechers
    Peers            []PeerInfo
}
```

### 4.8 DHT网络 (DHT)

```go
type DHTNode struct {
    nodeID           [20]byte
    routingTable     *RoutingTable
    
    // 网络
    conn             *net.UDPConn
    address          string
    
    // 缓存
    peerCache        map[[20]byte][]PeerInfo  // info_hash -> peers
    nodeCache        map[[20]byte][]NodeInfo  // info_hash -> nodes
    
    // 状态
    bootstrapNodes   []NodeInfo
    isBooted         bool
}

type RoutingTable struct {
    buckets          [160]*Bucket           // k-buckets (160 for 160-bit IDs)
    k                int                    // bucket capacity (20)
}

type Bucket struct {
    nodes            []*NodeInfo
    lastChanged      time.Time
}

type NodeInfo struct {
    ID               [20]byte
    Address          string                // ip:port
    Port             uint16
    LastSeen         time.Time
}
```

**DHT操作**:
- **ping**: 检查节点活性
- **find_node**: 根据ID查找最近的节点
- **get_peers**: 查询拥有info_hash的peers
- **announce_peer**: 声称拥有某个info_hash

---

## 5. 协议消息格式

### 5.1 握手 (Handshake)

```
长度：68字节

┌─────────────────────────────────────────────────────────────┐
│ pstrlen(1) │ pstr(19) │ reserved(8) │ info_hash(20) │ peer_id(20) │
│      1     │   "BitTorrent protocol"  │   8bytes    │     20      │     20      │
└─────────────────────────────────────────────────────────────┘

例:
\x13BitTorrent protocol\x00\x00\x00\x00\x00\x00\x00\x00
<info_hash_20bytes><peer_id_20bytes>
```

### 5.2 消息格式

```
┌──────────────────┬──────────────┐
│  message length  │   message    │
│   (4 bytes)      │   (variable) │
│   big-endian     │   msg_id(1)  │
└──────────────────┴──────────────┘
                   │ + payload
```

**消息ID**:
- 0: choke
- 1: unchoke
- 2: interested
- 3: not interested
- 4: have
- 5: bitfield
- 6: request
- 7: piece
- 8: cancel
- 9: port

---

## 6. 数据结构设计

### 6.1 Bitfield

```go
type Bitfield struct {
    data        []byte
    bitCount    int
}

// 表示哪些pieces已下载
// 每bit对应一个piece
// MSB优先 (most significant bit first)
```

### 6.2 Info Hash

```go
// 计算方式
info_dict := torrent["info"]
info_bencoded := bencode.Encode(info_dict)
info_hash := sha1.Sum(info_bencoded)  // 20字节
```

### 6.3 Peer ID

```
格式: -GO4200-xxxxxxxxxx
      2 + 5 + 12 = 20 bytes

GO      = Go language identifier
4200    = Version 4.2.0 (major.minor)
xxxxxxxx = Random bytes
```

---

## 7. 并发模型

### 7.1 Goroutine设计

```
main goroutine
├─ Session Loop (core/session.go)
│  ├─ Tracker Update Loop (discovery/tracker/manager.go)
│  ├─ DHT Loop (discovery/dht/node.go)
│  ├─ Download Manager Loop (core/download_manager.go)
│  ├─ Upload Manager Loop (core/upload_manager.go)
│  └─ Piece Scheduler Loop (core/piece_scheduler.go)
│
├─ Peer Connection Handlers (network/peer_manager.go)
│  └─ Per-Peer Handler (network/connection.go)
│     ├─ Read Loop
│     └─ Write Loop
│
├─ Network Listen (network/protocol_handler.go)
│  └─ Incoming Connection Handler
│
└─ API/CLI Interface
   └─ Command Processor
```

### 7.2 同步机制

- **Channel**: goroutine间的消息传递
- **Mutex/RWMutex**: 保护共享状态
- **Context**: 优雅关闭
- **WaitGroup**: 同步等待

---

## 8. 错误处理策略

### 8.1 分类

```
网络错误
├─ 连接超时 → 重试，更换peer
├─ 连接拒绝 → 标记peer不可用
└─ I/O错误 → 重试或跳过

协议错误
├─ 握手失败 → 断开连接
├─ 消息格式错误 → 断开连接
└─ 状态转移违反 → 日志记录+断开

数据错误
├─ Hash验证失败 → 丢弃+重新请求
├─ Piece不完整 → 超时重新请求
└─ 文件写入失败 → 暂停+重试

系统错误
├─ 磁盘满 → 暂停下载
├─ 内存不足 → 清理缓存+降速
└─ 文件权限 → 终止并报错
```

### 8.2 重试策略

```go
const (
    MaxRetries           = 5
    InitialBackoff       = 1 * time.Second
    MaxBackoff           = 1 * time.Minute
    BackoffMultiplier    = 2.0
)

// Exponential backoff with jitter
retryAfter := InitialBackoff * time.Duration(math.Pow(BackoffMultiplier, float64(retries)))
retryAfter += time.Duration(rand.Intn(1000)) * time.Millisecond
```

---

## 9. 性能优化

### 9.1 I/O优化
- 批量读写
- 内存映射文件
- 异步I/O
- 缓冲区池

### 9.2 网络优化
- TCP_NODELAY (禁用Nagle算法)
- 连接复用
- 请求批处理
- 流量整形

### 9.3 CPU优化
- Goroutine池
- 对象池（减少GC）
- 算法优化（Kademlia树）
- SIMD优化（可选）

### 9.4 内存优化
- 对象复用
- 及时释放
- 限制缓冲区大小
- GC友好的数据结构

---

## 10. 配置系统

```yaml
# gobt.yaml

network:
  max_connections: 200
  listen_port: 6881-6889
  max_peers_request: 50
  
download:
  max_concurrent_requests: 16
  piece_batch_size: 4
  request_timeout: 15s
  
upload:
  max_unchoked: 4
  upload_rate_limit: 0  # 0 = unlimited
  choke_interval: 10s
  
storage:
  output_dir: "./downloads"
  cache_size: 100MB
  resume_support: true
  
tracker:
  retry_interval: 30s
  max_retries: 3
  
dht:
  enabled: true
  bootstrap_nodes:
    - "router.bittorrent.com:6881"
    - "dht.transmissionbt.com:6881"
  
logging:
  level: "info"  # debug, info, warn, error
  output: "stdout"
```

---

## 11. 状态转移图

```
          ┌─────────┐
          │ Created │
          └────┬────┘
               │ Init()
               ▼
          ┌─────────────┐
          │ Initialized │
          └────┬────────┘
               │ Start()
               ▼
       ┌──────────────────┐
       │  Downloading     │◄─────┐
       └──┬───────────────┘      │
          │                      │ Progress
          ├─ Pause()─────────────┐
          │                      │
          ▼                      │
       ┌──────────┐              │
       │  Paused  │──────────────┘
       └──┬───────┘
          │ Resume() or All downloaded
          ▼
       ┌────────┐
       │ Seeding│
       └────┬───┘
            │ Stop()
            ▼
        ┌────────────┐
        │ Stopped    │
        └────────────┘
```

---

## 12. 测试策略

### 12.1 单元测试
- Bencode编解码
- Hash计算
- Bitfield操作
- 消息解析

### 12.2 集成测试
- Session生命周期
- Tracker通信
- DHT操作
- Peer连接

### 12.3 系统测试
- 小文件下载 (<100MB)
- 大文件下载 (>1GB)
- 多任务并行
- 网络中断恢复

### 12.4 性能测试
- 下载速度
- 内存占用
- CPU使用率
- 连接数量

---

## 13. 实现优先级

**Phase 1 (MVP)**:
- [x] Torrent文件解析
- [x] Bencode编解码
- [x] Tracker HTTP通信
- [x] Peer连接和握手
- [x] 基础消息交换
- [x] Piece下载和验证
- [x] 单线程下载完成

**Phase 2**:
- [ ] 多peer并发下载
- [ ] DHT支持
- [ ] Upload/Choking算法
- [ ] Resume支持
- [ ] 性能优化

**Phase 3**:
- [ ] PEX和LSD
- [ ] 协议扩展
- [ ] Web UI
- [ ] 日志和统计

**Phase 4+**:
- [ ] 加密支持
- [ ] Magnet link
- [ ] 插件系统

