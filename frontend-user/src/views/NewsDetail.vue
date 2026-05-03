<template>
  <div class="news-detail-page">
    <el-header class="header">
      <div class="container">
        <div class="logo" @click="router.push('/')">
          <el-icon size="24" color="#409EFF"><MapLocation /></el-icon>
          <span class="logo-text">智能交通管理系统</span>
        </div>
        <div class="nav-actions">
          <el-button link @click="router.push('/news')">返回列表</el-button>
          <el-button link @click="router.push('/')">返回首页</el-button>
        </div>
      </div>
    </el-header>

    <div class="main-content">
      <div class="container">
        <el-card v-loading="loading" class="detail-card">
          <template #header v-if="news">
            <div class="news-header">
              <h1>{{ news.title }}</h1>
              <div class="news-meta">
                <span><el-icon><User /></el-icon> 作者：{{ news.author || '系统' }}</span>
                <span><el-icon><View /></el-icon> 浏览：{{ news.view_count || 0 }}</span>
                <span><el-icon><Clock /></el-icon> 发布时间：{{ formatDate(news.created_at) }}</span>
              </div>
            </div>
          </template>

          <div v-if="news" class="news-content">
            <div v-if="news.image" class="news-image">
              <img :src="news.image" alt="资讯图片" />
            </div>
            <div v-html="news.content" class="content-text"></div>
          </div>

          <el-empty v-else description="资讯不存在或已被删除" />
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getNewsDetail } from '@/api/user'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const news = ref(null)

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const loadNewsDetail = async () => {
  const id = route.params.id
  if (!id) {
    router.push('/news')
    return
  }

  loading.value = true
  try {
    const res = await getNewsDetail(id)
    if (res.code === 200) {
      news.value = res.data
    }
  } catch (e) {
    console.error('加载资讯详情失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadNewsDetail()
})
</script>

<style scoped>
.news-detail-page {
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

.main-content {
  flex: 1;
  padding: 30px 0;
  background: #f5f7fa;
}

.main-content .container {
  flex-direction: column;
  align-items: stretch;
}

.detail-card {
  max-width: 900px;
  margin: 0 auto;
}

.news-header {
  text-align: center;
}

.news-header h1 {
  margin: 0 0 15px 0;
  font-size: 24px;
  color: #333;
  line-height: 1.5;
}

.news-meta {
  display: flex;
  justify-content: center;
  gap: 30px;
  font-size: 14px;
  color: #909399;
}

.news-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.news-image {
  text-align: center;
  margin-bottom: 20px;
}

.news-image img {
  max-width: 100%;
  border-radius: 8px;
}

.content-text {
  font-size: 16px;
  line-height: 2;
  color: #333;
}

.content-text :deep(p) {
  margin-bottom: 15px;
}

.content-text :deep(img) {
  max-width: 100%;
  border-radius: 8px;
}
</style>
