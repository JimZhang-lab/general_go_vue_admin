/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: Auth模块统一导出
 */

// 导出所有Auth相关组件
export { default as Login } from './Login.vue'
export { default as AdminManagement } from './AdminManagement.vue'
export { default as RoleManagement } from './RoleManagement.vue'
export { default as PermissionManagement } from './PermissionManagement.vue'
export { default as Profile } from './Profile.vue'
export { default as SystemLogs } from './SystemLogs.vue'

// 导出路由配置
export { authRoutes, checkPermission, filterMenusByPermission } from '@/router/auth'

// 导出权限相关工具函数
export { AuthUtils, routeGuards } from '@/utils/auth'

// Auth模块的类型定义
export interface AuthModule {
  // 管理员相关
  Admin: {
    id: number
    username: string
    phone: string
    email?: string
    deptId?: number
    postId?: number
    status: string
    remark?: string
    createTime?: string
  }
  
  // 角色相关
  Role: {
    id: number
    roleName: string
    roleKey: string
    sort: number
    status: string
    remark?: string
    createTime?: string
  }
  
  // 菜单权限相关
  Menu: {
    id: number
    menuName: string
    parentId: number
    menuType: string
    sort: number
    component?: string
    perms?: string
    icon?: string
    status: string
    remark?: string
    children?: AuthModule['Menu'][]
  }
  
  // 用户信息相关
  UserInfo: {
    id: number
    username: string
    phone?: string
    email?: string
    realName?: string
    avatar?: string
    status: string
    remark?: string
    lastLoginTime?: string
    createTime?: string
  }
  
  // 操作日志相关
  OperationLog: {
    id: number
    title: string
    operName: string
    requestMethod: string
    operUrl: string
    operIp: string
    status: number
    operTime: string
    operParam?: string
    jsonResult?: string
    errorMsg?: string
  }
  
  // 登录日志相关
  LoginLog: {
    id: number
    userName: string
    ipaddr: string
    loginLocation: string
    browser: string
    os: string
    status: number
    msg: string
    loginTime: string
  }
}

// Auth模块的常量定义
export const AUTH_CONSTANTS = {
  // 用户状态
  USER_STATUS: {
    ACTIVE: '1',
    INACTIVE: '2'
  },
  
  // 菜单类型
  MENU_TYPE: {
    DIRECTORY: 'M',  // 目录
    MENU: 'C',       // 菜单
    BUTTON: 'F'      // 按钮
  },
  
  // 请求方法
  REQUEST_METHOD: {
    GET: 'GET',
    POST: 'POST',
    PUT: 'PUT',
    DELETE: 'DELETE'
  },
  
  // 日志状态
  LOG_STATUS: {
    SUCCESS: 0,
    FAILURE: 1
  },
  
  // 登录状态
  LOGIN_STATUS: {
    SUCCESS: 1,
    FAILURE: 2
  },
  
  // 权限相关
  PERMISSIONS: {
    // 管理员管理
    ADMIN_LIST: 'system:admin:list',
    ADMIN_ADD: 'system:admin:add',
    ADMIN_EDIT: 'system:admin:edit',
    ADMIN_DELETE: 'system:admin:delete',
    ADMIN_RESET_PASSWORD: 'system:admin:resetPassword',
    
    // 角色管理
    ROLE_LIST: 'system:role:list',
    ROLE_ADD: 'system:role:add',
    ROLE_EDIT: 'system:role:edit',
    ROLE_DELETE: 'system:role:delete',
    ROLE_ASSIGN: 'system:role:assign',
    
    // 菜单管理
    MENU_LIST: 'system:menu:list',
    MENU_ADD: 'system:menu:add',
    MENU_EDIT: 'system:menu:edit',
    MENU_DELETE: 'system:menu:delete',
    
    // 日志管理
    LOG_LIST: 'system:log:list',
    LOG_EXPORT: 'system:log:export',
    
    // 个人资料
    PROFILE_EDIT: 'system:profile:edit',
    PROFILE_PASSWORD: 'system:profile:password'
  }
}

// Auth模块的工具函数
export const AuthUtils = {
  /**
   * 格式化用户状态
   */
  formatUserStatus(status: string): string {
    return status === AUTH_CONSTANTS.USER_STATUS.ACTIVE ? '启用' : '禁用'
  },
  
  /**
   * 格式化菜单类型
   */
  formatMenuType(type: string): string {
    switch (type) {
      case AUTH_CONSTANTS.MENU_TYPE.DIRECTORY:
        return '目录'
      case AUTH_CONSTANTS.MENU_TYPE.MENU:
        return '菜单'
      case AUTH_CONSTANTS.MENU_TYPE.BUTTON:
        return '按钮'
      default:
        return type
    }
  },
  
  /**
   * 格式化日期时间
   */
  formatDateTime(dateString?: string): string {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN')
  },
  
  /**
   * 格式化JSON字符串
   */
  formatJson(jsonString?: string): string {
    if (!jsonString) return ''
    try {
      return JSON.stringify(JSON.parse(jsonString), null, 2)
    } catch {
      return jsonString
    }
  },
  
  /**
   * 获取请求方法的样式类
   */
  getMethodClass(method: string): string {
    switch (method?.toUpperCase()) {
      case 'GET':
        return 'bg-green-100 text-green-800'
      case 'POST':
        return 'bg-blue-100 text-blue-800'
      case 'PUT':
        return 'bg-yellow-100 text-yellow-800'
      case 'DELETE':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  },
  
  /**
   * 生成随机密码
   */
  generateRandomPassword(length: number = 8): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
    let password = ''
    for (let i = 0; i < length; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return password
  },
  
  /**
   * 验证密码强度
   */
  validatePasswordStrength(password: string): {
    score: number
    message: string
    isValid: boolean
  } {
    let score = 0
    let message = ''
    
    if (password.length < 6) {
      return { score: 0, message: '密码长度至少6位', isValid: false }
    }
    
    if (password.length >= 8) score += 1
    if (/[a-z]/.test(password)) score += 1
    if (/[A-Z]/.test(password)) score += 1
    if (/[0-9]/.test(password)) score += 1
    if (/[^A-Za-z0-9]/.test(password)) score += 1
    
    switch (score) {
      case 0:
      case 1:
        message = '密码强度：弱'
        break
      case 2:
      case 3:
        message = '密码强度：中'
        break
      case 4:
      case 5:
        message = '密码强度：强'
        break
    }
    
    return { score, message, isValid: score >= 2 }
  }
}

// 默认导出Auth模块配置
export default {
  routes: authRoutes,
  constants: AUTH_CONSTANTS,
  utils: AuthUtils
}
