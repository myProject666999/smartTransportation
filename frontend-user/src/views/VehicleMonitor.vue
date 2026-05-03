<template>
  <div class="vehicle-monitor-page">
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
              <span>车辆监控查询</span>
              <div class="search-form">
                <el-input
                  v-model="searchForm.plate_number"
                  placeholder="请输入车牌号查询"
                  style="width: 250px"
                  clearable
                  @keyup.enter="searchByPlate"
                >
                  <template #append>
                    <el-button :icon="Search" @click="searchByPlate">查询</el-button>
                  </template>
                </el-input>
                <el-button @click="loadVehicleMonitors">刷新</el-button>
              </div>
            </div>
          </template>

          <el-table :data="vehicleMonitors" style="width: 100%" v-loading="loading" stripe>
            <el-table-column prop="plate_number" label="车牌号" width="120" />
            <el-table-column prop="location" label="当前位置" min-width="200" />
            <el-table-column prop="longitude" label="经度" width="120">
              <template #default="scope">
                {{ scope.row.longitude?.toFixed(6) || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="latitude" label="纬度" width="120">
              <template #default="scope">
                {{ scope.row.latitude?.toFixed(6) || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="speed" label="速度(km/h)" width="120">
              <template #default="scope">
                <el-tag :type="scope.row.speed > 80 ? 'danger' : 'success'">
                  {{ scope.row.speed || 0 }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="direction" label="方向" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
                  {{ scope.row.status === 1 ? '在线' : '离线' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="monitor_time" label="监控时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.monitor_time) }}
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadVehicleMonitors"
            @current-change="loadVehicleMonitors"
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
import { getVehicleMonitors, getVehicleMonitorByPlate } from '@/api/user'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const vehicleMonitors = ref([])
const searchForm = reactive({
  plate_number: ''
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

const loadVehicleMonitors = async () => {
  loading.value = true
  try {
    const res = await getVehicleMonitors({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.code === 200) {
      vehicleMonitors.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载车辆监控失败', e)
  } finally {
    loading.value = false
  }
}

const searchByPlate = async () => {
  if (!searchForm.plate_number.trim()) {
    pagination.page = 1
    loadVehicleMonitors()
    return
  }

  loading.value = true
  try {
    const res = await getVehicleMonitorByPlate(searchForm.plate_number)
    if (res.code === 200) {
      vehicleMonitors.value = res.data ? [res.data] : []
      pagination.total = res.data ? 1 : 0
    }
  } catch (e) {
    ElMessage.warning('未找到该车辆的监控信息')
    vehicleMonitors.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadVehicleMonitors()
})
</script>

<style scoped>
.vehicle-monitor-page {
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
