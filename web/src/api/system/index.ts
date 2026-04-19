/*
 * @Author: JimZhang
 * @Date: 2025-07-27 14:27:46
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-28 10:46:22
 * @FilePath: /go-vue-general-admin/web/src/api/system/index.ts
 * @Description: 后端 api 接口统一管理
 *
 */


import request from "@/utils/request";

// 定义通用响应接口
interface ApiResponse<T = any> {
  code: number;
  data: T;
  message: string;
}

// 定义分页参数接口
interface PageParams {
  pageNum?: number;
  pageSize?: number;
}

// 定义岗位相关接口
interface PostParams extends PageParams {
  postName?: string;
  postCode?: string;
  status?: string;
}

interface PostData {
  id?: number;
  postName: string;
  postCode: string;
  sort: number;
  status: string;
  remark?: string;
}

// 定义部门相关接口
interface DeptParams extends PageParams {
  deptName?: string;
  status?: string;
}

interface DeptData {
  id?: number;
  parentId: number;
  deptName: string;
  sort: number;
  status: string;
  remark?: string;
}

// 定义菜单相关接口
interface MenuParams extends PageParams {
  menuName?: string;
  status?: string;
}

interface MenuData {
  id?: number;
  parentId: number;
  menuName: string;
  menuType: string;
  sort: number;
  component?: string;
  perms?: string;
  icon?: string;
  status: string;
  remark?: string;
}

// 定义角色相关接口
interface RoleParams extends PageParams {
  roleName?: string;
  roleKey?: string;
  status?: string;
}

interface RoleData {
  id?: number;
  roleName: string;
  roleKey: string;
  sort: number;
  status: string;
  remark?: string;
}

// 定义管理员相关接口
interface AdminParams extends PageParams {
  username?: string;
  phone?: string;
  status?: string;
}

interface AdminData {
  id?: number;
  username: string;
  password?: string;
  nickname?: string;
  roleId?: number;
  phone: string;
  email?: string;
  deptId: number;
  postId: number;
  remark?: string;
  status: string | number;
}

// 定义登录相关接口
interface LoginData {
  username: string;
  password: string;
  image: string;  // 验证码图片内容
  idKey: string;  // 验证码ID
}

interface RegisterData {
  username: string;
  nickname?: string;
  email?: string;
  phone?: string;
  password: string;
  confirmPassword: string;
  image: string;
  idKey: string;
}

interface UpdateProfileData {
  username: string;
  nickname?: string;
  phone?: string;
  email?: string;
  note?: string;
  icon?: string;
}

interface ChangePasswordData {
  password: string;
  newPassword: string;
  resetPassword: string;
}

interface SettingItem {
  id?: number;
  groupKey: string;
  settingKey: string;
  settingName: string;
  settingValue: string;
  valueType?: string;
  optionsJson?: string;
  isEncrypted?: boolean;
  sort?: number;
  remark?: string;
}

interface NoticeParams extends PageParams {
  title?: string;
  noticeType?: string;
  status?: string | number;
}

interface NoticeData {
  id?: number;
  title: string;
  content: string;
  noticeType: string;
  noticeLevel: string;
  status: number;
  targetType: string;
}

interface LoginResponse {
  code: number;
  message: string;
  data: {
    token: string;
    sysAdmin: any;
    leftMenuList: any;
    permissionList: any;
  };
}

// 定义日志相关接口
interface LogParams extends PageParams {
  username?: string;
  beginTime?: string;
  endTime?: string;
}

interface LoginLogParams extends PageParams {
  username?: string;
  loginStatus?: string | number;
  title?: string;
  beginTime?: string;
  endTime?: string;
}

class AdminApi{
    // 验证码接口
    captcha() {
      return request({
        url: '/captcha',
        method: 'get',
        headers: {
          // 忽略全局的请求头
          isToken: false
        }
      })
    }

    // 登录接口
    login(data: LoginData): Promise<{ data: LoginResponse }> {
      return request({
        url: '/login',
        method: 'post',
        data,
        headers: {
          isToken: false
        }
      })
    }

    // 注册接口
    register(data: RegisterData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/register',
        method: 'post',
        data,
        headers: {
          isToken: false
        }
      })
    }

    // 登出接口
    logout() {
      return request({
        url: '/logout',
        method: 'post'
      })
    }

    // 获取管理员列表
    getAdminList(params: AdminParams): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/list',
        method: 'get',
        params
      })
    }

    // 添加管理员
    addAdmin(data: AdminData): Promise<{ data: ApiResponse }> {
      // 后端字段映射
      const payload = {
        username: data.username,
        password: data.password,
        nickname: data.nickname,
        phone: data.phone,
        email: data.email,
        deptId: data.deptId,
        postId: data.postId,
        roleId: data.roleId,
        note: data.remark,
        status: Number(data.status)
      }
      return request({
        url: '/admin/add',
        method: 'post',
        data: payload
      })
    }

    // 更新管理员
    updateAdmin(data: AdminData): Promise<{ data: ApiResponse }> {
      const payload = {
        id: data.id,
        username: data.username,
        nickname: data.nickname,
        phone: data.phone,
        email: data.email,
        deptId: data.deptId,
        postId: data.postId,
        roleId: data.roleId,
        note: data.remark,
        status: Number(data.status)
      }
      return request({
        url: '/admin/update',
        method: 'put',
        data: payload
      })
    }

    // 删除管理员
    deleteAdmin(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/delete',
        method: 'delete',
        data: { id }
      })
    }

    // 更新管理员状态
    updateAdminStatus(data: { id: number; status: number | string }): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/updateStatus',
        method: 'put',
        data
      })
    }

    // 重置管理员密码
    resetAdminPassword(data: { id: number; password: string }): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/updatePassword',
        method: 'put',
        data
      })
    }

    // 获取当前用户信息（不传id时后端默认返回当前登录用户）
    getCurrentUser(id?: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/info',
        method: 'get',
        params: id ? { id } : undefined
      })
    }

    // 更新个人信息
    updateProfile(data: UpdateProfileData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/updatePersonal',
        method: 'put',
        data
      })
    }

    // 修改个人密码
    changePassword(data: ChangePasswordData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/updatePersonalPassword',
        method: 'put',
        data
      })
    }

    // 上传文件
    uploadFile(data: FormData): Promise<{ data: ApiResponse<string> }> {
      return request({
        url: '/upload',
        method: 'post',
        data,
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
    }

    // 获取角色列表
    getRoleList(params: RoleParams): Promise<{ data: ApiResponse }> {
      return request({
        url: '/role/list',
        method: 'get',
        params
      })
    }

    // 添加角色
    addRole(data: RoleData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/role/add',
        method: 'post',
        data
      })
    }

    // 更新角色
    updateRole(data: RoleData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/role/update',
        method: 'put',
        data
      })
    }

    // 删除角色
    deleteRole(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/role/delete',
        method: 'delete',
        data: { id }
      })
    }

    // 角色下拉列表
    getRoleSelectList(): Promise<{ data: ApiResponse }> {
      return request({
        url: '/role/vo/list',
        method: 'get'
      })
    }

    // 查询角色已有权限ID
    getRolePermissionIds(id: number): Promise<{ data: ApiResponse<number[]> }> {
      return request({
        url: '/role/vo/idList',
        method: 'get',
        params: { id }
      })
    }

    // 分配角色权限
    assignRolePermissions(data: { id: number; menuIds: number[] }): Promise<{ data: ApiResponse }> {
      return request({
        url: '/role/assignPermissions',
        method: 'put',
        data
      })
    }

    // 获取菜单列表
    getMenuList(params: any): Promise<{ data: ApiResponse }> {
      return request({
        url: '/menu/list',
        method: 'get',
        params
      })
    }

    // 添加菜单
    addMenu(data: MenuData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/menu/add',
        method: 'post',
        data
      })
    }

    // 更新菜单
    updateMenu(data: MenuData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/menu/update',
        method: 'put',
        data
      })
    }

    // 删除菜单
    deleteMenu(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/menu/delete',
        method: 'delete',
        data: { id }
      })
    }

    // 获取部门列表
    getDeptList(params: DeptParams): Promise<{ data: ApiResponse }> {
      return request({
        url: '/dept/list',
        method: 'get',
        params
      })
    }

    // 添加部门
    addDept(data: DeptData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/dept/add',
        method: 'post',
        data
      })
    }

    // 更新部门
    updateDept(data: DeptData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/dept/update',
        method: 'put',
        data
      })
    }

    // 删除部门
    deleteDept(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/dept/delete',
        method: 'delete',
        data: { id }
      })
    }

    // 获取岗位列表
    getPostList(params: PostParams): Promise<{ data: ApiResponse }> {
      return request({
        url: '/post/list',
        method: 'get',
        params
      })
    }

    // 添加岗位
    addPost(data: PostData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/post/add',
        method: 'post',
        data
      })
    }

    // 更新岗位
    updatePost(data: PostData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/post/update',
        method: 'put',
        data
      })
    }

    // 删除岗位
    deletePost(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/post/delete',
        method: 'delete',
        data: { id }
      })
    }

    // 获取操作日志
    getOperationLogs(params: LogParams): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysOperationLog/list',
        method: 'get',
        params
      })
    }

    // 删除单条操作日志
    deleteSysOperationLog(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysOperationLog/delete',
        method: 'delete',
        params: { id }
      })
    }

    // -------------------------
    // 系统监控管理
    // -------------------------

    // 获取服务器监控信息
    getServerMonitorInfo(): Promise<{ data: ApiResponse }> {
      return request({
        url: '/monitor/server',
        method: 'get'
      });
    }

    // 获取 Dashboard 大盘数据
    getDashboardStats(): Promise<{ data: ApiResponse }> {
      return request({
        url: '/monitor/dashboard',
        method: 'get'
      });
    }

    // 批量删除操作日志
    batchDeleteOperationLogs(ids: number[]): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysOperationLog/batch/delete',
        method: 'delete',
        data: { ids }
      })
    }

    // 导出演作日志
    exportOperationLogs(params: LogParams): Promise<Blob> {
      return request({
        url: '/sysOperationLog/export',
        method: 'get',
        params,
        responseType: 'blob'
      })
    }

    // 清空操作日志
    cleanOperationLogs(): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysOperationLog/clean',
        method: 'delete'
      })
    }

    // 获取登录日志
    getLoginLogs(params: LoginLogParams): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysLoginInfo/list',
        method: 'get',
        params
      })
    }

    // 导出登录日志
    exportLoginLogs(params: LoginLogParams): Promise<Blob> {
      return request({
        url: '/sysLoginInfo/export',
        method: 'get',
        params,
        responseType: 'blob'
      })
    }

    // 删除单条登录日志
    deleteLoginLog(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysLoginInfo/delete',
        method: 'delete',
        data: { id }
      })
    }

    // 批量删除登录日志
    batchDeleteLoginLogs(ids: number[]): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysLoginInfo/batch/delete',
        method: 'delete',
        data: { ids }
      })
    }

    // 清空登录日志
    cleanLoginLogs(): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysLoginInfo/clean',
        method: 'delete'
      })
    }

    // 获取系统设置
    getSettingList(groupKey: string): Promise<{ data: ApiResponse<SettingItem[]> }> {
      return request({
        url: '/setting/list',
        method: 'get',
        params: { groupKey }
      })
    }

    // 批量更新系统设置
    batchUpdateSettings(items: SettingItem[]): Promise<{ data: ApiResponse }> {
      return request({
        url: '/setting/batch/update',
        method: 'put',
        data: { items }
      })
    }

    // 获取通知列表
    getNoticeList(params: NoticeParams): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/list',
        method: 'get',
        params
      })
    }

    // 获取当前用户通知
    getCurrentNotices(params: { limit?: number; unreadOnly?: boolean }): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/current',
        method: 'get',
        params
      })
    }

    // 新增通知
    addNotice(data: NoticeData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/add',
        method: 'post',
        data
      })
    }

    // 更新通知
    updateNotice(data: NoticeData): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/update',
        method: 'put',
        data
      })
    }

    // 更新通知状态
    updateNoticeStatus(data: { id: number; status: number }): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/updateStatus',
        method: 'put',
        data
      })
    }

    // 删除通知
    deleteNotice(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/delete',
        method: 'delete',
        data: { id }
      })
    }

    // 批量删除通知
    batchDeleteNotices(ids: number[]): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/batch/delete',
        method: 'delete',
        data: { ids }
      })
    }

    // 标记通知已读
    markNoticeRead(id: number): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/read',
        method: 'put',
        data: { id }
      })
    }

    // 全部标记已读
    markAllNoticeRead(): Promise<{ data: ApiResponse }> {
      return request({
        url: '/notice/readAll',
        method: 'put'
      })
    }

}

export default new AdminApi();
