import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { adminLogin, getAdminProfile, changePassword } from '@/api/admin'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const adminInfo = ref(JSON.parse(localStorage.getItem('admin_info') || '{}'))

  const isLoggedIn = computed(() => !!token.value)

  async function login(username, password) {
    const res = await adminLogin({ username, password })
    if (res.code === 200) {
      token.value = res.data.token
      adminInfo.value = res.data.user_info
      localStorage.setItem('admin_token', token.value)
      localStorage.setItem('admin_info', JSON.stringify(adminInfo.value))
    }
    return res
  }

  async function fetchProfile() {
    const res = await getAdminProfile()
    if (res.code === 200) {
      adminInfo.value = res.data
      localStorage.setItem('admin_info', JSON.stringify(adminInfo.value))
    }
    return res
  }

  async function changePwd(oldPassword, newPassword) {
    const res = await changePassword({ old_password: oldPassword, new_password: newPassword })
    return res
  }

  function logout() {
    token.value = ''
    adminInfo.value = {}
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_info')
  }

  return {
    token,
    adminInfo,
    isLoggedIn,
    login,
    fetchProfile,
    changePwd,
    logout
  }
})
