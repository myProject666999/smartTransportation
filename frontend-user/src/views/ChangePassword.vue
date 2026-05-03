<template>
  <div class="change-password-page">
    <el-header class="header">
      <div class="container">
        <div class="logo" @click="router.push('/')">
          <el-icon size="24" color="#409EFF"><MapLocation /></el-icon>
          <span class="logo-text">智能交通管理系统</span>
        </div>
        <div class="nav-actions">
          <el-button link @click="router.push('/profile')">个人信息</el-button>
          <el-button link @click="router.push('/')">返回首页</el-button>
        </div>
      </div>
    </el-header>

    <div class="main-content">
      <div class="container">
        <el-card class="form-card">
          <template #header>
            <span>修改密码</span>
          </template>

          <el-form
            :model="form"
            :rules="rules"
            ref="formRef"
            label-width="120px"
            class="password-form"
          >
            <el-form-item label="原密码" prop="old_password">
              <el-input
                v-model="form.old_password"
                type="password"
                placeholder="请输入原密码"
                show-password
                style="width: 350px"
              />
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
              <el-input
                v-model="form.new_password"
                type="password"
                placeholder="请输入新密码（6-20位）"
                show-password
                style="width: 350px"
              />
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirm_password">
              <el-input
                v-model="form.confirm_password"
                type="password"
                placeholder="请再次输入新密码"
                show-password
                style="width: 350px"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="handleSubmit">
                确认修改
              </el-button>
              <el-button @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card class="tips-card">
          <template #header>
            <span><el-icon><Warning /></el-icon> 安全提示</span>
          </template>
          <ul class="tips-list">
            <li>密码长度应为6-20位字符</li>
            <li>建议使用字母、数字和特殊字符的组合</li>
            <li>不要使用与其他网站相同的密码</li>
            <li>定期更换密码可以提高账户安全性</li>
            <li>如果忘记密码，请联系管理员重置</li>
          </ul>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { changePassword } from '@/api/user'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const formRef = ref(null)

const form = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  old_password: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度应为6-20位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await changePassword({
      old_password: form.old_password,
      new_password: form.new_password,
      confirm_password: form.confirm_password
    })
    if (res.code === 200) {
      ElMessage.success('密码修改成功！')
      handleReset()
    }
  } catch (e) {
    console.error('密码修改失败', e)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  form.old_password = ''
  form.new_password = ''
  form.confirm_password = ''
  formRef.value?.resetFields()
}
</script>

<style scoped>
.change-password-page {
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
  align-items: center;
  gap: 20px;
}

.form-card {
  width: 100%;
  max-width: 600px;
}

.password-form {
  padding: 20px 0;
}

.tips-card {
  width: 100%;
  max-width: 600px;
}

.tips-list {
  margin: 0;
  padding-left: 20px;
}

.tips-list li {
  margin-bottom: 10px;
  color: #606266;
  font-size: 14px;
}

.tips-list li:last-child {
  margin-bottom: 0;
}
</style>
