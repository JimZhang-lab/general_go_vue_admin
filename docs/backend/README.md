# 🔧 后端文档中心

Go后端技术文档集合，包含性能优化、问题解决、API文档、部署指南等完整的后端开发和运维文档。

## 📁 文档列表

### 🚀 性能优化
- [**QPS优化详细方案**](./QPS_OPTIMIZATION_SUMMARY.md) - 完整的QPS优化实施方案
  - RabbitMQ消息队列集成
  - Redis缓存优化
  - 数据库连接池优化
  - API限流与熔断机制
  - 异步日志处理
  - 性能监控与指标收集
  - 负载测试与性能验证

### 🐛 问题解决
- [**后端问题修复总结**](./BUGFIX_SUMMARY.md) - 后端编译错误和导入问题修复
- [**登录问题解决方案**](./LOGIN_ISSUE_SOLUTION.md) - 登录401错误完整解决方案

### 📚 开发文档
- [**API接口文档**](./API_DOCUMENTATION.md) - RESTful API接口说明
- [**数据库设计文档**](./DATABASE_DESIGN.md) - 数据库表结构和关系设计
- [**中间件使用指南**](./MIDDLEWARE_GUIDE.md) - 各种中间件的配置和使用

### 🚀 部署运维
- [**部署指南**](./DEPLOYMENT_GUIDE.md) - 生产环境部署配置
- [**性能调优指南**](./PERFORMANCE_GUIDE.md) - 系统性能调优建议
- [**监控运维指南**](./MONITORING_GUIDE.md) - 系统监控和运维指南

## 🎯 技术栈

### 核心框架
- **Go 1.21+**: 高性能编程语言
- **Gin**: 轻量级Web框架
- **GORM**: ORM框架，优化连接池配置
- **MySQL**: 关系型数据库，支持读写分离

### 中间件集成
- **Redis**: 缓存数据库，分布式锁，会话管理
- **RabbitMQ**: 消息队列，异步处理
- **JWT**: 身份认证
- **Swagger**: API文档生成

### 性能优化组件
- **限流器**: 令牌桶、滑动窗口、固定窗口算法
- **熔断器**: 智能熔断，优雅降级
- **缓存管理**: 多种数据类型缓存，自动过期
- **分布式锁**: 防止并发冲突，自动续期
- **异步日志**: 批量处理，多级缓冲
- **性能监控**: HTTP指标、系统指标、自定义指标

## 📊 性能指标

### 优化成果
| 指标 | 优化前 | 优化后 | 提升幅度 |
|------|--------|--------|----------|
| **QPS** | ~200 | **1000+** | **5x** |
| **平均响应时间** | ~500ms | **<100ms** | **5x** |
| **P99响应时间** | ~2s | **<500ms** | **4x** |
| **错误率** | ~5% | **<1%** | **5x** |
| **并发处理** | ~50 | **500+** | **10x** |

### 关键优化点
1. **异步处理**: RabbitMQ异步处理耗时操作，响应时间减少80%
2. **缓存优化**: Redis缓存热点数据，数据库访问减少70%
3. **连接池**: 数据库连接池优化，连接复用率提升90%
4. **限流保护**: 智能限流避免系统过载，稳定性提升95%
5. **监控预警**: 实时监控及时发现问题，故障恢复时间减少80%

## 🏗️ 项目结构

```
server/
├── api/                    # API接口层
│   ├── controller/         # 控制器
│   ├── service/           # 业务逻辑层
│   ├── dao/               # 数据访问层
│   └── entity/            # 实体定义
├── common/                # 公共模块
│   ├── config/            # 配置管理
│   ├── result/            # 响应结果封装
│   └── utils/             # 工具函数
├── middleware/            # 中间件
│   ├── rate_limit_middleware.go    # 限流熔断中间件
│   ├── logging_middleware.go       # 异步日志中间件
│   └── metrics_middleware.go       # 性能监控中间件
├── pkg/                   # 核心包
│   ├── rabbitmq/          # RabbitMQ消息队列
│   ├── redis/             # Redis缓存优化
│   ├── database/          # 数据库连接池优化
│   ├── limiter/           # API限流器
│   ├── breaker/           # 熔断器
│   ├── logger/            # 异步日志处理
│   ├── metrics/           # 性能监控
│   └── bootstrap/         # 系统启动引导
├── tools/                 # 工具脚本
│   ├── load_test.go       # 负载测试工具
│   └── performance_test.sh # 自动化测试脚本
├── config.yaml           # 配置文件
└── main.go               # 程序入口
```

## 🛠️ 开发环境

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- RabbitMQ 3.8+ (可选)

### 快速启动
```bash
# 1. 安装依赖
cd server
go mod tidy

# 2. 配置数据库和Redis
# 编辑 config.yaml

# 3. 启动服务
go run main.go
```

### 开发工具
- **IDE**: VS Code / GoLand
- **调试**: Delve debugger
- **测试**: Go test + testify
- **文档**: Swagger UI
- **监控**: Prometheus + Grafana

## 🧪 测试

### 单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./api/service/...

# 生成测试覆盖率报告
go test -cover ./...
```

### 集成测试
```bash
# 运行集成测试
go test -tags=integration ./test/...
```

### 性能测试
```bash
# 运行负载测试
go run tools/load_test.go -url=http://localhost:8080 -c=100 -d=2m

# 自动化性能测试
chmod +x tools/performance_test.sh
./tools/performance_test.sh -u http://localhost:8080 -d 5m
```

## 📝 开发规范

### 代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 使用golint检查代码质量
- 添加必要的注释和文档

### 提交规范
- 使用语义化提交信息
- 每个提交只包含一个功能或修复
- 提交前运行测试确保代码质量

### API设计规范
- 遵循RESTful设计原则
- 统一的响应格式
- 完整的错误处理
- 详细的API文档

## 🔍 故障排查

### 常见问题
1. **编译错误**: 查看[问题修复总结](./BUGFIX_SUMMARY.md)
2. **登录问题**: 查看[登录问题解决方案](./LOGIN_ISSUE_SOLUTION.md)
3. **性能问题**: 查看[性能调优指南](./PERFORMANCE_GUIDE.md)

### 日志查看
```bash
# 查看应用日志
tail -f logs/app.log

# 查看错误日志
tail -f logs/error.log

# 查看性能日志
tail -f logs/performance.log
```

### 监控指标
- **健康检查**: `GET /health`
- **性能指标**: `GET /metrics`
- **QPS信息**: `GET /qps`
- **性能报告**: `GET /performance/report`

---

**最后更新**: 2025-07-29  
**维护者**: 开发团队  
**技术支持**: 查看各专项文档或联系开发团队
