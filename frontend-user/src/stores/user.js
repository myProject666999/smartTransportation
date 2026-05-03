import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userLogin, userRegister, getUserProfile, updateUserProfile, changePassword } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  const isLoggedIn = computed(() => !!token.value)

  async function login(username, password) {
    const res = await userLogin({ username, password })
    if (res.code === 200) {
      token.value = res.data.token
      userInfo.value = res.data.user_info
      localStorage.setItem('token', token.value)
      localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
    }
    return res
  }

  async function register(userData) {
    const res = await userRegister(userData)
    return res
  }

  async function fetchUserProfile() {
    const res = await getUserProfile()
    if (res.code === 200) {
      userInfo.value = res.data
      localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
    }
    return res
  }

  async function updateProfile(data) {
    const res = await updateUserProfile(data)
    if (res.code === 200) {
      userInfo.value = { ...userInfo.value, ...data }
      localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
    }
    return res
  }

  async function changePwd(oldPassword, newPassword) {
    const res = await changePassword({ old_password: oldPassword, new_password: newPassword })
    return res
  }

  function logout() {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    login,
    register,
    fetchUserProfile,
    updateProfile,
    changePwd,
    logout
  }
})
