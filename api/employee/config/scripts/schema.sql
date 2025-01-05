CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    gender CHAR(1) CHECK (gender IN ('M', 'F')),
    phone VARCHAR(15),
    email VARCHAR(100) UNIQUE NOT NULL,
    profile_picture VARCHAR(255),
    profile_type VARCHAR(50) CHECK (profile_type IN ('Medic', 'Technical Staff')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE shifts (
    id SERIAL PRIMARY KEY,
    shift_date DATE NOT NULL,
    shift_type INT NOT NULL, -- 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE employee_shifts (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    shift_id INT NOT NULL REFERENCES shifts(id) ON DELETE CASCADE,
    profile_type VARCHAR(50) NOT NULL, -- "Medic", "Technical", "Administrator"
    UNIQUE (employee_id, shift_id)
);

CREATE INDEX IF NOT EXISTS idx_employee_username ON employee(username);
CREATE INDEX IF NOT EXISTS idx_employee_email ON employee(email);
