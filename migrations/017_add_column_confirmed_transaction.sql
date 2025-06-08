\c finaidb;

ALTER TABLE transactions
ADD COLUMN confirmed BOOLEAN DEFAULT FALSE;
