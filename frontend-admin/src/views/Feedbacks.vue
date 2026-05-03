<template>
  <div class="feedbacks-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>投诉反馈管理</h3>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-card>
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="标题/用户名"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="类型">
            <el-select v-model="searchForm.type" placeholder="全部" clearable style="width: 100px">
              <el-option label="投诉" value="投诉" />
              <el-option label="建议" value="建议" />
              <el-option label="咨询" value="咨询" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="待处理" :value="0" />
              <el-option label="已回复" :value="1" />
              <el-option label="已关闭" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
          <el-table-column prop="username" label="用户" width="100" />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="scope">
              <el-tag :type="getTypeTagType(scope.row.type)">
                {{ scope.row.type }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusTagType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="提交时间" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleView(scope.row)">查看</el-button>
              <el-button
                v-if="scope.row.status === 0"
                type="success"
                link
                size="small"
                @click="handleReply(scope.row)"
              >
                回复
              </el-button>
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

    <el-dialog v-model="detailDialogVisible" title="反馈详情" width="600px">
      <el-descriptions :column="2" border v-if="currentItem">
        <el-descriptions-item label="标题">{{ currentItem.title }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentItem.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getTypeTagType(currentItem.type)">{{ currentItem.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTagType(currentItem.status)">{{ getStatusText(currentItem.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="联系方式">{{ currentItem.contact || '-' }}</el-descriptions-item>
        <el-descriptions-item label="提交时间">{{ formatDate(currentItem.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="反馈内容" :span="2">{{ currentItem.content }}</el-descriptions-item>
        <el-descriptions-item v-if="currentItem.reply_content" label="回复内容" :span="2">
          {{ currentItem.reply_content }}
        </el-descriptions-item>
        <el-descriptions-item v-if="currentItem.reply_time" label="回复时间" :span="2">
          {{ formatDate(currentItem.reply_time) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <el-dialog v-model="replyDialogVisible" title="回复反馈" width="500px">
      <el-form :model="replyForm" :rules="replyRules" ref="replyFormRef" label-width="100px">
        <el-form-item label="回复内容" prop="reply">
          <el-input
            v-model="replyForm.reply"
            type="textarea"
            :rows="5"
            placeholder="请输入回复内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="replyDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="replying" @click="handleSubmitReply">发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getFeedbacks, replyFeedback, deleteFeedback } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const replying = ref(false)
const tableData = ref([])
const detailDialogVisible = ref(false)
const replyDialogVisible = ref(false)
const currentItem = ref({})
const replyFormRef = ref(null)

const searchForm = reactive({
  keyword: '',
  type: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const replyForm = reactive({
  reply: ''
})

const replyRules = {
  reply: [
    { required: true, message: '请输入回复内容', trigger: 'blur' },
    { min: 5, message: '回复内容至少5个字', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const getTypeTagType = (type) => {
  const types = { '投诉': 'danger', '建议': 'warning', '咨询': 'primary' }
  return types[type] || 'info'
}

const getStatusTagType = (status) => {
  const types = { 0: 'warning', 1: 'success', 2: 'info' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 0: '待处理', 1: '已回复', 2: '已关闭' }
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
    if (searchForm.type) {
      params.type = searchForm.type
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getFeedbacks(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载反馈列表失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetSearch = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const handleView = (row) => {
  currentItem.value = row
  detailDialogVisible.value = true
}

const handleReply = (row) => {
  currentItem.value = row
  replyForm.reply = ''
  replyFormRef.value?.resetFields()
  replyDialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除该反馈记录吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteFeedback(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadData()
      }
    } catch (e) {
      console.error('删除失败', e)
    }
  }).catch(() => {})
}

const handleSubmitReply = async () => {
  const valid = await replyFormRef.value.validate().catch(() => false)
  if (!valid) return

  replying.value = true
  try {
    const res = await replyFeedback(currentItem.value.id, replyForm.reply)
    if (res.code === 200) {
      ElMessage.success('回复成功')
      replyDialogVisible.value = false
      loadData()
    }
  } catch (e) {
    console.error('回复失败', e)
  } finally {
    replying.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.feedbacks-page {
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
