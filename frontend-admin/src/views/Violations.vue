<template>
  <div class="violations-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>违规管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加违规
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
              placeholder="车牌号/违规类型"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 120px">
              <el-option label="待处理" :value="0" />
              <el-option label="已处理" :value="1" />
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
          <el-table-column prop="violation_type" label="违规类型" min-width="150" />
          <el-table-column prop="location" label="违规地点" min-width="150" />
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
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                {{ scope.row.status === 1 ? '已处理' : '待处理' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="violation_time" label="违规时间" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.violation_time) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleView(scope.row)">查看</el-button>
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
            <el-form-item label="车牌号" prop="plate_number">
              <el-input v-model="form.plate_number" placeholder="请输入车牌号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="违规类型" prop="violation_type">
              <el-select v-model="form.violation_type" placeholder="请选择违规类型" style="width: 100%">
                <el-option label="闯红灯" value="闯红灯" />
                <el-option label="超速行驶" value="超速行驶" />
                <el-option label="违规停车" value="违规停车" />
                <el-option label="不按车道行驶" value="不按车道行驶" />
                <el-option label="逆行" value="逆行" />
                <el-option label="酒后驾驶" value="酒后驾驶" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="违规地点" prop="location">
          <el-input v-model="form.location" placeholder="请输入违规地点" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="罚款金额" prop="fine_amount">
              <el-input-number v-model="form.fine_amount" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="扣分" prop="points">
              <el-input-number v-model="form.points" :min="0" :max="12" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="违规时间" prop="violation_time">
          <el-date-picker
            v-model="form.violation_time"
            type="datetime"
            placeholder="选择违规时间"
            style="width: 100%"
            value-format="YYYY-MM-DDTHH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入违规描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">待处理</el-radio>
            <el-radio :label="1">已处理</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="违规详情" width="500px">
      <el-descriptions :column="2" border v-if="currentItem">
        <el-descriptions-item label="车牌号">{{ currentItem.plate_number }}</el-descriptions-item>
        <el-descriptions-item label="违规类型">{{ currentItem.violation_type }}</el-descriptions-item>
        <el-descriptions-item label="违规地点">{{ currentItem.location }}</el-descriptions-item>
        <el-descriptions-item label="违规时间">{{ formatDate(currentItem.violation_time) }}</el-descriptions-item>
        <el-descriptions-item label="罚款金额">
          <span style="color: #F56C6C; font-weight: bold;">¥{{ currentItem.fine_amount }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="扣分">
          <el-tag :type="currentItem.points > 3 ? 'danger' : 'warning'" size="small">
            {{ currentItem.points }}分
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态" :span="2">
          <el-tag :type="currentItem.status === 1 ? 'success' : 'danger'">
            {{ currentItem.status === 1 ? '已处理' : '待处理' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentItem.description || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getViolations, createViolation, updateViolation, deleteViolation } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
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
  plate_number: '',
  violation_type: '',
  location: '',
  fine_amount: 0,
  points: 0,
  violation_time: '',
  description: '',
  status: 0
})

const dialogTitle = computed(() => isEdit.value ? '编辑违规' : '添加违规')

const rules = {
  plate_number: [
    { required: true, message: '请输入车牌号', trigger: 'blur' }
  ],
  violation_type: [
    { required: true, message: '请选择违规类型', trigger: 'change' }
  ],
  location: [
    { required: true, message: '请输入违规地点', trigger: 'blur' }
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
    if (searchForm.keyword) {
      params.keyword = searchForm.keyword
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getViolations(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载违规列表失败', e)
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

const handleAdd = () => {
  isEdit.value = false
  formRef.value?.resetFields()
  form.id = null
  form.plate_number = ''
  form.violation_type = ''
  form.location = ''
  form.fine_amount = 0
  form.points = 0
  form.violation_time = ''
  form.description = ''
  form.status = 0
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.plate_number = row.plate_number
  form.violation_type = row.violation_type
  form.location = row.location
  form.fine_amount = row.fine_amount
  form.points = row.points
  form.violation_time = row.violation_time
  form.description = row.description
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除车牌号为"${row.plate_number}"的违规记录吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteViolation(row.id)
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
      violation_type: form.violation_type,
      location: form.location,
      fine_amount: form.fine_amount,
      points: form.points,
      violation_time: form.violation_time,
      description: form.description,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await updateViolation(form.id, data)
    } else {
      res = await createViolation(data)
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
.violations-page {
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
