<template>
  <div class="routes-page">
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
            <div class="card-header">
              <span>交通路线查询</span>
              <div class="search-form">
                <el-input
                  v-model="searchForm.keyword"
                  placeholder="搜索路线名称、起点、终点"
                  style="width: 300px"
                  clearable
                  @clear="handleSearch"
                  @keyup.enter="handleSearch"
                >
                  <template #append>
                    <el-button :icon="Search" @click="handleSearch" />
                  </template>
                </el-input>
              </div>
            </div>
          </template>

          <el-table :data="routes" style="width: 100%" v-loading="loading">
            <el-table-column prop="route_no" label="路线编号" width="120" />
            <el-table-column prop="route_name" label="路线名称" />
            <el-table-column prop="start_station" label="起点站" width="150" />
            <el-table-column prop="end_station" label="终点站" width="150" />
            <el-table-column prop="distance" label="距离(km)" width="100">
              <template #default="scope">
                {{ scope.row.distance || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="duration" label="预计时长(分钟)" width="120">
              <template #default="scope">
                {{ scope.row.duration || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="price" label="票价(元)" width="100">
              <template #default="scope">
                ¥{{ scope.row.price || 0 }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button type="primary" link @click="viewDetail(scope.row)">查看详情</el-button>
                <el-button type="success" link @click="bookTicket(scope.row)">购票</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadRoutes"
            @current-change="loadRoutes"
            style="margin-top: 20px; justify-content: flex-end"
          />
        </el-card>
      </div>
    </div>

    <el-dialog v-model="detailVisible" title="路线详情" width="600px">
      <el-descriptions :column="1" border v-if="currentRoute">
        <el-descriptions-item label="路线编号">{{ currentRoute.route_no }}</el-descriptions-item>
        <el-descriptions-item label="路线名称">{{ currentRoute.route_name }}</el-descriptions-item>
        <el-descriptions-item label="起点站">{{ currentRoute.start_station }}</el-descriptions-item>
        <el-descriptions-item label="终点站">{{ currentRoute.end_station }}</el-descriptions-item>
        <el-descriptions-item label="途经站点">{{ currentRoute.stations || '暂无' }}</el-descriptions-item>
        <el-descriptions-item label="距离">{{ currentRoute.distance }} km</el-descriptions-item>
        <el-descriptions-item label="预计时长">{{ currentRoute.duration }} 分钟</el-descriptions-item>
        <el-descriptions-item label="票价">¥{{ currentRoute.price }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="primary" @click="bookTicket(currentRoute)">立即购票</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getRoutes } from '@/api/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const routes = ref([])
const detailVisible = ref(false)
const currentRoute = ref(null)

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadRoutes = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) {
      params.keyword = searchForm.keyword
    }

    const res = await getRoutes(params)
    if (res.code === 200) {
      routes.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载路线失败', e)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadRoutes()
}

const viewDetail = (row) => {
  currentRoute.value = row
  detailVisible.value = true
}

const bookTicket = (route) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  router.push('/tickets')
}

onMounted(() => {
  loadRoutes()
})
</script>

<style scoped>
.routes-page {
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  display: flex;
  gap: 10px;
}
</style>
