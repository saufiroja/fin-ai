\c finaidb;

DROP TABLE IF EXISTS budgets;
CREATE TABLE budgets (
    budget_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250),
    category_id VARCHAR(250),
    amount_limit INTEGER NOT NULL,
    month INTEGER,
    year INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_budgets_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_budgets_category FOREIGN KEY (category_id) REFERENCES categories(category_id)
);
