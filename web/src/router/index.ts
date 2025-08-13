import { createRouter, createWebHistory } from 'vue-router'
import { AuthUtils, routeGuards } from '@/utils/auth'
import { authRoutes } from './auth'

// 用户端路由
const userRoutes = [
  //   {
  //   path: '/home',
  //   name: 'Home',
  //   component: () => import('../views/UsersViews/Home.vue'),
  //   meta: {
  //     title: 'Home',
  //   },
  // },
  {
    path: '/signin',
    name: 'Signin',
    component: () => import('../views/UsersViews/Signin.vue'),
    meta: {
      title: 'Signin',
    },
    children: [

    ]
  },
  {
    path: '/',
    name: 'Signup',
    component: () => import('../views/UsersViews/Signup.vue'),
    meta: {
      title: 'Signup',
    },
  },
  {
    path: '/error-404',
    name: '404 Error',
    component: () => import('../views/Errors/FourZeroFour.vue'),
    meta: {
      title: '404 Error',
    },
  },
]

// 运营端路由
const adminRoutes = [
  // 根路径重定向到登录页面
  {
    path: '/',
    redirect: '/adminLogin'
  },
  {
    path: '/dashboard',
    name: 'Ecommerce',
    component: () => import('../views/Ecommerce.vue'),
    meta: {
      title: 'eCommerce Dashboard',
      requiresAuth: true,
    },
  },
  {
    path: '/calendar',
    name: 'Calendar',
    component: () => import('../views/Others/Calendar.vue'),
    meta: {
      title: 'Calendar',
      requiresAuth: true,
    },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Others/UserProfile.vue'),
    meta: {
      title: 'Profile',
      requiresAuth: true,
    },
  },
  {
    path: '/form-elements',
    name: 'Form Elements',
    component: () => import('../views/Forms/FormElements.vue'),
    meta: {
      title: 'Form Elements',
      requiresAuth: true,
    },
  },
  {
    path: '/basic-tables',
    name: 'Basic Tables',
    component: () => import('../views/Tables/BasicTables.vue'),
    meta: {
      title: 'Basic Tables',
      requiresAuth: true,
    },
  },
  {
    path: '/line-chart',
    name: 'Line Chart',
    component: () => import('../views/Chart/LineChart/LineChart.vue'),
    meta: {
      title: 'Line Chart',
      requiresAuth: true,
    },
  },
  {
    path: '/bar-chart',
    name: 'Bar Chart',
    component: () => import('../views/Chart/BarChart/BarChart.vue'),
    meta: {
      title: 'Bar Chart',
      requiresAuth: true,
    },
  },
  {
    path: '/alerts',
    name: 'Alerts',
    component: () => import('../views/UiElements/Alerts.vue'),
    meta: {
      title: 'Alerts',
      requiresAuth: true,
    },
  },
  {
    path: '/avatars',
    name: 'Avatars',
    component: () => import('../views/UiElements/Avatars.vue'),
    meta: {
      title: 'Avatars',
      requiresAuth: true,
    },
  },
  {
    path: '/badge',
    name: 'Badge',
    component: () => import('../views/UiElements/Badges.vue'),
    meta: {
      title: 'Badge',
      requiresAuth: true,
    },
  },

  {
    path: '/buttons',
    name: 'Buttons',
    component: () => import('../views/UiElements/Buttons.vue'),
    meta: {
      title: 'Buttons',
      requiresAuth: true,
    },
  },

  {
    path: '/images',
    name: 'Images',
    component: () => import('../views/UiElements/Images.vue'),
    meta: {
      title: 'Images',
      requiresAuth: true,
    },
  },
  {
    path: '/videos',
    name: 'Videos',
    component: () => import('../views/UiElements/Videos.vue'),
    meta: {
      title: 'Videos',
      requiresAuth: true,
    },
  },
  {
    path: '/blank',
    name: 'Blank',
    component: () => import('../views/Pages/BlankPage.vue'),
    meta: {
      title: 'Blank',
      requiresAuth: true,
    },
  },
  {
    path: '/test-toast',
    name: 'TestToast',
    component: () => import('../views/TestToast.vue'),
    meta: {
      title: 'Toast Test',
      requiresAuth: true,
    },
  },
  {
    path: '/test-auth',
    name: 'TestAuthModule',
    component: () => import('../views/Auth/TestAuth.vue'),
    meta: {
      title: 'Auth模块测试',
      requiresAuth: false,
    },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(_to, _from, savedPosition) {
    return savedPosition || { left: 0, top: 0 }
  },
  routes: [
    // ...userRoutes,
    ...(authRoutes as any),
    ...adminRoutes,
    // 可以在这里添加更多路由组
  ],
})

export default router

router.beforeEach((to, _from, next) => {
  // 设置页面标题
  document.title = `Vue.js ${to.meta.title} | TailAdmin - Vue.js Tailwind CSS Dashboard Template`

  // 检查路由是否需要认证
  if (to.meta.requiresAuth) {
    routeGuards.requireAuth(to, _from, next)
  }
  // 检查已登录用户是否访问登录页
  else if (to.meta.redirectIfAuthenticated) {
    routeGuards.redirectIfAuthenticated(to, _from, next)
  }
  // 其他路由直接通过
  else {
    next()
  }
})