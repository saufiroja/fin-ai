\c finaidb;

DROP TABLE IF EXISTS chat_sessions;
CREATE TABLE chat_sessions (
    chat_session_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    title TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_chat_sessions_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);