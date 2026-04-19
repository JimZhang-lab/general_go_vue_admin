# Auth 模块文档

## 概述

Auth 模块是通用后台管理系统的核心权限管理模块，提供了完整的用户认证、授权、角色管理、权限控制等功能。

## 功能特性

### 🔐 用户认证
- 管理员登录/登出
- 验证码验证
- 会话管理
- 自动登录（记住我）
- 登录状态检查

### 👥 用户管理
- 管理员账户管理
- 用户信息维护
- 状态控制（启用/禁用）
- 密码重置
- 个人资料管理

### 🎭 角色管理
- 角色创建和编辑
- 角色权限分配
- 角色状态管理
- 角色层级控制

### 🔑 权限管理
- 菜单权限管理
- 按钮权限控制
- 接口权限验证
- 权限树形结构

### 📊 系统日志
- 操作日志记录
- 登录日志追踪
- 日志查询和导出
- 安全审计

## 目录结构

```
src/views/Auth/
├── Login.vue                    # 登录页面
├── AdminManagement.vue          # 管理员管理
├── RoleManagement.vue           # 角色管理
├── PermissionManagement.vue     # 权限管理
├── Profile.vue                  # 个人资料
├── SystemLogs.vue               # 系统日志
├── index.ts                     # 模块导出
└── README.md                    # 文档说明
```

## 组件说明

### 1. Login.vue - 登录组件
**功能：**
- 用户名/密码登录
- 验证码验证
- 记住登录状态
- 登录状态检查

**主要特性：**
- 智能数据绑定
- 表单验证
- 错误处理
- 自动跳转

### 2. AdminManagement.vue - 管理员管理
**功能：**
- 管理员列表查看
- 添加/编辑管理员
- 状态控制
- 密码重置
- 批量操作

**主要特性：**
- 分页查询
- 高级搜索
- 表单验证
- 权限控制

### 3. RoleManagement.vue - 角色管理
**功能：**
- 角色列表管理
- 角色权限分配
- 角色状态控制
- 权限树展示

**主要特性：**
- 树形权限选择
- 批量权限操作
- 角色继承
- 动态权限更新

### 4. PermissionManagement.vue - 权限管理
**功能：**
- 菜单权限管理
- 权限树形结构
- 权限层级控制
- 动态菜单生成

**主要特性：**
- 树形表格展示
- 拖拽排序
- 权限继承
- 实时预览

### 5. Profile.vue - 个人资料
**功能：**
- 个人信息编辑
- 密码修改
- 头像上传
- 登录记录查看

**主要特性：**
- 实时验证
- 安全密码策略
- 文件上传
- 操作日志

### 6. SystemLogs.vue - 系统日志
**功能：**
- 操作日志查看
- 登录日志追踪
- 日志搜索过滤
- 日志导出

**主要特性：**
- 多维度搜索
- 日志详情查看
- 数据导出
- 实时更新

## API 接口

### 认证相关
```typescript
// 获取验证码
adminApi.captcha()

// 用户登录
adminApi.login(data: LoginData)

// 用户登出
adminApi.logout()
```

### 管理员管理
```typescript
// 获取管理员列表
adminApi.getAdminList(params: AdminParams)

// 添加管理员
adminApi.addAdmin(data: AdminData)

// 更新管理员
adminApi.updateAdmin(data: AdminData)

// 删除管理员
adminApi.deleteAdmin(id: number)

// 更新管理员状态
adminApi.updateAdminStatus(data: { id: number; status: string })

// 重置管理员密码
adminApi.resetAdminPassword(data: { id: number; password: string })
```

### 角色管理
```typescript
// 获取角色列表
adminApi.getRoleList(params: RoleParams)

// 添加角色
adminApi.addRole(data: RoleData)

// 更新角色
adminApi.updateRole(data: RoleData)

// 删除角色
adminApi.deleteRole(id: number)
```

### 权限管理
```typescript
// 获取菜单列表
adminApi.getMenuList(params: MenuParams)

// 添加菜单
adminApi.addMenu(data: MenuData)

// 更新菜单
adminApi.updateMenu(data: MenuData)

// 删除菜单
adminApi.deleteMenu(id: number)
```

### 个人资料
```typescript
// 获取当前用户信息
adminApi.getCurrentUser()

// 更新个人信息
adminApi.updateProfile(data: any)

// 修改个人密码
adminApi.changePassword(data: { oldPassword: string; newPassword: string })
```

### 系统日志
```typescript
// 获取操作日志
adminApi.getOperationLogs(params: LogParams)

// 获取登录日志
adminApi.getLoginLogs(params: LogParams)
```

## 权限控制

### 路由权限
```typescript
// 在路由配置中设置权限
{
  path: '/auth/admin',
  meta: {
    permissions: ['system:admin:list']
  }
}
```

### 组件权限
```vue
<!-- 使用 v-permission 指令 -->
<button v-permission="['system:admin:add']">添加用户</button>

<!-- 使用权限检查函数 -->
<button v-if="checkPermission(['system:admin:edit'])">编辑</button>
```

### 角色权限
```vue
<!-- 使用 v-role 指令 -->
<div v-role="['admin', 'manager']">管理员内容</div>

<!-- 使用角色检查函数 -->
<div v-if="checkRole(['admin'])">超级管理员内容</div>
```

## 使用示例

### 1. 基本使用
```vue
<template>
  <div>
    <!-- 管理员管理组件 -->
    <AdminManagement />
  </div>
</template>

<script setup>
import { AdminManagement } from '@/views/Auth'
</script>
```

### 2. 权限检查
```vue
<template>
  <div>
    <!-- 只有有权限的用户才能看到 -->
    <button 
      v-permission="['system:admin:add']"
      @click="addAdmin"
    >
      添加管理员
    </button>
  </div>
</template>

<script setup>
import { checkPermission } from '@/views/Auth'

const canEdit = checkPermission(['system:admin:edit'])
</script>
```

### 3. 路由守卫
```typescript
import { routeGuards } from '@/utils/auth'

// 在路由配置中使用
{
  path: '/admin',
  beforeEnter: routeGuards.requireAuth
}
```

## 配置选项

### 权限常量
```typescript
import { AUTH_CONSTANTS } from '@/views/Auth'

// 用户状态
AUTH_CONSTANTS.USER_STATUS.ACTIVE    // '1'
AUTH_CONSTANTS.USER_STATUS.INACTIVE  // '2'

// 菜单类型
AUTH_CONSTANTS.MENU_TYPE.DIRECTORY   // 'M'
AUTH_CONSTANTS.MENU_TYPE.MENU        // 'C'
AUTH_CONSTANTS.MENU_TYPE.BUTTON      // 'F'
```

### 工具函数
```typescript
import { AuthUtils } from '@/views/Auth'

// 格式化用户状态
AuthUtils.formatUserStatus('1')  // '启用'

// 格式化日期时间
AuthUtils.formatDateTime('2023-01-01T00:00:00Z')

// 生成随机密码
AuthUtils.generateRandomPassword(8)

// 验证密码强度
AuthUtils.validatePasswordStrength('password123')
```

## 最佳实践

### 1. 权限设计
- 使用细粒度权限控制
- 权限命名规范：`模块:功能:操作`
- 合理设计角色层级
- 定期审查权限分配

### 2. 安全考虑
- 密码强度验证
- 登录失败限制
- 会话超时控制
- 操作日志记录

### 3. 用户体验
- 友好的错误提示
- 加载状态显示
- 操作确认对话框
- 快捷键支持

### 4. 性能优化
- 权限数据缓存
- 懒加载组件
- 虚拟滚动列表
- 分页查询

## 故障排除

### 常见问题

1. **登录失败**
   - 检查验证码是否正确
   - 确认用户名密码
   - 查看网络连接

2. **权限不生效**
   - 检查权限配置
   - 确认用户角色
   - 刷新权限缓存

3. **页面无法访问**
   - 检查路由配置
   - 确认登录状态
   - 查看权限设置

### 调试技巧

1. 开启浏览器开发者工具
2. 查看网络请求状态
3. 检查控制台错误信息
4. 使用Vue DevTools调试

## 更新日志

### v1.0.0 (2025-07-27)
- 初始版本发布
- 完整的权限管理功能
- 支持多种数据格式绑定
- 高并发处理优化
- 完善的错误处理机制

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License
