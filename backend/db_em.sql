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

-- EMPLOYEES (❌ BỎ EMAIL)
CREATE TABLE employees (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  gender ENUM('male', 'female', 'other') DEFAULT NULL COMMENT 'male=Nam, female=Nữ, other=Khác',
  date_of_birth DATE,

  phone VARCHAR(20) UNIQUE, -- thêm unique cho sạch data

  department_id INT UNSIGNED NULL,
  position VARCHAR(100),
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

  INDEX idx_department_id (department_id),
  INDEX idx_employee_name (name),
  INDEX idx_employee_phone (phone)
);

-- USERS (🔥 ENFORCE 1-1)
CREATE TABLE users (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,

  role_id TINYINT UNSIGNED NOT NULL DEFAULT 2,

  employee_id INT UNSIGNED NULL,
  UNIQUE KEY uq_user_employee (employee_id), -- 🔥 QUAN TRỌNG

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