# RedQuill Frontend

基于 Vue3 + Vite + TypeScript + Ant Design Vue 的 AI 小说生成平台前端应用。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vite** - 下一代前端构建工具
- **TypeScript** - JavaScript 的超集
- **Ant Design Vue** - 企业级 UI 设计语言
- **Vue Router** - 官方路由管理器
- **Pinia** - Vue 状态管理库
- **Axios** - HTTP 客户端

## 功能特性

### 🔐 用户管理
- 用户注册/登录
- 个人信息管理
- JWT 身份验证

### 🤖 LLM 模型管理
- 多厂商 LLM 模型支持（OpenAI、DeepSeek、豆包、千问、文心一言等）
- 模型配置和测试
- 使用统计

### 📝 Prompt 模板管理
- 模板创建和编辑
- 变量支持
- 分类和标签管理

### 📚 小说管理
- 小说项目创建
- 项目蓝图管理
- AI 上下文设置

### 🎭 AI 内容生成
- **故事核心生成** - 基于用户想法生成故事核心
- **世界观构建** - 生成完整的世界观设定
- **角色塑造** - 创建有深度的角色档案
- **章节生成** - 自动生成小说章节内容
- **流式生成** - 实时显示生成过程

### 🚀 流式生成支持
- Server-Sent Events (SSE) 实时流式响应
- 避免超时问题
- 实时内容展示

## 项目结构

```
redquill-frontend/
├── src/
│   ├── components/          # 组件
│   │   ├── Layout.vue       # 主布局
│   │   └── generate/        # 生成组件
│   ├── views/               # 页面
│   │   ├── Login.vue        # 登录页
│   │   ├── Register.vue     # 注册页
│   │   ├── Dashboard.vue    # 仪表盘
│   │   ├── Novels.vue       # 小说管理
│   │   ├── NovelDetail.vue # 小说详情
│   │   ├── NovelGenerate.vue # AI生成
│   │   ├── LLMModels.vue    # LLM模型管理
│   │   ├── Prompts.vue      # Prompt管理
│   │   └── Users.vue        # 用户管理
│   ├── stores/              # 状态管理
│   │   ├── auth.ts          # 认证状态
│   │   └── novel.ts         # 小说状态
│   ├── utils/               # 工具类
│   │   └── api.ts           # API 封装
│   ├── router/              # 路由配置
│   ├── styles/              # 样式文件
│   ├── App.vue              # 根组件
│   └── main.ts              # 入口文件
├── package.json
├── vite.config.ts
├── tsconfig.json
└── README.md
```

## 快速开始

### 环境要求

- Node.js >= 18.0.0
- npm >= 8.0.0

### 安装依赖

```bash
npm install
```

### 开发环境

```bash
npm run dev
```

应用将在 `http://localhost:3000` 启动

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## API 集成

前端通过代理配置连接到后端 API：

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

## 流式生成使用

### 基本用法

```typescript
import { streamGenerate } from '@/utils/api'

// 流式生成故事核心
await streamGenerate(
  '/generate/story-core',
  {
    novel_id: 'novel_id',
    llm_model_id: 'model_id',
    input_data: { /* 输入数据 */ },
    stream: true
  },
  (content: string) => {
    // 处理流式内容
    console.log('收到内容:', content)
  },
  () => {
    // 生成完成
    console.log('生成完成')
  },
  (error: string) => {
    // 处理错误
    console.error('生成失败:', error)
  }
)
```

### 组件中使用

```vue
<template>
  <div class="stream-content">
    <pre>{{ streamContent }}</pre>
    <span v-if="streaming" class="stream-cursor">|</span>
  </div>
</template>

<script setup lang="ts">
import { streamGenerate } from '@/utils/api'

const streaming = ref(false)
const streamContent = ref('')

const handleStreamGenerate = async () => {
  streaming.value = true
  streamContent.value = ''
  
  await streamGenerate(
    '/generate/story-core',
    data,
    (content: string) => {
      streamContent.value += content
    },
    () => {
      streaming.value = false
    },
    (error: string) => {
      streaming.value = false
      console.error(error)
    }
  )
}
</script>
```

## 开发指南

### 添加新页面

1. 在 `src/views/` 创建 Vue 组件
2. 在 `src/router/index.ts` 添加路由
3. 在 `src/components/Layout.vue` 添加导航菜单

### 添加新 API

1. 在 `src/utils/api.ts` 添加 API 方法
2. 在对应的 store 中添加状态管理
3. 在组件中调用 API

### 状态管理

使用 Pinia 进行状态管理：

```typescript
// stores/example.ts
import { defineStore } from 'pinia'

export const useExampleStore = defineStore('example', () => {
  const data = ref([])
  const loading = ref(false)
  
  const fetchData = async () => {
    // 获取数据逻辑
  }
  
  return {
    data,
    loading,
    fetchData
  }
})
```

## 部署

### 构建

```bash
npm run build
```

构建产物将输出到 `dist/` 目录

### 部署到 Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        root /path/to/dist;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

MIT License
