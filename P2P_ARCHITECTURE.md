# P2P 网络通信整体架构

下面这张图把一个 BitTorrent / P2P 客户端拆成几个核心层次，方便理解 Week 4 之后要做的内容。

```mermaid
flowchart TB
    subgraph Input[输入层]
        A1[.torrent 文件]
        A2[magnet link]
    end

    subgraph Metadata[元数据层]
        B1[bencode 解析]
        B2[torrent 元数据]
        B3[info hash 计算]
    end

    subgraph Discovery[Peer 发现层]
        C1[Tracker]
        C2[DHT]
        C3[本地缓存的 peer 列表]
    end

    subgraph Protocol[Peer 通信层]
        D1[TCP 连接]
        D2[Handshake]
        D3[Bitfield / Interested / Unchoke]
        D4[Request / Piece / Cancel]
    end

    subgraph Scheduler[分片调度层]
        E1[Piece 选择]
        E2[优先级策略]
        E3[多 peer 并发请求]
    end

    subgraph Storage[校验与存储层]
        F1[Piece SHA1 校验]
        F2[文件映射]
        F3[磁盘写入]
        F4[续传状态]
    end

    subgraph Upload[上传与状态层]
        G1[已下载块统计]
        G2[上传统计]
        G3[Seeder / Leecher 状态]
    end

    A1 --> B1
    A2 --> B1
    B1 --> B2
    B2 --> B3
    B3 --> C1
    B3 --> C2
    C1 --> C3
    C2 --> C3
    C3 --> D1
    D1 --> D2
    D2 --> D3
    D3 --> D4
    D4 --> E1
    E1 --> E2
    E2 --> E3
    E3 --> F1
    F1 --> F2
    F2 --> F3
    F3 --> F4
    F4 --> G1
    F4 --> G2
    G1 --> G3
    G2 --> G3
```

## 这张图怎么理解

- 前两层负责“知道自己要什么”和“知道去哪里找人”
- 中间两层负责“真正跟别人说话”和“决定先下什么”
- 后两层负责“确认数据没坏”和“把文件落到磁盘”
- 最后一层负责上传状态和统计

## 和 Week 的对应关系

- Week 2 / Week 3：元数据层
- Week 4：Peer 发现层，重点是 Tracker
- Week 5：DHT 层
- Week 6：Peer Wire 协议层
- 后续：分片调度、校验、存储、上传
