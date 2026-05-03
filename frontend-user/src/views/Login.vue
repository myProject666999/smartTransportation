<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <el-icon size="40" color="#409EFF"><MapLocation /></el-icon>
        <h2>用户登录</h2>
        <p>智能交通管理系统</p>
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
        <div class="login-links">
          <span>还没有账号？</span>
          <router-link to="/register">立即注册</router-link>
        </div>
        <div class="login-links" style="margin-top: 10px;">
          <span>管理员入口：</span>
          <a href="http://localhost:3001" target="_blank">后台管理</a>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

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
    const res = await userStore.login(form.username, form.password)
    if (res.code === 200) {
      ElMessage.success('登录成功')
      const redirect = route.query.redirect || '/'
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-container {
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
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

.login-links {
  text-align: center;
  font-size: 14px;
}

.login-links span {
  color: #909399;
}

.login-links a, .login-links router-link {
  color: #409EFF;
  text-decoration: none;
}

.login-links a:hover, .login-links router-link:hover {
  text-decoration: underline;
}
</style>
