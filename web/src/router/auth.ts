/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: Auth模块路由配置
 */

import { routeGuards } from '@/utils/auth'

// 定义路由类型
interface RouteConfig {
  path: string
  name?: string
  component?: any
  redirect?: string
  meta?: {
    title?: string
    icon?: string
    requiresAuth?: boolean
    permissions?: string[]
    hideInMenu?: boolean
    hideWhenNoChildren?: boolean
    hideInBreadcrumb?: boolean
  }
  children?: RouteConfig[]
  beforeEnter?: any
}

// Auth模块路由配置
export const authRoutes: RouteConfig[] = [
  {
    path: '/adminLogin',
    name: 'AdminLogin',
    component: () => import('@/views/Auth/Login.vue'),
    meta: {
      title: '管理员登录',
      requiresAuth: false,
      hideInMenu: true
    },
    beforeEnter: routeGuards.redirectIfAuthenticated
  },
  {
    path: '/auth',
    name: 'AuthManagement',
    redirect: '/auth/dashboard',
    meta: {
      title: '权限管理',
      icon: 'shield-check',
      requiresAuth: true,
      permissions: ['system:auth:view']
    },
    children: [
      {
        path: 'dashboard',
        name: 'AuthDashboard',
        component: () => import('@/views/Auth/Dashboard.vue'),
        meta: {
          title: '权限总览',
          icon: 'chart-bar',
          requiresAuth: true,
          permissions: ['system:auth:view']
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'test-components',
        name: 'TestAuthComponents',
        component: () => import('@/views/Auth/TestAuthComponents.vue'),
        meta: {
          title: 'Auth组件测试',
          requiresAuth: true,
          permissions: ['system:auth:view'],
          hideInMenu: true
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'layout-test',
        name: 'LayoutTest',
        component: () => import('@/views/Auth/LayoutTest.vue'),
        meta: {
          title: '布局测试',
          requiresAuth: true,
          permissions: ['system:auth:view'],
          hideInMenu: true
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'admin',
        name: 'AdminManagement',
        component: () => import('@/views/Auth/AdminManagement.vue'),
        meta: {
          title: '管理员管理',
          icon: 'users',
          requiresAuth: true,
          permissions: ['system:admin:list']
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'role',
        name: 'RoleManagement',
        component: () => import('@/views/Auth/RoleManagement.vue'),
        meta: {
          title: '角色管理',
          icon: 'user-group',
          requiresAuth: true,
          permissions: ['system:role:list']
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'permission',
        name: 'PermissionManagement',
        component: () => import('@/views/Auth/PermissionManagement.vue'),
        meta: {
          title: '权限管理',
          icon: 'key',
          requiresAuth: true,
          permissions: ['system:menu:list']
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Auth/Profile.vue'),
        meta: {
          title: '个人资料',
          icon: 'user-circle',
          requiresAuth: true,
          hideInMenu: false
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'logs',
        name: 'SystemLogs',
        component: () => import('@/views/Auth/SystemLogs.vue'),
        meta: {
          title: '系统日志',
          icon: 'document-text',
          requiresAuth: true,
          permissions: ['system:log:list']
        },
        beforeEnter: routeGuards.requireAuth
      },
      {
        path: 'test',
        name: 'TestAuth',
        component: () => import('@/views/Auth/TestAuth.vue'),
        meta: {
          title: 'Auth模块测试',
          icon: 'beaker',
          requiresAuth: false,
          hideInMenu: false
        }
      },
      {
        path: 'toast-test',
        name: 'ToastTest',
        component: () => import('@/views/Auth/ToastTest.vue'),
        meta: {
          title: 'Toast测试',
          icon: 'chat-bubble-left-ellipsis',
          requiresAuth: false,
          hideInMenu: false
        }
      },
      {
        path: 'logout-test',
        name: 'LogoutTest',
        component: () => import('@/views/Auth/LogoutTest.vue'),
        meta: {
          title: '登出测试',
          icon: 'arrow-right-on-rectangle',
          requiresAuth: false,
          hideInMenu: false
        }
      },
      {
        path: 'confirm-button-test',
        name: 'ConfirmButtonTest',
        component: () => import('@/views/Auth/ConfirmButtonTest.vue'),
        meta: {
          title: '确认按钮测试',
          icon: 'check-circle',
          requiresAuth: false,
          hideInMenu: false
        }
      },
      {
        path: 'captcha-debug',
        name: 'CaptchaDebug',
        component: () => import('@/views/Auth/CaptchaDebug.vue'),
        meta: {
          title: '验证码调试',
          icon: 'shield-check',
          requiresAuth: false,
          hideInMenu: false
        }
      }
    ]
  }
]

// 权限检查函数
export const checkPermission = (permissions: string[]): boolean => {
  // 从store获取用户权限列表
  const userPermissions = JSON.parse(localStorage.getItem('permissionList') || '[]')
  
  // 如果没有权限要求，直接返回true
  if (!permissions || permissions.length === 0) {
    return true
  }
  
  // 检查用户是否有任一所需权限
  return permissions.some(permission => userPermissions.includes(permission))
}

// 菜单过滤函数 - 根据权限过滤菜单
export const filterMenusByPermission = (routes: RouteConfig[]): RouteConfig[] => {
  return routes.filter(route => {
    // 检查当前路由权限
    if (route.meta?.permissions && !checkPermission(route.meta.permissions)) {
      return false
    }
    
    // 递归过滤子路由
    if (route.children && route.children.length > 0) {
      route.children = filterMenusByPermission(route.children)
      // 如果所有子路由都被过滤掉了，则隐藏父路由
      if (route.children.length === 0 && route.meta?.hideWhenNoChildren) {
        return false
      }
    }
    
    return true
  })
}

// 生成面包屑导航
export const generateBreadcrumb = (route: any): Array<{ name: string; path?: string }> => {
  const breadcrumb = []
  const matched = route.matched
  
  for (const match of matched) {
    if (match.meta?.title && !match.meta?.hideInBreadcrumb) {
      breadcrumb.push({
        name: match.meta.title,
        path: match.path === route.path ? undefined : match.path
      })
    }
  }
  
  return breadcrumb
}

// 获取页面标题
export const getPageTitle = (route: any): string => {
  const title = route.meta?.title || '管理系统'
  return `${title} - 通用管理系统`
}

// 权限指令 - 用于在模板中控制元素显示
export const vPermission = {
  mounted(el: HTMLElement, binding: any) {
    const { value } = binding
    
    if (value && value.length > 0) {
      const hasPermission = checkPermission(value)
      
      if (!hasPermission) {
        el.style.display = 'none'
        // 或者直接移除元素
        // el.parentNode?.removeChild(el)
      }
    }
  },
  updated(el: HTMLElement, binding: any) {
    const { value } = binding
    
    if (value && value.length > 0) {
      const hasPermission = checkPermission(value)
      
      if (!hasPermission) {
        el.style.display = 'none'
      } else {
        el.style.display = ''
      }
    }
  }
}

// 角色检查函数
export const checkRole = (roles: string[]): boolean => {
  // 从store获取用户角色列表
  const userRoles = JSON.parse(localStorage.getItem('userRoles') || '[]')
  
  // 如果没有角色要求，直接返回true
  if (!roles || roles.length === 0) {
    return true
  }
  
  // 检查用户是否有任一所需角色
  return roles.some(role => userRoles.includes(role))
}

// 角色指令 - 用于在模板中根据角色控制元素显示
export const vRole = {
  mounted(el: HTMLElement, binding: any) {
    const { value } = binding
    
    if (value && value.length > 0) {
      const hasRole = checkRole(value)
      
      if (!hasRole) {
        el.style.display = 'none'
      }
    }
  },
  updated(el: HTMLElement, binding: any) {
    const { value } = binding
    
    if (value && value.length > 0) {
      const hasRole = checkRole(value)
      
      if (!hasRole) {
        el.style.display = 'none'
      } else {
        el.style.display = ''
      }
    }
  }
}

export default authRoutes
