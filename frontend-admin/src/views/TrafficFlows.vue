<template>
  <div class="traffic-flows-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>交通流量管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加数据
          </el-button>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-card>
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="地点/道路名称"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="location" label="地点" min-width="150" />
          <el-table-column prop="road_name" label="道路名称" min-width="150" />
          <el-table-column prop="vehicle_count" label="车辆数" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.vehicle_count > 500 ? 'danger' : scope.row.vehicle_count > 200 ? 'warning' : 'success'" size="small">
                {{ scope.row.vehicle_count }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="flow_rate" label="流量(辆/小时)" width="120">
            <template #default="scope">
              <span style="font-weight: bold;">{{ scope.row.flow_rate || 0 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="congestion_level" label="拥堵等级" width="100">
            <template #default="scope">
              <el-tag :type="getCongestionType(scope.row.congestion_level)">
                {{ getCongestionText(scope.row.congestion_level) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="avg_speed" label="平均速度" width="100">
            <template #default="scope">
              {{ scope.row.avg_speed || 0 }} km/h
            </template>
          </el-table-column>
          <el-table-column prop="record_time" label="记录时间" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.record_time) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button type="danger" link size="small" @click="handleDelete(scope.row)">删除</el-button>
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

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="550px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="地点" prop="location">
              <el-input v-model="form.location" placeholder="请输入地点" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="道路名称" prop="road_name">
              <el-input v-model="form.road_name" placeholder="请输入道路名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="车辆数" prop="vehicle_count">
              <el-input-number v-model="form.vehicle_count" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="流量(辆/小时)" prop="flow_rate">
              <el-input-number v-model="form.flow_rate" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="拥堵等级" prop="congestion_level">
              <el-select v-model="form.congestion_level" placeholder="请选择拥堵等级" style="width: 100%">
                <el-option label="畅通" :value="1" />
                <el-option label="缓行" :value="2" />
                <el-option label="拥堵" :value="3" />
                <el-option label="严重拥堵" :value="4" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="平均速度" prop="avg_speed">
              <el-input-number v-model="form.avg_speed" :min="0" :max="120" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="记录时间" prop="record_time">
          <el-date-picker
            v-model="form.record_time"
            type="datetime"
            placeholder="选择记录时间"
            style="width: 100%"
            value-format="YYYY-MM-DDTHH:mm:ss"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getTrafficFlows, createTrafficFlow, updateTrafficFlow, deleteTrafficFlow } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const tableData = ref([])
const formRef = ref(null)

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  location: '',
  road_name: '',
  vehicle_count: 0,
  flow_rate: 0,
  congestion_level: 1,
  avg_speed: 0,
  record_time: ''
})

const dialogTitle = computed(() => isEdit.value ? '编辑流量数据' : '添加流量数据')

const rules = {
  location: [
    { required: true, message: '请输入地点', trigger: 'blur' }
  ],
  road_name: [
    { required: true, message: '请输入道路名称', trigger: 'blur' }
  ],
  congestion_level: [
    { required: true, message: '请选择拥堵等级', trigger: 'change' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const getCongestionType = (level) => {
  const types = { 1: 'success', 2: 'warning', 3: 'danger', 4: 'danger' }
  return types[level] || 'info'
}

const getCongestionText = (level) => {
  const texts = { 1: '畅通', 2: '缓行', 3: '拥堵', 4: '严重拥堵' }
  return texts[level] || '未知'
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

    const res = await getTrafficFlows(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载交通流量列表失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetSearch = () => {
  searchForm.keyword = ''
  pagination.page = 1
  loadData()
}

const handleAdd = () => {
  isEdit.value = false
  formRef.value?.resetFields()
  form.id = null
  form.location = ''
  form.road_name = ''
  form.vehicle_count = 0
  form.flow_rate = 0
  form.congestion_level = 1
  form.avg_speed = 0
  form.record_time = ''
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.location = row.location
  form.road_name = row.road_name
  form.vehicle_count = row.vehicle_count
  form.flow_rate = row.flow_rate
  form.congestion_level = row.congestion_level
  form.avg_speed = row.avg_speed
  form.record_time = row.record_time
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除该交通流量记录吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteTrafficFlow(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadData()
      }
    } catch (e) {
      console.error('删除失败', e)
    }
  }).catch(() => {})
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = {
      location: form.location,
      road_name: form.road_name,
      vehicle_count: form.vehicle_count,
      flow_rate: form.flow_rate,
      congestion_level: form.congestion_level,
      avg_speed: form.avg_speed,
      record_time: form.record_time
    }

    let res
    if (isEdit.value) {
      res = await updateTrafficFlow(form.id, data)
    } else {
      res = await createTrafficFlow(data)
    }

    if (res.code === 200) {
      ElMessage.success(isEdit.value ? '修改成功' : '添加成功')
      dialogVisible.value = false
      loadData()
    }
  } catch (e) {
    console.error('提交失败', e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.traffic-flows-page {
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
