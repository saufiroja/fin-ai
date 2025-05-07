\c finaidb;

DROP TABLE IF EXISTS transactions;
CREATE TABLE transactions (
    transaction_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250),
    category_id VARCHAR(250),
    type VARCHAR(20) CHECK (type IN ('income', 'expense')),
    description TEXT NOT NULL,
    amount INTEGER NOT NULL,
    source VARCHAR(100) NOT NULL,
    transaction_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_transactions_category FOREIGN KEY (category_id) REFERENCES categories(category_id)
);
