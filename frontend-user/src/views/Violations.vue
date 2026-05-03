<template>
  <div class="violations-page">
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
              <span>违规管理</span>
              <div class="search-form">
                <el-select v-model="searchForm.status" placeholder="处理状态" clearable style="width: 120px" @change="loadViolations">
                  <el-option label="待处理" :value="0" />
                  <el-option label="已处理" :value="1" />
                </el-select>
              </div>
            </div>
          </template>

          <el-table :data="violations" style="width: 100%" v-loading="loading" stripe>
            <el-table-column prop="plate_number" label="车牌号" width="120" />
            <el-table-column prop="violation_type" label="违规类型" min-width="150" />
            <el-table-column prop="location" label="违规地点" min-width="200" />
            <el-table-column prop="violation_time" label="违规时间" width="160">
              <template #default="scope">
                {{ formatDate(scope.row.violation_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="fine_amount" label="罚款金额" width="100">
              <template #default="scope">
                <span style="color: #F56C6C; font-weight: bold;">¥{{ scope.row.fine_amount }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="points" label="扣分" width="80">
              <template #default="scope">
                <el-tag :type="scope.row.points > 3 ? 'danger' : 'warning'" size="small">
                  {{ scope.row.points }}分
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                  {{ scope.row.status === 1 ? '已处理' : '待处理' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="scope">
                <el-button
                  v-if="scope.row.status === 0"
                  type="primary"
                  link
                  @click="handlePay(scope.row)"
                >
                  处理缴费
                </el-button>
                <span v-else style="color: #909399;">-</span>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadViolations"
            @current-change="loadViolations"
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
import { getViolations, payViolation } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const violations = ref([])
const searchForm = reactive({
  status: ''
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

const loadViolations = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getViolations(params)
    if (res.code === 200) {
      violations.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载违规记录失败', e)
  } finally {
    loading.value = false
  }
}

const handlePay = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要处理该违规记录吗？\n罚款金额：¥${row.fine_amount}\n扣分：${row.points}分`,
      '确认处理',
      {
        confirmButtonText: '确认缴费',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await payViolation(row.id)
    if (res.code === 200) {
      ElMessage.success('处理成功！')
      loadViolations()
    }
  } catch (e) {
    if (e !== 'cancel') {
      console.error('处理失败', e)
    }
  }
}

onMounted(() => {
  loadViolations()
})
</script>

<style scoped>
.violations-page {
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
