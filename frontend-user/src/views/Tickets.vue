<template>
  <div class="tickets-page">
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
            <span>交通购票</span>
          </template>

          <el-table :data="tickets" style="width: 100%" v-loading="loading">
            <el-table-column prop="route_name" label="路线名称" />
            <el-table-column prop="start_station" label="起点站" />
            <el-table-column prop="end_station" label="终点站" />
            <el-table-column prop="departure_time" label="发车时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.departure_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="arrival_time" label="到达时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.arrival_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="price" label="票价" width="100">
              <template #default="scope">
                ¥{{ scope.row.price }}
              </template>
            </el-table-column>
            <el-table-column prop="available_seats" label="余票" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.available_seats > 0 ? 'success' : 'danger'">
                  {{ scope.row.available_seats }}张
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="scope">
                <el-button type="primary" size="small" :disabled="scope.row.available_seats <= 0" @click="handleBook(scope.row)">
                  购票
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadTickets"
            @current-change="loadTickets"
            style="margin-top: 20px; justify-content: flex-end"
          />
        </el-card>
      </div>
    </div>

    <el-dialog v-model="bookDialogVisible" title="购买车票" width="500px">
      <el-form ref="bookFormRef" :model="bookForm" :rules="bookRules" label-width="100px">
        <el-form-item label="路线">
          <el-input :value="currentTicket?.route_name" disabled />
        </el-form-item>
        <el-form-item label="票价">
          <el-input :value="'¥' + (currentTicket?.price || 0)" disabled />
        </el-form-item>
        <el-form-item label="购票数量" prop="quantity">
          <el-input-number v-model="bookForm.quantity" :min="1" :max="currentTicket?.available_seats || 1" />
        </el-form-item>
        <el-form-item label="乘客姓名" prop="passenger_name">
          <el-input v-model="bookForm.passenger_name" placeholder="请输入乘客姓名" />
        </el-form-item>
        <el-form-item label="身份证号" prop="passenger_id_card">
          <el-input v-model="bookForm.passenger_id_card" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="联系电话" prop="passenger_phone">
          <el-input v-model="bookForm.passenger_phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="总价">
          <el-text type="danger" size="large">
            ¥{{ (currentTicket?.price || 0) * bookForm.quantity }}
          </el-text>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bookDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="booking" @click="submitBook">确认购买</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getTickets, bookTicket, payBooking } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const tickets = ref([])
const bookDialogVisible = ref(false)
const bookFormRef = ref(null)
const currentTicket = ref(null)
const booking = ref(false)

const bookForm = reactive({
  quantity: 1,
  passenger_name: '',
  passenger_id_card: '',
  passenger_phone: ''
})

const bookRules = {
  quantity: [{ required: true, message: '请选择购票数量', trigger: 'change' }],
  passenger_name: [{ required: true, message: '请输入乘客姓名', trigger: 'blur' }],
  passenger_id_card: [{ required: true, message: '请输入身份证号', trigger: 'blur' }],
  passenger_phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }]
}

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

const loadTickets = async () => {
  loading.value = true
  try {
    const res = await getTickets({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.code === 200) {
      tickets.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载车票失败', e)
  } finally {
    loading.value = false
  }
}

const handleBook = (ticket) => {
  currentTicket.value = ticket
  bookForm.quantity = 1
  bookForm.passenger_name = ''
  bookForm.passenger_id_card = ''
  bookForm.passenger_phone = ''
  bookDialogVisible.value = true
}

const submitBook = async () => {
  const valid = await bookFormRef.value.validate().catch(() => false)
  if (!valid) return

  booking.value = true
  try {
    const res = await bookTicket({
      ticket_id: currentTicket.value.id,
      passenger_name: bookForm.passenger_name,
      passenger_id_card: bookForm.passenger_id_card,
      passenger_phone: bookForm.passenger_phone,
      quantity: bookForm.quantity
    })

    if (res.code === 200) {
      try {
        await ElMessageBox.confirm(`购票成功！订单号：${res.data.order_no}，总金额：¥${res.data.total_amount}，是否立即支付？`, '提示', {
          confirmButtonText: '立即支付',
          cancelButtonText: '稍后支付',
          type: 'success'
        })

        const payRes = await payBooking({ booking_id: res.data.booking_id })
        if (payRes.code === 200) {
          ElMessage.success('支付成功！')
          bookDialogVisible.value = false
          loadTickets()
        }
      } catch (e) {
        if (e !== 'cancel') {
          throw e
        }
        ElMessage.info('您可以在购票管理中完成支付')
        bookDialogVisible.value = false
      }
    }
  } catch (e) {
    console.error('购票失败', e)
  } finally {
    booking.value = false
  }
}

onMounted(() => {
  loadTickets()
})
</script>

<style scoped>
.tickets-page {
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
</style>
