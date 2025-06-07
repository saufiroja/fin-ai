\c finaidb;

DROP TABLE IF EXISTS receipts;
CREATE TABLE receipts (
    receipt_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250),
    file_path TEXT,
    extracted_receipt TEXT,
    transaction_id VARCHAR(250),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_receipts_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_receipts_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

CREATE INDEX idx_receipts_user ON receipts(user_id);
CREATE INDEX idx_receipts_transaction ON receipts(transaction_id);
CREATE INDEX idx_receipts_receipt_id ON receipts(receipt_id);
