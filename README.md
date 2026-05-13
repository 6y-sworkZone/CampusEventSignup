# 校园活动报名系统

基于 Vue3 + Element Plus + Go + Gin + GORM + SQLite 的全栈校园活动报名系统。

## 功能特性

### 用户管理
- 学生注册（学号唯一、手机号11位、学院选择）
- 管理员注册（工号+用户管理权限）
- 登录失败3次锁定5分钟
- 记住我7天功能
- 头像上传
- 修改密码

### 活动管理
- 创建活动（标题2-50字、富文本描述、封面图上传、时间先后校验、报名截止、人数上限、多标签）
- 编辑活动（人数保护）
- 删除活动（二次确认）
- 活动列表（多条件筛选、搜索）
- 管理员强制结束活动

### 报名模块
- 报名弹窗确认
- 满员后候补报名（20%上限）
- 有人取消自动递补
- 签到需4位签到码
- 活动前30分钟开放签到
- 活动前2小时不可取消报名
- 我的报名按状态筛选展示

### 数据统计
- 管理员仪表盘（饼图展示学院分布、折线图展示7天趋势）
- 活动详情页展示报名统计
- 导出报名名单为Excel

## 技术栈

### 前端
- Vue 3
- Vue Router 4
- Pinia 状态管理
- Element Plus UI 组件库
- ECharts 图表库
- Axios HTTP 客户端
- Vite 构建工具

### 后端
- Go 1.21+
- Gin Web 框架
- GORM ORM
- SQLite 数据库
- JWT 认证
- excelize Excel 导出

## 项目结构

```
CampusEventSignup/
├── backend/                 # 后端项目
│   ├── config/             # 配置文件
│   │   └── config.go
│   ├── controllers/        # 控制器
│   │   ├── auth.go
│   │   ├── event.go
│   │   ├── registration.go
│   │   ├── stats.go
│   │   └── upload.go
│   ├── middleware/         # 中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── models/             # 数据模型
│   │   ├── user.go
│   │   ├── event.go
│   │   └── registration.go
│   ├── routes/             # 路由
│   │   └── routes.go
│   ├── utils/              # 工具函数
│   │   ├── jwt.go
│   │   └── response.go
│   ├── main.go             # 入口文件
│   ├── go.mod
│   └── go.sum
├── frontend/               # 前端项目
│   ├── src/
│   │   ├── views/          # 页面组件
│   │   │   ├── Login.vue
│   │   │   ├── Register.vue
│   │   │   ├── Home.vue
│   │   │   ├── Events.vue
│   │   │   ├── EventDetail.vue
│   │   │   ├── MyRegistrations.vue
│   │   │   ├── Profile.vue
│   │   │   ├── Layout.vue
│   │   │   └── admin/
│   │   │       ├── Dashboard.vue
│   │   │       ├── Events.vue
│   │   │       └── EventForm.vue
│   │   ├── stores/         # Pinia 状态
│   │   │   └── user.js
│   │   ├── router/         # 路由配置
│   │   │   └── index.js
│   │   ├── utils/          # 工具函数
│   │   │   └── request.js
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
└── README.md
```

## 快速开始

### 后端启动

1. 进入后端目录：
```bash
cd backend
```

2. 安装依赖：
```bash
go mod tidy
```

3. 启动后端服务：
```bash
go run main.go
```

后端将在随机端口（10000-60000之间）启动，控制台会输出端口号。

### 前端启动

1. 进入前端目录：
```bash
cd frontend
```

2. 安装依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm run dev
```

4. 修改 `vite.config.js` 中的代理配置，将目标端口改为后端实际运行的端口。

## 主要 API 接口

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/user/me` - 获取当前用户信息
- `POST /api/user/change-password` - 修改密码

### 活动相关
- `GET /api/events` - 获取活动列表
- `GET /api/events/:id` - 获取活动详情
- `POST /api/events` - 创建活动（管理员）
- `PUT /api/events/:id` - 编辑活动（管理员）
- `DELETE /api/events/:id` - 删除活动（管理员）
- `POST /api/events/:id/end` - 结束活动（管理员）
- `GET /api/events/:id/stats` - 获取活动统计

### 报名相关
- `POST /api/events/:id/register` - 报名活动
- `POST /api/events/:id/cancel` - 取消报名
- `POST /api/events/:id/signin` - 签到
- `GET /api/registrations/my` - 获取我的报名列表
- `GET /api/events/:id/registrations` - 获取活动报名列表（管理员）

### 统计相关
- `GET /api/stats/dashboard` - 获取仪表盘统计（管理员）
- `GET /api/stats/export` - 导出报名数据（管理员）

### 文件上传
- `POST /api/upload/avatar` - 上传头像
- `POST /api/upload/cover` - 上传活动封面（管理员）

## 核心特性说明

### 安全机制
- 密码使用 bcrypt 加密存储
- JWT Token 认证
- 登录失败锁定机制（3次锁定5分钟）
- CORS 跨域配置
- 权限控制（管理员/普通用户）

### 报名逻辑
- 报名人数限制
- 候补队列机制（上限为正式名额的20%）
- 自动递补机制（有人取消时候补自动转正）
- 活动前2小时不可取消报名
- 签到码验证（活动前30分钟开放）

### 数据导出
- 支持导出所有报名数据为 Excel
- 包含活动名称、姓名、学号/工号、学院、手机号、报名时间、签到状态等信息

## 注意事项

1. 后端启动时会自动创建 SQLite 数据库文件和上传目录
2. 前端开发时需要根据后端实际端口修改 Vite 代理配置
3. 生产环境建议配置 HTTPS 和适当的 JWT 密钥
4. 建议定期备份数据库文件

## 开发说明

本项目采用前后端分离架构，后端提供 RESTful API，前端通过 Axios 调用接口实现完整功能。
