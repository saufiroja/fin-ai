\c finaidb;

DROP TABLE IF EXISTS chat_messages;
CREATE TYPE chat_message_sender AS ENUM ('user', 'ai');
CREATE TABLE chat_messages (
    chat_message_id VARCHAR(250) PRIMARY KEY,
    chat_session_id VARCHAR(250) NOT NULL,
    user_id VARCHAR(250) NOT NULL,
    sender chat_message_sender NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_chat_messages_chat_session FOREIGN KEY (chat_session_id) REFERENCES chat_sessions(chat_session_id),
    CONSTRAINT fk_chat_messages_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);