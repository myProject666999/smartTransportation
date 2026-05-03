<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <el-icon size="28" color="#409EFF"><MapLocation /></el-icon>
        <span class="logo-text">交通管理后台</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
        unique-opened
      >
        <el-menu-item index="/dashboard">
          <template #title>
            <el-icon><DataAnalysis /></el-icon>
            <span>控制台</span>
          </template>
        </el-menu-item>
        
        <el-sub-menu index="system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/admins">
            <el-icon><Avatar /></el-icon>
            <span>管理员管理</span>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="traffic">
          <template #title>
            <el-icon><Van /></el-icon>
            <span>交通管理</span>
          </template>
          <el-menu-item index="/vehicles">
            <el-icon><Van /></el-icon>
            <span>车辆管理</span>
          </el-menu-item>
          <el-menu-item index="/vehicle-monitors">
            <el-icon><Monitor /></el-icon>
            <span>车辆监控</span>
          </el-menu-item>
          <el-menu-item index="/road-conditions">
            <el-icon><MapLocation /></el-icon>
            <span>路况管理</span>
          </el-menu-item>
          <el-menu-item index="/routes">
            <el-icon><Guide /></el-icon>
            <span>路线管理</span>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="business">
          <template #title>
            <el-icon><OfficeBuilding /></el-icon>
            <span>业务管理</span>
          </template>
          <el-menu-item index="/tickets">
            <el-icon><Ticket /></el-icon>
            <span>票务管理</span>
          </el-menu-item>
          <el-menu-item index="/bookings">
            <el-icon><Document /></el-icon>
            <span>购票管理</span>
          </el-menu-item>
          <el-menu-item index="/violations">
            <el-icon><Warning /></el-icon>
            <span>违规管理</span>
          </el-menu-item>
          <el-menu-item index="/feedbacks">
            <el-icon><ChatDotRound /></el-icon>
            <span>投诉反馈</span>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="monitor">
          <template #title>
            <el-icon><Monitor /></el-icon>
            <span>监控管理</span>
          </template>
          <el-menu-item index="/traffic-flows">
            <el-icon><TrendCharts /></el-icon>
            <span>交通流量</span>
          </el-menu-item>
          <el-menu-item index="/signal-lights">
            <el-icon><Lightning /></el-icon>
            <span>信号灯控制</span>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="content">
          <template #title>
            <el-icon><Reading /></el-icon>
            <span>内容管理</span>
          </template>
          <el-menu-item index="/news">
            <el-icon><Reading /></el-icon>
            <span>资讯管理</span>
          </el-menu-item>
          <el-menu-item index="/announcements">
            <el-icon><Bell /></el-icon>
            <span>公告管理</span>
          </el-menu-item>
          <el-menu-item index="/carousels">
            <el-icon><Picture /></el-icon>
            <span>轮播图管理</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute.meta?.title">{{ currentRoute.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="username">
              <el-icon><User /></el-icon>
              {{ adminStore.adminInfo.username || '管理员' }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="change-password">修改密码</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminStore()

const activeMenu = computed(() => route.path)
const currentRoute = computed(() => route)

const handleCommand = async (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'change-password':
      router.push('/change-password')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        adminStore.logout()
        ElMessage.success('已退出登录')
        router.push('/login')
      } catch (e) {}
      break
  }
}

onMounted(() => {
  if (adminStore.isLoggedIn) {
    adminStore.fetchProfile()
  }
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #304156;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid #3a4a5e;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: bold;
}

.header {
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.username {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  color: #409EFF;
}

.main {
  background: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
