# 🏢 Employee Management System

> Hệ thống Quản lý Nhân sự toàn diện với **Vue 3** + **Golang (Gin)** + **MySQL**

---

## 📋 Mục lục

1. [Tổng quan dự án](#1-tổng-quan-dự-án)
2. [Công nghệ sử dụng](#2-công-nghệ-sử-dụng)
3. [Cấu trúc thư mục](#3-cấu-trúc-thư-mục)
4. [Database Schema](#4-database-schema)
5. [Authentication & Authorization](#5-authentication--authorization)
6. [API Documentation](#6-api-documentation)
7. [Frontend Documentation](#7-frontend-documentation)
8. [Feature Summary](#8-feature-summary)
9. [Setup & Run](#9-setup--run)
10. [Troubleshooting](#10-troubleshooting)

---

## 1. Tổng quan dự án

### 🎯 Mục tiêu

Xây dựng hệ thống quản lý nhân sự nội bộ cho doanh nghiệp, hỗ trợ quản trị viên (Admin) quản lý thông tin nhân viên, tài khoản người dùng, phòng ban và theo dõi thống kê tổng quan qua Dashboard.

### ✨ Chức năng chính

- **Đăng nhập / Xác thực JWT** — Bảo mật toàn bộ hệ thống
- **Dashboard thống kê** — Tổng quan nhân sự, phòng ban, biểu đồ phân bổ
- **Quản lý nhân viên (CRUD)** — Thêm, sửa, xoá, tìm kiếm, lọc, phân trang
- **Quản lý tài khoản (Admin)** — Tạo/sửa/xoá user, phân quyền, liên kết nhân viên
- **Hồ sơ cá nhân** — Xem thông tin bản thân (employee + user)
- **Phân quyền** — Admin vs User với middleware bảo vệ cả backend lẫn frontend

### 🏗️ Kiến trúc

```
┌─────────────────┐         ┌──────────────────┐         ┌─────────┐
│   Vue 3 (SPA)   │◄──API──►│  Golang (Gin)    │◄──ORM──►│  MySQL  │
│   Port: 3000    │         │  Port: 8080      │         │  :3306  │
│   Vite + Pinia  │         │  GORM + JWT      │         │         │
└─────────────────┘         └──────────────────┘         └─────────┘
```

---

## 2. Công nghệ sử dụng

### Frontend (`package.json`)

| Công nghệ              | Phiên bản | Vai trò                               |
| ---------------------- | --------- | ------------------------------------- |
| **Vue 3**              | ^3.5.32   | Framework SPA chính (Composition API) |
| **Vite**               | ^8.0.10   | Build tool & Dev server siêu nhanh    |
| **Vue Router**         | ^5.0.6    | Điều hướng SPA, navigation guard      |
| **Pinia**              | ^3.0.4    | State management (thay Vuex)          |
| **Axios**              | ^1.15.2   | HTTP client gọi API                   |
| **vue-toastification** | —         | Thông báo toast UI                    |

### Backend (`go.mod`)

| Công nghệ                          | Vai trò                                           |
| ---------------------------------- | ------------------------------------------------- |
| **Go 1.26**                        | Ngôn ngữ backend chính                            |
| **Gin v1.12**                      | Web framework HTTP                                |
| **GORM v1.31**                     | ORM cho MySQL (AutoMigrate, Preload, Soft Delete) |
| **golang-jwt/jwt v5**              | Tạo & xác thực JWT token                          |
| **bcrypt** (`golang.org/x/crypto`) | Hash mật khẩu an toàn                             |
| **godotenv**                       | Load biến môi trường từ `.env`                    |
| **gin-contrib/cors**               | Xử lý CORS cho frontend                           |
| **gorm.io/driver/mysql**           | Driver kết nối MySQL                              |

### Infrastructure

| Công nghệ          | Vai trò                       |
| ------------------ | ----------------------------- |
| **Docker Compose** | Container hóa MySQL + Backend |
| **MySQL 8.0**      | Database chính                |

---

## 3. Cấu trúc thư mục

### Backend (`backend/`)

```
backend/
├── .env                    # Biến môi trường (DB, JWT_SECRET, PORT)
├── main.go                 # Entry point — khởi tạo Gin, CORS, routes, seed
├── config/
│   └── config.go           # LoadEnv() + ConnectDB() — kết nối MySQL qua GORM
├── database/
│   └── database.go         # MigrateModels() + Seed() — tạo bảng & dữ liệu mẫu
├── models/
│   ├── role_model.go       # Role{ID, Name}
│   ├── department_model.go # Department{ID, Name}
│   ├── employee_model.go   # Employee{Name, Gender, Phone, Salary, Status...}
│   └── user_model.go       # User{Email, PasswordHash, RoleID, EmployeeID...}
├── dto/
│   ├── auth_dto.go         # LoginInput
│   ├── employee_dto.go     # CreateEmployeeInput, UpdateEmployeeInput, PublicResponse
│   ├── user_dto.go         # CreateUserInput, UpdateUserInput
│   └── dashboard_dto.go    # DashboardResponse, DeptStat
├── handlers/
│   ├── auth_handler.go     # LoginHandler — đăng nhập, tạo JWT
│   ├── me_handler.go       # MeHandler — lấy thông tin user hiện tại
│   ├── dashboard_handler.go# GetDashboardStats — thống kê tổng quan
│   ├── department_handler.go# GetDepartments — danh sách phòng ban
│   ├── employee_handler.go # CRUD nhân viên (phân quyền salary)
│   └── user_handler.go     # CRUD tài khoản (admin only)
├── middleware/
│   ├── jwt_middleware.go   # AuthMiddlewareJWT — xác thực Bearer token
│   ├── admin_middleware.go # AdminOnlyMiddleware — chỉ admin
│   └── error_middleware.go # ErrorHandler — xử lý lỗi tập trung
├── routes/
│   └── routes.go           # Định nghĩa toàn bộ API routes
├── utils/
│   └── response.go         # Helper: Success, BadRequest, NotFound...
├── db_em.sql               # SQL schema tham khảo
├── go.mod / go.sum         # Go dependencies
```

### Frontend (`frontend/`)

```
frontend/
├── .env                    # VITE_API_URL=http://localhost:8080/api
├── index.html              # HTML entry point
├── vite.config.js          # Vite config (alias @, port 3000)
├── package.json            # Dependencies
└── src/
    ├── main.js             # Khởi tạo Vue app + Pinia + Router + Toast
    ├── App.vue             # Root component (RouterView)
    ├── style.css           # Global styles
    ├── api/
    │   └── index.js        # Axios instance + interceptors (auto token, 401 redirect)
    ├── router/
    │   └── index.js        # Routes + Navigation Guards (auth, admin)
    ├── stores/
    │   ├── auth.js          # Login, logout, verifyToken, token/user state
    │   ├── dashboard.js     # fetchDashboardData
    │   ├── department.js    # fetchDepartments
    │   ├── employee.js      # CRUD employees + pagination
    │   ├── user.js          # CRUD users (admin)
    │   └── ui.js            # Sidebar state (collapse, mobile)
    ├── layouts/
    │   └── MainLayout.vue   # Sidebar + Header + RouterView
    ├── components/
    │   ├── AppHeader.vue    # Header: tiêu đề, user info, logout
    │   ├── AppSidebar.vue   # Sidebar: menu động theo role
    │   └── EmployeeForm.vue # Modal form thêm/sửa nhân viên
    ├── views/
    │   ├── LoginView.vue         # Trang đăng nhập
    │   ├── DashboardView.vue     # Trang tổng quan thống kê
    │   ├── EmployeeListView.vue  # Danh sách + CRUD nhân viên
    │   ├── ProfileView.vue       # Hồ sơ cá nhân
    │   ├── UserManagementView.vue# Quản lý tài khoản (admin)
    │   └── DepartmentView.vue    # Phòng ban (placeholder)
    └── assets/icons/             # SVG icons
```

---

## 4. Database Schema

### Sơ đồ quan hệ (ER Diagram)

```mermaid
erDiagram
    ROLES ||--o{ USERS : "1-N"
    DEPARTMENTS ||--o{ EMPLOYEES : "1-N"
    EMPLOYEES ||--o| USERS : "1-1"

    ROLES {
        TINYINT id PK
        VARCHAR name UK "admin, user"
        TIMESTAMP created_at
    }

    DEPARTMENTS {
        INT id PK
        VARCHAR name UK
        TIMESTAMP created_at
    }

    EMPLOYEES {
        INT id PK
        VARCHAR name "NOT NULL"
        ENUM gender "male, female, other"
        DATE date_of_birth
        VARCHAR phone UK
        INT department_id FK "→ departments.id"
        VARCHAR position
        DECIMAL salary "default 0"
        DATE hire_date
        TINYINT status "1=active, 0=inactive"
        VARCHAR avatar_url
        TIMESTAMP created_at
        TIMESTAMP updated_at
        TIMESTAMP deleted_at "Soft delete"
    }

    USERS {
        INT id PK
        VARCHAR email UK "NOT NULL"
        VARCHAR password_hash "NOT NULL"
        TINYINT role_id FK "→ roles.id, default 2"
        INT employee_id FK_UK "→ employees.id, UNIQUE 1-1"
        BOOLEAN is_active "default true"
        TIMESTAMP created_at
        TIMESTAMP updated_at
        TIMESTAMP deleted_at "Soft delete"
    }
```

### Chi tiết quan hệ

| Quan hệ                                      | Loại         | Mô tả                             | On Delete |
| -------------------------------------------- | ------------ | --------------------------------- | --------- |
| `users.role_id` → `roles.id`                 | N:1          | Mỗi user có 1 role                | RESTRICT  |
| `users.employee_id` → `employees.id`         | 1:1 (UNIQUE) | Mỗi employee chỉ có tối đa 1 user | SET NULL  |
| `employees.department_id` → `departments.id` | N:1          | Nhiều employee thuộc 1 department | SET NULL  |

### Role Logic

| Role ID | Name    | Quyền hạn                                                             |
| ------- | ------- | --------------------------------------------------------------------- |
| 1       | `admin` | Toàn quyền: CRUD employee, CRUD user, xem salary, dashboard đầy đủ    |
| 2       | `user`  | Xem danh sách employee (ẩn salary), xem profile, xem dashboard cơ bản |

### Seed Data

- **Roles**: `admin`, `user`
- **Departments**: IT, HR, Finance, Marketing, Sales, Operations
- **Employees**: 40 nhân viên mẫu (tên tiếng Việt)
- **Admin**: `chiquoc64@admin.company.dev` / `admin123`
- **Users**: Tự động tạo từ tên nhân viên, mật khẩu mặc định `123456`

> ⚠️ Seed chỉ chạy khi `SEEDDATA=true` trong `.env`

---

## 5. Authentication & Authorization

### Login Flow

```
1. User nhập email + password
2. Frontend gửi POST /api/auth/login
3. Backend tìm user theo email (Preload Role + Employee)
4. Kiểm tra is_active → nếu false → 403 Forbidden
5. So sánh bcrypt hash password
6. Tạo JWT token (chứa user_id + role, hết hạn 48h)
7. Trả về { token, user: {id, email, name, avatar_url, role} }
8. Frontend lưu token + user vào localStorage
9. Mọi request sau đều gắn header: Authorization: Bearer <token>
```

### JWT Token

| Field      | Giá trị                           |
| ---------- | --------------------------------- |
| Algorithm  | HS256                             |
| Secret     | `JWT_SECRET` từ `.env`            |
| Expiration | 48 giờ                            |
| Claims     | `user_id` (uint), `role` (string) |

### Middleware Pipeline

```
Request → ErrorHandler → CORS → [Route]
                                    │
                          ┌─────────┴─────────┐
                     Public route         Protected route
                    (POST /auth/login)         │
                                        AuthMiddlewareJWT
                                        (verify token + DB check)
                                               │
                                    ┌──────────┴──────────┐
                               All users              Admin only
                            (GET employees,        AdminOnlyMiddleware
                             GET departments,      (check role="admin")
                             GET dashboard,             │
                             GET auth/me)      POST/PUT/DELETE employees
                                               GET/POST/PUT/DELETE users
```

### Middleware chi tiết

| Middleware            | File                  | Chức năng                                                                                       |
| --------------------- | --------------------- | ----------------------------------------------------------------------------------------------- |
| `ErrorHandler`        | `error_middleware.go` | Bắt lỗi chưa xử lý, trả JSON chuẩn                                                              |
| `AuthMiddlewareJWT`   | `jwt_middleware.go`   | Parse Bearer token, verify JWT, kiểm tra user tồn tại trong DB, set `user` + `role` vào context |
| `AdminOnlyMiddleware` | `admin_middleware.go` | Kiểm tra `role == "admin"`, reject nếu không phải                                               |

### Frontend Protected Routes

| Route          | Meta              | Guard Logic                          |
| -------------- | ----------------- | ------------------------------------ |
| `/login`       | `public: true`    | Nếu đã login → redirect `/dashboard` |
| `/dashboard`   | —                 | Phải login, verify token lần đầu     |
| `/employees`   | —                 | Phải login                           |
| `/profile`     | —                 | Phải login                           |
| `/departments` | —                 | Phải login                           |
| `/users`       | `adminOnly: true` | Phải login + role = admin            |

---

## 6. API Documentation

**Base URL**: `http://localhost:8080/api`

### 6.1 Authentication

#### `POST /auth/login` — Đăng nhập

|                  | Chi tiết                                                                                          |
| ---------------- | ------------------------------------------------------------------------------------------------- |
| **Middleware**   | Không (public)                                                                                    |
| **Handler**      | `handlers/auth_handler.go` → `LoginHandler`                                                       |
| **Request Body** | `{ "email": "string (required)", "password": "string (required)" }`                               |
| **Success 200**  | `{ "success": true, "token": "jwt...", "user": { "id", "email", "name", "avatar_url", "role" } }` |
| **Error 401**    | Sai email hoặc password                                                                           |
| **Error 403**    | Tài khoản bị vô hiệu hoá (`is_active = false`)                                                    |
| **Logic**        | Tìm user by email → check is_active → bcrypt compare → tạo JWT 48h → trả token + user info        |

#### `GET /auth/me` — Thông tin user hiện tại

|                 | Chi tiết                                                                                      |
| --------------- | --------------------------------------------------------------------------------------------- |
| **Middleware**  | `AuthMiddlewareJWT`                                                                           |
| **Handler**     | `handlers/me_handler.go` → `MeHandler`                                                        |
| **Success 200** | `{ "success": true, "user": { "id", "email", "name", "avatar", "role", "employee": {...} } }` |
| **Logic**       | Lấy user_id từ context → Preload Role + Employee.Department → trả JSON                        |

### 6.2 Dashboard

#### `GET /dashboard` — Thống kê tổng quan

|                 | Chi tiết                                                                                                                                                         |
| --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Middleware**  | `AuthMiddlewareJWT`                                                                                                                                              |
| **Handler**     | `handlers/dashboard_handler.go` → `GetDashboardStats`                                                                                                            |
| **Success 200** | `{ "totalEmployees", "activeEmployees", "inactiveEmployees", "totalDepartments", "totalUsers", "totalAdminRole", "employeesByDepartment": [{"name", "count"}] }` |
| **Logic**       | Đếm employees (total/active/inactive), departments, users, admins + group by department                                                                          |

### 6.3 Departments

#### `GET /departments` — Danh sách phòng ban

|                 | Chi tiết                                                      |
| --------------- | ------------------------------------------------------------- |
| **Middleware**  | `AuthMiddlewareJWT`                                           |
| **Handler**     | `handlers/department_handler.go` → `GetDepartments`           |
| **Success 200** | `{ "success": true, "data": [{"id", "name", "created_at"}] }` |

### 6.4 Employees

#### `GET /employees` — Danh sách nhân viên (có phân trang, tìm kiếm, lọc)

|                  | Chi tiết                                                                            |
| ---------------- | ----------------------------------------------------------------------------------- |
| **Middleware**   | `AuthMiddlewareJWT`                                                                 |
| **Role**         | Tất cả (admin thấy salary, user thấy public response ẩn salary)                     |
| **Handler**      | `handlers/employee_handler.go` → `GetEmployeeList`                                  |
| **Query Params** | `page` (default 1), `limit` (default 4), `search` (tên/phone), `department_id`      |
| **Success 200**  | `{ "success": true, "data": [...], "total": N, "page": N, "limit": N }`             |
| **Logic**        | Preload Department + User, search LIKE multi-word, filter dept, phân quyền response |

#### `GET /employees/:id` — Chi tiết nhân viên

|                 | Chi tiết                                                            |
| --------------- | ------------------------------------------------------------------- |
| **Middleware**  | `AuthMiddlewareJWT`                                                 |
| **Handler**     | `GetEmployeeID`                                                     |
| **Success 200** | `{ "success": true, "data": {...} }` — Admin: full, User: ẩn salary |
| **Error 404**   | Không tìm thấy                                                      |

#### `POST /employees` — Tạo nhân viên mới

|                  | Chi tiết                                                                                                     |
| ---------------- | ------------------------------------------------------------------------------------------------------------ |
| **Middleware**   | `AuthMiddlewareJWT` + `AdminOnlyMiddleware`                                                                  |
| **Handler**      | `CreateEmployee`                                                                                             |
| **Request Body** | `{ "name*", "phone*", "salary* (>0)", "gender", "date_of_birth", "department_id", "position", "hire_date" }` |
| **Success 201**  | `{ "success": true, "data": {...} }`                                                                         |
| **Error 409**    | Phone đã tồn tại                                                                                             |
| **Logic**        | Validate phone unique → validate department exists → parse dates → create → preload dept                     |

#### `PUT /employees/:id` — Cập nhật nhân viên

|                  | Chi tiết                                                                                                                         |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| **Middleware**   | `AuthMiddlewareJWT` + `AdminOnlyMiddleware`                                                                                      |
| **Handler**      | `UpdateEmployee`                                                                                                                 |
| **Request Body** | Tất cả field optional (partial update): `name, gender, phone, department_id, position, salary, hire_date, date_of_birth, status` |
| **Logic**        | Chỉ update field được gửi → validate phone unique (trừ chính nó) → parse dates                                                   |

#### `DELETE /employees/:id` — Xoá nhân viên (soft delete)

|                | Chi tiết                                                                 |
| -------------- | ------------------------------------------------------------------------ |
| **Middleware** | `AuthMiddlewareJWT` + `AdminOnlyMiddleware`                              |
| **Handler**    | `DeleteEmployee`                                                         |
| **Logic**      | Xoá user liên kết → set status = inactive → soft delete (GORM DeletedAt) |

### 6.5 Users (Admin Only)

> Tất cả endpoints trong group này yêu cầu `AuthMiddlewareJWT` + `AdminOnlyMiddleware`

#### `GET /users` — Danh sách tài khoản

|                 | Chi tiết                                                |
| --------------- | ------------------------------------------------------- |
| **Handler**     | `handlers/user_handler.go` → `GetUsers`                 |
| **Success 200** | `{ "success": true, "data": [user + role + employee] }` |

#### `GET /users/:id` — Chi tiết tài khoản

|                 | Chi tiết                                                |
| --------------- | ------------------------------------------------------- |
| **Handler**     | `GetUserByID`                                           |
| **Success 200** | `{ "success": true, "data": {user + role + employee} }` |

#### `POST /users` — Tạo tài khoản

|                  | Chi tiết                                                                                         |
| ---------------- | ------------------------------------------------------------------------------------------------ |
| **Handler**      | `CreateUser`                                                                                     |
| **Request Body** | `{ "email*", "password* (min 6)", "role_id* (1 or 2)", "employee_id" }`                          |
| **Validation**   | User role (2) bắt buộc có employee_id; Kiểm tra employee tồn tại; Kiểm tra employee chưa có user |
| **Error 409**    | Email đã tồn tại hoặc employee đã có tài khoản                                                   |

#### `PUT /users/:id` — Cập nhật tài khoản

|                  | Chi tiết                                                                |
| ---------------- | ----------------------------------------------------------------------- |
| **Handler**      | `UpdateUser`                                                            |
| **Request Body** | `{ "email", "password", "role_id", "employee_id" }` — tất cả optional   |
| **Bảo vệ**       | Admin không thể sửa admin khác; Không thể tự đổi role_id của chính mình |
| **Logic**        | Partial update, hash password mới nếu có, check employee_id conflict    |

#### `DELETE /users/:id` — Xoá tài khoản (soft delete)

|             | Chi tiết                                           |
| ----------- | -------------------------------------------------- |
| **Handler** | `DeleteUser`                                       |
| **Bảo vệ**  | Không thể xoá chính mình; Không thể xoá admin khác |

---

## 7. Frontend Documentation

### 7.1 Axios Instance (`src/api/index.js`)

```javascript
// Base URL từ .env hoặc fallback
baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api";
timeout: 10000;
```

**Request Interceptor**: Tự động gắn `Authorization: Bearer <token>` từ localStorage.

**Response Interceptor**:

- Tự động unwrap `response.data` (chỉ trả data, không trả wrapper axios)
- Nếu **401**: Xoá `token` + `user` khỏi storage → redirect `/login`
- Chuẩn hoá error message → `Promise.reject({ status, message })`

### 7.2 Router (`src/router/index.js`)

| Route          | Component            | Meta                  | Mô tả                           |
| -------------- | -------------------- | --------------------- | ------------------------------- |
| `/login`       | `LoginView`          | `public: true`        | Trang đăng nhập                 |
| `/`            | `MainLayout`         | redirect `/dashboard` | Layout chính (Sidebar + Header) |
| `/dashboard`   | `DashboardView`      | —                     | Trang tổng quan                 |
| `/employees`   | `EmployeeListView`   | —                     | Quản lý nhân viên               |
| `/profile`     | `ProfileView`        | —                     | Hồ sơ cá nhân                   |
| `/users`       | `UserManagementView` | `adminOnly: true`     | Quản lý tài khoản               |
| `/departments` | `DepartmentView`     | —                     | Phòng ban (placeholder)         |

**Navigation Guards** (theo thứ tự):

1. **Token verification**: Lần đầu vào protected route → gọi `GET /auth/me` verify token
2. **Auth check**: Route protected mà chưa login → redirect `/login`
3. **Login guard**: Đã login mà vào `/login` → redirect `/dashboard`
4. **Admin guard**: Route `adminOnly` mà không phải admin → redirect `/dashboard`

### 7.3 Stores (Pinia)

#### `auth.js` — Quản lý xác thực

| State/Action             | Mô tả                                                       |
| ------------------------ | ----------------------------------------------------------- |
| `token`                  | JWT token (khởi tạo từ localStorage)                        |
| `user`                   | Thông tin user (khởi tạo từ localStorage)                   |
| `isLoggedIn`             | Computed: `!!token`                                         |
| `login(email, password)` | Gọi API login → lưu token + user vào state + localStorage   |
| `logout()`               | Xoá token + user khỏi state + localStorage + sessionStorage |
| `verifyToken()`          | Gọi `GET /auth/me` → cập nhật user → trả `true/false`       |

#### `employee.js` — Quản lý nhân viên

| State/Action               | Mô tả                                                     |
| -------------------------- | --------------------------------------------------------- |
| `employees`                | Mảng nhân viên hiện tại                                   |
| `total`                    | Tổng số nhân viên (cho phân trang)                        |
| `loading`                  | Trạng thái loading                                        |
| `fetchEmployees(params)`   | `GET /employees` với `page, limit, search, department_id` |
| `createEmployee(data)`     | `POST /employees`                                         |
| `updateEmployee(id, data)` | `PUT /employees/:id`                                      |
| `deleteEmployee(id)`       | `DELETE /employees/:id`                                   |

#### `user.js` — Quản lý tài khoản

| State/Action           | Mô tả               |
| ---------------------- | ------------------- |
| `users`                | Mảng user           |
| `loading`              | Trạng thái loading  |
| `fetchUsers()`         | `GET /users`        |
| `createUser(data)`     | `POST /users`       |
| `updateUser(id, data)` | `PUT /users/:id`    |
| `deleteUser(id)`       | `DELETE /users/:id` |

#### `dashboard.js` — Dashboard

| State/Action           | Mô tả                       |
| ---------------------- | --------------------------- |
| `dashboardData`        | Object chứa tất cả thống kê |
| `fetchDashboardData()` | `GET /dashboard`            |

#### `department.js` — Phòng ban

| State/Action         | Mô tả              |
| -------------------- | ------------------ |
| `departments`        | Mảng phòng ban     |
| `fetchDepartments()` | `GET /departments` |

#### `ui.js` — Giao diện

| State/Action            | Mô tả                        |
| ----------------------- | ---------------------------- |
| `mobileOpen`            | Sidebar mở trên mobile       |
| `isCollapsed`           | Sidebar thu gọn trên desktop |
| `toggleMobileSidebar()` | Bật/tắt sidebar mobile       |
| `toggleCollapse()`      | Bật/tắt thu gọn sidebar      |

### 7.4 Components

#### `MainLayout.vue`

- Layout chính: `AppSidebar` (fixed left) + `AppHeader` (sticky top) + `<RouterView>`
- Responsive: sidebar ẩn trên mobile, body margin-left thay đổi khi collapse

#### `AppHeader.vue`

- Hiển thị tiêu đề trang (auto mapping từ route path)
- User info: avatar, email, role badge
- Nút đăng xuất + hamburger menu (mobile)

#### `AppSidebar.vue`

- Menu items: Dashboard, Nhân viên, Hồ sơ, Phòng Ban
- Menu "Tài khoản" chỉ hiện khi `role === "admin"`
- Hỗ trợ collapse (ẩn text, chỉ icon) + responsive mobile overlay

#### `EmployeeForm.vue`

- Modal form thêm mới / cập nhật nhân viên
- Các field: Họ tên, Giới tính, Ngày sinh, SĐT, Phòng ban (select), Chức vụ, Lương, Ngày vào làm
- Validate trước khi submit
- Emit `@saved` và `@close`

### 7.5 Views

#### `LoginView.vue`

- Form email + password với icon, toggle password visibility
- Loading spinner khi đang xử lý (800ms delay UX)
- Toast thông báo thành công/thất bại
- Redirect `/dashboard` sau login

#### `DashboardView.vue`

- 4 stat cards: Tổng NV, Đang làm việc, Phòng ban, Tổng Admin (admin only)
- Biểu đồ phân bổ nhân viên theo phòng ban (progress bar)
- Hover animation trên stat cards

#### `EmployeeListView.vue`

- Bảng dữ liệu: Avatar, Tên, Ngày sinh, Giới tính, SĐT, Phòng ban/Chức vụ, Lương (admin), Trạng thái
- Tìm kiếm debounce 300ms (tên/SĐT)
- Lọc theo phòng ban (select)
- Phân trang (7 items/page)
- Nút Thêm/Sửa/Xoá (admin only)
- Modal xác nhận xoá
- Skeleton loading

#### `ProfileView.vue`

- 2 cột: Sidebar (avatar, tên, role, status, thống kê nhanh) + Main (thông tin cá nhân + công việc)
- Dữ liệu lấy từ `authStore.user` (bao gồm employee data)
- Format tiền VND, ngày tiếng Việt

#### `UserManagementView.vue`

- 2 tab: "Tài khoản" (list users) + "Chưa có tài khoản" (employees without user)
- Tìm kiếm theo email/tên
- Modal form tạo/sửa tài khoản (email, password, confirm password, role, employee liên kết)
- Nút "Cấp tài khoản" nhanh cho employee chưa có
- Bảo vệ: Admin không thể sửa/xoá admin khác, không thể tự đổi role, icon 🔒 cho tài khoản bị khoá
- Modal xác nhận xoá

#### `DepartmentView.vue`

- Placeholder (chưa implement)

---

## 8. Feature Summary

### Tính năng theo Role

| #   | Tính năng                    |     Admin      |      User      |
| --- | ---------------------------- | :------------: | :------------: |
| 1   | Đăng nhập / Đăng xuất        |       ✅       |       ✅       |
| 2   | Xem Dashboard thống kê       |   ✅ (full)    |  ✅ (cơ bản)   |
| 3   | Xem danh sách nhân viên      | ✅ (có salary) | ✅ (ẩn salary) |
| 4   | Tìm kiếm nhân viên (tên/SĐT) |       ✅       |       ✅       |
| 5   | Lọc nhân viên theo phòng ban |       ✅       |       ✅       |
| 6   | Phân trang danh sách         |       ✅       |       ✅       |
| 7   | Thêm nhân viên mới           |       ✅       |       ❌       |
| 8   | Sửa thông tin nhân viên      |       ✅       |       ❌       |
| 9   | Xoá nhân viên (soft delete)  |       ✅       |       ❌       |
| 10  | Xem hồ sơ cá nhân            |       ✅       |       ✅       |
| 11  | Xem danh sách phòng ban      |       ✅       |       ✅       |
| 12  | Quản lý tài khoản (CRUD)     |       ✅       |       ❌       |
| 13  | Phân quyền tài khoản         |       ✅       |       ❌       |
| 14  | Cấp tài khoản nhanh cho NV   |       ✅       |       ❌       |
| 15  | Xem NV chưa có tài khoản     |       ✅       |       ❌       |

### Tính năng kỹ thuật

- ✅ JWT Authentication (48h expiry)
- ✅ Bcrypt password hashing
- ✅ GORM Soft Delete (DeletedAt)
- ✅ CORS configuration
- ✅ Auto DB migration
- ✅ Database seeding (40 employees + admin + users)
- ✅ Axios interceptors (auto token, 401 redirect)
- ✅ Vue Router navigation guards
- ✅ Responsive design (mobile sidebar overlay)
- ✅ Collapsible sidebar
- ✅ Toast notifications
- ✅ Debounced search
- ✅ Skeleton loading states
- ✅ Confirm modal before delete
- ✅ Vietnamese date & currency formatting
- ✅ Docker Compose support

---

## 9. Setup & Run

### Yêu cầu

- **Go** ≥ 1.26
- **Node.js** ≥ 18
- **MySQL** 8.0
- **Docker** (optional)

### Clone project

```bash
git clone https://github.com/Quocdev03/Employee_Management.git
cd Employee_Management
```

### Backend

```bash
cd backend

# Cấu hình .env (đã có sẵn, chỉnh sửa nếu cần)
# DB_HOST=127.0.0.1
# DB_PORT=3306
# DB_USER=quocdt6423
# DB_PASSWORD=quoc@123456
# DB_NAME=employee_db
# PORT=8080
# JWT_SECRET=...
# SEEDDATA=false  ← đổi thành true để seed dữ liệu mẫu lần đầu

# Cài dependencies
go mod tidy

# Chạy server
go run main.go
```

> Server chạy tại `http://localhost:8080`

### Frontend

```bash
cd frontend

# Cài dependencies
npm install

# Chạy dev server
npm run dev
```

> App chạy tại `http://localhost:3000`

### Docker

```bash
# Từ thư mục gốc
docker-compose up --build
```

Sẽ khởi tạo:

- **MySQL** container (`mysql-employee_db`) — port 3306
- **Go App** container (`go-app`) — port 8080

### Build Production

```bash
# Frontend
cd frontend
npm run build    # Output: dist/

# Backend
cd backend
go build -o server main.go
./server
```

### Tài khoản mặc định (sau khi seed)

| Loại      | Email                         | Password   |
| --------- | ----------------------------- | ---------- |
| **Admin** | `chiquoc64@admin.company.dev` | `admin123` |

> 💡 Để seed data: đặt `SEEDDATA=true` trong `backend/.env` rồi chạy lại server. Chỉ cần làm 1 lần.

---

## 10. Troubleshooting

### Port conflict

```
# Kiểm tra process đang dùng port
# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :3000

# Kill process
taskkill /PID <PID> /F

# Hoặc đổi port trong .env / vite.config.js
```

### MySQL connection refused

```
# Kiểm tra MySQL đang chạy
# Kiểm tra thông tin trong backend/.env:
#   DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
# Đảm bảo database employee_db đã được tạo:
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS employee_db;"
```

### `npm install` lỗi

```bash
# Xoá cache và cài lại
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### `go mod tidy` lỗi

```bash
# Xoá cache Go modules
go clean -modcache
go mod tidy
```

### Docker cache gây lỗi build

```bash
# Build lại không dùng cache
docker-compose down -v
docker-compose build --no-cache
docker-compose up
```

### Frontend không gọi được API (CORS)

- Kiểm tra `VITE_API_URL` trong `frontend/.env` trỏ đúng backend
- Backend đã cấu hình CORS cho origins: `localhost:3000`, `localhost:5173`, `localhost:4200`
- Đảm bảo backend đang chạy trước khi mở frontend

### Token hết hạn / 401 liên tục

- JWT token hết hạn sau 48h → đăng nhập lại
- Nếu DB bị reset mà token cũ vẫn còn → xoá localStorage trong browser DevTools

---

> **Tác giả**: Cao Chí Quốc  
> **Cập nhật lần cuối**: Tháng 05/2026
