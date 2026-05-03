import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册', guest: true }
  },
  {
    path: '/routes',
    name: 'Routes',
    component: () => import('@/views/Routes.vue'),
    meta: { title: '交通路线' }
  },
  {
    path: '/tickets',
    name: 'Tickets',
    component: () => import('@/views/Tickets.vue'),
    meta: { title: '交通购票', requiresAuth: true }
  },
  {
    path: '/news',
    name: 'News',
    component: () => import('@/views/News.vue'),
    meta: { title: '资讯' }
  },
  {
    path: '/news/:id',
    name: 'NewsDetail',
    component: () => import('@/views/NewsDetail.vue'),
    meta: { title: '资讯详情' }
  },
  {
    path: '/announcements',
    name: 'Announcements',
    component: () => import('@/views/Announcements.vue'),
    meta: { title: '交通公告' }
  },
  {
    path: '/vehicle-monitor',
    name: 'VehicleMonitor',
    component: () => import('@/views/VehicleMonitor.vue'),
    meta: { title: '车辆监控查询', requiresAuth: true }
  },
  {
    path: '/road-conditions',
    name: 'RoadConditions',
    component: () => import('@/views/RoadConditions.vue'),
    meta: { title: '路况查询' }
  },
  {
    path: '/violations',
    name: 'Violations',
    component: () => import('@/views/Violations.vue'),
    meta: { title: '违规管理', requiresAuth: true }
  },
  {
    path: '/bookings',
    name: 'Bookings',
    component: () => import('@/views/Bookings.vue'),
    meta: { title: '购票管理', requiresAuth: true }
  },
  {
    path: '/feedbacks',
    name: 'Feedbacks',
    component: () => import('@/views/Feedbacks.vue'),
    meta: { title: '投诉反馈', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '个人信息', requiresAuth: true }
  },
  {
    path: '/change-password',
    name: 'ChangePassword',
    component: () => import('@/views/ChangePassword.vue'),
    meta: { title: '密码修改', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 智能交通管理系统` : '智能交通管理系统'
  
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.guest && userStore.isLoggedIn) {
    next({ name: 'Home' })
  } else {
    next()
  }
})

export default router
