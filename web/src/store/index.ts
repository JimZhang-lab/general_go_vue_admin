import { defineStore } from 'pinia'
import storage from '@/utils/storage'

/**
 * 状态管理接口定义
 * 定义了应用中使用的各种状态类型
 */
export interface State {
  /**
   * 系统管理员信息
   */
  sysAdmin: any

  /**
   * 认证令牌
   */
  token: string | null

  /**
   * 左侧菜单列表
   */
  leftMenuList: any

  /**
   * 权限列表
   */
  permissionList: any

  /**
   * 当前激活路径
   */
  activePath: string | null
}

/**
 * 获取初始状态函数
 * 从本地存储中获取初始状态值，如果不存在则使用空字符串作为默认值
 * @returns 初始化的状态对象
 */
const getInitialState = (): State => ({
  sysAdmin: storage.getItem("sysAdmin") || "",
  token: storage.getItem("token") || "",
  leftMenuList: storage.getItem("leftMenuList") || "",
  permissionList: storage.getItem("permissionList") || "",
  activePath: storage.getItem("activePath") || ""
})

/**
 * 创建主存储实例
 * 使用 pinia 的 defineStore 创建一个名为 'main' 的存储实例
 */
export const useMainStore = defineStore('main', {
  state: getInitialState,
  actions: {
    /**
     * 保存系统管理员信息到状态和本地存储中
     * @param sysAdmin 系统管理员信息
     */
    saveSysAdmin(sysAdmin: any) {
      this.sysAdmin = sysAdmin
      storage.setItem('sysAdmin', sysAdmin)
    },
    /**
     * 保存认证令牌到状态和本地存储中
     * @param token 认证令牌
     */
    saveToken(token: string) {
      this.token = token
      storage.setItem('token', token)
    },
    /**
     * 保存左侧菜单列表到状态和本地存储中
     * @param leftMenuList 左侧菜单列表
     */
    saveLeftMenuList(leftMenuList: any) {
      this.leftMenuList = leftMenuList
      storage.setItem('leftMenuList', leftMenuList)
    },
    /**
     * 保存权限列表到状态和本地存储中
     * @param permissionList 权限列表
     */
    savePermissionList(permissionList: any) {
      this.permissionList = permissionList
      storage.setItem('permissionList', permissionList)
    },
    /**
     * 保存当前激活路径到状态和本地存储中
     * @param activePath 激活路径
     */
    saveActivePath(activePath: string) {
      this.activePath = activePath
      storage.setItem('activePath', activePath)
    },

    /**
     * 重置 store 状态到初始值
     */
    $reset() {
      this.sysAdmin = ""
      this.token = ""
      this.leftMenuList = ""
      this.permissionList = ""
      this.activePath = ""
    },

    /**
     * 登出用户，清除所有状态和存储
     */
    logout() {
      // 清除本地存储
      storage.removeItem('token')
      storage.removeItem('sysAdmin')
      storage.removeItem('leftMenuList')
      storage.removeItem('permissionList')
      storage.removeItem('activePath')
      storage.removeItem('keepLoggedIn')
      storage.removeItem('lastActivity')

      // 重置状态
      this.$reset()
    }
  }
})

export default useMainStore