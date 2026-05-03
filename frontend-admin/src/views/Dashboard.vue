<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409EFF;">
              <el-icon size="32"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">注册用户</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67C23A;">
              <el-icon size="32"><Ticket /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.tickets }}</div>
              <div class="stat-label">车次数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #E6A23C;">
              <el-icon size="32"><Van /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.vehicles }}</div>
              <div class="stat-label">车辆数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #F56C6C;">
              <el-icon size="32"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.violations }}</div>
              <div class="stat-label">待处理违规</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="charts-row">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>交通流量趋势</span>
          </template>
          <div ref="flowChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>状态分布</span>
          </template>
          <div ref="pieChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="list-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最新订单</span>
              <el-button type="primary" link @click="router.push('/bookings')">查看更多</el-button>
            </div>
          </template>
          <el-table :data="recentBookings" style="width: 100%">
            <el-table-column prop="order_no" label="订单号" width="180" />
            <el-table-column prop="passenger_name" label="乘客" />
            <el-table-column prop="total_amount" label="金额" width="100">
              <template #default="scope">
                ¥{{ scope.row.total_amount }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最新反馈</span>
              <el-button type="primary" link @click="router.push('/feedbacks')">查看更多</el-button>
            </div>
          </template>
          <el-table :data="recentFeedbacks" style="width: 100%">
            <el-table-column prop="title" label="标题" show-overflow-tooltip />
            <el-table-column prop="feedback_type" label="类型" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
                  {{ scope.row.status === 1 ? '已回复' : '待处理' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { getUsers, getBookings, getFeedbacks, getVehicles, getTickets, getViolations } from '@/api/admin'

const router = useRouter()

const flowChartRef = ref(null)
const pieChartRef = ref(null)
let flowChart = null
let pieChart = null

const stats = reactive({
  users: 0,
  tickets: 0,
  vehicles: 0,
  violations: 0
})

const recentBookings = ref([])
const recentFeedbacks = ref([])

const getStatusType = (status) => {
  const types = { 0: 'warning', 1: 'success', 2: 'info' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 0: '待支付', 1: '已完成', 2: '已取消' }
  return texts[status] || '未知'
}

const initCharts = () => {
  if (flowChartRef.value) {
    flowChart = echarts.init(flowChartRef.value)
    const flowOption = {
      tooltip: { trigger: 'axis' },
      legend: { data: ['车流量', '平均速度'] },
      xAxis: {
        type: 'category',
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      },
      yAxis: [
        { type: 'value', name: '车流量' },
        { type: 'value', name: '速度(km/h)' }
      ],
      series: [
        {
          name: '车流量',
          type: 'bar',
          data: [1200, 1500, 1800, 1600, 2000, 1400, 1000],
          itemStyle: { color: '#409EFF' }
        },
        {
          name: '平均速度',
          type: 'line',
          yAxisIndex: 1,
          data: [45, 42, 38, 40, 35, 48, 50],
          itemStyle: { color: '#67C23A' }
        }
      ]
    }
    flowChart.setOption(flowOption)
  }

  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    const pieOption = {
      tooltip: { trigger: 'item' },
      legend: {
        orient: 'vertical',
        left: 'left'
      },
      series: [
        {
          name: '订单状态',
          type: 'pie',
          radius: '60%',
          data: [
            { value: 10, name: '待支付', itemStyle: { color: '#E6A23C' } },
            { value: 85, name: '已完成', itemStyle: { color: '#67C23A' } },
            { value: 5, name: '已取消', itemStyle: { color: '#909399' } }
          ],
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    }
    pieChart.setOption(pieOption)
  }
}

const loadData = async () => {
  try {
    const [usersRes, bookingsRes, feedbacksRes, vehiclesRes, ticketsRes, violationsRes] = await Promise.all([
      getUsers({ page: 1, page_size: 1 }),
      getBookings({ page: 1, page_size: 5 }),
      getFeedbacks({ page: 1, page_size: 5 }),
      getVehicles({ page: 1, page_size: 1 }),
      getTickets({ page: 1, page_size: 1 }),
      getViolations({ page: 1, page_size: 1, status: 0 })
    ])

    stats.users = usersRes.data?.total || 0
    stats.vehicles = vehiclesRes.data?.total || 0
    stats.tickets = ticketsRes.data?.total || 0
    stats.violations = violationsRes.data?.total || 0

    recentBookings.value = bookingsRes.data?.list || []
    recentFeedbacks.value = feedbacksRes.data?.list || []
  } catch (e) {
    console.error('加载数据失败', e)
  }
}

const handleResize = () => {
  flowChart?.resize()
  pieChart?.resize()
}

onMounted(() => {
  initCharts()
  loadData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  flowChart?.dispose()
  pieChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard {
  min-height: 100%;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  cursor: pointer;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.charts-row {
  margin-bottom: 20px;
}

.chart-container {
  height: 300px;
}

.list-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
