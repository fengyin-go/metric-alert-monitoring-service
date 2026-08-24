# 监控告警（monitoring）

纯 Go 标准库实现的监控告警后端服务。零第三方依赖，开箱即跑。

## 功能特性

- **指标管理**：上报 gauge/counter 类指标，同名指标自动更新值。
- **告警规则**：配置指标阈值规则（gt/gte/lt/lte/eq），支持 info/warning/critical 严重程度。
- **规则评估**：一键评估全部启用规则，命中阈值即触发告警事件。
- **静默期**：为指标设置静默时间窗，静默期间不触发告警。
- **告警事件**：事件状态机（firing → resolved），支持统计。

## 技术栈

- Go 1.22+，仅使用标准库 `net/http`
- 标准工程分层：`cmd` / `internal`（app/config/model/store/service/handler）/ `pkg`
- 内存存储，线程安全（`sync.RWMutex`）

## 运行

```bash
# 在 origin/ 目录下
go run ./cmd/server

# 或指定端口
PORT=9090 go run ./cmd/server
```

服务默认监听 `:8080`。环境变量：`PORT`、`ADDR`、`MAX_PAGE_SIZE`、`EVAL_INTERVAL`、`LOG_LEVEL`。

## API 一览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/metrics/upsert` | 上报/更新指标 |
| GET | `/api/metrics` | 分页查询指标（`type`/`keyword`） |
| GET | `/api/metrics/stats` | 指标统计 |
| GET | `/api/metrics/rules` | 查询指标关联的规则（`name`） |
| GET | `/api/metrics/{id}` | 获取指标详情 |
| DELETE | `/api/metrics/{id}` | 删除指标 |
| POST | `/api/rules` | 创建告警规则 |
| GET | `/api/rules` | 分页查询规则（`severity`/`keyword`） |
| GET | `/api/rules/stats` | 规则统计 |
| GET | `/api/rules/{id}` | 获取规则详情 |
| PUT | `/api/rules/{id}` | 更新规则 |
| DELETE | `/api/rules/{id}` | 删除规则 |
| POST | `/api/events/evaluate` | 评估规则触发告警 |
| GET | `/api/events` | 分页查询事件（`metric_name`/`severity`/`status`） |
| GET | `/api/events/stats` | 告警事件统计 |
| GET | `/api/events/active` | 当前 firing 事件 |
| GET | `/api/events/count` | 按指标统计事件数（`metric_name`） |
| GET | `/api/events/rule` | 查询规则触发的事件（`rule_id`） |
| GET | `/api/events/{id}` | 获取事件详情 |
| POST | `/api/events/{id}/resolve` | 解决告警事件 |
| POST | `/api/silences` | 创建静默期 |
| GET | `/api/silences` | 分页查询静默期（`metric_name`） |
| GET | `/api/silences/active` | 当前生效的静默期 |
| GET | `/api/silences/stats` | 静默期统计 |
| DELETE | `/api/silences/metric` | 按指标删除静默期（`metric_name`） |
| DELETE | `/api/silences/{id}` | 删除静默期 |

统一响应结构：`{"code":0,"message":"ok","data":...}`。

## 测试

```bash
go test ./...
```
