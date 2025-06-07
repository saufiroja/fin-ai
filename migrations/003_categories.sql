\c finaidb;

DROP TABLE IF EXISTS categories;
CREATE TABLE categories (
    category_id VARCHAR(250) PRIMARY KEY,
    name VARCHAR(250) NOT NULL,
    name_embedding vector(1536) NOT NULL, -- for OpenAI embeddings
    type VARCHAR(20) CHECK (type IN ('income', 'expense')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_category_id ON categories(category_id);
CREATE INDEX idx_categories_name_embedding 
ON categories USING hnsw (name_embedding vector_cosine_ops);
