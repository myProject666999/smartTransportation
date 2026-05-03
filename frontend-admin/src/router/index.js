import { createRouter, createWebHistory } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import Layout from '@/views/Layout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '管理员登录', guest: true }
  },
  {
    path: '/',
    name: 'Layout',
    component: Layout,
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '控制台' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/Users.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'admins',
        name: 'Admins',
        component: () => import('@/views/Admins.vue'),
        meta: { title: '管理员管理' }
      },
      {
        path: 'vehicles',
        name: 'Vehicles',
        component: () => import('@/views/Vehicles.vue'),
        meta: { title: '车辆管理' }
      },
      {
        path: 'vehicle-monitors',
        name: 'VehicleMonitors',
        component: () => import('@/views/VehicleMonitors.vue'),
        meta: { title: '车辆监控管理' }
      },
      {
        path: 'road-conditions',
        name: 'RoadConditions',
        component: () => import('@/views/RoadConditions.vue'),
        meta: { title: '路况管理' }
      },
      {
        path: 'violations',
        name: 'Violations',
        component: () => import('@/views/Violations.vue'),
        meta: { title: '违规管理' }
      },
      {
        path: 'tickets',
        name: 'Tickets',
        component: () => import('@/views/Tickets.vue'),
        meta: { title: '票务管理' }
      },
      {
        path: 'bookings',
        name: 'Bookings',
        component: () => import('@/views/Bookings.vue'),
        meta: { title: '购票管理' }
      },
      {
        path: 'feedbacks',
        name: 'Feedbacks',
        component: () => import('@/views/Feedbacks.vue'),
        meta: { title: '投诉反馈管理' }
      },
      {
        path: 'traffic-flows',
        name: 'TrafficFlows',
        component: () => import('@/views/TrafficFlows.vue'),
        meta: { title: '交通流量管理' }
      },
      {
        path: 'signal-lights',
        name: 'SignalLights',
        component: () => import('@/views/SignalLights.vue'),
        meta: { title: '信号灯控制管理' }
      },
      {
        path: 'routes',
        name: 'Routes',
        component: () => import('@/views/Routes.vue'),
        meta: { title: '路线管理' }
      },
      {
        path: 'news',
        name: 'News',
        component: () => import('@/views/News.vue'),
        meta: { title: '资讯管理' }
      },
      {
        path: 'announcements',
        name: 'Announcements',
        component: () => import('@/views/Announcements.vue'),
        meta: { title: '公告管理' }
      },
      {
        path: 'carousels',
        name: 'Carousels',
        component: () => import('@/views/Carousels.vue'),
        meta: { title: '轮播图管理' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人信息' }
      },
      {
        path: 'change-password',
        name: 'ChangePassword',
        component: () => import('@/views/ChangePassword.vue'),
        meta: { title: '修改密码' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 智能交通管理后台` : '智能交通管理后台'
  
  const adminStore = useAdminStore()
  
  if (to.meta.requiresAuth && !adminStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.guest && adminStore.isLoggedIn) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
