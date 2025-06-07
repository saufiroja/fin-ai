\c finaidb;

DROP TABLE IF EXISTS ai_summaries;
CREATE TABLE ai_summaries (
    summary_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    period_type VARCHAR(20) CHECK (period_type IN ('weekly', 'monthly', 'yearly')),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    summary_data JSONB NOT NULL, -- insights, recommendations, highlights
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_summaries_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);