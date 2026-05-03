<template>
  <div class="feedbacks-page">
    <el-header class="header">
      <div class="container">
        <div class="logo" @click="router.push('/')">
          <el-icon size="24" color="#409EFF"><MapLocation /></el-icon>
          <span class="logo-text">智能交通管理系统</span>
        </div>
        <div class="nav-actions">
          <el-button type="primary" @click="showFeedbackDialog = true">
            <el-icon><Plus /></el-icon> 提交反馈
          </el-button>
          <el-button link @click="router.push('/')">返回首页</el-button>
        </div>
      </div>
    </el-header>

    <div class="main-content">
      <div class="container">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>投诉反馈</span>
              <div class="search-form">
                <el-select v-model="searchForm.type" placeholder="反馈类型" clearable style="width: 120px" @change="loadFeedbacks">
                  <el-option label="投诉" value="投诉" />
                  <el-option label="建议" value="建议" />
                  <el-option label="咨询" value="咨询" />
                </el-select>
                <el-select v-model="searchForm.status" placeholder="处理状态" clearable style="width: 120px" @change="loadFeedbacks">
                  <el-option label="待处理" :value="0" />
                  <el-option label="已回复" :value="1" />
                  <el-option label="已关闭" :value="2" />
                </el-select>
              </div>
            </div>
          </template>

          <el-table :data="feedbacks" style="width: 100%" v-loading="loading" stripe>
            <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="scope">
                <el-tag :type="getTypeTagType(scope.row.type)">
                  {{ scope.row.type }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="content" label="内容" min-width="250" show-overflow-tooltip />
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
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="scope">
                <el-button type="primary" link size="small" @click="handleView(scope.row)">
                  查看详情
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
            @size-change="loadFeedbacks"
            @current-change="loadFeedbacks"
            style="margin-top: 20px; justify-content: flex-end"
          />
        </el-card>
      </div>
    </div>

    <el-dialog
      v-model="showFeedbackDialog"
      title="提交投诉反馈"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="feedbackForm" :rules="feedbackRules" ref="feedbackFormRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="feedbackForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="feedbackForm.type" placeholder="请选择反馈类型" style="width: 100%">
            <el-option label="投诉" value="投诉" />
            <el-option label="建议" value="建议" />
            <el-option label="咨询" value="咨询" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="feedbackForm.content"
            type="textarea"
            :rows="5"
            placeholder="请详细描述您的问题或建议"
          />
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="feedbackForm.contact" placeholder="请输入您的联系方式（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFeedbackDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmitFeedback">提交</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="showDetailDialog"
      title="反馈详情"
      width="600px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="标题">{{ currentFeedback.title }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getTypeTagType(currentFeedback.type)">{{ currentFeedback.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTagType(currentFeedback.status)">{{ getStatusText(currentFeedback.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="提交时间">{{ formatDate(currentFeedback.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="内容" :span="2">
          {{ currentFeedback.content }}
        </el-descriptions-item>
        <el-descriptions-item v-if="currentFeedback.contact" label="联系方式">{{ currentFeedback.contact }}</el-descriptions-item>
        <el-descriptions-item v-if="currentFeedback.reply_time" label="回复时间">{{ formatDate(currentFeedback.reply_time) }}</el-descriptions-item>
        <el-descriptions-item v-if="currentFeedback.reply_content" label="回复内容" :span="2">
          {{ currentFeedback.reply_content }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getFeedbacks, submitFeedback } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const feedbacks = ref([])
const showFeedbackDialog = ref(false)
const showDetailDialog = ref(false)
const currentFeedback = ref({})
const feedbackFormRef = ref(null)

const searchForm = reactive({
  type: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const feedbackForm = reactive({
  title: '',
  type: '',
  content: '',
  contact: ''
})

const feedbackRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择反馈类型', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入反馈内容', trigger: 'blur' },
    { min: 10, message: '内容至少10个字', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const getTypeTagType = (type) => {
  const types = {
    '投诉': 'danger',
    '建议': 'warning',
    '咨询': 'primary'
  }
  return types[type] || 'info'
}

const getStatusTagType = (status) => {
  const types = {
    0: 'warning',
    1: 'success',
    2: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    0: '待处理',
    1: '已回复',
    2: '已关闭'
  }
  return texts[status] || '未知'
}

const loadFeedbacks = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.type) {
      params.type = searchForm.type
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getFeedbacks(params)
    if (res.code === 200) {
      feedbacks.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载反馈失败', e)
  } finally {
    loading.value = false
  }
}

const handleView = (row) => {
  currentFeedback.value = row
  showDetailDialog.value = true
}

const handleSubmitFeedback = async () => {
  const valid = await feedbackFormRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const res = await submitFeedback(feedbackForm)
    if (res.code === 200) {
      ElMessage.success('提交成功！')
      showFeedbackDialog.value = false
      resetForm()
      loadFeedbacks()
    }
  } catch (e) {
    console.error('提交失败', e)
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  feedbackForm.title = ''
  feedbackForm.type = ''
  feedbackForm.content = ''
  feedbackForm.contact = ''
  feedbackFormRef.value?.resetFields()
}

onMounted(() => {
  loadFeedbacks()
})
</script>

<style scoped>
.feedbacks-page {
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
</style>
