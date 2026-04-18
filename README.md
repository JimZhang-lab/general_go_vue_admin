# General Go Vue Admin

一个面向生产的后台管理系统，基于 **Go + Vue3 + TypeScript + TailwindCSS**。  
目标不是“能跑”，而是长期可维护、可演进、可观测、可持续优化。

## 1. 项目定位

- 面向中后台场景的通用管理底座（用户、角色、菜单、组织、日志）
- 强调安全性（鉴权、验证码、登录防暴力破解、统一错误处理）
- 强调稳定性（中间件治理、连接池配置、构建与自检闭环）
- 强调体验一致（统一表单反馈、浅色动态登录/注册页、组件复用）

## 2. 核心能力

### 2.1 业务能力

- 认证模块：登录、注册、退出、验证码
- 账号体系：管理员管理、个人资料、密码修改、头像上传
- 权限体系：角色管理、菜单管理、权限分配
- 组织体系：部门管理、岗位管理
- 审计体系：登录日志、操作日志

### 2.2 安全能力

- JWT 鉴权与路由保护
- 验证码校验（已修复为一次性消费，防止重复复用）
- 登录失败防暴力破解：
  - 连续失败达到阈值自动锁定
  - 锁定时间到期自动恢复
  - 阈值与时长可配置（`security.loginFailedAttemptLimit` / `security.loginLockMinutes`）
- 统一错误模型与 TraceID 支撑

### 2.3 稳定性与工程能力

- 后端自检：`go test ./...` + `go build ./...` + `go vet ./...`
- 前端自检：`npm run build-only`
- 一键全链路自检脚本：`server/tools/self_check.sh`
- 前端分包策略（manualChunks）降低主包体积与首屏压力

## 3. 技术栈

### 后端（server）

- Go 1.21+
- Gin
- GORM + MySQL
- Redis
- JWT
- Swagger

### 前端（web）

- Vue 3 + TypeScript
- Vite
- Pinia
- Vue Router
- TailwindCSS + DaisyUI
- Axios

## 4. 项目结构

```text
general_go_vue_admin/
├── docs/                       # 项目文档
│   ├── EVOLUTION_MATRIX.md     # 演进矩阵（对标能力清单）
│   ├── backend/
│   └── frontend/
├── server/                     # Go 后端
│   ├── api/
│   │   ├── controller/
│   │   ├── service/
│   │   ├── dao/
│   │   └── entity/
│   ├── common/
│   ├── middleware/
│   ├── pkg/
│   ├── router/
│   ├── test/
│   ├── tools/
│   │   └── self_check.sh       # 一键自检
│   ├── config.yaml
│   └── main.go
└── web/                        # Vue 前端
    ├── src/
    │   ├── api/
    │   ├── views/
    │   ├── components/
    │   ├── router/
    │   ├── store/
    │   └── utils/
    ├── vite.config.ts
    └── package.json
```

## 5. 快速启动

## 5.1 环境准备

- Go 1.21+
- Node.js 18+
- MySQL 8+
- Redis 6+

## 5.2 启动依赖（示例）

```bash
# Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# MySQL
docker run -d --name mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=admin1234 \
  mysql:8.0
```

## 5.3 启动后端

```bash
cd server
go mod tidy
go run main.go
```

默认地址：`http://127.0.0.1:8080`

## 5.4 启动前端

```bash
cd web
npm install
npm run dev
```

默认地址：`http://localhost:3000`

## 6. 初始化数据与默认账号

项目支持首启种子数据（见 `server/config.yaml` 的 `seed` 配置）。  
默认会初始化管理员账号（可在配置中修改）：

- username: `admin`
- password: `admin123`

> 如果数据库中已存在管理员数据，则不会重复初始化。

## 7. 关键配置说明（server/config.yaml）

## 7.1 服务与数据库

- `server.port` / `server.host` / `server.model`
- `db.*`（连接、连接池、慢查询阈值等）
- `redis.*`

## 7.2 登录安全策略

```yaml
security:
  loginFailedAttemptLimit: 5
  loginLockMinutes: 15
```

说明：

- 连续失败达到 `loginFailedAttemptLimit` 时触发锁定
- 锁定时长为 `loginLockMinutes`

## 7.3 种子数据

```yaml
seed:
  enable: true
  admin:
    username: admin
    password: admin123
    nickname: 系统管理员
    email: admin@example.com
    phone: 13800138000
```

## 8. API 与调试入口

- Swagger：`/swagger/index.html`
- 验证码：`GET /api/captcha`
- 登录：`POST /api/login`
- 注册：`POST /api/register`

## 9. 质量保障与自检

## 9.1 一键自检（推荐）

```bash
server/tools/self_check.sh
```

执行内容：

1. `go test ./...`
2. `go build ./...`
3. `npm run build-only`

## 9.2 手动检查（可选）

```bash
cd server && go vet ./...
```

## 10. 已完成优化（近期）

- 修复配置文件路径依赖启动目录的问题，增强测试/部署稳定性
- 修复若干结构体 tag 与不可达代码问题（`go vet` 清零）
- 修复验证码可复用漏洞（改为一次性消费）
- 增强登录安全：失败次数限制 + 临时锁定 + 可配置策略
- 个人中心链路打通：资料映射、改密参数、头像上传闭环
- 上传接口异常分支补全，避免失败后继续执行
- 前端构建分包优化，降低主包体积

## 11. 演进路线

详见：[docs/EVOLUTION_MATRIX.md](./docs/EVOLUTION_MATRIX.md)

高优先级方向：

- 字典管理
- 系统参数中心
- 定时任务中心
- 数据权限（按组织维度）
- 审计增强（变更前后 diff / 风险分级）

## 12. 开发约定

- 新增功能默认补充错误处理和边界验证
- 新增后端功能至少通过 `go test` + `go build` + `go vet`
- 前端提交前至少通过 `npm run build-only`
- 重要改动建议同步更新 `docs/` 文档

## 13. License

MIT
