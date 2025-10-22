import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/Layout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/app/dashboard'
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/app',
      component: Layout,
      children: [
        {
          path: '',
          redirect: '/app/dashboard'
        },
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/Users.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'llm-models',
          name: 'LLMModels',
          component: () => import('@/views/LLMModels.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'prompts',
          name: 'Prompts',
          component: () => import('@/views/Prompts.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'novels',
          name: 'Novels',
          component: () => import('@/views/Novels.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'novel/:id',
          name: 'NovelDetail',
          component: () => import('@/views/NovelDetail.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'novel/:id/generate',
          name: 'NovelGenerate',
          component: () => import('@/views/NovelGenerate.vue'),
          meta: { requiresAuth: true }
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      // 未登录，重定向到登录页
      next('/login')
      return
    }
  }
  
  // 如果已登录且访问登录/注册页面，重定向到仪表盘
  if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
    next('/app/dashboard')
    return
  }
  
  // 其他情况正常导航
  next()
})

export default router
