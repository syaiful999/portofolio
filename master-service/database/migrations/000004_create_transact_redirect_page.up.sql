CREATE TABLE transaction.transact_redirect_page (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_date TIMESTAMP NOT NULL DEFAULT now(),
    created_by VARCHAR(255) NOT NULL,
    modified_date TIMESTAMP NOT NULL DEFAULT now(),
    modified_by VARCHAR(255) NOT NULL,
    expired_date TIMESTAMP NOT NULL DEFAULT (now() + INTERVAL '24 hours'),
    page_name VARCHAR(100) NOT NULL
);

CREATE INDEX idx_redirect_page_active ON transaction.transact_redirect_page (is_active);
CREATE INDEX idx_redirect_page_created_by ON transaction.transact_redirect_page (created_by);
