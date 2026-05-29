CREATE TABLE master.master_enum (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by VARCHAR(255) NOT NULL,
    created_date TIMESTAMP NOT NULL DEFAULT now(),
    modified_by VARCHAR(255) NOT NULL,
    modified_date TIMESTAMP NOT NULL DEFAULT now(),
    enum_value VARCHAR(255) NOT NULL,
    enum_type VARCHAR(100) NOT NULL,
    enum_code VARCHAR(100)
);

CREATE INDEX idx_master_enum_type ON master.master_enum (enum_type);
CREATE INDEX idx_master_enum_code ON master.master_enum (enum_code);
CREATE INDEX idx_master_enum_active ON master.master_enum (is_active);
