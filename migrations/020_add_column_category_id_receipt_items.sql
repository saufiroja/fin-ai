\c finaidb;

ALTER TABLE receipt_items ADD COLUMN category_id VARCHAR, ADD CONSTRAINT fk_receipt_items_category FOREIGN KEY (category_id) REFERENCES categories(category_id);
ALTER TABLE receipt_items
ADD COLUMN ai_category_confidence DECIMAL(3,2) DEFAULT 0.00;
