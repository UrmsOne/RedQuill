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


