<template>
  <div class="home">
    <el-header class="header">
      <div class="container">
        <div class="logo">
          <el-icon size="24" color="#409EFF"><MapLocation /></el-icon>
          <span class="logo-text">智能交通管理系统</span>
        </div>
        <el-menu :default-active="activeMenu" mode="horizontal" background-color="transparent" text-color="#333" active-text-color="#409EFF" router>
          <el-menu-item index="/">首页</el-menu-item>
          <el-menu-item index="/routes">交通路线</el-menu-item>
          <el-menu-item index="/tickets">交通购票</el-menu-item>
          <el-menu-item index="/news">资讯</el-menu-item>
          <el-menu-item index="/announcements">公告</el-menu-item>
          <el-menu-item index="/road-conditions">路况查询</el-menu-item>
        </el-menu>
        <div class="user-actions">
          <template v-if="userStore.isLoggedIn">
            <el-dropdown @command="handleCommand">
              <span class="username">
                <el-icon><User /></el-icon>
                {{ userStore.userInfo.username || '用户' }}
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                  <el-dropdown-item command="bookings">购票管理</el-dropdown-item>
                  <el-dropdown-item command="violations">违规管理</el-dropdown-item>
                  <el-dropdown-item command="feedbacks">投诉反馈</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" link @click="router.push('/login')">登录</el-button>
            <el-button type="primary" @click="router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>

    <el-carousel :interval="4000" height="300px" class="carousel">
      <el-carousel-item v-for="item in carousels" :key="item.id">
        <div class="carousel-item" :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
          <div class="carousel-content">
            <h2>{{ item.title }}</h2>
            <p>{{ item.link || '点击查看详情' }}</p>
          </div>
        </div>
      </el-carousel-item>
      <el-carousel-item v-if="carousels.length === 0">
        <div class="carousel-item" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
          <div class="carousel-content">
            <h2>智能交通管理系统</h2>
            <p>便捷出行，智能导航</p>
          </div>
        </div>
      </el-carousel-item>
    </el-carousel>

    <div class="main-content">
      <div class="container">
        <div class="quick-access">
          <h3>快捷入口</h3>
          <div class="access-grid">
            <div class="access-item" @click="router.push('/routes')">
              <el-icon size="40" color="#409EFF"><Guide /></el-icon>
              <span>交通路线</span>
            </div>
            <div class="access-item" @click="router.push('/tickets')">
              <el-icon size="40" color="#67C23A"><Ticket /></el-icon>
              <span>交通购票</span>
            </div>
            <div class="access-item" @click="router.push('/vehicle-monitor')">
              <el-icon size="40" color="#E6A23C"><Monitor /></el-icon>
              <span>车辆监控</span>
            </div>
            <div class="access-item" @click="router.push('/road-conditions')">
              <el-icon size="40" color="#F56C6C"><Warning /></el-icon>
              <span>路况查询</span>
            </div>
            <div class="access-item" @click="router.push('/violations')">
              <el-icon size="40" color="#909399"><Document /></el-icon>
              <span>违规查询</span>
            </div>
            <div class="access-item" @click="router.push('/news')">
              <el-icon size="40" color="#409EFF"><Reading /></el-icon>
              <span>交通资讯</span>
            </div>
          </div>
        </div>

        <div class="content-row">
          <div class="news-section">
            <div class="section-header">
              <h3>最新资讯</h3>
              <el-button type="primary" link @click="router.push('/news')">查看更多</el-button>
            </div>
            <div class="news-list" v-loading="newsLoading">
              <div class="news-item" v-for="item in newsList" :key="item.id" @click="router.push(`/news/${item.id}`)">
                <h4>{{ item.title }}</h4>
                <p class="summary">{{ item.summary || item.content?.substring(0, 100) + '...' }}</p>
                <span class="date">{{ formatDate(item.created_at) }}</span>
              </div>
              <el-empty v-if="newsList.length === 0 && !newsLoading" description="暂无资讯" />
            </div>
          </div>

          <div class="announcements-section">
            <div class="section-header">
              <h3>交通公告</h3>
              <el-button type="primary" link @click="router.push('/announcements')">查看更多</el-button>
            </div>
            <div class="announcement-list" v-loading="announcementsLoading">
              <div class="announcement-item" v-for="item in announcements" :key="item.id">
                <el-tag :type="item.priority > 0 ? 'danger' : 'primary'" size="small">
                  {{ item.priority > 0 ? '重要' : '公告' }}
                </el-tag>
                <h4>{{ item.title }}</h4>
                <span class="date">{{ formatDate(item.publish_time || item.created_at) }}</span>
              </div>
              <el-empty v-if="announcements.length === 0 && !announcementsLoading" description="暂无公告" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-footer class="footer">
      <div class="container">
        <p>© 2024 智能交通管理系统 版权所有</p>
      </div>
    </el-footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getCarousels, getNews, getAnnouncements } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const carousels = ref([])
const newsList = ref([])
const announcements = ref([])
const newsLoading = ref(false)
const announcementsLoading = ref(false)

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const fetchCarousels = async () => {
  try {
    const res = await getCarousels()
    if (res.code === 200) {
      carousels.value = res.data || []
    }
  } catch (e) {
    console.error('获取轮播图失败', e)
  }
}

const fetchNews = async () => {
  newsLoading.value = true
  try {
    const res = await getNews({ page: 1, page_size: 5 })
    if (res.code === 200) {
      newsList.value = res.data.list || []
    }
  } catch (e) {
    console.error('获取资讯失败', e)
  } finally {
    newsLoading.value = false
  }
}

const fetchAnnouncements = async () => {
  announcementsLoading.value = true
  try {
    const res = await getAnnouncements({ page: 1, page_size: 5 })
    if (res.code === 200) {
      announcements.value = res.data.list || []
    }
  } catch (e) {
    console.error('获取公告失败', e)
  } finally {
    announcementsLoading.value = false
  }
}

const handleCommand = async (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'bookings':
      router.push('/bookings')
      break
    case 'violations':
      router.push('/violations')
      break
    case 'feedbacks':
      router.push('/feedbacks')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        userStore.logout()
        ElMessage.success('已退出登录')
        router.push('/')
      } catch (e) {}
      break
  }
}

onMounted(() => {
  fetchCarousels()
  fetchNews()
  fetchAnnouncements()
})
</script>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  height: 60px;
  line-height: 60px;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.username {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  color: #409EFF;
}

.carousel {
  width: 100%;
}

.carousel-item {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.carousel-content {
  text-align: center;
}

.carousel-content h2 {
  font-size: 36px;
  margin-bottom: 10px;
}

.carousel-content p {
  font-size: 18px;
  opacity: 0.9;
}

.main-content {
  flex: 1;
  padding: 30px 0;
  background: #f5f7fa;
}

.main-content .container {
  flex-direction: column;
  align-items: stretch;
}

.quick-access {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.quick-access h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  color: #333;
}

.access-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 20px;
}

.access-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.access-item:hover {
  background: #f5f7fa;
  transform: translateY(-2px);
}

.content-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
}

.news-section, .announcements-section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.news-item, .announcement-item {
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s;
}

.news-item:last-child, .announcement-item:last-child {
  border-bottom: none;
}

.news-item:hover, .announcement-item:hover {
  background: #f9f9f9;
  padding-left: 10px;
  padding-right: 10px;
  margin-left: -10px;
  margin-right: -10px;
}

.news-item h4, .announcement-item h4 {
  margin: 0 0 8px 0;
  font-size: 15px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.news-item .summary {
  margin: 0 0 8px 0;
  font-size: 13px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.announcement-item h4 {
  margin-top: 8px;
}

.date {
  font-size: 12px;
  color: #c0c4cc;
}

.footer {
  background: #303133;
  color: #fff;
  text-align: center;
  padding: 20px 0;
}

.footer p {
  margin: 0;
}
</style>
