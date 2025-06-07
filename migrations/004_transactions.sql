\c finaidb;

DROP TABLE IF EXISTS transactions;
CREATE TABLE transactions (
    transaction_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250),
    category_id VARCHAR(250),
    type VARCHAR(20) CHECK (type IN ('income', 'expense')),
    description TEXT NOT NULL,
    description_embedding vector(1536) NOT NULL, -- for OpenAI embeddings
    amount INTEGER NOT NULL,
    source VARCHAR(100) NOT NULL,
    transaction_date TIMESTAMP,
    ai_category_confidence DECIMAL(3,2) DEFAULT 0.00,
    is_auto_categorized BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_transactions_category FOREIGN KEY (category_id) REFERENCES categories(category_id)
);

CREATE INDEX idx_transactions_transaction_id ON transactions(transaction_id);
CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date DESC);
CREATE INDEX idx_transactions_category_type ON transactions(category_id, type);
CREATE INDEX idx_transactions_description_embedding 
ON transactions USING hnsw (description_embedding vector_cosine_ops);
CREATE INDEX idx_transactions_user_id ON transactions (user_id);
CREATE INDEX idx_transactions_user_date_type
ON transactions (user_id, transaction_date DESC, type)
INCLUDE (transaction_id, amount, category_id);
CREATE INDEX idx_transactions_ai_categorization
ON transactions (user_id, is_auto_categorized, ai_category_confidence DESC)
WHERE ai_category_confidence > 0.5;
CREATE INDEX idx_transactions_user_month_category
ON transactions (user_id, date_trunc('month', transaction_date), category_id)
INCLUDE (amount, type);
CREATE INDEX idx_transactions_user_transaction_date
ON transactions (user_id, transaction_date DESC);
