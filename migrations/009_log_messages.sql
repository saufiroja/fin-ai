\c finaidb;

DROP TABLE IF EXISTS log_messages;
CREATE TABLE log_messages (
    log_messages_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    message TEXT,
    response TEXT,
    input_token INT DEFAULT 0,
    output_token INT DEFAULT 0,
    topic VARCHAR(50) NOT NULL,
    model VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_log_messages_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);