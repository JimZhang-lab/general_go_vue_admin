<template>
  <AuthLayout>
    <PageBreadcrumb :pageTitle="currentPageTitle" />
    <div class="space-y-5 sm:space-y-6">

      <!-- 搜索和操作栏 -->
      <AuthCard title="搜索条件">
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between" style="gap: 1rem;">
          <!-- 搜索表单 -->
          <div class="flex flex-col sm:flex-row flex-1" style="gap: 1rem;">
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.username"
                type="text"
                placeholder="搜索用户名..."
                icon="search"
              />
            </div>
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.phone"
                type="text"
                placeholder="搜索手机号..."
                icon="phone"
              />
            </div>
            <div class="flex-1">
              <select
                v-model="searchForm.status"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:border-gray-700 dark:bg-gray-800 dark:text-white dark:placeholder-gray-500"
              >
                <option value="">全部状态</option>
                <option value="1">启用</option>
                <option value="2">禁用</option>
              </select>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex" style="gap: 0.5rem;">
            <AuthButton
              @click="searchAdmins"
              variant="primary"
              size="md"
              text="搜索"
            />
            <AuthButton
              @click="resetSearch"
              variant="secondary"
              size="md"
              text="重置"
            />
            <AuthButton
              @click="showAddModal = true"
              variant="success"
              size="md"
              text="添加管理员"
            />
          </div>
        </div>
      </AuthCard>

      <!-- 管理员列表 -->
      <AuthCard title="管理员列表">
        <div class="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-white/[0.03]">
          <div class="max-w-full overflow-x-auto custom-scrollbar">
            <table class="min-w-full">
              <thead>
                <tr class="border-b border-gray-200 dark:border-gray-700">
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">ID</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">用户信息</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">联系方式</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">部门岗位</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">状态</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">创建时间</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">操作</p>
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                <tr v-for="admin in adminList" :key="admin.id" class="border-t border-gray-100 dark:border-gray-800">
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">#{{ admin.id }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-3">
                      <div class="w-10 h-10 overflow-hidden rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
                        <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                      </div>
                      <div>
                        <span class="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                          {{ admin.username }}
                        </span>
                        <span class="block text-gray-500 text-theme-xs dark:text-gray-400">
                          {{ admin.remark || '管理员' }}
                        </span>
                      </div>
                    </div>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div>
                      <p class="text-gray-800 text-theme-sm dark:text-white/90">{{ admin.phone }}</p>
                      <p class="text-gray-500 text-theme-xs dark:text-gray-400">{{ admin.email || '-' }}</p>
                    </div>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div>
                      <p class="text-gray-800 text-theme-sm dark:text-white/90">{{ admin.deptName || '-' }}</p>
                      <p class="text-gray-500 text-theme-xs dark:text-gray-400">{{ admin.postName || '-' }}</p>
                    </div>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <span
                      :class="[
                        'rounded-full px-2 py-0.5 text-theme-xs font-medium',
                        {
                          'bg-success-50 text-success-700 dark:bg-success-500/15 dark:text-success-500':
                            admin.status === '1',
                          'bg-error-50 text-error-700 dark:bg-error-500/15 dark:text-error-500':
                            admin.status === '2',
                        },
                      ]"
                    >
                      {{ admin.status === '1' ? '启用' : '禁用' }}
                    </span>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ formatDate(admin.createTime) }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-2">
                      <button
                        @click="editAdmin(admin)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 dark:bg-blue-500/15 dark:text-blue-400 dark:hover:bg-blue-500/25"
                      >
                        编辑
                      </button>
                      <button
                        @click="toggleAdminStatus(admin)"
                        :class="[
                          'inline-flex items-center px-2 py-1 text-xs font-medium rounded',
                          admin.status === '1'
                            ? 'text-red-600 bg-red-50 hover:bg-red-100 dark:bg-red-500/15 dark:text-red-400 dark:hover:bg-red-500/25'
                            : 'text-green-600 bg-green-50 hover:bg-green-100 dark:bg-green-500/15 dark:text-green-400 dark:hover:bg-green-500/25'
                        ]"
                      >
                        {{ admin.status === '1' ? '禁用' : '启用' }}
                      </button>
                      <button
                        @click="resetPassword(admin)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-orange-600 bg-orange-50 rounded hover:bg-orange-100 dark:bg-orange-500/15 dark:text-orange-400 dark:hover:bg-orange-500/25"
                      >
                        重置密码
                      </button>
                      <button
                        @click="deleteAdmin(admin)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 dark:bg-red-500/15 dark:text-red-400 dark:hover:bg-red-500/25"
                      >
                        删除
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </AuthCard>

      <!-- 分页 -->
      <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="prevPage"
            :disabled="pagination.page <= 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            上一页
          </button>
          <button
            @click="nextPage"
            :disabled="pagination.page >= pagination.pages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            下一页
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              显示第 <span class="font-medium">{{ (pagination.page - 1) * pagination.pageSize + 1 }}</span> 到
              <span class="font-medium">{{ Math.min(pagination.page * pagination.pageSize, pagination.total) }}</span> 条，
              共 <span class="font-medium">{{ pagination.total }}</span> 条记录
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
              <button
                @click="prevPage"
                :disabled="pagination.page <= 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
              >
                上一页
              </button>
              <button
                v-for="page in visiblePages"
                :key="page"
                @click="goToPage(page)"
                :class="page === pagination.page ? 'bg-blue-50 border-blue-500 text-blue-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
                class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
              >
                {{ page }}
              </button>
              <button
                @click="nextPage"
                :disabled="pagination.page >= pagination.pages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
              >
                下一页
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑管理员模态框 -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="fixed inset-0 bg-black/50" @click="closeModal"></div>
      <div class="relative w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-xl dark:border-gray-800 dark:bg-white/[0.03]">
          <!-- Modal Header -->
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-white/90">
              {{ showAddModal ? '添加管理员' : '编辑管理员' }}
            </h3>
            <button
              @click="closeModal"
              class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Modal Body -->
          <form @submit.prevent="submitForm" class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- 用户名 -->
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  用户名 *
                </label>
                <input
                  v-model="adminForm.username"
                  type="text"
                  required
                  class="h-11 w-full rounded-lg border border-gray-300 bg-transparent px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                  placeholder="请输入用户名"
                />
              </div>

              <!-- 密码 -->
              <div v-if="showAddModal">
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  密码 *
                </label>
                <input
                  v-model="adminForm.password"
                  type="password"
                  required
                  class="h-11 w-full rounded-lg border border-gray-300 bg-transparent px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                  placeholder="请输入密码"
                />
              </div>

              <!-- 手机号 -->
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  手机号 *
                </label>
                <input
                  v-model="adminForm.phone"
                  type="tel"
                  required
                  class="h-11 w-full rounded-lg border border-gray-300 bg-transparent px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                  placeholder="请输入手机号"
                />
              </div>

              <!-- 邮箱 -->
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  邮箱
                </label>
                <input
                  v-model="adminForm.email"
                  type="email"
                  class="h-11 w-full rounded-lg border border-gray-300 bg-transparent px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                  placeholder="请输入邮箱"
                />
              </div>

              <!-- 部门 -->
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  部门
                </label>
                <div class="relative z-20 bg-transparent">
                  <select
                    v-model="adminForm.deptId"
                    class="h-11 w-full appearance-none rounded-lg border border-gray-300 bg-transparent bg-none px-4 py-2.5 pr-11 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                  >
                    <option value="">请选择部门</option>
                    <option v-for="dept in deptList" :key="dept.id" :value="dept.id">
                      {{ dept.deptName }}
                    </option>
                  </select>
                  <span class="absolute z-30 text-gray-500 -translate-y-1/2 pointer-events-none right-4 top-1/2 dark:text-gray-400">
                    <svg class="stroke-current" width="20" height="20" viewBox="0 0 20 20" fill="none">
                      <path d="M4.79175 7.396L10.0001 12.6043L15.2084 7.396" stroke="" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </span>
                </div>
              </div>

              <!-- 岗位 -->
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  岗位
                </label>
                <div class="relative z-20 bg-transparent">
                  <select
                    v-model="adminForm.postId"
                    class="h-11 w-full appearance-none rounded-lg border border-gray-300 bg-transparent bg-none px-4 py-2.5 pr-11 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                  >
                    <option value="">请选择岗位</option>
                    <option v-for="post in postList" :key="post.id" :value="post.id">
                      {{ post.postName }}
                    </option>
                  </select>
                  <span class="absolute z-30 text-gray-500 -translate-y-1/2 pointer-events-none right-4 top-1/2 dark:text-gray-400">
                    <svg class="stroke-current" width="20" height="20" viewBox="0 0 20 20" fill="none">
                      <path d="M4.79175 7.396L10.0001 12.6043L15.2084 7.396" stroke="" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </span>
                </div>
              </div>

              <!-- 状态 -->
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  状态
                </label>
                <div class="relative z-20 bg-transparent">
                  <select
                    v-model="adminForm.status"
                    class="h-11 w-full appearance-none rounded-lg border border-gray-300 bg-transparent bg-none px-4 py-2.5 pr-11 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                  >
                    <option value="1">启用</option>
                    <option value="2">禁用</option>
                  </select>
                  <span class="absolute z-30 text-gray-500 -translate-y-1/2 pointer-events-none right-4 top-1/2 dark:text-gray-400">
                    <svg class="stroke-current" width="20" height="20" viewBox="0 0 20 20" fill="none">
                      <path d="M4.79175 7.396L10.0001 12.6043L15.2084 7.396" stroke="" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </span>
                </div>
              </div>
            </div>

            <!-- 备注 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                备注
              </label>
              <textarea
                v-model="adminForm.remark"
                rows="3"
                class="w-full rounded-lg border border-gray-300 bg-transparent px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                placeholder="请输入备注信息"
              ></textarea>
            </div>

            <!-- Modal Footer -->
            <div class="flex justify-end gap-3 pt-6 border-t border-gray-200 dark:border-gray-700">
              <button
                type="button"
                @click="closeModal"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700"
              >
                取消
              </button>
              <button
                type="submit"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-brand-600 border border-transparent rounded-lg hover:bg-brand-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500"
              >
                {{ showAddModal ? '添加' : '更新' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 重置密码模态框 -->
    <div v-if="showPasswordModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="fixed inset-0 bg-black/50" @click="showPasswordModal = false"></div>
      <div class="relative w-full max-w-md">
        <div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-xl dark:border-gray-800 dark:bg-white/[0.03]">
          <!-- Modal Header -->
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-white/90">
              重置密码
            </h3>
            <button
              @click="showPasswordModal = false"
              class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Modal Body -->
          <form @submit.prevent="submitPasswordReset" class="space-y-6">
            <!-- 新密码 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                新密码 *
              </label>
              <input
                v-model="passwordForm.password"
                type="password"
                required
                class="h-11 w-full rounded-lg border border-gray-300 bg-transparent px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                placeholder="请输入新密码"
              />
            </div>

            <!-- 确认密码 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-400">
                确认密码 *
              </label>
              <input
                v-model="passwordForm.confirmPassword"
                type="password"
                required
                class="h-11 w-full rounded-lg border border-gray-300 bg-transparent px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
                placeholder="请再次输入密码"
              />
            </div>

            <!-- Modal Footer -->
            <div class="flex justify-end gap-3 pt-6 border-t border-gray-200 dark:border-gray-700">
              <button
                type="button"
                @click="showPasswordModal = false"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700"
              >
                取消
              </button>
              <button
                type="submit"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-brand-600 border border-transparent rounded-lg hover:bg-brand-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500"
              >
                重置密码
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { AuthLayout, AuthCard, AuthButton, AuthInput } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'

const currentPageTitle = ref('管理员管理')

// 定义接口类型
interface Admin {
  id: number
  username: string
  phone: string
  email?: string
  deptId?: number
  postId?: number
  deptName?: string
  postName?: string
  status: string
  remark?: string
  createTime?: string
}

interface Dept {
  id: number
  deptName: string
}

interface Post {
  id: number
  postName: string
}

// 响应式数据
const adminList = ref<Admin[]>([])
const deptList = ref<Dept[]>([])
const postList = ref<Post[]>([])
const loading = ref(false)

// 搜索表单
const searchForm = reactive({
  username: '',
  phone: '',
  status: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  pages: 0
})

// 模态框状态
const showAddModal = ref(false)
const showEditModal = ref(false)
const showPasswordModal = ref(false)

// 表单数据
const adminForm = reactive({
  id: undefined as number | undefined,
  username: '',
  password: '',
  phone: '',
  email: '',
  deptId: 0,
  postId: 0,
  status: '1',
  remark: ''
})

// 密码重置表单
const passwordForm = reactive({
  id: 0,
  password: '',
  confirmPassword: ''
})

// 计算属性 - 可见页码
const visiblePages = computed(() => {
  const pages = []
  const total = pagination.pages
  const current = pagination.page

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }

  return pages
})

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN')
}

// 获取管理员列表
const getAdminList = async () => {
  try {
    loading.value = true
    const params = {
      pageNum: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.username,
      status: searchForm.status,
      phone: searchForm.phone
    }

    const { data: res } = await adminApi.getAdminList(params)

    if (res.code === 200) {
      adminList.value = res.data.list || []
      pagination.total = res.data.total || 0
      pagination.pages = Math.ceil(pagination.total / pagination.pageSize)

      ToastAlert.success({
        title: '获取成功',
        message: '管理员列表已加载'
      })
    } else {
      ToastAlert.error({
        title: '获取失败',
        message: res.message || '获取管理员列表失败'
      })
    }
  } catch (error) {
    console.error('获取管理员列表失败:', error)
    ToastAlert.error({
      title: '获取管理员列表失败',
      message: '网络异常，请重试'
    })
  } finally {
    loading.value = false
  }
}

// 获取部门列表
const getDeptList = async () => {
  try {
    const { data: res } = await adminApi.getDeptList({})

    if (res.code === 200) {
      deptList.value = res.data || []
    } else {
      console.error('获取部门列表失败:', res.message)
    }
  } catch (error) {
    console.error('获取部门列表失败:', error)
    // 使用模拟数据作为后备
    deptList.value = [
      { id: 1, deptName: '技术部' },
      { id: 2, deptName: '运营部' },
      { id: 3, deptName: '市场部' },
      { id: 4, deptName: '人事部' }
    ]
  }
}

// 获取岗位列表
const getPostList = async () => {
  try {
    const { data: res } = await adminApi.getPostList({})

    if (res.code === 200) {
      postList.value = res.data || []
    } else {
      console.error('获取岗位列表失败:', res.message)
    }
  } catch (error) {
    console.error('获取岗位列表失败:', error)
    // 使用模拟数据作为后备
    postList.value = [
      { id: 1, postName: '系统管理员' },
      { id: 2, postName: '运营专员' },
      { id: 3, postName: '市场专员' },
      { id: 4, postName: '人事专员' }
    ]
  }
}

// 搜索管理员
const searchAdmins = () => {
  pagination.page = 1
  getAdminList()
}

// 重置搜索
const resetSearch = () => {
  Object.assign(searchForm, {
    username: '',
    phone: '',
    status: ''
  })
  pagination.page = 1
  getAdminList()
}

// 编辑管理员
const editAdmin = (admin: Admin) => {
  Object.assign(adminForm, {
    id: admin.id,
    username: admin.username,
    phone: admin.phone,
    email: admin.email || '',
    deptId: admin.deptId || 0,
    postId: admin.postId || 0,
    status: admin.status,
    remark: admin.remark || ''
  })
  showEditModal.value = true
}

// 切换管理员状态
const toggleAdminStatus = async (admin: Admin) => {
  try {
    const newStatus = admin.status === '1' ? '2' : '1'
    const { data: res } = await adminApi.updateAdminStatus({
      id: admin.id,
      status: newStatus
    })

    if (res.code === 200) {
      ToastAlert.success({
        title: '操作成功',
        message: `管理员已${newStatus === '1' ? '启用' : '禁用'}`
      })
      getAdminList()
    } else {
      ToastAlert.error({
        title: '操作失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('切换管理员状态失败:', error)
    ToastAlert.error({
      title: '操作失败',
      message: '网络异常，请重试'
    })
  }
}

// 重置密码
const resetPassword = (admin: Admin) => {
  passwordForm.id = admin.id
  passwordForm.password = ''
  passwordForm.confirmPassword = ''
  showPasswordModal.value = true
}

// 删除管理员
const deleteAdmin = async (admin: Admin) => {
  if (!confirm(`确定要删除管理员 "${admin.username}" 吗？`)) {
    return
  }

  try {
    const { data: res } = await adminApi.deleteAdmin(admin.id)

    if (res.code === 200) {
      ToastAlert.success({
        title: '删除成功',
        message: '管理员已删除'
      })
      getAdminList()
    } else {
      ToastAlert.error({
        title: '删除失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('删除管理员失败:', error)
    ToastAlert.error({
      title: '删除失败',
      message: '网络异常，请重试'
    })
  }
}

// 提交表单
const submitForm = async () => {
  try {
    const isAdd = showAddModal.value
    const apiMethod = isAdd ? adminApi.addAdmin : adminApi.updateAdmin

    const { data: res } = await apiMethod(adminForm)

    if (res.code === 200) {
      ToastAlert.success({
        title: '操作成功',
        message: `管理员已${isAdd ? '添加' : '更新'}`
      })
      closeModal()
      getAdminList()
    } else {
      ToastAlert.error({
        title: '操作失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('提交表单失败:', error)
    ToastAlert.error({
      title: '操作失败',
      message: '网络异常，请重试'
    })
  }
}

// 提交密码重置
const submitPasswordReset = async () => {
  if (passwordForm.password !== passwordForm.confirmPassword) {
    ToastAlert.error({
      title: '密码不匹配',
      message: '两次输入的密码不一致'
    })
    return
  }

  try {
    const { data: res } = await adminApi.resetAdminPassword({
      id: passwordForm.id,
      password: passwordForm.password
    })

    if (res.code === 200) {
      ToastAlert.success({
        title: '重置成功',
        message: '密码已重置'
      })
      showPasswordModal.value = false
    } else {
      ToastAlert.error({
        title: '重置失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('重置密码失败:', error)
    ToastAlert.error({
      title: '重置失败',
      message: '网络异常，请重试'
    })
  }
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  Object.assign(adminForm, {
    id: undefined,
    username: '',
    password: '',
    phone: '',
    email: '',
    deptId: 0,
    postId: 0,
    status: '1',
    remark: ''
  })
}

// 分页操作
const prevPage = () => {
  if (pagination.page > 1) {
    pagination.page--
    getAdminList()
  }
}

const nextPage = () => {
  if (pagination.page < pagination.pages) {
    pagination.page++
    getAdminList()
  }
}

const goToPage = (page: number | string) => {
  if (typeof page === 'number' && page !== pagination.page) {
    pagination.page = page
    getAdminList()
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getAdminList()
  getDeptList()
  getPostList()
})
</script>
