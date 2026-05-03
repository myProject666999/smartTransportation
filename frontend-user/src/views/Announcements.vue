<template>
  <div class="announcements-page">
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
            <span>交通公告</span>
          </template>

          <div class="announcements-list" v-loading="loading">
            <div class="announcement-item" v-for="item in announcements" :key="item.id">
              <div class="announcement-header">
                <el-tag :type="item.priority > 0 ? 'danger' : 'primary'" size="small">
                  {{ item.priority > 0 ? '重要公告' : '普通公告' }}
                </el-tag>
                <h3>{{ item.title }}</h3>
              </div>
              <div class="announcement-content" v-html="item.content"></div>
              <div class="announcement-footer">
                <span>发布时间：{{ formatDate(item.publish_time || item.created_at) }}</span>
                <span v-if="item.expire_time">有效期至：{{ formatDate(item.expire_time) }}</span>
              </div>
            </div>
            <el-empty v-if="announcements.length === 0 && !loading" description="暂无公告" />
          </div>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadAnnouncements"
            @current-change="loadAnnouncements"
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
import { getAnnouncements } from '@/api/user'

const router = useRouter()

const loading = ref(false)
const announcements = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const loadAnnouncements = async () => {
  loading.value = true
  try {
    const res = await getAnnouncements({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.code === 200) {
      announcements.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载公告失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAnnouncements()
})
</script>

<style scoped>
.announcements-page {
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

.announcements-list {
  max-width: 900px;
  margin: 0 auto;
}

.announcement-item {
  padding: 20px;
  margin-bottom: 20px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  background: #fff;
}

.announcement-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.announcement-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
  flex: 1;
}

.announcement-content {
  font-size: 15px;
  line-height: 1.8;
  color: #666;
  margin-bottom: 15px;
}

.announcement-content :deep(p) {
  margin-bottom: 10px;
}

.announcement-footer {
  display: flex;
  gap: 30px;
  font-size: 13px;
  color: #909399;
  padding-top: 15px;
  border-top: 1px dashed #ebeef5;
}
</style>
