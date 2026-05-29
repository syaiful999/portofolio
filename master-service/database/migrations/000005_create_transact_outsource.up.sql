CREATE TABLE transaction.transact_outsource (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by VARCHAR(255) NOT NULL,
    created_date TIMESTAMP NOT NULL DEFAULT now(),
    modified_by VARCHAR(255) NOT NULL,
    modified_date TIMESTAMP NOT NULL DEFAULT now(),
    outsource_id UUID NOT NULL REFERENCES master.master_enum(id),
    department_id UUID NOT NULL REFERENCES master.master_enum(id)
);

CREATE INDEX idx_transact_outsource_outsource ON transaction.transact_outsource (outsource_id);
CREATE INDEX idx_transact_outsource_department ON transaction.transact_outsource (department_id);
CREATE INDEX idx_transact_outsource_active ON transaction.transact_outsource (is_active);
