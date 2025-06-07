\c finaidb;

DROP TABLE IF EXISTS categories;
CREATE TABLE categories (
    category_id VARCHAR(250) PRIMARY KEY,
    name VARCHAR(250) NOT NULL,
    type VARCHAR(20) CHECK (type IN ('income', 'expense')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_category_id ON categories(category_id);

INSERT INTO categories (category_id, name, type, created_at, updated_at) VALUES
('cat001', 'Salary', 'income', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat002', 'Freelance Work', 'income', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat003', 'Investment Income', 'income', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat004', 'Transportation', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat005', 'Groceries', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat006', 'Food & Dining', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat007', 'Rent', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat008', 'Utilities', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat009', 'Entertainment', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat010', 'Healthcare', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat011', 'Education', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('cat012', 'Miscellaneous', 'expense', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
