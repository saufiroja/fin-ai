\c finaidb;

DROP TABLE IF EXISTS financial_goals;
CREATE TABLE financial_goals (
    goal_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    target_amount INTEGER NOT NULL,
    current_amount INTEGER DEFAULT 0,
    target_date DATE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'paused')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_goals_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);