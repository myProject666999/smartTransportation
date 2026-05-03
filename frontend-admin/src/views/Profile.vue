<template>
  <div class="profile-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>个人信息</h3>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card>
            <div class="avatar-card">
              <el-avatar :size="120" class="avatar">
                <el-icon v-if="!adminInfo.avatar"><User /></el-icon>
                <img v-else :src="adminInfo.avatar" alt="avatar" />
              </el-avatar>
              <h3 class="username">{{ adminInfo.username }}</h3>
              <el-tag :type="adminInfo.role === 1 ? 'danger' : 'primary'" class="role-tag">
                {{ adminInfo.role === 1 ? '超级管理员' : '普通管理员' }}
              </el-tag>
              <div class="info-item">
                <span class="label">管理员ID：</span>
                <span class="value">{{ adminInfo.id }}</span>
              </div>
              <div class="info-item">
                <span class="label">账号状态：</span>
                <el-tag :type="adminInfo.status === 1 ? 'success' : 'danger'" size="small">
                  {{ adminInfo.status === 1 ? '正常' : '禁用' }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">创建时间：</span>
                <span class="value">{{ formatDate(adminInfo.created_at) }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="16">
          <el-card>
            <template #header>
              <span>编辑个人信息</span>
            </template>
            <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
              <el-form-item label="用户名" prop="username">
                <el-input v-model="form.username" placeholder="请输入用户名" />
              </el-form-item>
              <el-form-item label="姓名" prop="name">
                <el-input v-model="form.name" placeholder="请输入姓名" />
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="form.email" placeholder="请输入邮箱" />
              </el-form-item>
              <el-form-item label="手机号" prop="phone">
                <el-input v-model="form.phone" placeholder="请输入手机号" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="submitting" @click="handleSubmit">
                  保存修改
                </el-button>
                <el-button @click="handleReset">重置</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { updateAdminProfile, getAdminProfile } from '@/api/admin'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const submitting = ref(false)
const formRef = ref(null)

const adminInfo = computed(() => userStore.admin || {})

const form = reactive({
  username: '',
  name: '',
  email: '',
  phone: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const loadProfile = async () => {
  try {
    const res = await getAdminProfile()
    if (res.code === 200 && res.data) {
      form.username = res.data.username || ''
      form.name = res.data.name || ''
      form.email = res.data.email || ''
      form.phone = res.data.phone || ''
    }
  } catch (e) {
    console.error('加载个人信息失败', e)
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  loadProfile()
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = {
      username: form.username,
      name: form.name,
      email: form.email,
      phone: form.phone
    }

    const res = await updateAdminProfile(data)
    if (res.code === 200) {
      ElMessage.success('修改成功')
      userStore.setAdmin({
        ...adminInfo.value,
        ...data
      })
    }
  } catch (e) {
    console.error('保存失败', e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-page {
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

.avatar-card {
  text-align: center;
  padding: 20px;
}

.avatar {
  margin-bottom: 20px;
  background: #F5F7FA;
}

.username {
  margin: 0 0 10px 0;
  font-size: 20px;
  color: #303133;
}

.role-tag {
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 10px;
  font-size: 14px;
}

.info-item .label {
  color: #909399;
}

.info-item .value {
  color: #606266;
}
</style>
