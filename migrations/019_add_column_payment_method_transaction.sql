\c finaidb;

ALTER TABLE transactions
ADD COLUMN payment_method VARCHAR DEFAULT 'cash';
