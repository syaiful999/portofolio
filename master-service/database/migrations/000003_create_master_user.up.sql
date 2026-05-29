CREATE TABLE master.master_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by VARCHAR(255) NOT NULL,
    created_date TIMESTAMP NOT NULL DEFAULT now(),
    modified_by VARCHAR(255) NOT NULL,
    modified_date TIMESTAMP NOT NULL DEFAULT now(),
    user_name VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255),
    role_id UUID REFERENCES master.master_enum(id),
    department_id UUID REFERENCES master.master_enum(id),
    outsource_id UUID REFERENCES master.master_enum(id),
    location_id UUID REFERENCES master.master_enum(id),
    picture TEXT,
    is_status BOOLEAN NOT NULL DEFAULT FALSE,
    password_last_changed TIMESTAMP,
    password_expiry_date TIMESTAMP,
    must_change_password BOOLEAN NOT NULL DEFAULT TRUE,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    last_failed_login TIMESTAMP,
    account_locked_until TIMESTAMP
);

CREATE UNIQUE INDEX idx_master_user_email ON master.master_user (LOWER(email)) WHERE is_active = TRUE;
CREATE UNIQUE INDEX idx_master_user_username ON master.master_user (LOWER(user_name)) WHERE is_active = TRUE;
CREATE INDEX idx_master_user_role ON master.master_user (role_id);
CREATE INDEX idx_master_user_department ON master.master_user (department_id);
CREATE INDEX idx_master_user_active ON master.master_user (is_active);
