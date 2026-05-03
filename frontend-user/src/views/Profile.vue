<template>
  <div class="profile-page">
    <el-header class="header">
      <div class="container">
        <div class="logo" @click="router.push('/')">
          <el-icon size="24" color="#409EFF"><MapLocation /></el-icon>
          <span class="logo-text">智能交通管理系统</span>
        </div>
        <div class="nav-actions">
          <el-button link @click="router.push('/')">返回首页</el-button>
        </div>
      </div>
    </el-header>

    <div class="main-content">
      <div class="container">
        <el-card class="profile-card">
          <template #header>
            <span>个人信息</span>
          </template>
          <el-tabs v-model="activeTab">
            <el-tab-pane label="基本信息" name="info">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="用户名">
                  {{ userInfo.username }}
                </el-descriptions-item>
                <el-descriptions-item label="邮箱">
                  {{ userInfo.email || '未设置' }}
                </el-descriptions-item>
                <el-descriptions-item label="手机号">
                  {{ userInfo.phone || '未设置' }}
                </el-descriptions-item>
                <el-descriptions-item label="真实姓名">
                  {{ userInfo.real_name || '未设置' }}
                </el-descriptions-item>
                <el-descriptions-item label="身份证号">
                  {{ userInfo.id_card || '未设置' }}
                </el-descriptions-item>
                <el-descriptions-item label="注册时间">
                  {{ formatDate(userInfo.created_at) }}
                </el-descriptions-item>
              </el-descriptions>
              <div style="margin-top: 20px;">
                <el-button type="primary" @click="showEdit = true">编辑信息</el-button>
                <el-button @click="router.push('/change-password')">修改密码</el-button>
              </div>
            </el-tab-pane>
            <el-tab-pane label="快捷操作" name="quick">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-card shadow="hover" class="quick-card" @click="router.push('/bookings')">
                    <el-icon size="36" color="#409EFF"><Ticket /></el-icon>
                    <span>购票管理</span>
                  </el-card>
                </el-col>
                <el-col :span="6">
                  <el-card shadow="hover" class="quick-card" @click="router.push('/violations')">
                    <el-icon size="36" color="#F56C6C"><Document /></el-icon>
                    <span>违规管理</span>
                  </el-card>
                </el-col>
                <el-col :span="6">
                  <el-card shadow="hover" class="quick-card" @click="router.push('/feedbacks')">
                    <el-icon size="36" color="#E6A23C"><ChatDotRound /></el-icon>
                    <span>投诉反馈</span>
                  </el-card>
                </el-col>
                <el-col :span="6">
                  <el-card shadow="hover" class="quick-card" @click="router.push('/tickets')">
                    <el-icon size="36" color="#67C23A"><ShoppingCart /></el-icon>
                    <span>购买车票</span>
                  </el-card>
                </el-col>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="showEdit" title="编辑个人信息" width="500px">
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="100px">
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="editForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="editForm.real_name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="身份证号" prop="id_card">
          <el-input v-model="editForm.id_card" placeholder="请输入身份证号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEdit = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('info')
const showEdit = ref(false)
const saving = ref(false)
const editFormRef = ref(null)

const userInfo = computed(() => userStore.userInfo)

const editForm = reactive({
  email: '',
  phone: '',
  real_name: '',
  id_card: ''
})

const editRules = {
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const openEdit = () => {
  editForm.email = userInfo.value.email || ''
  editForm.phone = userInfo.value.phone || ''
  editForm.real_name = userInfo.value.real_name || ''
  editForm.id_card = userInfo.value.id_card || ''
  showEdit.value = true
}

const handleSave = async () => {
  const valid = await editFormRef.value.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    const res = await userStore.updateProfile(editForm)
    if (res.code === 200) {
      ElMessage.success('保存成功')
      showEdit.value = false
      userStore.fetchUserProfile()
    }
  } catch (e) {
    console.error('保存失败', e)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (userStore.isLoggedIn) {
    userStore.fetchUserProfile()
  }
})
</script>

<style scoped>
.profile-page {
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

.profile-card {
  max-width: 900px;
  margin: 0 auto;
}

.quick-card {
  text-align: center;
  padding: 30px 0;
  cursor: pointer;
  transition: all 0.3s;
}

.quick-card:hover {
  transform: translateY(-5px);
}

.quick-card span {
  display: block;
  margin-top: 10px;
  font-size: 14px;
  color: #333;
}
</style>
