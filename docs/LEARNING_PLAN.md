# gop2p 学习计划 - 从零开始

> 目标：通过逐步实现BitTorrent客户端，系统学习网络编程、分布式系统和Go语言最佳实践

**总体周期**: 8-12周（每周投入15-25小时）  
**学习方式**: 先学习 → 然后实践 → 最后优化

---

## 第一阶段：基础准备 (1周)

> **学习目标**: 理解BitTorrent协议基础，熟悉Go开发环境

### Day 1-2: 环境和工具 (~2-3小时)

```
学习任务:
□ 安装Go 1.25+，验证 go version
□ 选择IDE: VS Code + Go Extension 或 GoLand
□ 创建GitHub账户并配置git
□ 理解Go的包管理 (go.mod/go.sum)
□ 学习basic Go syntax: 变量、函数、struct

关键命令:
$ go version
$ go mod init gop2p
$ go run main.go
$ go test ./...

输出:
└─ 项目已初始化，能运行 Hello World
```

### Day 3-4: BitTorrent核心概念 (~4-5小时)

```
学习任务:
□ 阅读 docs/QUICKREF.md 的第一部分 (协议基础)
□ 理解5个关键术语:
  - Torrent文件和Info Hash
  - Piece和Block
  - Peer/Seeder/Leecher
  - Tracker
  - DHT (基础概念)
□ 查看一个.torrent文件的内容
□ 理解Bencode编码格式

学习资源:
- docs/QUICKREF.md 中的"术语速览"部分
- 在线工具: Bencode解码器

输出:
└─ 能解释"为什么需要Tracker"等基本问题
```

### Day 5-7: 代码框架搭建 (~4-5小时)

```
学习任务:
□ 阅读 docs/ARCHITECTURE.md 的目录结构部分
□ 创建以下基础目录:
  pkg/bencode/
  pkg/torrent/
  pkg/protocol/
  pkg/network/
  cmd/
  tests/
  
□ 创建第一个Go模块:
  pkg/utils/version.go - 包含版本信息
  
□ 学习Go package和module最佳实践

代码示例:
// pkg/utils/version.go
package utils

const (
    MajorVersion = 4
    MinorVersion = 2
    PatchVersion = 0
)

func Version() string {
    return fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)
}

□ 学会编写基础单元测试

输出:
└─ 项目基础框架完成，能跑测试
```

---

## 第二阶段：Bencode编解码 (1-2周)

> **学习目标**: 理解序列化，实现Bencode编解码

### Week 2-1: Bencode编码器 (~5-6小时)

```
学习任务:
□ 理解Bencode格式 (docs/QUICKREF.md)
  - 整数: i42e
  - 字符串: 5:hello
  - 列表: li1e4:spame
  - 字典: d3:agei27ee

□ 实现编码器逻辑:
  pkg/bencode/encoder.go
  - encodeInteger()
  - encodeString()
  - encodeList()
  - encodeDict()

□ 学习Go的interface和error handling

代码框架:
type Encoder struct {
    buf bytes.Buffer
}

func (e *Encoder) Encode(v interface{}) ([]byte, error) {
    switch v.(type) {
    case int64:
        return e.encodeInteger(v.(int64))
    case string:
        return e.encodeString(v.(string))
    // ...
    }
}

学习点:
- 递归数据结构处理
- Go的type assertion
- error作为返回值
- buffer/bytes库使用

输出:
└─ 能编码各种数据类型
```

### Week 2-2: Bencode解码器 (~5-6小时)

```
学习任务:
□ 实现解码器逻辑:
  pkg/bencode/decoder.go
  - 字符流解析
  - 递归结构处理
  - 错误检测

□ 学习io.Reader接口使用
□ 实现错误处理和边界检查

代码框架:
type Decoder struct {
    r io.Reader
}

func (d *Decoder) Decode() (interface{}, error) {
    // 读取第一个字节判断类型
    // 调用对应的decode函数
}

学习点:
- io interface设计
- 状态机式解析
- 错误恢复
- 内存管理

输出:
└─ 能解析.torrent文件的raw data
```

### Week 2-3: 测试和优化 (~3-4小时)

```
学习任务:
□ 编写全面的单元测试:
  pkg/bencode/bencode_test.go
  - 测试边界情况
  - 测试嵌套结构
  - 性能测试 (benchmark)

□ 学习Go的testing框架
□ 使用 go test -v, -bench, -cover

□ 优化性能:
  - buffer复用
  - 少量内存分配

输出:
└─ bencode模块完成，测试覆盖>90%
```

---

## 第三阶段：Torrent文件解析 (1-2周)

> **学习目标**: 理解文件格式，实现元数据提取

### Week 3-1: Torrent元数据结构 (~4-5小时)

```
学习任务:
□ 理解.torrent文件结构:
  {
    "announce": "http://tracker.example.com/announce",
    "announce-list": [["tracker1"], ["tracker2"]],
    "info": {
      "name": "filename",
      "length": 1024,
      "piece length": 262144,
      "pieces": "<20-byte hashes>"
    }
  }

□ 定义Go数据结构:
  pkg/torrent/metadata.go

代码示例:
type Torrent struct {
    Announce     string
    AnnounceList [][]string
    CreationDate int64
    Comment      string
    Info         *Info
}

type Info struct {
    Name       string
    Length     int64
    PieceLength int
    Pieces     string // raw binary
}

学习点:
- struct tag在JSON中的应用
- 可选字段处理
- 二进制数据处理

输出:
└─ 定义好数据结构
```

### Week 3-2: Info Hash计算 (~3-4小时)

```
学习任务:
□ 理解Info Hash:
  SHA1(bencode(torrent["info"])) = 20 bytes

□ 实现hash计算:
  pkg/torrent/parser.go

代码框架:
import (
    "crypto/sha1"
)

func CalculateInfoHash(torrent *Torrent) [20]byte {
    // 1. bencode torrent.Info
    infoEncoded := bencode.Encode(torrent.Info)
    
    // 2. SHA1 hash
    hash := sha1.Sum(infoEncoded)
    
    return hash
}

□ 验证计算结果

学习点:
- crypto/sha1包使用
- 数组vs切片
- bencode integration

输出:
└─ 能计算Info Hash
```

### Week 3-3: 完整Torrent解析 (~4-5小时)

```
学习任务:
□ 实现完整的.torrent文件解析:
  pkg/torrent/parser.go

代码框架:
func Parse(filename string) (*Torrent, error) {
    // 1. 读取文件
    data, err := ioutil.ReadFile(filename)
    
    // 2. Bencode解码
    value, err := bencode.Decode(data)
    
    // 3. 转换为Torrent struct
    torrent := convertToTorrent(value)
    
    // 4. 验证
    err = validate(torrent)
    
    return torrent, nil
}

□ 实现验证逻辑:
  - 必要字段检查
  - 大小合理性检查
  - Info Hash计算

□ 完整的错误处理

学习点:
- 文件I/O
- 类型转换
- 验证逻辑
- 错误链

输出:
└─ 能解析任何.torrent文件
```

### Week 3-4: 测试 (~3-4小时)

```
学习任务:
□ 准备测试.torrent文件:
  tests/fixtures/sample.torrent
  
□ 编写测试:
  pkg/torrent/parser_test.go
  - 解析测试
  - info hash验证
  - 错误处理

□ 性能测试

输出:
└─ Torrent解析模块完成
```

---

## 第四阶段：网络基础 (1-2周)

> **学习目标**: 理解网络编程，实现Peer连接

### Week 4-1: TCP连接基础 (~5-6小时)

```
学习任务:
□ 学习Go的net库:
  - net.Dial 连接
  - net.Listen 监听
  - bufio 读写

□ 实现简单的TCP连接:
  pkg/network/connection.go

代码框架:
type Connection struct {
    conn   net.Conn
    reader *bufio.Reader
    writer *bufio.Writer
    
    sendCh chan []byte
    recvCh chan []byte
    errCh  chan error
}

func Dial(addr string) (*Connection, error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, err
    }
    
    return &Connection{
        conn:   conn,
        reader: bufio.NewReader(conn),
        writer: bufio.NewWriter(conn),
    }, nil
}

□ 实现close逻辑

学习点:
- TCP连接生命周期
- bufio buffer优化
- 并发安全性

输出:
└─ 能建立TCP连接
```

### Week 4-2: BitTorrent握手 (~5-6小时)

```
学习任务:
□ 理解握手协议 (docs/QUICKREF.md):
  [1] 协议长度: 0x13
  [19] "BitTorrent protocol"
  [8] reserved flags
  [20] info_hash
  [20] peer_id

总长: 68字节

□ 生成Peer ID:
  pkg/utils/peer_id.go

代码框架:
func GeneratePeerID() [20]byte {
    var id [20]byte
    // "-GO4200-" + 12 random bytes
    copy(id[:], "-GO4200-")
    rand.Read(id[8:])
    return id
}

□ 实现握手:
  pkg/protocol/handshake.go

type Handshake struct {
    InfoHash [20]byte
    PeerID   [20]byte
    Reserved [8]byte
}

func (h *Handshake) Marshal() []byte {
    // 编码握手信息
}

func (h *Handshake) Unmarshal(data []byte) error {
    // 解码握手信息
}

学习点:
- 二进制编码/解码
- 大端序处理
- 固定长度数组

输出:
└─ 能发送和接收握手
```

### Week 4-3: 基础消息交换 (~4-5小时)

```
学习任务:
□ 理解消息格式:
  [4] 消息长度 (大端序)
  [1] 消息ID
  [N] 有效负载

□ 实现消息编解码:
  pkg/protocol/messages.go

代码框架:
type Message struct {
    ID      uint8
    Payload []byte
}

func (m *Message) Marshal() []byte {
    // 编码消息
}

func (m *Message) Unmarshal(data []byte) error {
    // 解码消息
}

□ 实现4个基础消息:
  - Bitfield (0x05)
  - Have (0x04)
  - Interested (0x02)
  - NotInterested (0x03)

学习点:
- 变长数据处理
- 消息设计模式
- 二进制协议

输出:
└─ 能交换基础消息
```

### Week 4-4: 完整连接管理 (~3-4小时)

```
学习任务:
□ 整合握手和消息:
  pkg/network/peer_manager.go

□ 实现连接生命周期:
  - 连接
  - 握手
  - 消息循环
  - 断开

□ 错误处理和超时

输出:
└─ 能建立完整的peer连接
```

---

## 第五阶段：Tracker通信 (1周)

> **学习目标**: 理解tracker协议，获取peer列表

### Week 5-1: HTTP请求构建 (~4-5小时)

```
学习任务:
□ 学习net/http库:
  - url.Values构建query
  - http.Get请求
  - 响应处理

□ 理解announce请求:
  GET /announce?
    info_hash=...&
    peer_id=...&
    port=...&
    uploaded=...&
    downloaded=...&
    left=...&
    event=started

□ 实现请求构建:
  pkg/tracker/http.go

代码框架:
type AnnounceRequest struct {
    InfoHash   [20]byte
    PeerID     [20]byte
    Port       uint16
    Uploaded   int64
    Downloaded int64
    Left       int64
    Event      string
}

func (r *AnnounceRequest) BuildURL(tracker string) string {
    // 构建URL
}

学习点:
- URL编码
- Query参数
- HTTP客户端

输出:
└─ 能构建announce请求
```

### Week 5-2: 响应解析 (~4-5小时)

```
学习任务:
□ 理解tracker响应格式:
  Bencode字典:
  {
    "interval": 1800,
    "peers": [
      {"ip": "...", "port": ...},
      ...
    ]
  }

□ 实现响应解析:
  pkg/tracker/http.go

代码框架:
type AnnounceResponse struct {
    Interval int
    Peers    []PeerInfo
}

type PeerInfo struct {
    IP   string
    Port uint16
}

func ParseAnnounceResponse(body []byte) (*AnnounceResponse, error) {
    // 1. Bencode解码
    // 2. 提取字段
    // 3. 解析peer列表
}

□ 错误处理:
  - Tracker返回failure
  - 无效响应
  - 连接超时

学习点:
- JSON/Bencode响应处理
- 错误处理
- 数据验证

输出:
└─ 能获取peer列表
```

### Week 5-3: 完整Tracker客户端 (~3-4小时)

```
学习任务:
□ 整合请求和响应:
  pkg/tracker/manager.go

代码框架:
type TrackerManager struct {
    tracker string
    session *Session
}

func (tm *TrackerManager) Announce(event string) ([]PeerInfo, error) {
    req := AnnounceRequest{...}
    resp, err := tm.sendRequest(req)
    return resp.Peers, nil
}

□ 重试逻辑
□ 定期更新

输出:
└─ 能与tracker通信
```

---

## 第六阶段：下载引擎基础 (2-3周)

> **学习目标**: 理解下载流程，实现基本的piece下载

### Week 6-1: Bitfield管理 (~3-4小时)

```
学习任务:
□ 理解Bitfield:
  - 每个bit代表一个piece
  - 1 = 已有, 0 = 缺少
  - MSB优先 (most significant bit first)

□ 实现Bitfield:
  pkg/core/bitfield.go

代码框架:
type Bitfield struct {
    data     []byte
    bitCount int
}

func (b *Bitfield) Has(index int) bool {
    byteIndex := index / 8
    bitIndex := 7 - (index % 8)
    return (b.data[byteIndex] >> uint(bitIndex)) & 1 == 1
}

func (b *Bitfield) Set(index int) {
    byteIndex := index / 8
    bitIndex := 7 - (index % 8)
    b.data[byteIndex] |= 1 << uint(bitIndex)
}

□ 完整的操作接口

学习点:
- 位操作
- 位序问题
- 性能优化

输出:
└─ 能管理piece状态
```

### Week 6-2: 请求和接收 (~6-7小时)

```
学习任务:
□ 理解Request/Piece消息:
  Request: [4]piece_index [4]offset [4]length
  Piece: [4]piece_index [4]offset [N]data

□ 实现request发送:
  pkg/protocol/messages.go (扩展)

□ 实现piece接收和缓冲:
  pkg/core/download_manager.go

代码框架:
type DownloadManager struct {
    pieces        []*Piece
    activeRequests map[string]chan []byte
}

type Piece struct {
    Index      int
    Size       int
    Hash       [20]byte
    Data       []byte
    Downloaded int
}

func (dm *DownloadManager) RequestPiece(
    peer *Peer,
    pieceIndex int,
    offset int,
    length int,
) error {
    // 构建request消息
    // 发送给peer
    // 等待响应
}

□ 超时处理

学习点:
- 异步操作
- 缓冲区管理
- 协议消息

输出:
└─ 能请求和接收数据
```

### Week 6-3: 数据验证 (~4-5小时)

```
学习任务:
□ 实现SHA1验证:
  pkg/hash/verifier.go

代码框架:
func VerifyPiece(data []byte, expectedHash [20]byte) bool {
    actualHash := sha1.Sum(data)
    return actualHash == expectedHash
}

□ 集成到下载流程:
  - 接收完整piece
  - 计算SHA1
  - 对比hash
  - 验证失败重新请求

学习点:
- crypto库使用
- 错误恢复
- 流程控制

输出:
└─ 能验证数据完整性
```

### Week 6-4: 文件写入 (~4-5小时)

```
学习任务:
□ 实现文件写入:
  pkg/storage/file_manager.go

代码框架:
type FileManager struct {
    files []*File
}

type File struct {
    path   string
    handle *os.File
    size   int64
}

func (fm *FileManager) WritePiece(
    pieceIndex int,
    data []byte,
) error {
    // 计算piece在文件中的位置
    // 写入数据
    // 更新进度
}

□ 处理多文件torrent
□ 错误处理

学习点:
- 文件I/O
- 位置计算
- 并发安全

输出:
└─ 能将下载数据保存到磁盘
```

### Week 6-5: 主下载循环 (~5-6小时)

```
学习任务:
□ 整合所有部分:
  pkg/core/session.go

代码框架:
type Session struct {
    peers          []*Peer
    downloadMgr    *DownloadManager
    fileManager    *FileManager
    progress       float64
}

func (s *Session) Start() error {
    // 1. 连接tracker，获取peers
    // 2. 连接多个peers
    // 3. 发送握手
    // 4. 主循环:
    //    - 选择piece
    //    - 发送request
    //    - 接收piece
    //    - 验证
    //    - 保存
    //    - 更新进度
}

□ 进度跟踪
□ 状态管理

学习点:
- 系统集成
- 状态机
- 事件循环

输出:
└─ 能完成基本的文件下载
```

---

## 第七阶段：优化和完善 (2-3周)

> **学习目标**: 性能优化，错误处理，完整功能

### Week 7-1: 并发优化 (~5-6小时)

```
学习任务:
□ 学习Go并发:
  - Goroutines
  - Channels
  - sync包 (Mutex, WaitGroup)

□ 优化连接管理:
  - 并发连接多个peers
  - 连接池

□ 优化消息处理:
  - 每个peer一个goroutine处理消息
  - channel通信

学习点:
- Goroutine生命周期
- Channel设计
- 竞态条件检测

输出:
└─ 能支持100+并发连接
```

### Week 7-2: 错误处理和恢复 (~4-5小时)

```
学习任务:
□ 实现错误处理:
  pkg/core/errors.go

□ 实现重试逻辑:
  - 连接失败重试
  - 数据验证失败重新请求
  - Tracker超时更换tracker

□ Graceful shutdown

学习点:
- 错误定义
- 错误链
- 恢复策略

输出:
└─ 系统更稳定
```

### Week 7-3: 配置系统 (~3-4小时)

```
学习任务:
□ 实现可配置参数:
  pkg/config/config.go

□ 支持命令行标志:
  - 输出目录
  - 最大连接数
  - 带宽限制
  
□ 配置文件支持 (YAML)

学习点:
- flag包使用
- YAML解析
- 配置管理

输出:
└─ 能配置运行参数
```

### Week 7-4: CLI和测试 (~5-6小时)

```
学习任务:
□ 实现命令行工具:
  cmd/main.go

代码框架:
func main() {
    app := cli.NewApp()
    app.Commands = []cli.Command{
        {
            Name: "start",
            Action: startDownload,
        },
        {
            Name: "info",
            Action: showTorrentInfo,
        },
    }
    app.Run(os.Args)
}

□ 显示下载进度
□ 实时统计

□ 编写全面的测试:
  - 单元测试 (>80%覆盖率)
  - 集成测试
  - 端到端测试

输出:
└─ 可用的CLI工具
```

---

## 第八阶段：高级功能 (可选，2-3周)

> 如果前面阶段都完成，可继续学习

### Week 8-1: 上传和做种

```
学习任务:
□ 实现Choking算法 (BEP 6)
□ 实现上传管理
□ 做种功能

学习点:
- 激励机制设计
- 优化策略
```

### Week 8-2: DHT网络

```
学习任务:
□ 学习DHT和Kademlia算法
□ 实现DHT节点
□ Peer发现

学习点:
- 分布式系统
- 路由算法
- 网络拓扑
```

---

## 学习路线总结

```
第一阶段: 基础 (1周)
  目标: 理解协议，搭建框架
  
  第二阶段: Bencode (1-2周)
    目标: 序列化和反序列化
    ↓
    第三阶段: Torrent解析 (1-2周)
      目标: 读取元数据
      ↓
      第四阶段: 网络基础 (1-2周)
        目标: TCP通信和握手
        ↓
        第五阶段: Tracker (1周)
          目标: 获取peer列表
          ↓
          第六阶段: 下载引擎 (2-3周)
            目标: 完成第一次下载
            ↓
            第七阶段: 优化 (2-3周)
              目标: 稳定可用的客户端
              ↓
              (可选) 第八阶段: 高级功能 (2-3周)
```

---

## 每个阶段的输出物

| 阶段 | 输出 | 代码行数 |
|------|------|--------|
| 1 | 项目框架 | ~100 |
| 2 | Bencode编解码 | ~300 |
| 3 | Torrent解析 | ~200 |
| 4 | 网络基础 | ~400 |
| 5 | Tracker通信 | ~300 |
| 6 | 下载引擎 | ~1000 |
| 7 | 优化完善 | ~500 |
| 8 | 高级功能 | ~800 |
| | **总计** | **~3600** |

---

## 学习建议

### ✅ 推荐做法
```
□ 每个阶段完成后，编写单元测试
□ 定期运行测试，确保功能正常
□ 阅读参考文档中的相应部分
□ 参考存在的开源项目理解细节
□ 记录遇到的问题和解决方案
□ 定期提交代码到git
□ 每周回顾学到的新概念
```

### ❌ 要避免
```
□ 一次性写太多代码而不测试
□ 跳过某些阶段直接进入高级功能
□ 忽视错误处理
□ 不写测试就认为完成
□ 复制粘贴代码而不理解
```

---

## 资源和参考

### 关键文档 (按阅读顺序)
1. **docs/QUICKREF.md** - 快速参考，学习各个阶段
2. **docs/DESIGN.md** - 架构细节，深入理解
3. **docs/ARCHITECTURE.md** - 代码框架，实现参考

### 外部资源
- [BitTorrent BEP 3](http://www.bittorrent.org/beps/bep_0003.html) - 协议规范
- [Go官方文档](https://golang.org/doc/) - Go语言学习
- [Effective Go](https://golang.org/doc/effective_go) - Go最佳实践
- [Transmission源码](https://github.com/transmission/transmission) - 参考实现

---

## 时间估算

```
每周投入时间: 15-25小时

阶段1: 1周   × 20小时 = 20小时
阶段2: 1.5周 × 20小时 = 30小时
阶段3: 1.5周 × 20小时 = 30小时
阶段4: 1.5周 × 22小时 = 33小时
阶段5: 1周   × 20小时 = 20小时
阶段6: 2.5周 × 24小时 = 60小时
阶段7: 2.5周 × 24小时 = 60小时
─────────────────────────────────
总计: 12周   × ~21小时 = ~253小时

换算: 
- 每周25小时，约10周完成
- 每周15小时，约17周完成
- 全职投入，约1.5个月完成
```

---

## 学习成果

完成这个计划后，你将学到：

### Go语言
- ✅ 包管理和模块化设计
- ✅ 并发编程 (Goroutines/Channels)
- ✅ 接口和类型系统
- ✅ 错误处理最佳实践
- ✅ 测试驱动开发 (TDD)

### 网络编程
- ✅ TCP/IP编程
- ✅ 二进制协议设计
- ✅ HTTP客户端开发
- ✅ 并发连接管理
- ✅ 错误恢复和超时处理

### 系统设计
- ✅ 模块化架构
- ✅ 状态管理
- ✅ 事件驱动编程
- ✅ 性能优化
- ✅ 可靠性设计

### BitTorrent协议
- ✅ 协议深度理解
- ✅ P2P架构
- ✅ 分布式系统基础
- ✅ 激励机制设计

---

## 建议的第一步

```
今天开始:

1. 阅读本文档的"第一阶段" (30分钟)
2. 阅读 docs/QUICKREF.md 的基础部分 (30分钟)
3. 按照Day 1-2的任务安装Go环境 (1-2小时)
4. 创建项目框架 (1小时)
5. 运行第一个单元测试 (30分钟)

→ 第一天：认识项目，搭建环境
```

---

**开始日期**: 2026-07-01  
**预计完成**: 2026-08 至 2026-09  
**难度**: ⭐⭐⭐ (中等难度)  
**收获**: 🎓 系统的Go学习 + 🚀 可用的BitTorrent客户端

**准备好开始了吗？** 👍
