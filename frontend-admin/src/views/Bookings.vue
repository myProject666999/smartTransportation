<template>
  <div class="bookings-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>购票管理</h3>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-card>
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="订单号/用户名"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 120px">
              <el-option label="待支付" :value="0" />
              <el-option label="已支付" :value="1" />
              <el-option label="已取消" :value="2" />
              <el-option label="已使用" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="booking_no" label="订单号" width="180" />
          <el-table-column prop="username" label="用户名" width="100" />
          <el-table-column prop="ticket_name" label="车票名称" min-width="150">
            <template #default="scope">
              {{ scope.row.ticket?.name || scope.row.ticket_name || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="route_name" label="路线" min-width="150">
            <template #default="scope">
              {{ scope.row.ticket?.route_name || scope.row.route_name || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="quantity" label="数量" width="80" />
          <el-table-column prop="total_price" label="总价" width="100">
            <template #default="scope">
              <span style="color: #F56C6C; font-weight: bold;">¥{{ scope.row.total_price }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="下单时间" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleView(scope.row)">查看</el-button>
              <el-button
                v-if="scope.row.status === 0"
                type="danger"
                link
                size="small"
                @click="handleCancel(scope.row)"
              >
                取消订单
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
          style="margin-top: 20px; justify-content: flex-end"
        />
      </el-card>
    </div>

    <el-dialog v-model="detailDialogVisible" title="订单详情" width="500px">
      <el-descriptions :column="2" border v-if="currentItem">
        <el-descriptions-item label="订单号">{{ currentItem.booking_no }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentItem.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="车票名称">{{ currentItem.ticket?.name || currentItem.ticket_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="路线">{{ currentItem.ticket?.route_name || currentItem.route_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="数量">{{ currentItem.quantity }}</el-descriptions-item>
        <el-descriptions-item label="总价">
          <span style="color: #F56C6C; font-weight: bold;">¥{{ currentItem.total_price }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="状态" :span="2">
          <el-tag :type="getStatusType(currentItem.status)">{{ getStatusText(currentItem.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="乘车日期" :span="2">{{ currentItem.travel_date || '-' }}</el-descriptions-item>
        <el-descriptions-item label="下单时间" :span="2">{{ formatDate(currentItem.created_at) }}</el-descriptions-item>
        <el-descriptions-item v-if="currentItem.payment_time" label="支付时间" :span="2">{{ formatDate(currentItem.payment_time) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getBookings, cancelBooking } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const tableData = ref([])
const detailDialogVisible = ref(false)
const currentItem = ref({})

const searchForm = reactive({
  keyword: '',
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

const getStatusType = (status) => {
  const types = { 0: 'warning', 1: 'success', 2: 'info', 3: 'primary' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 0: '待支付', 1: '已支付', 2: '已取消', 3: '已使用' }
  return texts[status] || '未知'
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) {
      params.keyword = searchForm.keyword
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getBookings(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载订单列表失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const handleView = (row) => {
  currentItem.value = row
  detailDialogVisible.value = true
}

const handleCancel = (row) => {
  ElMessageBox.confirm(`确认要取消订单"${row.booking_no}"吗？`, '确认取消', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await cancelBooking(row.id)
      if (res.code === 200) {
        ElMessage.success('订单已取消')
        loadData()
      }
    } catch (e) {
      console.error('取消失败', e)
    }
  }).catch(() => {})
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.bookings-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h3 {
  margin: 0;
  color: #303133;
}

.search-form {
  margin-bottom: 20px;
}
</style>
