\c finaidb;

DROP TABLE IF EXISTS ai_recommendations;
CREATE TYPE recommendation_priority AS ENUM ('low', 'medium', 'high');
CREATE TYPE recommendation_type AS ENUM ('budget alert', 'saving tip', 'spending warning');
CREATE TABLE ai_recommendations (
    recommendation_id VARCHAR(250) PRIMARY KEY,
    user_id VARCHAR(250) NOT NULL,
    recommendation_type recommendation_type NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    content_embedding vector(1536) NOT NULL, -- for OpenAI embeddings
    priority recommendation_priority DEFAULT 'medium',
    is_read BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_recommendations_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX idx_recommendations_embedding 
ON ai_recommendations USING hnsw (content_embedding vector_cosine_ops);
