# 🎨 前端文档中心

Vue3前端技术文档集合，包含组件开发、项目配置、构建部署等完整的前端开发文档。

## 📁 文档列表

### 🚀 快速开始
- [**开发环境搭建**](#开发环境搭建) - 前端开发环境配置
- [**项目启动指南**](#项目启动) - 快速启动前端项目
- [**目录结构说明**](#项目结构) - 前端项目目录结构

### 📚 开发文档
- [**组件开发指南**](./COMPONENT_GUIDE.md) - Vue3组件开发规范
- [**路由配置指南**](./ROUTING_GUIDE.md) - Vue Router配置说明
- [**状态管理指南**](./STATE_MANAGEMENT.md) - Pinia状态管理
- [**API调用指南**](./API_GUIDE.md) - 前后端API对接

### 🎨 UI/UX文档
- [**设计规范**](./DESIGN_GUIDE.md) - UI设计规范和组件库
- [**主题配置**](./THEME_GUIDE.md) - 主题色和暗黑模式配置
- [**响应式设计**](./RESPONSIVE_GUIDE.md) - 移动端适配指南

### 🛠️ 构建部署
- [**构建配置**](./BUILD_GUIDE.md) - Vite构建配置和优化
- [**部署指南**](./DEPLOYMENT_GUIDE.md) - 生产环境部署
- [**性能优化**](./PERFORMANCE_GUIDE.md) - 前端性能优化

## 🎯 技术栈

### 核心框架
- **Vue 3**: 渐进式JavaScript框架
- **TypeScript**: JavaScript超集，提供类型安全
- **Vite**: 现代化前端构建工具
- **Vue Router**: 官方路由管理器
- **Pinia**: 新一代状态管理库

### UI组件库
- **Element Plus**: Vue 3组件库
- **Tailwind CSS**: 原子化CSS框架
- **Heroicons**: 精美的SVG图标库

### 开发工具
- **ESLint**: 代码质量检查
- **Prettier**: 代码格式化
- **Husky**: Git hooks管理
- **Commitizen**: 规范化提交信息

## 🏗️ 项目结构

```
web/
├── public/                 # 静态资源
│   ├── favicon.ico
│   └── logo.png
├── src/                   # 源代码
│   ├── api/               # API接口
│   │   ├── admin.ts       # 管理员相关API
│   │   ├── auth.ts        # 认证相关API
│   │   └── index.ts       # API统一导出
│   ├── assets/            # 资源文件
│   │   ├── images/        # 图片资源
│   │   └── styles/        # 样式文件
│   ├── components/        # 公共组件
│   │   ├── common/        # 通用组件
│   │   ├── form/          # 表单组件
│   │   └── layout/        # 布局组件
│   ├── composables/       # 组合式函数
│   │   ├── useAuth.ts     # 认证相关
│   │   └── useApi.ts      # API调用
│   ├── layouts/           # 页面布局
│   │   ├── DefaultLayout.vue
│   │   └── AuthLayout.vue
│   ├── router/            # 路由配置
│   │   ├── index.ts       # 路由主文件
│   │   └── guards.ts      # 路由守卫
│   ├── stores/            # 状态管理
│   │   ├── auth.ts        # 认证状态
│   │   ├── user.ts        # 用户状态
│   │   └── index.ts       # Store统一导出
│   ├── types/             # TypeScript类型定义
│   │   ├── api.ts         # API类型
│   │   ├── user.ts        # 用户类型
│   │   └── common.ts      # 通用类型
│   ├── utils/             # 工具函数
│   │   ├── request.ts     # HTTP请求封装
│   │   ├── storage.ts     # 本地存储
│   │   └── helpers.ts     # 辅助函数
│   ├── views/             # 页面组件
│   │   ├── Auth/          # 认证页面
│   │   ├── Dashboard/     # 仪表板
│   │   ├── Admin/         # 管理员管理
│   │   └── System/        # 系统管理
│   ├── App.vue            # 根组件
│   └── main.ts            # 应用入口
├── index.html             # HTML模板
├── package.json           # 项目配置
├── tsconfig.json          # TypeScript配置
├── vite.config.ts         # Vite配置
├── tailwind.config.ts     # Tailwind配置
└── eslint.config.ts       # ESLint配置
```

## 🛠️ 开发环境搭建

### 环境要求
- **Node.js**: 16.0+
- **npm**: 8.0+ 或 **yarn**: 1.22+
- **Git**: 2.0+

### 安装依赖
```bash
# 使用npm
npm install

# 或使用yarn
yarn install

# 或使用pnpm
pnpm install
```

## 🚀 项目启动

### 开发模式
```bash
# 启动开发服务器
npm run dev

# 或
yarn dev

# 访问地址
http://localhost:5173
```

### 构建生产版本
```bash
# 构建生产版本
npm run build

# 预览构建结果
npm run preview
```

### 代码检查
```bash
# ESLint检查
npm run lint

# 修复ESLint问题
npm run lint:fix

# TypeScript类型检查
npm run type-check
```

## 🎨 核心功能

### 认证系统
- **登录/登出**: JWT Token认证
- **权限控制**: 基于角色的权限管理
- **路由守卫**: 自动权限验证

### 用户界面
- **响应式设计**: 支持PC和移动端
- **主题切换**: 明亮/暗黑主题
- **国际化**: 多语言支持
- **组件库**: 基于Element Plus

### 数据管理
- **状态管理**: Pinia集中状态管理
- **API调用**: 统一的HTTP请求封装
- **数据缓存**: 智能数据缓存策略

### 开发体验
- **热重载**: 开发时实时更新
- **TypeScript**: 完整的类型支持
- **代码规范**: ESLint + Prettier
- **Git规范**: Husky + Commitizen

## 📱 页面功能

### 认证页面
- **登录页**: 用户名密码登录，验证码验证
- **注册页**: 用户注册功能
- **忘记密码**: 密码重置功能

### 管理页面
- **仪表板**: 系统概览和统计信息
- **用户管理**: 用户增删改查
- **角色管理**: 角色权限配置
- **菜单管理**: 动态菜单配置
- **系统设置**: 系统参数配置

### 功能页面
- **个人中心**: 个人信息管理
- **操作日志**: 用户操作记录
- **系统监控**: 系统性能监控
- **文件管理**: 文件上传下载

## 🔧 配置说明

### Vite配置
```typescript
// vite.config.ts
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

### 路由配置
```typescript
// router/index.ts
const routes = [
  {
    path: '/',
    component: DefaultLayout,
    children: [
      { path: 'dashboard', component: Dashboard },
      { path: 'admin', component: AdminManagement }
    ]
  }
]
```

### API配置
```typescript
// utils/request.ts
const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(config => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
```

## 🧪 测试

### 单元测试
```bash
# 运行单元测试
npm run test

# 测试覆盖率
npm run test:coverage
```

### E2E测试
```bash
# 运行端到端测试
npm run test:e2e
```

## 📝 开发规范

### 组件开发
- 使用组合式API (Composition API)
- 遵循Vue 3最佳实践
- 组件命名使用PascalCase
- 文件命名使用kebab-case

### 代码风格
- 使用TypeScript进行类型检查
- 遵循ESLint规则
- 使用Prettier格式化代码
- 添加必要的注释

### 提交规范
```bash
# 使用Commitizen提交
npm run commit

# 提交格式
feat: 添加新功能
fix: 修复bug
docs: 更新文档
style: 代码格式调整
refactor: 代码重构
test: 添加测试
chore: 构建过程或辅助工具的变动
```

## 🔍 故障排查

### 常见问题
1. **依赖安装失败**: 清除node_modules重新安装
2. **端口占用**: 修改vite.config.ts中的端口配置
3. **API调用失败**: 检查后端服务是否启动
4. **路由跳转异常**: 检查路由配置和权限设置

### 调试工具
- **Vue DevTools**: Vue官方调试工具
- **Network面板**: 查看API请求
- **Console**: 查看错误信息
- **Sources**: 断点调试

---

**最后更新**: 2025-07-29
**维护者**: 前端开发团队
**技术支持**: 查看各专项文档或联系开发团队
