<template>
  <div class="change-password-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>密码修改</h3>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-card>
        <div class="password-form-container">
          <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
            <el-form-item label="当前密码" prop="oldPassword">
              <el-input
                v-model="form.oldPassword"
                type="password"
                placeholder="请输入当前密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="form.newPassword"
                type="password"
                placeholder="请输入新密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirmPassword">
              <el-input
                v-model="form.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="submitting" @click="handleSubmit">
                确认修改
              </el-button>
              <el-button @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>

          <el-divider>安全提示</el-divider>
          <div class="tips">
            <h4>密码安全建议：</h4>
            <ul>
              <li>密码长度至少 8 位字符</li>
              <li>建议包含大小写字母、数字和特殊字符</li>
              <li>避免使用与用户名相同或相似的密码</li>
              <li>定期更换密码以确保账户安全</li>
              <li>不要在多个平台使用相同的密码</li>
            </ul>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { changeAdminPassword } from '@/api/admin'
import { ElMessage } from 'element-plus'

const submitting = ref(false)
const formRef = ref(null)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入新密码'))
  } else if (value !== form.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleReset = () => {
  formRef.value?.resetFields()
  form.oldPassword = ''
  form.newPassword = ''
  form.confirmPassword = ''
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const res = await changeAdminPassword({
      old_password: form.oldPassword,
      new_password: form.newPassword
    })

    if (res.code === 200) {
      ElMessage.success('密码修改成功')
      handleReset()
    }
  } catch (e) {
    console.error('密码修改失败', e)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.change-password-page {
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

.password-form-container {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
}

.tips {
  padding: 0 20px;
}

.tips h4 {
  margin: 0 0 15px 0;
  color: #303133;
}

.tips ul {
  margin: 0;
  padding-left: 20px;
}

.tips li {
  color: #606266;
  line-height: 2;
  font-size: 14px;
}
</style>
