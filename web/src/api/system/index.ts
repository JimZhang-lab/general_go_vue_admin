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
  title?: string;
  operName?: string;
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
        data
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

    // 获取当前用户信息(后端未提供此接口，保留占位)
    getCurrentUser(): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/info',
        method: 'get'
      })
    }

    // 更新个人信息
    updateProfile(data: any): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/updatePersonal',
        method: 'put',
        data
      })
    }

    // 修改个人密码
    changePassword(data: { password: string; newPassword: string; resetPassword: string }): Promise<{ data: ApiResponse }> {
      return request({
        url: '/admin/updatePersonalPassword',
        method: 'put',
        data
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
        url: `/api/dept/delete/${id}`,
        method: 'delete'
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
    getOperationLogs(params: any): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysOperationLog/list',
        method: 'get',
        params
      })
    }

    // 获取登录日志
    getLoginLogs(params: any): Promise<{ data: ApiResponse }> {
      return request({
        url: '/sysLoginInfo/list',
        method: 'get',
        params
      })
    }

}

export default new AdminApi();