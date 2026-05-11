-- =============================================
-- Employee Management System - Database Schema
-- Cập nhật: Thêm bảng positions (chức vụ theo phòng ban)
-- =============================================

-- ROLES
CREATE TABLE roles (
  id TINYINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(50) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- DEPARTMENTS
CREATE TABLE departments (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- POSITIONS (Chức vụ thuộc về một phòng ban cụ thể)
CREATE TABLE positions (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(150) NOT NULL,
  department_id INT UNSIGNED NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_pos_department
    FOREIGN KEY (department_id)
    REFERENCES departments(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

  INDEX idx_position_dept (department_id),
  INDEX idx_position_name (name)
);

-- EMPLOYEES
CREATE TABLE employees (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  gender ENUM('male', 'female', 'other') DEFAULT NULL COMMENT 'male=Nam, female=Nữ, other=Khác',
  date_of_birth DATE,
  phone VARCHAR(20) UNIQUE NOT NULL,

  department_id INT UNSIGNED NULL,
  position_id   INT UNSIGNED NULL,

  salary DECIMAL(15,2) DEFAULT 0,
  hire_date DATE,
  status TINYINT DEFAULT 1 COMMENT '1=active, 0=inactive',
  avatar_url VARCHAR(500),

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_emp_department
    FOREIGN KEY (department_id)
    REFERENCES departments(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE,

  CONSTRAINT fk_emp_position
    FOREIGN KEY (position_id)
    REFERENCES positions(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE,

  INDEX idx_department_id (department_id),
  INDEX idx_position_id (position_id),
  INDEX idx_employee_name (name),
  INDEX idx_employee_phone (phone)
);

-- USERS (1-1 với employees)
CREATE TABLE users (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,

  role_id TINYINT UNSIGNED NOT NULL DEFAULT 2,

  employee_id INT UNSIGNED NULL,
  UNIQUE KEY uq_user_employee (employee_id), -- Enforce 1-1

  is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,

  CONSTRAINT fk_user_role
    FOREIGN KEY (role_id)
    REFERENCES roles(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,

  CONSTRAINT fk_user_employee
    FOREIGN KEY (employee_id)
    REFERENCES employees(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE,

  INDEX idx_role_id (role_id)
);

-- =============================================
-- SEED DATA
-- =============================================

INSERT INTO roles (name) VALUES ('admin'), ('employee');

INSERT INTO departments (name) VALUES
  ('IT'), ('HR'), ('Finance'), ('Marketing'), ('Sales'), ('Operations');

-- Positions cho từng phòng ban
INSERT INTO positions (name, department_id) VALUES
  -- IT (id=1)
  ('IT Manager',          1),
  ('System Admin',        1),
  ('Developer',           1),
  ('Help Desk',           1),
  -- HR (id=2)
  ('HR Manager',          2),
  ('Recruiter',           2),
  ('HR Officer',          2),
  ('Training Specialist', 2),
  -- Finance (id=3)
  ('CFO',                 3),
  ('Accountant',          3),
  ('Financial Analyst',   3),
  ('Auditor',             3),
  -- Marketing (id=4)
  ('Marketing Manager',   4),
  ('Content Writer',      4),
  ('Graphic Designer',    4),
  ('SEO Specialist',      4),
  -- Sales (id=5)
  ('Sales Manager',       5),
  ('Sales Rep',           5),
  ('Account Executive',   5),
  ('Business Dev',        5),
  -- Operations (id=6)
  ('Operations Manager',  6),
  ('Logistics Officer',   6),
  ('Process Analyst',     6),
  ('Supply Chain',        6);