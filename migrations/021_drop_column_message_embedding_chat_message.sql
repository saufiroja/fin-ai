\c finaidb;

ALTER TABLE chat_messages
DROP COLUMN IF EXISTS message_embedding;
