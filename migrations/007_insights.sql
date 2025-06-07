\c finaidb;

DROP TABLE IF EXISTS insights;
CREATE TABLE insights (
    insight_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    insight_type VARCHAR(100),
    content JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_insights_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX idx_insights_user ON insights(user_id);
CREATE INDEX idx_insights_insight_id ON insights(insight_id);
CREATE INDEX idx_insights_type ON insights(insight_type);
