-- Seed: Roles
INSERT INTO master.master_enum (created_by, modified_by, enum_value, enum_type, enum_code) VALUES
('system', 'system', 'Super Admin', 'role', 'SUPER_ADMIN'),
('system', 'system', 'Admin', 'role', 'ADMIN'),
('system', 'system', 'User', 'role', 'USER');

-- Seed: Password Policy
INSERT INTO master.master_enum (created_by, modified_by, enum_value, enum_type, enum_code) VALUES
('system', 'system', '90', 'PASSWORD_POLICY', 'EXPIRY_DAYS'),
('system', 'system', '5', 'PASSWORD_POLICY', 'MAX_FAILED_ATTEMPTS'),
('system', 'system', '30', 'PASSWORD_POLICY', 'LOCKOUT_DURATION_MINUTES');

-- Seed: Departments
INSERT INTO master.master_enum (created_by, modified_by, enum_value, enum_type, enum_code) VALUES
('system', 'system', 'IT Department', 'department', 'IT'),
('system', 'system', 'HR Department', 'department', 'HR'),
('system', 'system', 'Finance Department', 'department', 'FINANCE');

-- Seed: Locations
INSERT INTO master.master_enum (created_by, modified_by, enum_value, enum_type, enum_code) VALUES
('system', 'system', 'Jakarta', 'location', 'JKT'),
('system', 'system', 'Surabaya', 'location', 'SBY'),
('system', 'system', 'Bandung', 'location', 'BDG');

-- Seed: Default admin user (password: admin123 - bcrypt hashed)
INSERT INTO master.master_user (created_by, modified_by, user_name, name, email, password, role_id, is_status)
SELECT
    'system', 'system', 'admin', 'Administrator', 'admin@moyo.local',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    e.id, TRUE
FROM master.master_enum e
WHERE e.enum_type = 'role' AND e.enum_code = 'SUPER_ADMIN';
