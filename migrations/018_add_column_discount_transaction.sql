\c finaidb;

ALTER TABLE transactions
ADD COLUMN discount INT DEFAULT 0;
