<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <el-icon size="48" color="#409EFF"><MapLocation /></el-icon>
        <h1>智能交通管理系统</h1>
        <p>管理员后台</p>
      </div>
      <el-form ref="formRef" :model="form" :rules="rules" class="login-form" label-width="0">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" prefix-icon="Lock" size="large" show-password @keyup.enter="handleLogin" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" class="login-btn" @click="handleLogin">登录</el-button>
        </el-form-item>
        <div class="login-tips">
          <el-text type="info">默认账号：admin / admin123</el-text>
        </div>
        <div class="login-links">
          <span>用户入口：</span>
          <a href="http://localhost:3000" target="_blank">用户前台</a>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminStore()

const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await adminStore.login(form.username, form.password)
    if (res.code === 200) {
      ElMessage.success('登录成功')
      const redirect = route.query.redirect || '/dashboard'
      router.push(redirect)
    }
  } catch (e) {
    console.error('登录失败', e)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
}

.login-container {
  background: #fff;
  border-radius: 12px;
  padding: 50px 40px;
  width: 420px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h1 {
  margin: 15px 0 5px 0;
  color: #333;
  font-size: 24px;
}

.login-header p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.login-form {
  margin-bottom: 0;
}

.login-btn {
  width: 100%;
}

.login-tips {
  text-align: center;
  margin-bottom: 15px;
}

.login-links {
  text-align: center;
  font-size: 14px;
}

.login-links span {
  color: #909399;
}

.login-links a {
  color: #409EFF;
  text-decoration: none;
}

.login-links a:hover {
  text-decoration: underline;
}
</style>
