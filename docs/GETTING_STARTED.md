# 🚀 第一步：现在就开始

> 这是一个**手把手指南**，让你在30分钟内完成第一步

---

## Step 1: 验证Go环境 (5分钟)

### 1.1 检查Go是否安装

```bash
# 打开终端/命令行，运行:
go version

# 预期输出类似:
# go version go1.25.5 windows/amd64

# 如果显示"command not found"，需要安装Go:
# https://golang.org/dl/
```

### 1.2 检查GOPATH

```bash
go env GOPATH

# 这会显示你的Go工作目录
# 记住这个路径，稍后会用到
```

---

## Step 2: 项目初始化 (10分钟)

### 2.1 进入项目目录

```bash
# Windows:
cd d:\CODE\github\go\gobt

# Linux/Mac:
cd ~/path/to/gobt
```

### 2.2 检查现有文件

```bash
# 你应该看到:
ls

# 输出:
# README.md
# go.mod
# docs/
# .idea/

# 如果go.mod不存在，创建它:
go mod init gobt

# 这会创建 go.mod 文件
```

### 2.3 现有go.mod

```bash
cat go.mod

# 输出应该是:
# module gobt
# go 1.25
```

---

## Step 3: 创建第一个程序 (10分钟)

### 3.1 创建cmd目录

```bash
# 创建cmd文件夹
mkdir -p cmd

# 验证:
ls cmd
```

### 3.2 创建main.go

创建文件 `cmd/main.go`，内容如下：

```go
package main

import (
	"fmt"
	"gobt/pkg/utils"
)

func main() {
	fmt.Println("=== gobt BitTorrent Client ===")
	fmt.Printf("Version: %s\n", utils.Version())
	fmt.Println("Ready to download torrents!")
}
```

### 3.3 创建utils包

创建 `pkg/utils/version.go`：

```go
package utils

import "fmt"

const (
	MajorVersion = 4
	MinorVersion = 2
	PatchVersion = 0
)

func Version() string {
	return fmt.Sprintf("%d.%d.%d",
		MajorVersion, MinorVersion, PatchVersion)
}
```

### 3.4 运行程序

```bash
# 在项目根目录运行:
go run ./cmd

# 预期输出:
# === gobt BitTorrent Client ===
# Version: 4.2.0
# Ready to download torrents!
```

**🎉 恭喜！你的第一个程序运行了！**

---

## Step 4: 写第一个测试 (5分钟)

### 4.1 创建测试文件

创建 `pkg/utils/version_test.go`：

```go
package utils

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	version := Version()
	
	// 验证版本格式
	expected := "4.2.0"
	if version != expected {
		t.Errorf("Expected %s, got %s", expected, version)
	}
}

func TestVersionFormat(t *testing.T) {
	version := Version()
	
	// 验证包含版本号
	if !strings.Contains(version, "4") {
		t.Error("Version should contain major version")
	}
}
```

### 4.2 运行测试

```bash
# 运行测试
go test ./pkg/utils -v

# 预期输出:
# === RUN   TestVersion
# --- PASS: TestVersion (0.00s)
# === RUN   TestVersionFormat
# --- PASS: TestVersionFormat (0.00s)
# PASS
# ok      gobt/pkg/utils  0.123s
```

**✅ 测试通过！**

---

## Step 5: 提交到Git (5分钟)

### 5.1 初始化Git（如果还没有）

```bash
# 检查是否已初始化
git status

# 如果未初始化:
git init
git config user.name "Your Name"
git config user.email "your.email@example.com"
```

### 5.2 提交第一个版本

```bash
# 查看改动
git status

# 添加所有文件
git add .

# 提交
git commit -m "feat: init project structure with version module"

# 查看日志
git log --oneline

# 输出应该显示:
# abc1234 feat: init project structure with version module
```

---

## 检查清单 ✅

完成以下所有项后，你已准备好开始学习：

- [ ] Go环境验证 (`go version` 有输出)
- [ ] 进入项目目录
- [ ] go.mod 存在且正确
- [ ] `cmd/main.go` 创建并能运行
- [ ] `pkg/utils/version.go` 创建
- [ ] 程序输出正确的版本号 (4.2.0)
- [ ] `pkg/utils/version_test.go` 创建
- [ ] 测试通过 (`go test ./pkg/utils -v`)
- [ ] 代码提交到Git
- [ ] 能运行 `go run ./cmd` 命令

---

## 现在的项目结构

```
gobt/
├── cmd/
│   └── main.go                    ✅ 已创建
├── pkg/
│   └── utils/
│       ├── version.go             ✅ 已创建
│       └── version_test.go        ✅ 已创建
├── docs/
│   ├── REQUIREMENTS.md            (已有)
│   ├── DESIGN.md                  (已有)
│   ├── ARCHITECTURE.md            (已有)
│   ├── ROADMAP.md                 (已有)
│   ├── SUMMARY.md                 (已有)
│   ├── QUICKREF.md                (已有)
│   ├── INDEX.md                   (已有)
│   └── archive/                   (历史文档归档)
├── go.mod                         ✅ 已有
├── README.md                      ✅ 已有
└── .git/                          ✅ 已初始化
```

---

## 下一步怎么做？

### 选项1: 跟随项目计划 (推荐)
```
继续学习第一阶段，进入Bencode编解码

阅读:
1. README.md - 当前周计划与执行路线
2. docs/QUICKREF.md - Bencode部分

实现:
1. 创建 pkg/bencode/encoder.go
2. 创建 pkg/bencode/decoder.go
3. 编写单元测试
4. 测试通过后提交

预计时间: 5-6小时
```

### 选项2: 了解更多细节
```
深入学习BitTorrent协议

阅读:
1. docs/QUICKREF.md - 协议基础
2. docs/DESIGN.md - 系统架构

理解:
1. Torrent文件格式
2. Info Hash计算
3. Peer通信流程
```

### 选项3: 快速浏览
```
看看已有的设计文档

阅读:
1. docs/SUMMARY.md - 项目总览
2. docs/INDEX.md - 文档导航
3. docs/ROADMAP.md - 开发计划
```

---

## 常见问题

### Q: 我运行 `go run ./cmd` 出错了

**A**: 
```bash
# 1. 确认你在项目根目录
pwd  # 应该显示 ...gobt

# 2. 确认 main.go 存在
ls cmd/main.go

# 3. 检查导入路径
# main.go 中的 "gobt/pkg/utils" 必须正确

# 4. 尝试重新加载依赖
go mod tidy
go run ./cmd
```

### Q: 测试无法运行

**A**:
```bash
# 1. 确认测试文件存在
ls pkg/utils/version_test.go

# 2. 确认文件名以 _test.go 结尾
# 文件名必须是 *_test.go

# 3. 测试函数必须以 Test 开头
# func TestXxx(t *testing.T)

# 4. 重新运行
go test ./pkg/utils -v
```

### Q: Git 提交失败

**A**:
```bash
# 1. 配置git用户信息
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# 2. 重新提交
git add .
git commit -m "feat: init project"
```

### Q: 如何继续下一步？

**A**:
参考 `README.md` 的周计划部分，并结合 `docs/QUICKREF.md` 开始实现。

---

## 需要帮助？

### 查看相关文档

| 问题 | 文档 |
|------|------|
| 不理解项目结构 | docs/INDEX.md |
| 想了解下一步做什么 | README.md |
| 对BitTorrent有疑问 | docs/QUICKREF.md |
| 想看代码框架 | docs/ARCHITECTURE.md |
| 想看整体计划 | docs/ROADMAP.md |

### 命令参考

```bash
# 运行程序
go run ./cmd

# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/utils -v

# 生成测试覆盖率
go test ./... -cover

# 代码格式化
go fmt ./...

# 代码分析
go vet ./...

# 安装依赖
go mod tidy

# 查看项目结构
tree

# 列出目录
ls -la
```

---

## 下一个里程碑

```
🎯 目标: 在一周内完成 Bencode 编解码

□ 学习 Bencode 格式 (1小时)
□ 实现编码器 (2-3小时)
□ 实现解码器 (2-3小时)
□ 编写全面测试 (1-2小时)
□ 代码审查和优化 (1小时)

预计时间: 7-10小时
```

---

## 记住

- 💡 **从小处开始** - 先完成一个小功能，再逐步扩展
- ✍️ **边写代码边学习** - 不要只读文档，动手实现
- 🧪 **多写测试** - 测试会帮你发现问题
- 📝 **记录学习笔记** - 记下遇到的问题和解决方案
- 🔄 **定期提交** - 每完成一个小功能就提交到git
- 🤔 **理解每一行代码** - 不要复制粘贴，要理解含义

---

**祝你学习顺利！** 🚀

如果这是你第一次开发Go项目，可能会遇到一些困惑，这都是正常的。
重要的是坚持下去，一步一步地完成每个阶段。

**立即开始第一步：** 打开终端，输入 `go version` 验证环境！

