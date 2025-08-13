import storage from './storage'
import { useMainStore } from '@/store'

/**
 * 认证工具类
 * 提供登录状态检查、token 验证等功能
 */
export class AuthUtils {
  
  /**
   * 检查用户是否已登录
   * @returns {boolean} 是否已登录
   */
  static isAuthenticated(): boolean {
    const token = storage.getItem('token')
    const sysAdmin = storage.getItem('sysAdmin')

    console.log('检查认证状态:', { token: token ? `${token.substring(0, 20)}...` : null, sysAdmin: !!sysAdmin })

    // 检查 token 和用户信息是否存在
    if (!token || !sysAdmin) {
      console.log('Token 或用户信息不存在')
      return false
    }

    // 暂时跳过 JWT 过期检查，让后端来验证 token 有效性
    // 这样可以避免前端解析 token 时的格式问题
    console.log('认证检查通过（跳过过期检查）')
    return true
  }
  
  /**
   * 检查 token 是否过期
   * @param {string} token JWT token
   * @returns {boolean} 是否过期
   */
  static isTokenExpired(token: string): boolean {
    try {
      // 检查 token 格式
      if (!token || typeof token !== 'string') {
        console.log('Token 为空或格式错误')
        return true
      }

      // 检查是否是 JWT 格式 (应该有3个部分，用.分隔)
      const parts = token.split('.')
      if (parts.length !== 3) {
        console.log('Token 不是标准的 JWT 格式，跳过过期检查:', token.substring(0, 50) + '...')
        // 对于非 JWT token，我们假设它是有效的，让后端来验证
        return false
      }

      // 解析 JWT token
      const payload = JSON.parse(atob(parts[1]))
      const currentTime = Math.floor(Date.now() / 1000)

      console.log('JWT Payload:', payload)
      console.log('当前时间:', currentTime, '过期时间:', payload.exp)

      // 检查是否有过期时间字段
      if (!payload.exp) {
        console.warn('Token 没有过期时间字段，假设有效')
        return false
      }

      // 检查是否过期
      const isExpired = payload.exp < currentTime
      console.log('Token 是否过期:', isExpired)
      return isExpired
    } catch (error) {
      console.error('Token 解析失败:', error)
      console.log('Token 内容:', token)
      // 如果解析失败，我们暂时假设 token 有效，让后端来验证
      return false
    }
  }
  
  /**
   * 获取当前用户信息
   * @returns {any} 用户信息
   */
  static getCurrentUser(): any {
    return storage.getItem('sysAdmin')
  }
  
  /**
   * 获取当前 token
   * @returns {string|null} token
   */
  static getToken(): string | null {
    return storage.getItem('token')
  }
  
  /**
   * 登出用户
   * 清除所有存储的用户信息
   */
  static logout(): void {
    // 清除本地存储
    storage.removeItem('token')
    storage.removeItem('sysAdmin')
    storage.removeItem('leftMenuList')
    storage.removeItem('permissionList')
    storage.removeItem('activePath')
    storage.removeItem('keepLoggedIn')
    
    // 清除 store 状态
    const store = useMainStore()
    store.$reset()
  }
  
  /**
   * 检查是否启用了"保持登录"
   * @returns {boolean} 是否保持登录
   */
  static isKeepLoggedIn(): boolean {
    return storage.getItem('keepLoggedIn') === 'true'
  }
  
  /**
   * 设置"保持登录"状态
   * @param {boolean} keep 是否保持登录
   */
  static setKeepLoggedIn(keep: boolean): void {
    storage.setItem('keepLoggedIn', keep.toString())
  }
  
  /**
   * 刷新用户会话
   * 如果启用了"保持登录"，则延长会话时间
   */
  static refreshSession(): void {
    if (this.isKeepLoggedIn() && this.isAuthenticated()) {
      // 这里可以调用后端 API 刷新 token
      // 暂时只更新本地时间戳
      const timestamp = Date.now().toString()
      storage.setItem('lastActivity', timestamp)
    }
  }
  
  /**
   * 检查会话是否活跃
   * @returns {boolean} 会话是否活跃
   */
  static isSessionActive(): boolean {
    const lastActivity = storage.getItem('lastActivity')
    if (!lastActivity) {
      // 如果没有活动记录，设置当前时间并返回 true
      this.refreshSession()
      return true
    }

    const now = Date.now()
    const lastTime = parseInt(lastActivity)
    const maxInactiveTime = this.isKeepLoggedIn() ? 7 * 24 * 60 * 60 * 1000 : 2 * 60 * 60 * 1000 // 保持登录7天，否则2小时

    return (now - lastTime) < maxInactiveTime
  }
}

/**
 * 路由守卫辅助函数
 */
export const routeGuards = {
  /**
   * 需要认证的路由守卫
   */
  requireAuth: (to: any, from: any, next: any) => {
    if (AuthUtils.isAuthenticated() && AuthUtils.isSessionActive()) {
      // 刷新会话
      AuthUtils.refreshSession()
      next()
    } else {
      // 未登录，重定向到登录页
      next({
        path: '/adminLogin',
        query: { redirect: to.fullPath }
      })
    }
  },
  
  /**
   * 已登录用户访问登录页的守卫
   */
  redirectIfAuthenticated: (to: any, from: any, next: any) => {
    if (AuthUtils.isAuthenticated() && AuthUtils.isSessionActive()) {
      // 已登录，重定向到权限管理模块的仪表板
      next('/auth/dashboard')
    } else {
      next()
    }
  }
}
