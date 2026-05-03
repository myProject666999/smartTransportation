<template>
  <div class="road-conditions-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>路况管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加路况
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
              placeholder="道路名称/位置"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="路况">
            <el-select v-model="searchForm.condition" placeholder="全部" clearable style="width: 100px">
              <el-option label="畅通" value="畅通" />
              <el-option label="缓行" value="缓行" />
              <el-option label="拥堵" value="拥堵" />
              <el-option label="施工" value="施工" />
              <el-option label="事故" value="事故" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="进行中" :value="1" />
              <el-option label="已结束" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="road_name" label="道路名称" width="150" />
          <el-table-column prop="location" label="位置" min-width="150" />
          <el-table-column prop="condition" label="路况" width="100">
            <template #default="scope">
              <el-tag :type="getConditionType(scope.row.condition)">
                {{ scope.row.condition }}
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
          <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="90">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'danger' : 'info'">
                {{ scope.row.status === 1 ? '进行中' : '已结束' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="start_time" label="开始时间" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.start_time) }}
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
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="道路名称" prop="road_name">
              <el-input v-model="form.road_name" placeholder="请输入道路名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="位置" prop="location">
              <el-input v-model="form.location" placeholder="请输入位置" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="路况" prop="condition">
              <el-select v-model="form.condition" placeholder="请选择路况" style="width: 100%">
                <el-option label="畅通" value="畅通" />
                <el-option label="缓行" value="缓行" />
                <el-option label="拥堵" value="拥堵" />
                <el-option label="施工" value="施工" />
                <el-option label="事故" value="事故" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="严重程度" prop="severity">
              <el-select v-model="form.severity" placeholder="请选择严重程度" style="width: 100%">
                <el-option label="轻微" value="轻微" />
                <el-option label="一般" value="一般" />
                <el-option label="严重" value="严重" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入路况描述"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始时间" prop="start_time">
              <el-date-picker
                v-model="form.start_time"
                type="datetime"
                placeholder="选择开始时间"
                style="width: 100%"
                value-format="YYYY-MM-DDTHH:mm:ss"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="预计结束时间" prop="end_time">
              <el-date-picker
                v-model="form.end_time"
                type="datetime"
                placeholder="选择预计结束时间"
                style="width: 100%"
                value-format="YYYY-MM-DDTHH:mm:ss"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">进行中</el-radio>
            <el-radio :label="0">已结束</el-radio>
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
import { getRoadConditions, createRoadCondition, updateRoadCondition, deleteRoadCondition } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const tableData = ref([])
const formRef = ref(null)

const searchForm = reactive({
  keyword: '',
  condition: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  road_name: '',
  location: '',
  condition: '',
  severity: '',
  description: '',
  start_time: '',
  end_time: '',
  status: 1
})

const dialogTitle = computed(() => isEdit.value ? '编辑路况' : '添加路况')

const rules = {
  road_name: [
    { required: true, message: '请输入道路名称', trigger: 'blur' }
  ],
  condition: [
    { required: true, message: '请选择路况', trigger: 'change' }
  ]
}

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
    if (searchForm.condition) {
      params.condition = searchForm.condition
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getRoadConditions(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载路况列表失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetSearch = () => {
  searchForm.keyword = ''
  searchForm.condition = ''
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const handleAdd = () => {
  isEdit.value = false
  formRef.value?.resetFields()
  form.id = null
  form.road_name = ''
  form.location = ''
  form.condition = ''
  form.severity = ''
  form.description = ''
  form.start_time = ''
  form.end_time = ''
  form.status = 1
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.road_name = row.road_name
  form.location = row.location
  form.condition = row.condition
  form.severity = row.severity
  form.description = row.description
  form.start_time = row.start_time
  form.end_time = row.end_time
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除路况"${row.road_name}"的记录吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteRoadCondition(row.id)
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
      road_name: form.road_name,
      location: form.location,
      condition: form.condition,
      severity: form.severity,
      description: form.description,
      start_time: form.start_time,
      end_time: form.end_time,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await updateRoadCondition(form.id, data)
    } else {
      res = await createRoadCondition(data)
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
.road-conditions-page {
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
