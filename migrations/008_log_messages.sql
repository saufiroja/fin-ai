\c finaidb;

DROP TABLE IF EXISTS log_messages;
CREATE TABLE log_messages (
    log_messages_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    role VARCHAR(100),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_log_messages_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);