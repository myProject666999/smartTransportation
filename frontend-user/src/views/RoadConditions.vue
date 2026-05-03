<template>
  <div class="road-conditions-page">
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
              <span>路况查询</span>
              <div class="search-form">
                <el-input
                  v-model="searchForm.keyword"
                  placeholder="搜索道路名称、位置"
                  style="width: 250px"
                  clearable
                  @keyup.enter="loadRoadConditions"
                >
                  <template #append>
                    <el-button :icon="Search" @click="loadRoadConditions">查询</el-button>
                  </template>
                </el-input>
                <el-select v-model="searchForm.condition" placeholder="路况状态" clearable style="width: 120px" @change="loadRoadConditions">
                  <el-option label="畅通" value="畅通" />
                  <el-option label="缓行" value="缓行" />
                  <el-option label="拥堵" value="拥堵" />
                  <el-option label="施工" value="施工" />
                  <el-option label="事故" value="事故" />
                </el-select>
              </div>
            </div>
          </template>

          <el-table :data="roadConditions" style="width: 100%" v-loading="loading" stripe>
            <el-table-column prop="road_name" label="道路名称" width="150" />
            <el-table-column prop="location" label="位置" min-width="200" />
            <el-table-column prop="condition" label="路况状态" width="120">
              <template #default="scope">
                <el-tag :type="getConditionType(scope.row.condition)">
                  {{ scope.row.condition || '未知' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="severity" label="严重程度" width="100">
              <template #default="scope">
                <el-tag :type="getSeverityType(scope.row.severity)" size="small">
                  {{ scope.row.severity || '一般' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
            <el-table-column prop="start_time" label="开始时间" width="160">
              <template #default="scope">
                {{ formatDate(scope.row.start_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="end_time" label="预计结束" width="160">
              <template #default="scope">
                {{ scope.row.end_time ? formatDate(scope.row.end_time) : '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'danger' : 'info'">
                  {{ scope.row.status === 1 ? '进行中' : '已结束' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadRoadConditions"
            @current-change="loadRoadConditions"
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
import { getRoadConditions } from '@/api/user'

const router = useRouter()

const loading = ref(false)
const roadConditions = ref([])
const searchForm = reactive({
  keyword: '',
  condition: ''
})

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

const getConditionType = (condition) => {
  const types = {
    '畅通': 'success',
    '缓行': 'warning',
    '拥堵': 'danger',
    '施工': 'info',
    '事故': 'danger'
  }
  return types[condition] || 'info'
}

const getSeverityType = (severity) => {
  const types = {
    '轻微': 'info',
    '一般': 'warning',
    '严重': 'danger'
  }
  return types[severity] || 'info'
}

const loadRoadConditions = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) {
      params.keyword = searchForm.keyword
    }
    if (searchForm.condition) {
      params.condition = searchForm.condition
    }

    const res = await getRoadConditions(params)
    if (res.code === 200) {
      roadConditions.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载路况失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRoadConditions()
})
</script>

<style scoped>
.road-conditions-page {
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
