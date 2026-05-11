# 🏢 Employee Management System

> Hệ thống Quản lý Nhân sự nội bộ với **Vue 3** + **Golang (Gin)** + **MySQL**

---

## 📋 Mục lục

1. [Tổng quan dự án](#1-tổng-quan-dự-án)
2. [Công nghệ sử dụng](#2-công-nghệ-sử-dụng)
3. [Cấu trúc thư mục](#3-cấu-trúc-thư-mục)
4. [Database Schema](#4-database-schema)
5. [API Documentation](#5-api-documentation)
6. [Setup & Run](#6-setup--run)
7. [Troubleshooting](#7-troubleshooting)

---

## 1. Tổng quan dự án

Xây dựng hệ thống quản lý nhân sự nội bộ cho doanh nghiệp, hỗ trợ Admin quản lý thông tin nhân viên, tài khoản người dùng, phòng ban, chức vụ và theo dõi thống kê qua Dashboard.

**Chức năng chính:**

| Tính năng                                |     Admin     |     User      |
| ---------------------------------------- | :-----------: | :-----------: |
| Đăng nhập / Đăng xuất (JWT)              |      ✅       |      ✅       |
| Dashboard thống kê                       |  ✅ (đầy đủ)  |  ✅ (cơ bản)  |
| Xem danh sách nhân viên                  | ✅ (có lương) | ✅ (ẩn lương) |
| Tìm kiếm & lọc theo phòng ban            |      ✅       |      ✅       |
| Thêm / Sửa / Xoá nhân viên (soft delete) |      ✅       |      ❌       |
| Quản lý phòng ban (CRUD)                 |      ✅       |      ❌       |
| Quản lý chức vụ theo phòng ban (CRUD)    |      ✅       |      ❌       |
| Xem hồ sơ cá nhân                        |      ✅       |      ✅       |
| Quản lý tài khoản (CRUD, phân quyền)     |      ✅       |      ❌       |
| Cấp tài khoản nhanh cho nhân viên        |      ✅       |      ❌       |

**Kiến trúc:**

```
┌─────────────────┐         ┌──────────────────┐         ┌─────────┐
│   Vue 3 (SPA)   │◄──API──►│  Golang (Gin)    │◄──ORM──►│  MySQL  │
│   Port: 3000    │         │  Port: 8080      │         │  :3306  │
│   Vite + Pinia  │         │  GORM + JWT      │         │         │
└─────────────────┘         └──────────────────┘         └─────────┘
```

---

## 2. Công nghệ sử dụng

### Frontend

| Công nghệ              | Vai trò                          |
| ---------------------- | -------------------------------- |
| **Vue 3** (^3.5)       | Framework SPA (Composition API)  |
| **Vite** (^8.0)        | Build tool & Dev server          |
| **Vue Router** (^5.0)  | Điều hướng SPA, navigation guard |
| **Pinia** (^3.0)       | State management                 |
| **Axios** (^1.15)      | HTTP client                      |
| **vue-toastification** | Thông báo toast UI               |

### Backend

| Công nghệ               | Vai trò                              |
| ----------------------- | ------------------------------------ |
| **Go 1.26 + Gin v1.12** | Web framework                        |
| **GORM v1.31**          | ORM cho MySQL (AutoMigrate, Preload) |
| **golang-jwt/jwt v5**   | Tạo & xác thực JWT token             |
| **bcrypt**              | Hash mật khẩu                        |
| **godotenv**            | Load biến môi trường từ `.env`       |
| **gin-contrib/cors**    | Xử lý CORS                           |

### Infrastructure

| Công nghệ          | Vai trò                           |
| ------------------ | --------------------------------- |
| **Docker Compose** | Container hóa MySQL + Backend     |
| **MySQL 8.0**      | Database chính                    |
| **Postman**        | Test & quản lý API endpoints      |
| **DBeaver**        | Quản lý và truy vấn cơ sở dữ liệu |

---

## 3. Cấu trúc thư mục

### Backend (`backend/`)

```
backend/
├── main.go                  # Entry point
├── config/config.go         # LoadEnv() + ConnectDB()
├── database/database.go     # MigrateModels() + Seed()
├── models/                  # Database Entities (Role, Department, Position, Employee, User)
├── dto/                     # Data Transfer Objects (Requests/Responses API)
├── handlers/                # auth, me, dashboard, department, position, employee, user
├── middleware/              # JWT, AdminOnly, ErrorHandler
├── routes/routes.go         # Định nghĩa API routes
├── utils/response.go        # Chuẩn hóa API response
└── db_em.sql                # SQL schema tham khảo
```

### Frontend (`frontend/`)

```
frontend/src/
├── api/index.js             # Axios instance + interceptors
├── router/index.js          # Routes + Navigation Guards
├── stores/                  # Pinia stores (auth, me, employee, user, dashboard, department)
├── assets/common.css        # CSS dùng chung (utilities, buttons, forms)
├── layouts/MainLayout.vue   # Sidebar + Header + RouterView
├── components/              # BaseModal, BaseTable, ConfirmModal, EmployeeForm
└── views/                   # Login, Dashboard, EmployeeList, Department, Profile, UserManagement
```

---

## 4. Database Schema

```mermaid
erDiagram
    ROLES ||--o{ USERS : "has"
    DEPARTMENTS ||--o{ POSITIONS : "has"
    DEPARTMENTS ||--o{ EMPLOYEES : "belongs to"
    POSITIONS ||--o{ EMPLOYEES : "holds"
    EMPLOYEES |o--o| USERS : "linked (1-1)"

    ROLES {
        TINYINT_UNSIGNED id PK
        VARCHAR50 name UK
        TIMESTAMP created_at
    }

    DEPARTMENTS {
        INT_UNSIGNED id PK
        VARCHAR100 name UK
        TIMESTAMP created_at
    }

    POSITIONS {
        INT_UNSIGNED id PK
        VARCHAR150 name
        INT_UNSIGNED department_id FK
        TIMESTAMP created_at
    }

    EMPLOYEES {
        INT_UNSIGNED id PK
        VARCHAR255 name
        ENUM gender
        DATE date_of_birth
        VARCHAR20 phone UK
        INT_UNSIGNED department_id FK
        INT_UNSIGNED position_id FK
        DECIMAL15_2 salary
        DATE hire_date
        TINYINT status
        VARCHAR500 avatar_url
        TIMESTAMP created_at
        TIMESTAMP updated_at
        TIMESTAMP deleted_at
    }

    USERS {
        INT_UNSIGNED id PK
        VARCHAR255 email UK
        VARCHAR255 password_hash
        TINYINT_UNSIGNED role_id FK
        INT_UNSIGNED employee_id FK_UK
        BOOLEAN is_active
        TIMESTAMP created_at
        TIMESTAMP updated_at
        TIMESTAMP deleted_at
    }
```

**Notes:**

| Column          | References       | Constraint                       |
| --------------- | ---------------- | -------------------------------- |
| `role_id`       | `roles.id`       | ON DELETE RESTRICT               |
| `department_id` | `departments.id` | ON DELETE SET NULL               |
| `position_id`   | `positions.id`   | ON DELETE SET NULL               |
| `employee_id`   | `employees.id`   | UNIQUE — 1-1, ON DELETE SET NULL |

- `deleted_at` — soft delete (GORM convention)
- `status` — `1` = active, `0` = inactive
- `gender` — `male` = Nam, `female` = Nữ, `other` = Khác
- `position_id` phải thuộc đúng `department_id` của employee → validate ở application layer

**Chi tiết quan hệ**

| Quan hệ                                      | Loại         | Mô tả                               | On Delete |
| -------------------------------------------- | ------------ | ----------------------------------- | --------- |
| `users.role_id` → `roles.id`                 | N:1          | Mỗi user có 1 role                  | RESTRICT  |
| `users.employee_id` → `employees.id`         | 1:1 (UNIQUE) | Mỗi employee chỉ có tối đa 1 user   | SET NULL  |
| `employees.department_id` → `departments.id` | N:1          | Nhiều employee thuộc 1 phòng ban    | SET NULL  |
| `employees.position_id` → `positions.id`     | N:1          | Nhiều employee có cùng chức vụ      | SET NULL  |
| `positions.department_id` → `departments.id` | N:1          | Chức vụ thuộc về 1 phòng ban cụ thể | CASCADE   |

**Phân quyền:**

| Role    | Quyền hạn                                                                           |
| ------- | ----------------------------------------------------------------------------------- |
| `admin` | Toàn quyền: CRUD employee, user, department, position; xem salary; dashboard đầy đủ |
| `user`  | Xem danh sách employee (ẩn salary), xem profile, dashboard cơ bản                   |

**Bảo vệ tài khoản Admin:**

- Admin không thể tự xoá tài khoản của chính mình
- Admin không thể xoá hoặc sửa quyền của Admin khác
- Xoá nhân viên sẽ bị chặn nếu nhân viên đó là Admin

**Seed data**

- **Roles**: `admin`, `user`
- **Departments**: IT, HR, Finance, Marketing, Sales, Operations
- **Positions**: 4 chức vụ cho mỗi phòng ban (24 chức vụ tổng cộng)
- **Employees**: 40 nhân viên mẫu (tên tiếng Việt)
- **Admin**: `chiquoc64@admin.company.dev` / `admin123`
- **Users**: Tự động tạo từ tên nhân viên, mật khẩu mặc định `123456`

> Seed chỉ chạy khi `SEEDDATA=true` trong `.env`

---

## 5. API Documentation

**Base URL**: `http://localhost:8080/api`  
Mọi route (trừ `/auth/login`) yêu cầu header `Authorization: Bearer <token>`.

**Chuẩn hóa Response:**

```json
// Thành công
{ "success": true, "data": { ... } }

// Phân trang
{ "success": true, "data": [...], "total": 40, "page": 1, "limit": 10 }

// Thông tin cá nhân (/auth/me)
{
  "success": true,
  "data": {
    "user": { "id": 1, "email": "...", "name": "...", "role": "admin", "avatar": "..." },
    "employee": { "id": 39, "name": "...", "department": {...}, "position": {...} }
  }
}

// Lỗi
{ "success": false, "error": "Mô tả lỗi" }
```

### Auth

| Method | Endpoint      | Mô tả                    | Auth |
| ------ | ------------- | ------------------------ | :--: |
| POST   | `/auth/login` | Đăng nhập, trả JWT token |  ❌  |
| GET    | `/auth/me`    | Thông tin user hiện tại  |  ✅  |

### Dashboard

| Method | Endpoint     | Mô tả              | Admin | User |
| ------ | ------------ | ------------------ | :---: | :--: |
| GET    | `/dashboard` | Thống kê tổng quan |  ✅   |  ✅  |

### Departments & Positions

| Method | Endpoint                         | Mô tả                             | Admin | User |
| ------ | -------------------------------- | --------------------------------- | :---: | :--: |
| GET    | `/departments`                   | Danh sách phòng ban (kèm chức vụ) |  ✅   |  ✅  |
| GET    | `/departments/:id`               | Chi tiết phòng ban                |  ✅   |  ✅  |
| GET    | `/departments/:deptId/positions` | Chức vụ theo phòng ban            |  ✅   |  ✅  |
| POST   | `/departments`                   | Tạo phòng ban                     |  ✅   |  ❌  |
| PUT    | `/departments/:id`               | Cập nhật phòng ban                |  ✅   |  ❌  |
| DELETE | `/departments/:id`               | Xoá phòng ban (khi không có NV)   |  ✅   |  ❌  |
| POST   | `/positions`                     | Tạo chức vụ                       |  ✅   |  ❌  |
| PUT    | `/positions/:id`                 | Cập nhật chức vụ                  |  ✅   |  ❌  |
| DELETE | `/positions/:id`                 | Xoá chức vụ (khi không có NV)     |  ✅   |  ❌  |

### Employees

| Method | Endpoint         | Mô tả                                           |     Admin      |      User      |
| ------ | ---------------- | ----------------------------------------------- | :------------: | :------------: |
| GET    | `/employees`     | Danh sách nhân viên (phân trang, tìm kiếm, lọc) | ✅ (có salary) | ✅ (ẩn salary) |
| GET    | `/employees/:id` | Chi tiết nhân viên                              |       ✅       |       ✅       |
| POST   | `/employees`     | Tạo nhân viên                                   |       ✅       |       ❌       |
| PUT    | `/employees/:id` | Cập nhật nhân viên                              |       ✅       |       ❌       |
| DELETE | `/employees/:id` | Xoá mềm nhân viên                               |       ✅       |       ❌       |

**Query params cho `GET /employees`:**

| Param           | Mô tả                            | Ví dụ              |
| --------------- | -------------------------------- | ------------------ |
| `page`          | Trang hiện tại (mặc định: 1)     | `?page=2`          |
| `limit`         | Số dòng mỗi trang (mặc định: 10) | `?limit=20`        |
| `search`        | Tìm theo tên / số điện thoại     | `?search=Quốc`     |
| `department_id` | Lọc theo phòng ban               | `?department_id=1` |

### Users

| Method | Endpoint     | Mô tả               | Admin | User |
| ------ | ------------ | ------------------- | :---: | :--: |
| GET    | `/users`     | Danh sách tài khoản |  ✅   |  ❌  |
| GET    | `/users/:id` | Chi tiết tài khoản  |  ✅   |  ❌  |
| POST   | `/users`     | Tạo tài khoản       |  ✅   |  ❌  |
| PUT    | `/users/:id` | Cập nhật tài khoản  |  ✅   |  ❌  |
| DELETE | `/users/:id` | Xoá tài khoản       |  ✅   |  ❌  |

---

## 6. Setup & Run

**Yêu cầu**

- **Go** ≥ 1.26
- **Node.js** ≥ 18
- **MySQL** 8.0
- **Docker** (optional)

```bash
git clone https://github.com/Quocdev03/Employee_Management.git
cd Employee_Management
```

### Backend

```bash
cd backend
# Chỉnh backend/.env nếu cần (DB_HOST, DB_USER, DB_PASSWORD, JWT_SECRET...)
# Đặt SEEDDATA=true để seed dữ liệu mẫu lần đầu
go mod tidy
go run main.go
# → http://localhost:8080
```

> **Lần đầu chạy**: GORM sẽ tự động tạo bảng `roles`, `departments`, `positions`, `employees`, `users`.  
> Nếu DB cũ đã có cột `position` (string) trong bảng `employees`, cần chạy migration thủ công hoặc drop/recreate DB.

### Frontend

```bash
cd frontend
npm install
npm run dev
# → http://localhost:3000
```

### Docker

```bash
# Từ thư mục gốc
docker-compose up --build
```

Sẽ khởi tạo:

- **MySQL** container (`mysql-employee_db`) — port 3306
- **Go App** container (`go-app`) — port 8080

### Tài khoản mặc định (sau khi seed)

| Loại      | Email                         | Password   |
| --------- | ----------------------------- | ---------- |
| **Admin** | `chiquoc64@admin.company.dev` | `admin123` |

---

## 7. Troubleshooting

**Port conflict**

```bash
netstat -ano | findstr :8080   # Windows
taskkill /PID <PID> /F
# Hoặc đổi port trong .env / vite.config.js
```

**MySQL connection refused**

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS employee_db;"
# Kiểm tra DB_HOST, DB_PORT, DB_USER, DB_PASSWORD trong backend/.env
```

**Migration lỗi sau khi cập nhật schema**

```bash
# Nếu DB cũ tồn tại, chạy lệnh sau để thêm bảng và cột mới:
mysql -u root -p employee_db < backend/db_em.sql

# Hoặc dùng migration thủ công:
ALTER TABLE employees
  DROP COLUMN IF EXISTS position,
  ADD COLUMN position_id INT UNSIGNED NULL,
  ADD CONSTRAINT fk_emp_position
    FOREIGN KEY (position_id) REFERENCES positions(id)
    ON DELETE SET NULL ON UPDATE CASCADE;
```

**`npm install` lỗi**

```bash
rm -rf node_modules package-lock.json
npm cache clean --force && npm install
```

**`go mod tidy` lỗi**

```bash
go clean -modcache && go mod tidy
```

**Docker cache lỗi**

```bash
docker-compose down -v
docker-compose build --no-cache && docker-compose up
```

**CORS / API không gọi được**

- Kiểm tra `VITE_API_URL` trong `frontend/.env` trỏ đúng backend
- Backend hỗ trợ CORS cho `localhost:3000`, `localhost:5173`, `localhost:4200`

**Token hết hạn / 401 liên tục**

- JWT hết hạn sau 48h → đăng nhập lại
- Nếu DB reset mà token cũ còn → xoá localStorage trong browser DevTools

---

> **Tác giả**: Cao Chí Quốc
