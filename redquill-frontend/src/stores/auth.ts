import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import { api } from '@/utils/api'

export interface User {
  id: string
  name: string
  email: string
  ctime: number
  mtime: number
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  // 初始化时检查token有效性
  const initAuth = () => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      token.value = storedToken
      // 这里可以添加验证token有效性的逻辑
    }
  }

  // 立即初始化
  initAuth()

  // 登录
  const login = async (email: string, password: string) => {
    try {
      loading.value = true
      const response = await api.post('/login', { email, password })
      
      if (response.data.token) {
        token.value = response.data.token
        user.value = response.data.user
        localStorage.setItem('token', response.data.token)
        message.success('登录成功')
        return true
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '登录失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (name: string, email: string, password: string) => {
    try {
      loading.value = true
      await api.post('/user', { name, email, password })
      message.success('注册成功，请登录')
      return true
    } catch (error: any) {
      message.error(error.response?.data?.error || '注册失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 登出
  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    message.success('已退出登录')
  }

  // 更新用户信息
  const updateUser = async (userData: Partial<User>) => {
    try {
      loading.value = true
      const response = await api.put(`/user/${user.value?.id}`, userData)
      user.value = { ...user.value!, ...response.data }
      message.success('更新成功')
      return true
    } catch (error: any) {
      message.error(error.response?.data?.error || '更新失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    if (!token.value) return
    
    try {
      const response = await api.get(`/user/${user.value?.id}`)
      user.value = response.data
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    login,
    register,
    logout,
    updateUser,
    fetchUserInfo
  }
})
