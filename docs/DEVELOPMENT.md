# OpenClaw 监控系统 - 阶段开发文档

## 项目完成状态

### ✅ 已完成

#### 阶段 1: 基础架构搭建
- [x] 项目结构设计
- [x] 数据库模型设计
- [x] API 接口设计
- [x] Go 项目初始化

#### 阶段 2: 后端开发
- [x] 数据库连接和模型实现 (`backend/internal/database/database.go`)
- [x] OpenClaw 命令执行器 (`backend/internal/collector/collector.go`)
- [x] 数据收集器实现
- [x] API 接口实现 (`backend/internal/api/handler.go`, `router.go`)
- [x] 定时任务调度器 (`backend/internal/scheduler/scheduler.go`)

#### 阶段 3: 前端开发
- [x] 页面布局设计 (`frontend/index.html`)
- [x] 实时数据展示 (`frontend/js/app.js`)
- [x] 图表可视化 (Chart.js)
- [x] 交互功能实现

#### 阶段 4: 配置和部署
- [x] 配置文件 (`config.yaml`)
- [x] 启动脚本 (`start.sh`, `start.bat`)
- [x] 文档编写 (`README.md`, `DEVELOPMENT.md`)

## 项目结构

```
monitor/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # 主入口
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler.go           # API 处理器
│   │   │   └── router.go            # 路由配置
│   │   ├── collector/
│   │   │   └── collector.go         # 数据收集器
│   │   ├── database/
│   │   │   └── database.go          # 数据库操作
│   │   ├── models/
│   │   │   └── models.go            # 数据模型
│   │   └── scheduler/
│   │       └── scheduler.go         # 定时任务
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── index.html                   # 主页面
│   ├── css/
│   │   └── style.css                # 样式
│   └── js/
│       └── app.js                   # 前端逻辑
├── data/                            # 数据库目录
├── docs/
│   └── DEVELOPMENT.md               # 开发文档
├── config.yaml                      # 配置文件
├── start.sh                         # Linux/Mac 启动脚本
├── start.bat                        # Windows 启动脚本
├── .gitignore
└── README.md
```

## 核心功能实现

### 1. 数据收集器 (Collector)

**文件**: `backend/internal/collector/collector.go`

**功能**:
- 执行 OpenClaw 命令 (`status`, `health`, `agents list`)
- 解析 JSON 输出
- 存储到数据库

**主要方法**:
- `CollectStatus()` - 收集状态信息
- `CollectHealth()` - 收集健康检查
- `CollectAgents()` - 收集代理信息
- `CollectAll()` - 收集所有数据

### 2. 数据库层 (Database)

**文件**: `backend/internal/database/database.go`

**功能**:
- SQLite 数据库初始化
- GORM 模型自动迁移
- CRUD 操作
- 数据清理

**主要表**:
- `status_snapshots` - 状态快照
- `health_checks` - 健康检查
- `sessions` - 会话记录
- `agents` - 代理信息
- `metrics` - 指标数据

### 3. API 接口 (API)

**文件**: `backend/internal/api/handler.go`, `router.go`

**端点**:
```
GET  /api/status/latest      - 最新状态
GET  /api/status/history     - 历史状态
POST /api/status/refresh     - 手动刷新
GET  /api/health/latest      - 最新健康检查
GET  /api/sessions/active    - 活动会话
GET  /api/agents             - 所有代理
GET  /api/system/info        - 系统信息
GET  /api/dashboard          - 仪表板数据
```

### 4. 定时调度器 (Scheduler)

**文件**: `backend/internal/scheduler/scheduler.go`

**功能**:
- 每 60 秒收集一次数据
- 每天凌晨 2 点清理旧数据
- 使用 cron 表达式

### 5. 前端界面 (Frontend)

**文件**: `frontend/index.html`, `js/app.js`, `css/style.css`

**功能**:
- 实时显示监控数据
- 24 小时趋势图表
- 手动刷新按钮
- 自动 60 秒刷新

## 使用说明

### 1. 配置

编辑 `config.yaml`:

```yaml
openclaw:
  path: "E:/openclaw/openclaw"  # OpenClaw 路径
  node_path: "node"              # Node.js 路径

server:
  port: "8080"
  host: "localhost"

database:
  path: "./data/monitor.db"

collector:
  interval: 60  # 收集间隔（秒）
  timeout: 30   # 超时时间（秒）

retention:
  days: 30  # 数据保留天数
```

### 2. 启动服务

**Windows**:
```bash
start.bat
```

**Linux/Mac**:
```bash
chmod +x start.sh
./start.sh
```

**手动启动**:
```bash
cd backend
go mod download
go run cmd/server/main.go -config ../config.yaml
```

### 3. 访问

打开浏览器访问: `http://localhost:8080`

## 数据流程

```
┌─────────────┐
│  Scheduler  │ (每60秒触发)
└──────┬──────┘
       │
       v
┌─────────────┐
│  Collector  │ 执行 OpenClaw 命令
└──────┬──────┘
       │
       v
┌─────────────┐
│   Parser    │ 解析 JSON 输出
└──────┬──────┘
       │
       v
┌─────────────┐
│  Database   │ 存储到 SQLite
└──────┬──────┘
       │
       v
┌─────────────┐
│  API Layer  │ 提供 REST API
└──────┬──────┘
       │
       v
┌─────────────┐
│  Frontend   │ 显示数据和图表
└─────────────┘
```

## OpenClaw 命令映射

| OpenClaw 命令 | 收集方法 | 存储表 | 更新频率 |
|--------------|---------|--------|---------|
|  status --json` | `CollectStatus()` | `status_snapshots` | 60秒 |
| `openclaw health --json` | `CollectHealth()` | `health_checks` | 60秒 |
| `openclaw agents list --json` | `CollectAgents()` | `agents` | 60秒 |

## 依赖项

### 后端 (Go)
- `github.com/gin-gonic/gin` - Web 框架
- `gorm.io/gorm` - ORM
- `gorm.io/driver/sqlite` - SQLite 驱动
- `github.com/robfig/cron/v3` - 定时任务
- `gopkg.in/yaml.v3` - YAML 解析
- `github.com/gin-contrib/cors` - CORS 支持

### 前端
- Chart.js 4.4.0 - 图表库

## 注意事项

1. **Node.js 版本**: OpenClaw 需要 Node.js v22.12+
2. **路径配置**: 确保 `config.yaml` 中的路径正确
3. **权限**: 确保有执行 OpenClaw 命令的权限
4. **端口**: 默认使用 8080 端口，确保未被占用
5. **数据库**: SQLite 文件存储在 `data/monitor.db`

## 故障排除

### 问题 1: 无法执行 OpenClaw 命令

**症状**: 日志显示 "command failed"

**解决方案**:
1. 检查 Node.js 版本: `node --version` (需要 >= 22.12)
2. 检查 OpenClaw 路径是否正确
3. 手动测试: `node E:/openclaw/openclaw/openclaw.mjs status --json`

### 问题 2: 数据库错误

**症状**: "failed to connect database"

**解决方案**:
1. 确保 `data` 目录存在: `mkdir data`
2. 检查文件权限
3. 删除旧数据库重试: `rm data/monitor.db`

### 问题 3: 前端无法加载数据

**症状**: 浏览器显示 "加载中..."

**解决方案**:
1. 检查后端是否运行: `curl http://localhost:8080/api/dashrd`
2. 查看浏览器控制台错误
3. 检查 CORS 配置

## 扩展建议

### 未来可以添加的功能:

1. **WebSocket 实时推送** - 替代轮询，实时推送数据更新
2. **告警系统** - 当健康状态异常时发送通知
3. **更多图表** - 添加饼图、柱状图等
4. **数据导出** - 导出 CSV/Excel 报表
5. **用户认证** - 添加登录功能
6. **多实例监控** - 同时监控多个 OpenClaw 实例
7. **性能指标** - CPU、内存使用率等
8. **日志查看** - 集成 OpenClaw 日志查看

## 开发环境

- Go 1.21+
- Node.js 22.12+ (用于运行 OpenClaw)
- 现代浏览器 (Chrome, Firefox, Edge)

## 许可证

MIT
