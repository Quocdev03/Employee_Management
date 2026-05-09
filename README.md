# 🏢 Employee Management System

> Hệ thống quản lý nhân viên nội bộ — **Vue 3** + **Golang** + **MySQL**

---

## 📋 Mục lục

- [Tổng quan](#tổng-quan)
- [Tech Stack](#tech-stack)
- [Tính năng](#tính-năng)
- [Cấu trúc Project](#cấu-trúc-project)
- [Cài đặt & Chạy](#cài-đặt--chạy)
- [API Reference](#api-reference)
- [Phân quyền](#phân-quyền)

---

## 🎯 Tổng quan

Hệ thống quản lý nhân viên dành cho **nội bộ doanh nghiệp** — không cho phép tự đăng ký, chỉ Admin tạo tài khoản mới.

| Đặc điểm      | Chi tiết                   |
| ------------- | -------------------------- |
| Loại ứng dụng | Web nội bộ (Internal Tool) |
| Kiến trúc     | REST API + SPA             |
| Xác thực      | JWT (24h expiry)           |
| Phân quyền    | Role-based: Admin / User   |

---

## 🛠 Tech Stack

### Backend

| Thư viện              | Mục đích                                      |
| --------------------- | --------------------------------------------- |
| **Golang 1.22**       | Ngôn ngữ chính                                |
| **Gin**               | HTTP Framework — routing, middleware, binding |
| **GORM**              | ORM — query, migrate, soft delete             |
| **golang-jwt/jwt v5** | Tạo & verify JWT token                        |
| **bcrypt**            | Hash mật khẩu (cost=10)                       |
| **godotenv**          | Load biến môi trường từ `.env`                |

### Frontend

| Thư viện            | Mục đích                                    |
| ------------------- | ------------------------------------------- |
| **Vue 3**           | Framework UI                                |
| **Composition API** | Logic tái sử dụng, gọn hơn Options API      |
| **Vite**            | Build tool, HMR nhanh                       |
| **Pinia**           | State management (thay thế Vuex)            |
| **Vue Router 4**    | Client-side routing + navigation guard      |
| **Axios**           | HTTP client + interceptor tự động gắn token |

### Infrastructure

| Công nghệ          | Mục đích                             |
| ------------------ | ------------------------------------ |
| **MySQL 8.0**      | Database chính                       |
| **Docker**         | Container hoá ứng dụng               |
| **Docker Compose** | Orchestrate app + db với healthcheck |

---

## ✨ Tính năng

### 🔐 Authentication

- Đăng nhập bằng email + mật khẩu
- JWT token tự động gắn vào mọi request (Axios interceptor)
- Token hết hạn (401) → tự động redirect về trang Login, xoá localStorage
- Đăng xuất xoá toàn bộ session

### 📊 Dashboard

- Thống kê: tổng nhân viên, đang làm việc, số phòng ban
- Số tài khoản hệ thống _(chỉ Admin thấy)_
- Biểu đồ cột phân bổ nhân viên theo từng phòng ban

### 👥 Quản lý nhân viên

- Danh sách dạng bảng với avatar chữ cái đầu
- **Tìm kiếm realtime** theo tên / email (debounce 300ms)
- **Lọc** theo phòng ban
- **Phân trang** — 10 bản ghi / trang
- **Thêm / Sửa** nhân viên — dùng chung 1 modal component (EmployeeForm)
- **Xoá** nhân viên — soft delete (có thể khôi phục)
- Loading skeleton khi đang tải, trạng thái rỗng khi không có dữ liệu

### 🔑 Quản lý tài khoản _(Admin only)_

- Danh sách tài khoản hệ thống
- Tạo tài khoản mới (chọn role Admin hoặc User)
- Cập nhật email, mật khẩu, role
- Xoá tài khoản (không thể tự xoá chính mình)

### 🧭 Layout & UX

- **Sidebar** thu gọn / mở rộng — thiết kế dễ thêm menu mới
- **Header** hiển thị email, role badge, nút đăng xuất
- **Toast notification** 3 loại: thành công (xanh), lỗi (đỏ), cảnh báo (vàng)
- Menu **"Tài khoản"** chỉ hiện với Admin (cả sidebar và route guard)

---

## 🔄 Luồng hoạt động

### 1. Luồng Đăng nhập

```
Nhập email + password
        │
        ▼
POST /api/auth/login
        │
        ├─ Không tìm thấy email ──► 401
        ├─ Sai mật khẩu ──────────► 401
        │
        ▼
Tạo JWT { user_id, role, exp: +24h }
Ký bằng JWT_SECRET → trả token
        │
        ▼
Lưu token + user vào localStorage
Pinia authStore cập nhật state
        │
        ▼
Router redirect → /dashboard
```

### 2. Luồng Gọi API (có xác thực)

```
Component gọi store action
        │
        ▼
Axios interceptor gắn tự động:
  Authorization: Bearer <token>
        │
        ▼
JWT Middleware: parse + verify token
        │
        ├─ Token lỗi / hết hạn ──► 401
        │                              │
        │                              ▼
        │                   Axios interceptor bắt 401
        │                   Xoá localStorage → /login
        ▼
Gắn user_id, role vào gin.Context
Handler xử lý → trả response
        │
        ▼
Pinia store cập nhật → UI re-render
```

### 3. Luồng Phân quyền Admin

```
Request → /api/users/*
        │
        ▼
JWT Middleware ✅
        │
        ▼
AdminOnly Middleware:
  Lấy role từ gin.Context
        │
        ├─ role != "admin" ──► 403 Forbidden
        │
        ▼
Handler thực thi bình thường
```

### 4. Luồng Tìm kiếm Nhân viên

```
Gõ vào ô search
        │
        ▼
debounce 300ms
        │
        ▼
GET /api/employees?search=...&department_id=...&page=1&limit=10
        │
        ▼
WHERE (name LIKE '%...%' OR email LIKE '%...%')
  AND department_id = ?
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 10 OFFSET 0
        │
        ▼
{ data: [...], total: N, page: 1, limit: 10 }
        │
        ▼
Pinia store cập nhật
Phân trang tính lại totalPages
```

### 5. Luồng Thêm / Sửa Nhân viên

```
Click "Thêm" hoặc "Sửa"
        │
        ▼
EmployeeForm modal mở
  isEdit = true nếu truyền employee prop
        │
        ▼
Submit form
        │
        ▼
validate(): name, email, salary > 0
        │
        ├─ Lỗi → hiển thị dưới field tương ứng
        │
        ▼
store.create() hoặc store.update()
→ POST hoặc PUT /api/employees[/:id]
        │
        ├─ Lỗi API → hiển thị message lỗi
        │
        ▼
Đóng modal → Toast thành công → Reload list
```

### 6. Luồng Xoá Nhân viên (Soft Delete)

```
Click "Xoá" → confirm dialog
        │
        ▼
DELETE /api/employees/:id
        │
        ▼
GORM: SET deleted_at = NOW()
  (bản ghi vẫn còn trong DB)
        │
        ▼
Mọi query tự động thêm
  WHERE deleted_at IS NULL
        │
        ▼
Toast thành công → Reload list
```

---

### Dữ liệu seed mặc định

| Bảng        | Dữ liệu                                       |
| ----------- | --------------------------------------------- |
| roles       | admin (id=1), user (id=2)                     |
| departments | IT, HR, Finance, Marketing, Sales, Operations |
| users       | admin@company.com / admin123                  |
| employees   | 8 nhân viên mẫu                               |

---

## 📁 Cấu trúc Project

```
employee-management/
├── docker-compose.yml
├── README.md
│
├── backend/
│   ├── Dockerfile
│   ├── .env
│   ├── go.mod
│   ├── main.go                  # Entry: LoadEnv → ConnectDB → Migrate → Seed → Gin
│   ├── config/
│   │   └── config.go            # MySQL DSN, connection pool
│   ├── models/
│   │   ├── role.go
│   │   ├── user.go              # FK: RoleID → roles
│   │   ├── department.go
│   │   └── employee.go          # FK: DepartmentID → departments, soft delete
│   ├── database/
│   │   └── database.go          # AutoMigrate + Seed (idempotent)
│   ├── middleware/
│   │   ├── jwt.go               # AuthMiddleware, Claims struct
│   │   ├── admin.go             # AdminOnly middleware
│   │   └── error.go             # Error handling middleware
│   ├── utils/
│   │   └── response.go          # Success / Created / BadRequest / Forbidden...
│   ├── handlers/
│   │   ├── auth.go              # POST /auth/login
│   │   ├── employee.go          # GET/POST/PUT/DELETE /employees
│   │   ├── user.go              # GET/POST/PUT/DELETE /users
│   │   └── department.go        # GET /departments
│   └── routes/
│       └── routes.go            # Group routes + apply middleware
│
└── frontend/
    ├── .env                     # VITE_API_URL
    ├── package.json
    └── src/
        ├── main.js              # createApp + Pinia + Router
        ├── App.vue              # RouterView + ToastNotification
        ├── api/
        │   └── index.js         # Axios + interceptors (token, 401 handler)
        ├── layouts/
        │   └── MainLayout.vue   # AppSidebar + AppHeader + RouterView
        ├── router/
        │   └── index.js         # Routes + beforeEach (isLoggedIn, adminOnly)
        ├── stores/
        │   ├── auth.js          # login(), logout(), isLoggedIn, user
        │   ├── employee.js      # fetchAll(), create(), update(), remove()
        │   ├── user.js          # fetchAll(), create(), update(), remove()
        │   ├── department.js    # fetchAll() — dùng cho dropdown form
        │   └── toast.js         # success(), error(), warning()
        ├── views/
        │   ├── LoginView.vue
        │   ├── DashboardView.vue
        │   ├── EmployeeListView.vue
        │   └── UserManagementView.vue
        └── components/
            ├── AppSidebar.vue          # Thu gọn/mở rộng, menu theo role
            ├── AppHeader.vue           # User info + logout
            ├── EmployeeForm.vue        # Modal dùng chung thêm + sửa
            └── ToastNotification.vue   # Teleport to body, TransitionGroup
```

## 🚀 Cài đặt & Chạy

### Yêu cầu

- [Docker](https://www.docker.com/) & Docker Compose

### 1. Clone & cấu hình

```bash
git clone https://github.com/yourname/employee-management
cd employee-management
```

Kiểm tra `backend/.env`:

### 2. Chạy bằng Docker

```bash
docker-compose up --build
```

> Tự động: tạo DB → migrate → seed dữ liệu mẫu

### 3. Chạy Frontend (dev mode)

```bash
cd frontend
npm install
npm run dev
```

| Service        | URL                   |
| -------------- | --------------------- |
| 🌐 Frontend    | http://localhost:5173 |
| ⚙️ Backend API | http://localhost:8080 |

### Reset toàn bộ (xoá DB)

```bash
docker-compose down -v
docker-compose up --build
```

---

## 📡 API Reference

### Header xác thực

```
Authorization: Bearer <jwt_token>
```

### Auth

| Method | Endpoint          | Body                  | Auth |
| ------ | ----------------- | --------------------- | ---- |
| POST   | `/api/auth/login` | `{ email, password }` | ❌   |

### Employees

| Method | Endpoint             | Params / Body                                                        | Auth |
| ------ | -------------------- | -------------------------------------------------------------------- | ---- |
| GET    | `/api/employees`     | `?page&limit&search&department_id&status`                            | ✅   |
| GET    | `/api/employees/:id` | —                                                                    | ✅   |
| POST   | `/api/employees`     | `{ name, email, phone, department_id, position, salary, hire_date }` | ✅   |
| PUT    | `/api/employees/:id` | Tương tự POST + `status`                                             | ✅   |
| DELETE | `/api/employees/:id` | —                                                                    | ✅   |

### Departments

| Method | Endpoint           | Auth |
| ------ | ------------------ | ---- |
| GET    | `/api/departments` | ✅   |

### Users _(Admin only)_

| Method | Endpoint         | Body                              | Auth  |
| ------ | ---------------- | --------------------------------- | ----- |
| GET    | `/api/users`     | —                                 | Admin |
| POST   | `/api/users`     | `{ email, password, role_id }`    | Admin |
| PUT    | `/api/users/:id` | `{ email?, password?, role_id? }` | Admin |
| DELETE | `/api/users/:id` | —                                 | Admin |

---

## 🛡 Phân quyền

| Tính năng                     | User | Admin |
| ----------------------------- | ---- | ----- |
| Đăng nhập / Đăng xuất         | ✅   | ✅    |
| Xem Dashboard                 | ✅   | ✅    |
| Xem / Tìm kiếm nhân viên      | ✅   | ✅    |
| Thêm / Sửa / Xoá nhân viên    | ❌   | ✅    |
| Xem thống kê tài khoản        | ❌   | ✅    |
| Quản lý tài khoản hệ thống    | ❌   | ✅    |
| Menu "Tài khoản" trên sidebar | ❌   | ✅    |

> **Double protection:** phân quyền kiểm tra ở cả 2 tầng:
>
> - **Frontend:** `router.beforeEach` + `v-if="auth.user?.role === 'admin'"`
> - **Backend:** `AdminOnly()` middleware trên toàn bộ `/api/users/*`

---

## 🔐 Tài khoản mặc định

| Email             | Password | Role  |
| ----------------- | -------- | ----- |
| admin@company.com | admin123 | Admin |

> ⚠️ Đổi `JWT_SECRET` và mật khẩu admin trước khi deploy production!
