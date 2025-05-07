\c finaidb;

DROP TABLE IF EXISTS categories;
CREATE TABLE categories (
    category_id VARCHAR(250) PRIMARY KEY,
    name VARCHAR(250) NOT NULL,
    type VARCHAR(20) CHECK (type IN ('income', 'expense')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
