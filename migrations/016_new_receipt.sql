\c finaidb;

DROP TABLE IF EXISTS receipts CASCADE;
CREATE TABLE receipts (
    receipt_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    merchant_name TEXT,
    sub_total INT NOT NULL,
    total_discount INT DEFAULT 0,
    total_shopping INT NOT NULL,
    metadata JSONB,
    extracted_receipt JSONB,
    extracted_receipt_embedding vector(1536) NOT NULL, -- for OpenAI embeddings
    confirmed BOOLEAN DEFAULT FALSE,
    transaction_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_receipts_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX idx_receipts_user ON receipts(user_id);
CREATE INDEX idx_receipts_receipt_id ON receipts(receipt_id);
CREATE INDEX idx_receipts_embedding
ON receipts USING hnsw (extracted_receipt_embedding vector_cosine_ops);

DROP TABLE IF EXISTS receipt_items;
CREATE TABLE receipt_items (
    receipt_item_id VARCHAR(250) PRIMARY KEY,
    receipt_id VARCHAR(250),
    item_name TEXT NOT NULL,
    item_quantity INT NOT NULL,
    item_price INT NOT NULL,
    item_price_total INT NOT NULL,
    item_discount INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_receipt_items_receipt FOREIGN KEY (receipt_id) REFERENCES receipts(receipt_id)
);

CREATE INDEX idx_receipt_items_receipt ON receipt_items(receipt_id);
CREATE INDEX idx_receipt_items_item_name ON receipt_items(item_name);
CREATE INDEX idx_receipt_items_receipt_item_id ON receipt_items(receipt_item_id);
