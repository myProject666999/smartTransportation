<template>
  <div class="news-page">
    <el-header class="header">
      <div class="container">
        <div class="logo" @click="router.push('/')">
          <el-icon size="24" color="#409EFF"><MapLocation /></el-icon>
          <span class="logo-text">智能交通管理系统</span>
        </div>
        <div class="nav-actions">
          <el-button link @click="router.push('/')">返回首页</el-button>
        </div>
      </div>
    </el-header>

    <div class="main-content">
      <div class="container">
        <el-card>
          <template #header>
            <span>交通资讯</span>
          </template>

          <div class="news-list" v-loading="loading">
            <div class="news-item" v-for="item in newsList" :key="item.id" @click="viewDetail(item.id)">
              <h3>{{ item.title }}</h3>
              <p class="summary">{{ item.summary || item.content?.substring(0, 150) + '...' }}</p>
              <div class="meta">
                <span class="author">作者：{{ item.author || '系统' }}</span>
                <span class="views">浏览：{{ item.view_count || 0 }}</span>
                <span class="date">{{ formatDate(item.created_at) }}</span>
              </div>
            </div>
            <el-empty v-if="newsList.length === 0 && !loading" description="暂无资讯" />
          </div>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadNews"
            @current-change="loadNews"
            style="margin-top: 20px; justify-content: flex-end"
          />
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getNews } from '@/api/user'

const router = useRouter()

const loading = ref(false)
const newsList = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const loadNews = async () => {
  loading.value = true
  try {
    const res = await getNews({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.code === 200) {
      newsList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载资讯失败', e)
  } finally {
    loading.value = false
  }
}

const viewDetail = (id) => {
  router.push(`/news/${id}`)
}

onMounted(() => {
  loadNews()
})
</script>

<style scoped>
.news-page {
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

.news-list {
  max-width: 900px;
}

.news-item {
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s;
}

.news-item:hover {
  padding-left: 10px;
  background: #f9f9f9;
}

.news-item:last-child {
  border-bottom: none;
}

.news-item h3 {
  margin: 0 0 10px 0;
  font-size: 18px;
  color: #333;
}

.news-item .summary {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}

.news-item .meta {
  display: flex;
  gap: 20px;
  font-size: 12px;
  color: #999;
}
</style>
