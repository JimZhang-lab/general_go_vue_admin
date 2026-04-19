# Auth组件库

这是专为权限管理系统优化的组件库，包含了完整的布局组件和功能组件。

## 📁 组件结构

### 布局组件
- **AuthLayout.vue** - 权限管理系统主布局
- **AuthSidebar.vue** - 优化的侧边栏导航
- **AuthHeader.vue** - 优化的顶部导航栏

### Header子组件
- **AuthHeaderLogo.vue** - 带页面信息的Logo组件
- **AuthSearchBar.vue** - 智能搜索栏，支持全局搜索
- **AuthNotificationMenu.vue** - 通知中心下拉菜单
- **AuthUserMenu.vue** - 用户菜单下拉组件

### Sidebar子组件
- **AuthSidebarWidget.vue** - 侧边栏底部小部件

### 通用组件
- **AuthCard.vue** - 统一的卡片组件
- **AuthButton.vue** - 多样式按钮组件
- **AuthInput.vue** - 功能丰富的输入框组件

## 🚀 主要特性

### AuthLayout
- 完整的权限管理系统布局
- 响应式设计，支持移动端
- 集成侧边栏和顶部导航
- 暗色模式支持

### AuthSidebar
- 权限管理专用菜单结构
- 支持展开/收缩
- 鼠标悬停展开
- 移动端友好
- 包含系统状态小部件

### AuthHeader
- 智能页面标题显示
- 全局搜索功能
- 通知中心
- 用户菜单
- 快速操作按钮

### AuthSearchBar
- 全局智能搜索
- 支持管理员、角色、权限搜索
- 键盘导航支持
- 搜索结果分类显示
- 实时搜索建议

### AuthNotificationMenu
- 实时通知显示
- 通知分类和状态管理
- 未读数量提醒
- 通知操作（标记已读、清空）
- 多种通知类型支持

### AuthUserMenu
- 用户信息展示
- 在线状态显示
- 快速设置入口
- 系统信息显示
- 安全退出功能

## 📖 使用方法

### 基本使用

```vue
<template>
  <AuthLayout>
    <!-- 页面内容 -->
    <div>
      <h1>权限管理页面</h1>
      <!-- 其他内容 -->
    </div>
  </AuthLayout>
</template>

<script setup>
import { AuthLayout } from '@/components/auth'
</script>
```

### 单独使用组件

```vue
<template>
  <div>
    <!-- 使用搜索栏 -->
    <AuthSearchBar />
    
    <!-- 使用通知菜单 -->
    <AuthNotificationMenu />
    
    <!-- 使用用户菜单 -->
    <AuthUserMenu />
  </div>
</template>

<script setup>
import { 
  AuthSearchBar, 
  AuthNotificationMenu, 
  AuthUserMenu 
} from '@/components/auth'
</script>
```

## 🎨 样式特性

### 设计语言
- 使用 `rounded-2xl` 大圆角设计
- `shadow-lg` 阴影效果
- 统一的颜色方案和间距
- 完整的暗色模式支持

### 响应式设计
- 移动端友好的布局
- 灵活的网格系统
- 自适应的组件尺寸

### 交互体验
- 流畅的过渡动画
- 直观的操作反馈
- 键盘导航支持
- 无障碍访问支持

## 🔧 配置选项

### AuthSidebar菜单配置
```javascript
const authMenuGroups = [
  {
    title: "权限管理",
    items: [
      {
        icon: ChartBarIcon,
        name: "权限总览",
        path: "/auth/dashboard",
      },
      // 更多菜单项...
    ],
  },
  // 更多菜单组...
]
```

### AuthSearchBar搜索数据
```javascript
const searchData = [
  {
    id: 1,
    type: 'admin',
    title: 'admin',
    description: '系统管理员',
    path: '/auth/admin'
  },
  // 更多搜索项...
]
```

### AuthNotificationMenu通知配置
```javascript
const notifications = [
  {
    id: 1,
    type: 'security',
    title: '安全警告',
    message: '检测到异常登录尝试',
    time: new Date(),
    read: false,
  },
  // 更多通知...
]
```

## 🎯 最佳实践

1. **布局使用**: 在权限管理相关页面使用 `AuthLayout` 替代 `AdminLayout`
2. **组件复用**: 充分利用通用组件 `AuthCard`、`AuthButton`、`AuthInput`
3. **搜索优化**: 根据实际数据结构配置搜索数据源
4. **通知管理**: 集成实时通知系统，提供良好的用户体验
5. **权限控制**: 结合路由权限控制，确保安全性

## 🔄 更新日志

### v2.1.0 (2025-07-28)
- ✨ 新增 AuthLayout 完整布局组件
- ✨ 新增 AuthSidebar 权限专用侧边栏
- ✨ 新增 AuthHeader 优化的顶部导航
- ✨ 新增 AuthSearchBar 全局智能搜索
- ✨ 新增 AuthNotificationMenu 通知中心
- ✨ 新增 AuthUserMenu 用户菜单
- ✨ 新增 AuthSidebarWidget 侧边栏小部件
- 🎨 统一设计语言和交互体验
- 📱 完整的响应式支持
- 🌙 完整的暗色模式支持

## 📞 技术支持

如有问题或建议，请联系开发团队或提交Issue。
