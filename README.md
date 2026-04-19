# 通用后台管理系统 (General Go-Vue Admin)

这是一个基于现代化技术栈开发的企业级通用后台管理系统。系统采纳前后端分离架构，专注于高可维护性、最佳实践与优美的界面交互。本系统内置了完善的角色权限控制 (RBAC)、动态路由菜单、精美的数据可视化概览面板以及完整的管理员与系统日志管理模块。

## 🛠️ 技术栈 (Tech Stack)

### 前端生态 (Frontend)
- **核心框架**: Vue 3 (Composition API) 
- **构建工具**: Vite 7
- **类型系统**: TypeScript
- **状态管理**: Pinia
- **路由控制**: Vue Router 4 (支持动态路由与权限拦截)
- **UI 样式**: Tailwind CSS v4, DaisyUI, 现代玻璃拟物态(Glassmorphism)设计
- **第三方组件**:
  - `axios` 统一网络请求封装与拦截
  - `ApexCharts` 动态数据可视化图表
  - `lucide-vue-next` 与 `heroicons` 提供海量精美矢量图

### 后端生态 (Backend)
- **核心框架**: Golang (1.20+)
- **Web 框架**: Gin (高性能路由分发)
- **ORM 映射**: GORM (兼容 MySQL)
- **认证授权**: JWT (JSON Web Token)
- **实时热更**: Air (开发环境自动化平滑重载)
- **数据库**: MySQL 8.x + Redis (可选，用于缓存)

---

## 🌟 核心特性 (Features)

1. **动态菜单与权限 (RBAC)**
   - 细粒度的路由守卫拦截
   - 根据用户角色和后端配置动态渲染侧边栏菜单
   - 指令级 (`v-permission`, `v-role`) 前端按钮权限控制
2. **极佳的视觉与动效体验**
   - 适配深色模式 (Dark Mode) 与亮色模式无缝切换
   - Tailwind 提供的响应式系统，完美适配桌面与移动端设备
3. **完备的基础运维功能**
   - 多级管理员账号与组织架构支持
   - 用户登入、登出行为跟踪及系统操作日志记录
4. **性能极致优化**
   - 后端针对关键登录和资源获取接口做了并发查询优化 (`goroutines`)
   - 彻底修复 GORM 模型校验问题与 JWT 失效逻辑

---

## 🚀 快速启动 (Getting Started)

### 1. 后端服务 (Go Server)
确保本地已安装 Go 环境以及相关数据库，然后修改 `server/config.yaml` 填入您的 MySQL 账号密码。

```bash
cd server
# 下载相关包
go mod tidy

# （推荐）使用 air 进行热重载开发
go install github.com/air-verse/air@latest
air

# 普通运行
go run main.go
```
*后端访问接口默认跑在 `http://127.0.0.1:8368`。*

### 2. 前端服务 (Vue Web)
确保本地已安装 Node.js(推荐 v20+) 与 Vite。

```bash
cd web
# 安装依赖
npm install  # 或 yarn install

# 启动冷更服务
npm run dev  # 或 yarn dev
```
*前端默认跑在 `http://localhost:3000`，内置 proxy 配置会自动将 `/api` 请求代理至后端。*

---

## 🛡️ 近期优化与 Bug Fixes (Changelog)

- **[修复]** 修复了 Vue Router 4 内部由相同 `name: 'Profile'` 导致的配置层隐式冲突覆盖，解决了点击个人资料菜单后白屏并提示 "No match found" 的致命 Bug。
- **[修复]** 重构后端 `server/pkg/jwt/jwt.go` 全局变量初始化陷阱，确保下发的 Token 根据最新 Config 热载入实时计算到期时间，根除“前端收到 Token 却立刻被识别伪过期”从而陷入认证死循环无法跳转 dashboard 的疑难杂症。
- **[修复]** 修正了 GORM 在扫描包含 slice 非映射字段 `LeftMenuVo.MenuSvoList` 时的结构体解析 panic。
- **[优化]** 清理了 `AuthUserMenu` 中未挂载的空连结（如 `/help`）。

---

## 📝 开源协议
MIT License.
