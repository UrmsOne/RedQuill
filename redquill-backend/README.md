## RedQuill Backend (Gin + MongoDB)

### Project Structure

```
redquill-backend/
  pkg/
    cmd/server/
      main.go
    config/
      config.go
    handlers/
      health.go
      user.go
    middleware/
      request_id.go
    routes/
      routes.go
    server/
      server.go
    services/
      user.go
    utils/
      mongo.go
  docker-compose.yml
  env.example
  go.mod
```

### Quick Start

- Copy env: `cp env.example .env`
- Start Mongo: `docker compose up -d`
- Run server:
  - `go run ./pkg/cmd/server`
  - Listens on `:8080` by default

### Endpoints

- Health: `GET /healthz`
- List users: `GET /api/v1/users`
- Register user: `POST /api/v1/user`
- Get user detail: `GET /api/v1/user/:id` (JWT)
- Update user: `PUT /api/v1/user/:id` (JWT)
- Delete user: `DELETE /api/v1/user/:id` (JWT)
- Login: `POST /api/v1/login` -> returns `{ token }`

### LLM Model Management (JWT Required)

- Create LLM model: `POST /api/v1/llm-model`
- List LLM models: `GET /api/v1/llm-models` (with pagination/sort/search)
- Get LLM model: `GET /api/v1/llm-model/:id`
- Update LLM model: `PUT /api/v1/llm-model/:id`
- Delete LLM model: `DELETE /api/v1/llm-model/:id`
- Test LLM model: `POST /api/v1/llm-model/:id/test`
- Use LLM model service: `POST /api/v1/llm-model/:id/service`

### Prompt Management (JWT Required)

- Create prompt: `POST /api/v1/prompt`
- List prompts: `GET /api/v1/prompts` (with pagination/sort/search)
- Get prompt: `GET /api/v1/prompt/:id`
- Update prompt: `PUT /api/v1/prompt/:id`
- Delete prompt: `DELETE /api/v1/prompt/:id`

### Novel Management (JWT Required)

- Create novel: `POST /api/v1/novel`
- List novels: `GET /api/v1/novels` (with pagination/sort/search)
- Get novel: `GET /api/v1/novel/:id`
- Update novel: `PUT /api/v1/novel/:id`
- Delete novel: `DELETE /api/v1/novel/:id`

### Story Development (JWT Required)

- Create story core: `POST /api/v1/story-core` (novel_id in request body)
- Get story cores: `GET /api/v1/story-cores/:novel_id`
- Create worldview: `POST /api/v1/worldview` (novel_id in request body)
- Get worldview: `GET /api/v1/worldview/:novel_id`
- Create character: `POST /api/v1/character` (novel_id in request body)
- Get characters: `GET /api/v1/characters/:novel_id`

### Chapter Management (JWT Required)

- Create chapter: `POST /api/v1/chapter` (novel_id in request body)
- Get chapters: `GET /api/v1/chapters/:novel_id`
- Get chapter: `GET /api/v1/chapter/:id`
- Create writing session: `POST /api/v1/writing-session` (novel_id in request body)
- Get writing session: `GET /api/v1/writing-session/:novel_id`

### AI Generation (JWT Required)

- Generate story core: `POST /api/v1/generate/story-core`
- Generate worldview: `POST /api/v1/generate/worldview`
- Generate character: `POST /api/v1/generate/character`
- Generate chapter: `POST /api/v1/generate/chapter`
- General LLM generation: `POST /api/v1/generate/llm`

### Auth
- JWT Bearer via `Authorization: Bearer <token>`
- Env:
  - `JWT_SECRET`
  - `JWT_TTL_MIN` minutes

### 通用查询能力（分页 / 排序 / 关键字）

- 目录：`pkg/common/query.go`
- 查询参数：
  - `page`: 第几页，默认 1
  - `pageSize`: 每页数量，默认 20，最大 100
  - `sort`: 排序字段，逗号分隔，前缀 `-` 表示倒序。例如：`name,-createdAt`
  - `q`: 关键字，模糊匹配（大小写不敏感）

- Handler 侧用法（示例 `users` 列表）：
  - 解析参数：`page, size, sortExpr, q := common.ParseCommonQueryParams(c.Request.URL.Query())`
  - 服务层应用：`ListPaged(ctx, page, size, sortExpr, q)`

- 返回结构：
  - `items`: 数据数组
  - `pagination`: `{ page, pageSize, total, totalPage }`

### Config

- `APP_ENV`: `development` | `production`
- `PORT`: default `8080`
- `MONGO_URI`: default `mongodb://localhost:27017`
- `MONGO_DB`: default `redquill`

### Notes

- Minimal layered scaffold: handlers -> services -> db
- Add new modules by following the `user` example


