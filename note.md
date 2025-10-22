# 开发笔记记录常见问题
## 前端
### 开发指南
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


## 后端
- 后端满足restful风格的设计，每种用户资源（数据模型），默认有两个属性参数，ctime(创建时间)、mtime(修改时间)；ctime和mtime为时间戳；
- 用户模型存放在pkg/models；业务逻辑在pkg/services；pkg/handlers为路由函数，主要作用为参数校验，权限控制等。pkg/routes为路由与路由函数的映射
- 当开发新功能时，先生成models，在写对应的services、再写对应的handlers、再更新routes。
- 中间件存放在pkg/middleware文件夹；pkg/config配置文件相关；pkg/common存放公共功能的代码；pkg/cmd为入口文件
- post/put方法的请求和响应的结构体，可以在models/xxx_model.go里面定义

### 接口设计规范
案例如下：
新增用户 :`Post: /api/v1/user`
查询用户列表: `Get /api/v1/users`,可以使用通用查询能力默认支持以下功能
查询用户详情: `Get /api/v1/user/<user:id>`
修复用户信息: `Put /api/v1/user/<user:id>`
删除用户: `Delete /api/v1/user/<user:id>`

每个新功能都需要实现以上接口，并在`Get /api/v1/users`api实现通用查询功能

### 通用查询能力（分页 / 排序 / 关键字）

- 目录：`pkg/common/query.go`
- 查询参数：
    - `page`: 第几页，默认 1
    - `pageSize`: 每页数量，默认 20，最大 100
    - `sort`: 排序字段，逗号分隔，前缀 `-` 表示倒序。例如：`name,-ctime`
    - `q`: 关键字，模糊匹配（大小写不敏感）

- Handler 侧用法（示例 `users` 列表）：
    - 解析参数：`page, size, sortExpr, q := common.ParseCommonQueryParams(c.Request.URL.Query())`
    - 服务层应用：`ListPaged(ctx, page, size, sortExpr, q)`

- 返回结构：
    - `items`: 数据数组
    - `pagination`: `{ page, pageSize, total, totalPage }`

### 命名规范
handlers命名：
- 新增用户：PostUsersHandler，与接口设计的新增用户一致
- 查询用户详情：GetUsersHandler，与接口设计的查询用户详情一致
- 分页查询用户列表：ListUsersHandler，与接口设计的分页查询用户列表详情一致
- 修改用户信息：PutUsersHandler，与接口设计的修改用户信息一致
- 删除用户：DeleteUsersHandler，与接口设计的删除用户一致

services命名：
也参考handlers命名

每个go文件开头的注视格式如下：
```
// Package {package name} 
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: {filename}.go
/@Description:
/*/

```

pkg/models文件夹下的文件命名规范：`xxx_model.go`
pkg/services文件夹下的文件命名规范：`xxx_service.go`
pkg/handlers：`xxx_handler.go`

### 新功能开发规范
- 开发业务新功能时，需要提供上述给出的完整接口设计和实现。命名规范，代码规范遵守上述要求。
- 通用组件开发时，命名规范，代码规范遵守上述要求，相关文件要放置到指定的文件夹下

