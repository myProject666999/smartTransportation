<template>
  <div class="signal-lights-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>信号灯控制管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加信号灯
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
              placeholder="路口/地点"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="绿灯" :value="1" />
              <el-option label="红灯" :value="2" />
              <el-option label="黄灯" :value="3" />
              <el-option label="故障" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="intersection" label="路口" min-width="180" />
          <el-table-column prop="location" label="地点" min-width="180" />
          <el-table-column prop="current_status" label="当前状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.current_status)">
                <span class="status-dot" :class="'dot-' + scope.row.current_status"></span>
                {{ getStatusText(scope.row.current_status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="灯时配置" width="180">
            <template #default="scope">
              <div class="timing-config">
                <span class="timing-item">绿: <span class="green">{{ scope.row.green_duration || 0 }}s</span></span>
                <span class="timing-item">红: <span class="red">{{ scope.row.red_duration || 0 }}s</span></span>
                <span class="timing-item">黄: <span class="yellow">{{ scope.row.yellow_duration || 0 }}s</span></span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="is_active" label="运行状态" width="100">
            <template #default="scope">
              <el-switch
                v-model="scope.row.is_active"
                active-color="#67C23A"
                inactive-color="#F56C6C"
                active-text="运行"
                inactive-text="停止"
                @change="(val) => handleToggleActive(scope.row, val)"
              />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleManualControl(scope.row)">手动控制</el-button>
              <el-button type="success" link size="small" @click="handleEdit(scope.row)">编辑</el-button>
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
        <el-form-item label="路口" prop="intersection">
          <el-input v-model="form.intersection" placeholder="请输入路口名称" />
        </el-form-item>
        <el-form-item label="地点" prop="location">
          <el-input v-model="form.location" placeholder="请输入地点描述" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="绿灯时长" prop="green_duration">
              <el-input-number v-model="form.green_duration" :min="0" :max="120" style="width: 100%" />
              <span class="unit">秒</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="红灯时长" prop="red_duration">
              <el-input-number v-model="form.red_duration" :min="0" :max="120" style="width: 100%" />
              <span class="unit">秒</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="黄灯时长" prop="yellow_duration">
              <el-input-number v-model="form.yellow_duration" :min="0" :max="10" style="width: 100%" />
              <span class="unit">秒</span>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="是否启用" prop="is_active">
          <el-switch v-model="form.is_active" active-text="是" inactive-text="否" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="controlDialogVisible" title="信号灯手动控制" width="400px">
      <div class="manual-control" v-if="currentItem">
        <div class="light-preview">
          <div class="light-circle red" :class="{ active: controlForm.status === 2 }"></div>
          <div class="light-circle yellow" :class="{ active: controlForm.status === 3 }"></div>
          <div class="light-circle green" :class="{ active: controlForm.status === 1 }"></div>
        </div>
        <div class="control-info">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="路口">{{ currentItem.intersection }}</el-descriptions-item>
            <el-descriptions-item label="当前状态">{{ getStatusText(currentItem.current_status) }}</el-descriptions-item>
          </el-descriptions>
        </div>
        <el-form :model="controlForm" label-width="100px" style="margin-top: 20px;">
          <el-form-item label="设置状态">
            <el-radio-group v-model="controlForm.status">
              <el-radio :value="1">
                <span style="color: #67C23A;">绿灯</span>
              </el-radio>
              <el-radio :value="2">
                <span style="color: #F56C6C;">红灯</span>
              </el-radio>
              <el-radio :value="3">
                <span style="color: #E6A23C;">黄灯</span>
              </el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="controlDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="controlling" @click="handleControlSubmit">执行</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getSignalLights, createSignalLight, updateSignalLight, deleteSignalLight, controlSignalLight } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const controlling = ref(false)
const dialogVisible = ref(false)
const controlDialogVisible = ref(false)
const isEdit = ref(false)
const tableData = ref([])
const formRef = ref(null)
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

const form = reactive({
  id: null,
  intersection: '',
  location: '',
  green_duration: 30,
  red_duration: 30,
  yellow_duration: 3,
  is_active: true,
  current_status: 1
})

const controlForm = reactive({
  status: 1
})

const dialogTitle = computed(() => isEdit.value ? '编辑信号灯' : '添加信号灯')

const rules = {
  intersection: [
    { required: true, message: '请输入路口名称', trigger: 'blur' }
  ],
  location: [
    { required: true, message: '请输入地点', trigger: 'blur' }
  ]
}

const getStatusType = (status) => {
  const types = { 0: 'info', 1: 'success', 2: 'danger', 3: 'warning' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 0: '故障', 1: '绿灯', 2: '红灯', 3: '黄灯' }
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

    const res = await getSignalLights(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载信号灯列表失败', e)
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

const handleAdd = () => {
  isEdit.value = false
  formRef.value?.resetFields()
  form.id = null
  form.intersection = ''
  form.location = ''
  form.green_duration = 30
  form.red_duration = 30
  form.yellow_duration = 3
  form.is_active = true
  form.current_status = 1
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.intersection = row.intersection
  form.location = row.location
  form.green_duration = row.green_duration
  form.red_duration = row.red_duration
  form.yellow_duration = row.yellow_duration
  form.is_active = row.is_active
  form.current_status = row.current_status
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除该信号灯记录吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteSignalLight(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadData()
      }
    } catch (e) {
      console.error('删除失败', e)
    }
  }).catch(() => {})
}

const handleToggleActive = (row, val) => {
  updateSignalLight(row.id, { is_active: val }).then(res => {
    if (res.code === 200) {
      ElMessage.success(val ? '已启用' : '已停止')
    } else {
      row.is_active = !val
    }
  }).catch(() => {
    row.is_active = !val
  })
}

const handleManualControl = (row) => {
  currentItem.value = row
  controlForm.status = row.current_status
  controlDialogVisible.value = true
}

const handleControlSubmit = async () => {
  controlling.value = true
  try {
    const res = await controlSignalLight(currentItem.value.id, {
      status: controlForm.status
    })
    if (res.code === 200) {
      ElMessage.success('控制成功')
      controlDialogVisible.value = false
      loadData()
    }
  } catch (e) {
    console.error('控制失败', e)
  } finally {
    controlling.value = false
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = {
      intersection: form.intersection,
      location: form.location,
      green_duration: form.green_duration,
      red_duration: form.red_duration,
      yellow_duration: form.yellow_duration,
      is_active: form.is_active
    }

    let res
    if (isEdit.value) {
      res = await updateSignalLight(form.id, data)
    } else {
      res = await createSignalLight(data)
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
.signal-lights-page {
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

.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 5px;
  animation: blink 1s infinite;
}

.dot-0 { background: #909399; }
.dot-1 { background: #67C23A; }
.dot-2 { background: #F56C6C; }
.dot-3 { background: #E6A23C; }

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.timing-config {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
}

.timing-item {
  color: #606266;
}

.timing-item .green { color: #67C23A; font-weight: bold; }
.timing-item .red { color: #F56C6C; font-weight: bold; }
.timing-item .yellow { color: #E6A23C; font-weight: bold; }

.unit {
  margin-left: 5px;
  font-size: 12px;
  color: #909399;
}

.manual-control {
  text-align: center;
}

.light-preview {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
}

.light-circle {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #EBEEF5;
  border: 3px solid #E4E7ED;
  transition: all 0.3s;
}

.light-circle.red { border-color: #F56C6C; }
.light-circle.yellow { border-color: #E6A23C; }
.light-circle.green { border-color: #67C23A; }

.light-circle.red.active {
  background: #F56C6C;
  box-shadow: 0 0 20px #F56C6C;
}

.light-circle.yellow.active {
  background: #E6A23C;
  box-shadow: 0 0 20px #E6A23C;
}

.light-circle.green.active {
  background: #67C23A;
  box-shadow: 0 0 20px #67C23A;
}
</style>
