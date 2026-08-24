# Advanced Flight Server v2

Advanced Flight Server v2 是一个使用 Go 编写的高性能飞行模拟网络 TCP 服务器。项目基于 `gnet` 的事件驱动网络模型，负责处理飞行员、管制员与 ATIS 客户端的连接、认证、状态同步和消息转发。

> 项目仍在持续开发中，协议和配置项可能发生变化。用于生产环境前，请自行完成安全审计、压力测试和数据备份。

## Features

- 基于 `gnet v2` 的多核、事件驱动 TCP 服务
- 飞行员、ATC 和 ATIS 会话管理
- 用户 CID、密码、账号状态与等级校验
- 飞行计划、位置、文本消息、移交、标签修改等协议包处理
- METAR 数据定时拉取与客户端查询响应
- 在线飞行员、ATC、ATIS 快照每 3 秒发布到 Redis
- IPv4/IPv6 与 CIDR IP 封禁，支持 `reject` 和 `silent` 两种处理方式
- 登录超时、空闲连接清理、粘包拆包与缓冲区保护
- 单连接 panic 恢复，避免异常数据包导致整个服务退出
- MySQL 和 PostgreSQL 数据库支持
- Zap 结构化日志、日志轮转和慢查询记录
- Windows/Linux、AMD64/ARM64 跨平台构建

## 环境要求

- Go 1.24 或更高版本
- MySQL 或 PostgreSQL
- Redis
- 一个与 `pkg/entity/user.go` 中 `User` 模型兼容的 `user` 表

项目目前不会自动创建或迁移数据库表，请提前准备账户数据库结构。

## 快速开始

### 1. 获取代码

```bash
git clone https://github.com/AFcPPe/advanced-flight-server-v2.git
cd advanced-flight-server-v2
```

### 2. 下载依赖

```bash
go mod download
```

### 3. 生成并修改配置

首次运行时，如果当前目录不存在 `config.yaml`，程序会自动生成默认配置。你也可以先运行一次：

```bash
go run .
```

随后编辑生成的 `config.yaml`，至少配置数据库和 Redis：

```yaml
app:
  name: advanced-flight-server
  version: 2.0.0
  env: prod

server:
  host: 0.0.0.0
  port: 6809
  motd: Welcome to Advanced Flight Server!
  auth_timeout: 5

database_accounts:
  account:
    driver: mysql # mysql 或 postgres
    host: 127.0.0.1
    port: 3306
    username: flight_server
    password: change-me
    database: flight_server
    charset: utf8mb4
    max_idle_conns: 10
    max_open_conns: 100
    conn_max_lifetime: 3600
    conn_max_idle_time: 600
    log_level: warn
    slow_threshold: 200

redis:
  addr: 127.0.0.1:6379
  password: ""
  db: 0

metar:
  url: http://metar.vatsim.net/metar.php?id=ALL
  interval: 10

ip_ban:
  enabled: true
  file: ip_ban.json
  interval: 1

logger:
  level: info
  filename: logs/app.log
  max_size: 100
  max_backups: 3
  max_age: 7
  compress: true
  console: true
  rotate_by_date: true
```

`config.yaml` 已被 Git 忽略，请不要把数据库密码、Redis 密码等凭据提交到仓库。

### 4. 启动服务

```bash
go run .
```

默认监听 `0.0.0.0:6809`。也可以先编译再运行：

```bash
go build -o advanced-flight-server .
./advanced-flight-server
```

Windows PowerShell：

```powershell
go build -o advanced-flight-server.exe .
.\advanced-flight-server.exe
```

## IP 封禁

启用 IP 封禁后，若规则文件不存在，程序会自动生成 `ip_ban.json` 示例文件。规则支持单个 IP 或 CIDR：

```json
{
  "rules": [
    {
      "cidr": "192.0.2.0/24",
      "action": "reject",
      "note": "连接后立即断开"
    },
    {
      "cidr": "198.51.100.7",
      "action": "silent",
      "note": "接受连接但丢弃输入且不响应"
    }
  ]
}
```

规则会按配置的分钟间隔自动重新加载，按文件顺序匹配，第一条命中的规则生效。

## Redis 快照

服务每 3 秒将在线用户快照写入 Redis：

- Key：`afs:snapshot:users`
- TTL：10 秒
- 内容：飞行员、ATC、ATIS 状态及时间戳的 JSON 数据

## 测试

```bash
go test ./...
```

如需检查并发数据竞争：

```bash
go test -race ./...
```

## 发布构建

项目使用 [GoReleaser](https://goreleaser.com/) 生成 Windows/Linux 的 AMD64 和 ARM64 构建产物。

安装 GoReleaser：

```bash
go install github.com/goreleaser/goreleaser/v2@latest
```

在 Windows PowerShell 中执行快照构建：

```powershell
.\build.ps1
```

正式发布需要先创建 Git tag，然后执行：

```powershell
.\build.ps1 -Release
```

构建产物会输出到 `dist/`。

## 安全说明

当前账户认证逻辑会直接比较客户端提交的密码与数据库中的密码字段。这意味着现有实现可能依赖明文密码，不建议在未经改造和审计的公网生产环境中直接使用。建议在部署前改用安全的密码哈希验证方案，并通过防火墙或私有网络限制数据库与 Redis 的访问。

如果发现安全问题，请避免在公开 Issue 中披露可直接利用的细节，优先联系项目维护者。

## 贡献

欢迎提交 Issue 和 Pull Request。提交代码前请确保：

```bash
gofmt -w .
go test ./...
```

## 开源协议

本项目基于 [MIT License](LICENSE) 开源。你可以自由使用、复制、修改、合并、发布和分发本项目，但需要保留原始版权声明和许可声明。软件按“原样”提供，不附带任何形式的担保。

Copyright (c) 2026 AFcPPe
