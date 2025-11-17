# RedQuill - AI 小说生成平台

<div align="center">

![RedQuill Logo](redquill-frontend/public/logo.svg)

**基于 AI 技术的智能小说创作平台**

[功能特性](#功能特性) • [快速开始](#快速开始) • [开发指南](#开发指南) • [API 文档](#api-文档)

</div>

---

## 📖 项目简介

RedQuill 是一个基于 AI 技术的小说创作平台，支持从故事核心、世界观、角色设定到章节内容的完整创作流程。平台集成了多厂商 LLM 模型，提供流式生成、Prompt 模板管理等核心功能，帮助作者高效创作高质量小说内容。

### 核心特性

- 🤖 **多模型支持** - 集成 OpenAI、DeepSeek、豆包、千问、文心一言等多个 LLM 服务
- 📝 **智能生成** - 支持故事核心、世界观、角色、大纲、章节的 AI 生成
- 🔄 **流式生成** - 基于 SSE 的实时流式内容生成，提升用户体验
- 📚 **Prompt 引擎** - 灵活的 Prompt 模板系统，支持变量替换和动态生成
- 🎯 **创作流程** - 完整的创作工作流，从构思到成稿一站式服务
- 🔐 **用户系统** - JWT 认证、用户管理、权限控制

---

## 🛠 技术栈

### 后端

- **语言**: Go 1.21+
- **框架**: Gin (HTTP Web Framework)
- **数据库**: MongoDB
- **认证**: JWT (JSON Web Token)
- **架构**: RESTful API

### 前端

- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件**: Ant Design Vue
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios

### AI 集成

- **多厂商支持**: OpenAI、DeepSeek、豆包、千问、文心一言、Azure OpenAI、Ollama
- **流式响应**: Server-Sent Events (SSE)
- **统一接口**: 抽象化的 LLM 客户端接口

---

## ✨ 功能特性

### 🔐 用户管理

- 用户注册/登录
- JWT 身份验证
- 个人信息管理
- 用户权限控制

### 🤖 LLM 模型管理

- 多厂商 LLM 模型配置
- 模型测试和验证
- API Key 管理
- 使用统计和监控

### 📝 Prompt 模板管理

- Prompt 模板创建和编辑
- 变量支持 (`{variable_name}`)
- 模板分类和标签
- 模板版本管理
- 默认模板初始化

### 📚 小说项目管理

- 小说项目创建和管理
- 项目蓝图设置
- AI 上下文配置
- 项目状态跟踪

### 🎭 AI 内容生成

#### 故事核心生成
- 基于用户想法生成故事核心
- 包含核心冲突、主题、创新点等
- 支持批量生成和选择

#### 世界观构建
- 生成完整的世界观设定
- 包含修炼体系、社会结构、地理环境等
- 支持流式生成和实时预览

#### 角色塑造
- 创建有深度的角色档案
- 包含灵魂档案（性格、背景、动机）
- 核心属性（境界、能力、物品）
- 支持批量生成

#### 大纲生成
- 生成完整的小说大纲
- 包含章节信息、故事弧线、关键主题
- 支持章节级别的详细规划

#### 章节生成
- 基于大纲和上下文生成章节内容
- 自动填充前情提要（上一章节摘要和正文）
- 支持章节目标、角色发展、质量评估
- 流式生成实时显示

### 🚀 流式生成支持

- Server-Sent Events (SSE) 实时流式响应
- 避免超时问题
- 实时内容展示和解析
- 支持 JSON 格式解析和内容提取

### 📊 数据管理

- 故事核心列表和详情查看
- 世界观详情查看
- 角色列表和详情管理
- 大纲详情查看（支持编辑和删除）
- 章节列表和内容查看

---

## 📁 项目结构

```
RedQuill/
├── redquill-backend/          # 后端服务
│   ├── pkg/
│   │   ├── cmd/server/        # 入口文件
│   │   ├── config/            # 配置管理
│   │   ├── models/            # 数据模型
│   │   ├── services/           # 业务逻辑层
│   │   ├── handlers/           # HTTP 处理器
│   │   ├── routes/             # 路由配置
│   │   ├── middleware/         # 中间件
│   │   ├── common/             # 公共功能
│   │   └── utils/              # 工具类
│   │       └── llm/            # LLM 客户端
│   ├── docker-compose.yml      # Docker 配置
│   ├── env.example             # 环境变量示例
│   └── README.md               # 后端文档
│
├── redquill-frontend/          # 前端应用
│   ├── src/
│   │   ├── components/         # 组件
│   │   │   ├── generate/       # 生成组件
│   │   │   └── Layout.vue     # 布局组件
│   │   ├── views/              # 页面
│   │   ├── stores/             # 状态管理
│   │   ├── utils/              # 工具类
│   │   ├── router/             # 路由配置
│   │   └── styles/             # 样式文件
│   ├── package.json
│   ├── vite.config.ts
│   └── README.md               # 前端文档
│
├── note.md                     # 开发规范文档
└── README.md                   # 本文档
```

---

## 🚀 快速开始

### 环境要求

**后端**
- Go 1.21 或更高版本
- MongoDB 5.0+ (或使用 Docker)
- 环境变量配置

**前端**
- Node.js 18.0+ （推荐使用v22.15.0）
- npm 8.0+ 或 yarn

### 后端启动

1. **克隆项目**
```bash
git clone <repository-url>
cd RedQuill/redquill-backend
```

2. **配置环境变量**
```bash
cp env.example .env
# 编辑 .env 文件，配置 MongoDB、JWT 等
```

3. **启动 MongoDB** (使用 Docker)
```bash
docker-compose up -d
```

4. **安装依赖并运行**
```bash
go mod download
go run ./pkg/cmd/server
```

后端服务将在 `http://localhost:8080` 启动

### 前端启动

1. **进入前端目录**
```bash
cd redquill-frontend
```

2. **安装依赖**
```bash
npm install
```

3. **启动开发服务器**
```bash
npm run dev
```

前端应用将在 `http://localhost:3000` 启动

### 环境变量配置

**后端 (.env)**
```env
APP_ENV=development
PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=redquill
JWT_SECRET=your-secret-key
JWT_TTL_MIN=1440
```

**前端 (vite.config.ts 代理配置)**
```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

---

## 📚 开发指南

### 后端开发规范

#### 项目结构规范

- **models**: 数据模型定义 (`xxx_model.go`)
- **services**: 业务逻辑层 (`xxx_service.go`)
- **handlers**: HTTP 处理器 (`xxx_handler.go`)
- **routes**: 路由配置 (`routes.go`)
- **middleware**: 中间件
- **common**: 公共功能（如通用查询）
- **utils**: 工具类

#### 接口设计规范

遵循 RESTful 风格：

- 新增资源: `POST /api/v1/{resource}`
- 查询列表: `GET /api/v1/{resources}` (支持分页/排序/搜索)
- 查询详情: `GET /api/v1/{resource}/:id`
- 更新资源: `PUT /api/v1/{resource}/:id`
- 删除资源: `DELETE /api/v1/{resource}/:id`

#### 通用查询能力

所有列表接口支持以下查询参数：

- `page`: 页码，默认 1
- `pageSize`: 每页数量，默认 20，最大 100
- `sort`: 排序字段，逗号分隔，前缀 `-` 表示倒序。例如：`name,-ctime`
- `q`: 关键字，模糊匹配（大小写不敏感）

**Handler 侧用法**:
```go
page, size, sortExpr, q := common.ParseCommonQueryParams(c.Request.URL.Query())
items, total, err := service.ListPaged(ctx, page, size, sortExpr, q)
```

**返回结构**:
```json
{
  "items": [...],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 100,
    "totalPage": 5
  }
}
```

#### 命名规范

- **Handlers**: `PostUsersHandler`, `GetUsersHandler`, `ListUsersHandler`, `PutUsersHandler`, `DeleteUsersHandler`
- **Services**: 参考 handlers 命名
- **Models**: `xxx_model.go`
- **Services**: `xxx_service.go`
- **Handlers**: `xxx_handler.go`

#### 代码注释格式

```go
// Package services
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: {filename}.go
/@Description: 功能描述
/*/
```

#### 数据模型规范

所有数据模型默认包含：
- `ctime`: 创建时间（时间戳）
- `mtime`: 修改时间（时间戳）

### 前端开发规范

#### 添加新页面

1. 在 `src/views/` 创建 Vue 组件
2. 在 `src/router/index.ts` 添加路由
3. 在 `src/components/Layout.vue` 添加导航菜单

#### 添加新 API

1. 在 `src/utils/api.ts` 添加 API 方法
2. 在对应的 store 中添加状态管理
3. 在组件中调用 API

#### 状态管理

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

---

## 📡 API 文档

### 认证

所有需要认证的接口都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

### 用户管理

- `POST /api/v1/user` - 注册用户
- `POST /api/v1/login` - 用户登录
- `GET /api/v1/users` - 查询用户列表（支持分页/排序/搜索）
- `GET /api/v1/user/:id` - 查询用户详情
- `PUT /api/v1/user/:id` - 更新用户信息
- `DELETE /api/v1/user/:id` - 删除用户

### LLM 模型管理

- `POST /api/v1/llm-model` - 创建 LLM 模型
- `GET /api/v1/llm-models` - 查询模型列表
- `GET /api/v1/llm-model/:id` - 查询模型详情
- `PUT /api/v1/llm-model/:id` - 更新模型
- `DELETE /api/v1/llm-model/:id` - 删除模型
- `POST /api/v1/llm-model/:id/test` - 测试模型
- `POST /api/v1/llm-model/:id/service` - 调用模型服务

### Prompt 模板管理

- `POST /api/v1/prompt` - 创建 Prompt 模板
- `GET /api/v1/prompts` - 查询模板列表
- `GET /api/v1/prompt/:id` - 查询模板详情
- `PUT /api/v1/prompt/:id` - 更新模板
- `DELETE /api/v1/prompt/:id` - 删除模板

### 小说管理

- `POST /api/v1/novel` - 创建小说
- `GET /api/v1/novels` - 查询小说列表
- `GET /api/v1/novel/:id` - 查询小说详情
- `PUT /api/v1/novel/:id` - 更新小说
- `DELETE /api/v1/novel/:id` - 删除小说

### 故事核心

- `POST /api/v1/story-core` - 创建故事核心
- `GET /api/v1/story-cores/:novel_id` - 查询故事核心列表

### 世界观

- `POST /api/v1/worldview` - 创建世界观
- `GET /api/v1/worldview/:novel_id` - 查询世界观

### 角色管理

- `POST /api/v1/character` - 创建角色
- `GET /api/v1/characters/:novel_id` - 查询角色列表

### 大纲管理

- `POST /api/v1/outline` - 创建大纲
- `GET /api/v1/outlines/:novel_id` - 查询大纲列表
- `GET /api/v1/outline/:id` - 查询大纲详情
- `PUT /api/v1/outline/:id` - 更新大纲
- `DELETE /api/v1/outline/:id` - 删除大纲

### 章节管理

- `POST /api/v1/chapter` - 创建章节
- `GET /api/v1/chapters/:novel_id` - 查询章节列表
- `GET /api/v1/chapter/:id` - 查询章节详情

### AI 生成接口

所有生成接口都支持流式生成（通过 `stream: true` 参数）：

- `POST /api/v1/generate/story-core` - 生成故事核心
- `POST /api/v1/generate/worldview` - 生成世界观
- `POST /api/v1/generate/character` - 生成角色
- `POST /api/v1/generate/outline` - 生成大纲
- `POST /api/v1/generate/chapter` - 生成章节
- `POST /api/v1/generate/llm` - 通用 LLM 生成

**流式生成示例**:
```json
{
  "novel_id": "novel_id",
  "llm_model_id": "model_id",
  "input_data": {
    "chapter_number": 1,
    "chapter_goal": "章节目标",
    "previous_summary": "前情提要",
    "characters_involved": [...],
    "outline_id": "outline_id",
    "characters_outline": {...}
  },
  "stream": true
}
```

---

## 🎯 核心功能说明

### Prompt 模板系统

Prompt 模板支持变量替换，使用 `{variable_name}` 格式：

```go
Content: `你是{novel_title}的御用写手。
根据以下信息生成章节：
- 故事核心：{story_core}
- 世界观：{worldview}
- 章节目标：{chapter_goal}`
```

系统会自动从 `input_data` 中提取变量值并替换。

### 章节生成流程

1. **数据准备**: 自动获取故事核心、世界观、大纲信息
2. **上下文构建**: 整合章节目标、前情提要、角色信息
3. **Prompt 渲染**: 使用模板和变量生成完整 Prompt
4. **LLM 调用**: 调用选定的 LLM 模型生成内容
5. **结果解析**: 解析 JSON 格式的元数据和正文内容
6. **数据保存**: 保存章节到数据库

### 流式生成机制

使用 Server-Sent Events (SSE) 实现流式生成：

- **后端**: 通过 `GenerateWithLLMStream` 方法返回流式响应
- **前端**: 使用 `streamGenerate` 工具函数处理 SSE 事件
- **优势**: 避免超时、实时反馈、更好的用户体验

---

## 🚢 部署指南

### 后端部署

1. **构建应用**
```bash
cd redquill-backend
go build -o redquill-server ./pkg/cmd/server
```

2. **配置环境变量**
```bash
export APP_ENV=production
export PORT=8080
export MONGO_URI=mongodb://your-mongo-uri
export JWT_SECRET=your-secret-key
```

3. **运行服务**
```bash
./redquill-server
```

### 前端部署

1. **构建生产版本**
```bash
cd redquill-frontend
npm run build
```

2. **部署到 Nginx**

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        root /path/to/redquill-frontend/dist;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # SSE 支持
        proxy_set_header Connection '';
        proxy_http_version 1.1;
        chunked_transfer_encoding off;
        proxy_buffering off;
        proxy_cache off;
    }
}
```

### Docker 部署

**后端 Dockerfile** (示例):
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o redquill-server ./pkg/cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/redquill-server .
CMD ["./redquill-server"]
```

---

## 📝 开发规范

详细开发规范请参考 [note.md](./note.md)

### 关键规范摘要

1. **后端 RESTful 设计**: 遵循 REST 风格，统一接口命名
2. **通用查询能力**: 所有列表接口支持分页、排序、关键字搜索
3. **命名规范**: Handlers、Services 统一命名风格
4. **代码注释**: 遵循统一的注释格式
5. **数据模型**: 默认包含 `ctime` 和 `mtime` 字段

---

## 🤝 贡献指南

### 贡献流程

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### Issue 规范

在提交 Issue 之前，请先搜索是否已有相关 Issue，避免重复。

#### Issue 类型

- **Bug 报告** (`bug`): 报告系统错误或异常行为
- **功能建议** (`feature`): 提出新功能或改进建议
- **文档问题** (`docs`): 文档错误或需要改进的地方
- **性能问题** (`performance`): 性能相关问题
- **问题咨询** (`question`): 使用问题或技术咨询

#### Issue 标题格式

```
[类型] 简短描述
```

示例：
- `[bug] 章节生成时前情提要未正确填充`
- `[feature] 支持批量生成章节`
- `[docs] 更新 API 文档中的示例代码`

#### Issue 内容模板

**Bug 报告模板**:
```markdown
## 问题描述
清晰简洁地描述问题

## 复现步骤
1. 执行操作 A
2. 执行操作 B
3. 看到错误 C

## 预期行为
描述期望的正确行为

## 实际行为
描述实际发生的错误行为

## 环境信息
- 操作系统: 
- Node.js 版本: 
- Go 版本: 
- 浏览器版本: (如果是前端问题)

## 截图/日志
如果有相关截图或错误日志，请附上
```

**功能建议模板**:
```markdown
## 功能描述
清晰描述想要的功能

## 使用场景
描述这个功能的使用场景和价值

## 实现建议
如果有实现思路，可以描述

## 相关 Issue
如果有相关的 Issue，请链接
```

### Pull Request 规范

#### PR 标题格式

```
[类型] 简短描述
```

示例：
- `[feat] 添加章节批量生成功能`
- `[fix] 修复前情提要自动填充问题`
- `[docs] 更新 README 中的部署说明`

#### PR 描述模板

```markdown
## 变更类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能优化
- [ ] 其他

## 变更描述
详细描述本次 PR 的变更内容

## 相关 Issue
关联的 Issue: #123

## 测试说明
描述如何测试这些变更

## 检查清单
- [ ] 代码已通过 linter 检查
- [ ] 已添加必要的测试
- [ ] 已更新相关文档
- [ ] 已检查向后兼容性
- [ ] 已测试功能是否正常工作
```

#### PR 提交规范

提交信息应遵循以下格式：

```
[类型] 简短描述

详细描述（可选）
```

**类型说明**:
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整（不影响代码运行）
- `refactor`: 代码重构（既不是新功能也不是修复 bug）
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具相关
- `ci`: CI/CD 相关

**提交示例**:
```bash
feat: 添加章节批量生成功能

- 支持选择多个章节号批量生成
- 添加批量生成进度显示
- 优化生成结果的展示方式

Closes #123
```

#### PR 审查要求

1. **代码质量**
   - 代码已通过 linter 检查
   - 遵循项目代码规范
   - 添加必要的注释

2. **功能完整性**
   - 功能已完整实现
   - 已处理边界情况
   - 已添加错误处理

3. **测试覆盖**
   - 新功能已添加测试
   - 修复的 bug 已添加回归测试
   - 所有测试通过

4. **文档更新**
   - 更新了相关文档
   - API 变更已更新文档
   - 添加了使用示例（如需要）

5. **向后兼容**
   - 不破坏现有功能
   - API 变更考虑兼容性
   - 数据库变更提供迁移方案

#### PR 合并规范

- **Squash and Merge**: 功能完整的 PR 使用 squash merge
- **Merge Commit**: 大型功能或多人协作的 PR 使用 merge commit
- **Rebase and Merge**: 简单的 bug 修复可以使用 rebase merge

### 代码提交规范

提交信息格式：

```
[类型] 简短描述（50字符以内）

详细描述（可选，72字符换行）

相关 Issue: #123
```

**类型说明**:
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具相关
- `ci`: CI/CD 相关

**提交示例**:
```bash
feat: 添加章节生成的前情提要自动填充

当用户选择章节号时，自动从上一个章节获取摘要和正文内容，
填充到前情提要字段中。

Closes #456
```

---

## 📄 许可证

MIT License

---

## 👥 作者

- **urmsone** - urmsone@163.com

---

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- Email: urmsone@163.com
- Issue: [GitHub Issues](https://github.com/your-repo/issues)

---

<div align="center">

**Made with ❤️ by RedQuill Team**

</div>

