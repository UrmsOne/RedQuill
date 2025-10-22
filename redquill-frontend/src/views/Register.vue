<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <h1>注册账号</h1>
        <p>加入RedQuill，开启AI创作之旅</p>
      </div>
      
      <a-form
        :model="form"
        :rules="rules"
        @finish="handleRegister"
        layout="vertical"
        class="register-form"
      >
        <a-form-item label="用户名" name="name">
          <a-input
            v-model:value="form.name"
            placeholder="请输入用户名"
            size="large"
          />
        </a-form-item>
        
        <a-form-item label="邮箱" name="email">
          <a-input
            v-model:value="form.email"
            placeholder="请输入邮箱"
            size="large"
          />
        </a-form-item>
        
        <a-form-item label="密码" name="password">
          <a-input-password
            v-model:value="form.password"
            placeholder="请输入密码"
            size="large"
          />
        </a-form-item>
        
        <a-form-item label="确认密码" name="confirmPassword">
          <a-input-password
            v-model:value="form.confirmPassword"
            placeholder="请再次输入密码"
            size="large"
          />
        </a-form-item>
        
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="loading"
          >
            注册
          </a-button>
        </a-form-item>
        
        <div class="register-footer">
          <a @click="$router.push('/login')">已有账号？立即登录</a>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  name: [
    { required: true, message: '请输入用户名' },
    { min: 2, message: '用户名至少2位' }
  ],
  email: [
    { required: true, message: '请输入邮箱' },
    { type: 'email', message: '请输入有效的邮箱地址' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6位' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    {
      validator: (_: any, value: string) => {
        if (value && value !== form.password) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      }
    }
  ]
}

const loading = computed(() => authStore.loading)

const handleRegister = async () => {
  const success = await authStore.register(form.name, form.email, form.password)
  if (success) {
    router.push('/login')
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-card {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-header h1 {
  font-size: 28px;
  font-weight: bold;
  color: #1890ff;
  margin-bottom: 8px;
}

.register-header p {
  color: #666;
  font-size: 14px;
}

.register-form {
  margin-top: 24px;
}

.register-footer {
  text-align: center;
  margin-top: 16px;
}

.register-footer a {
  color: #1890ff;
  text-decoration: none;
}

.register-footer a:hover {
  text-decoration: underline;
}
</style>
