<template>
  <a-layout class="layout">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      class="sider"
      :width="256"
    >
      <div class="logo">
        <img v-if="!collapsed" src="/logo.svg" alt="RedQuill" class="logo-img" />
        <img v-else src="/logo-mini.svg" alt="RedQuill" class="logo-img-mini" />
        <h1 v-if="!collapsed" class="logo-text">RedQuill</h1>
      </div>
      
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        theme="dark"
        class="sidebar-menu"
        @click="handleMenuClick"
      >
        <a-menu-item key="dashboard">
          <DashboardOutlined />
          <span>仪表盘</span>
        </a-menu-item>
        
        <a-sub-menu key="content">
          <template #title>
            <FileTextOutlined />
            <span>内容管理</span>
          </template>
          <a-menu-item key="novels">
            <BookOutlined />
            <span>小说管理</span>
          </a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="ai">
          <template #title>
            <RobotOutlined />
            <span>AI工具</span>
          </template>
          <a-menu-item key="llm-models">
            <ApiOutlined />
            <span>LLM模型</span>
          </a-menu-item>
          <a-menu-item key="prompts">
            <CodeOutlined />
            <span>Prompt管理</span>
          </a-menu-item>
        </a-sub-menu>
        
        <a-menu-item key="users">
          <UserOutlined />
          <span>用户管理</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    
    <!-- 主内容区域 -->
    <a-layout class="main-layout">
      <!-- 顶部导航栏 -->
      <a-layout-header class="header">
        <div class="header-content">
          <div class="header-left">
            <a-button
              type="text"
              @click="collapsed = !collapsed"
              class="trigger"
            >
              <MenuUnfoldOutlined v-if="collapsed" />
              <MenuFoldOutlined v-else />
            </a-button>
          </div>
          
          <div class="header-right">
            <a-dropdown>
              <a-button type="text" class="user-button">
                <a-avatar :size="32" class="user-avatar">
                  {{ user?.name?.charAt(0)?.toUpperCase() }}
                </a-avatar>
                <span class="user-name">{{ user?.name }}</span>
                <DownOutlined />
              </a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="handleProfile">
                    <UserOutlined />
                    个人资料
                  </a-menu-item>
                  <a-menu-item @click="handleLogout">
                    <LogoutOutlined />
                    退出登录
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </a-layout-header>
      
      <!-- 页面内容 -->
      <a-layout-content class="content">
        <div class="page-container">
          <router-view />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  UserOutlined,
  LogoutOutlined,
  DownOutlined,
  DashboardOutlined,
  FileTextOutlined,
  BookOutlined,
  RobotOutlined,
  ApiOutlined,
  CodeOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined
} from '@ant-design/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const collapsed = ref(false)
const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>(['content', 'ai'])

// 根据当前路由设置选中的菜单项
watch(() => route.path, (path) => {
  if (path.startsWith('/dashboard')) {
    selectedKeys.value = ['dashboard']
    openKeys.value = []
  } else if (path.startsWith('/novels')) {
    selectedKeys.value = ['novels']
    openKeys.value = ['content']
  } else if (path.startsWith('/llm-models')) {
    selectedKeys.value = ['llm-models']
    openKeys.value = ['ai']
  } else if (path.startsWith('/prompts')) {
    selectedKeys.value = ['prompts']
    openKeys.value = ['ai']
  } else if (path.startsWith('/users')) {
    selectedKeys.value = ['users']
    openKeys.value = []
  }
}, { immediate: true })

// 菜单点击处理
const handleMenuClick = ({ key }: { key: string }) => {
  switch (key) {
    case 'dashboard':
      router.push('/app/dashboard')
      break
    case 'novels':
      router.push('/app/novels')
      break
    case 'llm-models':
      router.push('/app/llm-models')
      break
    case 'prompts':
      router.push('/app/prompts')
      break
    case 'users':
      router.push('/app/users')
      break
  }
}

const handleProfile = () => {
  // TODO: 实现个人资料页面
  console.log('个人资料')
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
}

/* 侧边栏样式 */
.sider {
  background: #001529;
  box-shadow: 2px 0 8px 0 rgba(29, 35, 41, 0.05);
}

.logo {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 8px;
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 12px;
}

.logo-img-mini {
  width: 24px;
  height: 24px;
  margin: 0 auto;
}

.logo-text {
  color: white;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.sidebar-menu {
  border: none;
  background: transparent;
}

.sidebar-menu :deep(.ant-menu-item) {
  margin: 4px 8px;
  border-radius: 6px;
  height: 40px;
  line-height: 40px;
}

.sidebar-menu :deep(.ant-menu-submenu) {
  margin: 4px 8px;
  border-radius: 6px;
}

.sidebar-menu :deep(.ant-menu-submenu-title) {
  height: 40px;
  line-height: 40px;
  border-radius: 6px;
}

.sidebar-menu :deep(.ant-menu-item-selected) {
  background: #1890ff;
}

.sidebar-menu :deep(.ant-menu-item:hover) {
  background: rgba(255, 255, 255, 0.1);
}

/* 主布局样式 */
.main-layout {
  background: #f0f2f5;
}

/* 顶部导航栏样式 */
.header {
  background: white;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 10;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.header-left {
  display: flex;
  align-items: center;
}

.trigger {
  font-size: 18px;
  color: #666;
  border: none;
  background: transparent;
  padding: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.3s;
}

.trigger:hover {
  background: #f5f5f5;
  color: #1890ff;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-button {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border: none;
  background: transparent;
  border-radius: 6px;
  transition: all 0.3s;
}

.user-button:hover {
  background: #f5f5f5;
}

.user-avatar {
  margin-right: 8px;
  background: #1890ff;
  color: white;
  font-weight: 600;
}

.user-name {
  margin-right: 8px;
  font-weight: 500;
  color: #333;
}

/* 内容区域样式 */
.content {
  margin: 24px;
  padding: 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  min-height: calc(100vh - 112px);
}

.page-container {
  width: 100%;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sider {
    position: fixed;
    height: 100vh;
    z-index: 100;
  }
  
  .main-layout {
    margin-left: 0;
  }
  
  .content {
    margin: 16px;
    padding: 16px;
  }
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .main-layout {
    background: #141414;
  }
  
  .header {
    background: #1f1f1f;
    border-bottom: 1px solid #303030;
  }
  
  .content {
    background: #1f1f1f;
    border: 1px solid #303030;
  }
  
  .user-name {
    color: #fff;
  }
}
</style>
