<template>
  <div class="vehicle-monitors-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>车辆监控管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加监控
          </el-button>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-card>
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="车牌号">
            <el-input
              v-model="searchForm.plate_number"
              placeholder="请输入车牌号"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="在线" :value="1" />
              <el-option label="离线" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
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
          <el-table-column prop="speed" label="速度" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.speed > 80 ? 'danger' : 'success'" size="small">
                {{ scope.row.speed || 0 }} km/h
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="direction" label="方向" width="80" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
                {{ scope.row.status === 1 ? '在线' : '离线' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="monitor_time" label="监控时间" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.monitor_time) }}
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
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="车牌号" prop="plate_number">
          <el-input v-model="form.plate_number" placeholder="请输入车牌号" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="经度" prop="longitude">
              <el-input v-model="form.longitude" placeholder="请输入经度" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="纬度" prop="latitude">
              <el-input v-model="form.latitude" placeholder="请输入纬度" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="当前位置" prop="location">
          <el-input v-model="form.location" placeholder="请输入当前位置" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="速度(km/h)" prop="speed">
              <el-input-number v-model="form.speed" :min="0" :max="200" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="方向" prop="direction">
              <el-select v-model="form.direction" placeholder="请选择方向" style="width: 100%">
                <el-option label="东" value="东" />
                <el-option label="南" value="南" />
                <el-option label="西" value="西" />
                <el-option label="北" value="北" />
                <el-option label="东南" value="东南" />
                <el-option label="东北" value="东北" />
                <el-option label="西南" value="西南" />
                <el-option label="西北" value="西北" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">在线</el-radio>
            <el-radio :label="0">离线</el-radio>
          </el-radio-group>
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
import { getVehicleMonitors, createVehicleMonitor, updateVehicleMonitor, deleteVehicleMonitor } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const tableData = ref([])
const formRef = ref(null)

const searchForm = reactive({
  plate_number: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  plate_number: '',
  longitude: null,
  latitude: null,
  location: '',
  speed: 0,
  direction: '',
  status: 1
})

const dialogTitle = computed(() => isEdit.value ? '编辑监控' : '添加监控')

const rules = {
  plate_number: [
    { required: true, message: '请输入车牌号', trigger: 'blur' }
  ],
  location: [
    { required: true, message: '请输入当前位置', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.plate_number) {
      params.plate_number = searchForm.plate_number
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getVehicleMonitors(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载车辆监控列表失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetSearch = () => {
  searchForm.plate_number = ''
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const handleAdd = () => {
  isEdit.value = false
  formRef.value?.resetFields()
  form.id = null
  form.plate_number = ''
  form.longitude = null
  form.latitude = null
  form.location = ''
  form.speed = 0
  form.direction = ''
  form.status = 1
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.plate_number = row.plate_number
  form.longitude = row.longitude
  form.latitude = row.latitude
  form.location = row.location
  form.speed = row.speed
  form.direction = row.direction
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除车辆"${row.plate_number}"的监控记录吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteVehicleMonitor(row.id)
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
      plate_number: form.plate_number,
      longitude: form.longitude,
      latitude: form.latitude,
      location: form.location,
      speed: form.speed,
      direction: form.direction,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await updateVehicleMonitor(form.id, data)
    } else {
      res = await createVehicleMonitor(data)
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
.vehicle-monitors-page {
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
