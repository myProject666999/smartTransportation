<template>
  <div class="bookings-page">
    <el-header class="header">
      <div class="container">
        <div class="logo" @click="router.push('/')">
          <el-icon size="24" color="#409EFF"><MapLocation /></el-icon>
          <span class="logo-text">智能交通管理系统</span>
        </div>
        <div class="nav-actions">
          <el-button link @click="router.push('/tickets')">购票</el-button>
          <el-button link @click="router.push('/')">返回首页</el-button>
        </div>
      </div>
    </el-header>

    <div class="main-content">
      <div class="container">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>购票管理</span>
              <div class="search-form">
                <el-select v-model="searchForm.status" placeholder="订单状态" clearable style="width: 120px" @change="loadBookings">
                  <el-option label="待支付" :value="0" />
                  <el-option label="已支付" :value="1" />
                  <el-option label="已取消" :value="2" />
                  <el-option label="已使用" :value="3" />
                </el-select>
              </div>
            </div>
          </template>

          <el-table :data="bookings" style="width: 100%" v-loading="loading" stripe>
            <el-table-column prop="booking_no" label="订单号" width="180" />
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
            <el-table-column prop="travel_date" label="乘车日期" width="120">
              <template #default="scope">
                {{ scope.row.travel_date || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="下单时间" width="160">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <div class="action-btns">
                  <el-button
                    v-if="scope.row.status === 0"
                    type="primary"
                    link
                    size="small"
                    @click="handlePay(scope.row)"
                  >
                    去支付
                  </el-button>
                  <el-button
                    v-if="scope.row.status === 0"
                    type="danger"
                    link
                    size="small"
                    @click="handleCancel(scope.row)"
                  >
                    取消订单
                  </el-button>
                  <el-button
                    v-if="scope.row.status === 1"
                    type="primary"
                    link
                    size="small"
                    @click="handleView(scope.row)"
                  >
                    查看详情
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadBookings"
            @current-change="loadBookings"
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
import { getBookings, payBooking, cancelBooking } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const bookings = ref([])
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

const getStatusType = (status) => {
  const types = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'primary'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    0: '待支付',
    1: '已支付',
    2: '已取消',
    3: '已使用'
  }
  return texts[status] || '未知'
}

const loadBookings = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getBookings(params)
    if (res.code === 200) {
      bookings.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载订单失败', e)
  } finally {
    loading.value = false
  }
}

const handlePay = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要支付该订单吗？\n订单号：${row.booking_no}\n金额：¥${row.total_price}`,
      '确认支付',
      {
        confirmButtonText: '确认支付',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await payBooking({
      booking_id: row.id
    })
    if (res.code === 200) {
      ElMessage.success('支付成功！')
      loadBookings()
    }
  } catch (e) {
    if (e !== 'cancel') {
      console.error('支付失败', e)
    }
  }
}

const handleCancel = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认要取消该订单吗？\n订单号：${row.booking_no}`,
      '确认取消',
      {
        confirmButtonText: '确认取消',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await cancelBooking(row.id)
    if (res.code === 200) {
      ElMessage.success('订单已取消！')
      loadBookings()
    }
  } catch (e) {
    if (e !== 'cancel') {
      console.error('取消失败', e)
    }
  }
}

const handleView = (row) => {
  ElMessage.info(`订单详情：\n订单号：${row.booking_no}\n金额：¥${row.total_price}\n状态：${getStatusText(row.status)}`)
}

onMounted(() => {
  loadBookings()
})
</script>

<style scoped>
.bookings-page {
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

.action-btns {
  display: flex;
  gap: 5px;
}
</style>
