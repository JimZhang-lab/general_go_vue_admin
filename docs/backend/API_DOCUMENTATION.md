# 📚 API接口文档

Go后端RESTful API接口完整文档，包含认证、用户管理、系统管理等所有接口说明。

## 🔗 基础信息

### 服务地址
- **开发环境**: `http://localhost:8080`
- **生产环境**: `https://your-domain.com`

### 请求格式
- **Content-Type**: `application/json`
- **字符编码**: `UTF-8`
- **认证方式**: `Bearer Token (JWT)`

### 响应格式
```json
{
  "code": 200,
  "message": "成功",
  "data": {}
}
```

### 状态码说明
| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权或Token过期 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 🔐 认证接口

### 获取验证码
```http
GET /api/captcha
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "idKey": "captcha-id-123",
    "image": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
  }
}
```

### 用户登录
```http
POST /api/login
```

**请求参数**:
```json
{
  "username": "admin",
  "password": "admin123",
  "image": "1234",
  "idKey": "captcha-id-123"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "sysAdmin": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "email": "admin@example.com",
      "phone": "13800138000",
      "status": 1
    }
  }
}
```

### 用户登出
```http
POST /api/logout
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登出成功",
  "data": null
}
```

## 👥 用户管理接口

### 获取用户列表
```http
GET /api/admin/list?page=1&pageSize=10&username=admin
Authorization: Bearer {token}
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| pageSize | int | 否 | 每页数量，默认10 |
| username | string | 否 | 用户名筛选 |
| status | int | 否 | 状态筛选 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "管理员",
        "email": "admin@example.com",
        "phone": "13800138000",
        "status": 1,
        "createTime": "2025-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

### 获取用户详情
```http
GET /api/admin/{id}
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "管理员",
    "email": "admin@example.com",
    "phone": "13800138000",
    "status": 1,
    "deptId": 1,
    "postId": 1,
    "createTime": "2025-01-01T00:00:00Z"
  }
}
```

### 创建用户
```http
POST /api/admin
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "username": "newuser",
  "password": "password123",
  "nickname": "新用户",
  "email": "newuser@example.com",
  "phone": "13800138001",
  "deptId": 1,
  "postId": 1,
  "status": 1
}
```

### 更新用户
```http
PUT /api/admin/{id}
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "nickname": "更新的昵称",
  "email": "updated@example.com",
  "phone": "13800138002",
  "status": 1
}
```

### 删除用户
```http
DELETE /api/admin/{id}
Authorization: Bearer {token}
```

## 🏢 部门管理接口

### 获取部门列表
```http
GET /api/dept/list
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": [
    {
      "id": 1,
      "name": "技术部",
      "parentId": 0,
      "sort": 1,
      "status": 1,
      "children": [
        {
          "id": 2,
          "name": "前端组",
          "parentId": 1,
          "sort": 1,
          "status": 1
        }
      ]
    }
  ]
}
```

### 创建部门
```http
POST /api/dept
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "name": "新部门",
  "parentId": 1,
  "sort": 1,
  "status": 1
}
```

## 📋 角色管理接口

### 获取角色列表
```http
GET /api/role/list
Authorization: Bearer {token}
```

### 创建角色
```http
POST /api/role
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "name": "新角色",
  "code": "NEW_ROLE",
  "description": "角色描述",
  "status": 1,
  "menuIds": [1, 2, 3]
}
```

## 📊 系统监控接口

### 健康检查
```http
GET /health
```

**响应示例**:
```json
{
  "status": "healthy",
  "timestamp": "2025-07-29T18:00:00Z",
  "services": {
    "database": "healthy",
    "redis": "healthy",
    "rabbitmq": "healthy"
  }
}
```

### 性能指标
```http
GET /metrics
```

**响应示例**:
```json
{
  "http": {
    "totalRequests": 1000,
    "requestsPerSec": 50.5,
    "avgResponseTime": "85ms",
    "errorRate": 0.5
  },
  "system": {
    "memoryPercent": 45.2,
    "goroutineCount": 150,
    "gcCount": 25
  }
}
```

### QPS信息
```http
GET /qps
```

**响应示例**:
```json
{
  "current": 120.5,
  "peak": 500.0,
  "average": 85.2,
  "timestamp": "2025-07-29T18:00:00Z"
}
```

## 📝 日志管理接口

### 获取操作日志
```http
GET /api/logs/operation?page=1&pageSize=10
Authorization: Bearer {token}
```

### 获取登录日志
```http
GET /api/logs/login?page=1&pageSize=10
Authorization: Bearer {token}
```

## 🔧 错误处理

### 错误响应格式
```json
{
  "code": 400,
  "message": "请求参数错误",
  "data": null,
  "errors": [
    {
      "field": "username",
      "message": "用户名不能为空"
    }
  ]
}
```

### 常见错误码
| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 10001 | 验证码错误 | 重新获取验证码 |
| 10002 | 用户名或密码错误 | 检查登录信息 |
| 10003 | Token已过期 | 重新登录 |
| 10004 | 权限不足 | 联系管理员 |
| 10005 | 参数验证失败 | 检查请求参数 |

## 🧪 接口测试

### 使用curl测试
```bash
# 获取验证码
curl -X GET http://localhost:8080/api/captcha

# 用户登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","image":"1234","idKey":"test-key"}'

# 获取用户列表
curl -X GET http://localhost:8080/api/admin/list \
  -H "Authorization: Bearer your-token-here"
```

### 使用Postman
1. 导入API集合文件
2. 设置环境变量
3. 配置认证Token
4. 执行接口测试

---

**最后更新**: 2025-07-29  
**API版本**: v1.0  
**维护者**: 后端开发团队
