<template>
  <div class="announcements-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>公告管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 发布公告
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
              placeholder="标题/内容"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="已发布" :value="1" />
              <el-option label="草稿" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="title" label="标题" min-width="250" show-overflow-tooltip />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="scope">
              <el-tag :type="getTypeTagType(scope.row.type)">
                {{ scope.row.type || '普通' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="priority" label="优先级" width="100">
            <template #default="scope">
              <el-tag v-if="scope.row.priority === 1" type="danger" size="small">置顶</el-tag>
              <span v-else class="text-muted">普通</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
                {{ scope.row.status === 1 ? '已发布' : '草稿' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleView(scope.row)">查看</el-button>
              <el-button type="success" link size="small" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button v-if="scope.row.status === 0" type="warning" link size="small" @click="handlePublish(scope.row)">发布</el-button>
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
      width="650px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
                <el-option label="普通" value="普通" />
                <el-option label="重要" value="重要" />
                <el-option label="紧急" value="紧急" />
                <el-option label="通知" value="通知" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="优先级" prop="priority">
              <el-radio-group v-model="form.priority">
                <el-radio :value="0">普通</el-radio>
                <el-radio :value="1">置顶</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">发布</el-radio>
            <el-radio :value="0">存为草稿</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="10"
            placeholder="请输入公告内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="公告详情" width="600px">
      <div class="announcement-detail" v-if="currentItem">
        <div class="detail-header">
          <h2>{{ currentItem.title }}</h2>
          <div class="detail-meta">
            <el-tag :type="getTypeTagType(currentItem.type)">{{ currentItem.type || '普通' }}</el-tag>
            <span v-if="currentItem.priority === 1">
              <el-tag type="danger" size="small">置顶</el-tag>
            </span>
            <span class="time">{{ formatDate(currentItem.created_at) }}</span>
          </div>
        </div>
        <div class="detail-content">
          {{ currentItem.content }}
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getAnnouncements, createAnnouncement, updateAnnouncement, deleteAnnouncement } from '@/api/admin'
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
  title: '',
  type: '普通',
  priority: 0,
  content: '',
  status: 1
})

const dialogTitle = computed(() => isEdit.value ? '编辑公告' : '发布公告')

const rules = {
  title: [
    { required: true, message: '请输入公告标题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入公告内容', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const getTypeTagType = (type) => {
  const types = { '普通': 'info', '重要': 'warning', '紧急': 'danger', '通知': 'primary' }
  return types[type] || 'info'
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

    const res = await getAnnouncements(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载公告列表失败', e)
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
  form.title = ''
  form.type = '普通'
  form.priority = 0
  form.content = ''
  form.status = 1
  dialogVisible.value = true
}

const handleView = (row) => {
  currentItem.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.title = row.title
  form.type = row.type || '普通'
  form.priority = row.priority || 0
  form.content = row.content || ''
  form.status = row.status
  dialogVisible.value = true
}

const handlePublish = (row) => {
  ElMessageBox.confirm(`确认要发布公告"${row.title}"吗？`, '确认发布', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await updateAnnouncement(row.id, { status: 1 })
      if (res.code === 200) {
        ElMessage.success('发布成功')
        loadData()
      }
    } catch (e) {
      console.error('发布失败', e)
    }
  }).catch(() => {})
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除公告"${row.title}"吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteAnnouncement(row.id)
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
      title: form.title,
      type: form.type,
      priority: form.priority,
      content: form.content,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await updateAnnouncement(form.id, data)
    } else {
      res = await createAnnouncement(data)
    }

    if (res.code === 200) {
      ElMessage.success(isEdit.value ? '修改成功' : '发布成功')
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
.announcements-page {
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

.text-muted {
  color: #909399;
  font-size: 14px;
}

.announcement-detail {
  padding: 10px;
}

.detail-header {
  text-align: center;
  padding-bottom: 20px;
  border-bottom: 1px solid #EBEEF5;
  margin-bottom: 20px;
}

.detail-header h2 {
  margin: 0 0 15px 0;
  color: #303133;
  font-size: 20px;
}

.detail-meta {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  color: #909399;
  font-size: 14px;
}

.detail-content {
  line-height: 2;
  color: #606266;
  font-size: 15px;
  white-space: pre-wrap;
}
</style>
