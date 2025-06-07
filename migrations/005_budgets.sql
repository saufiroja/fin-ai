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

CREATE INDEX idx_budgets_user_month ON budgets(user_id, month, year);
CREATE INDEX idx_budgets_category ON budgets(category_id);
CREATE INDEX idx_budgets_budget_id ON budgets(budget_id);
