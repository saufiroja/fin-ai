\c finaidb;

DROP TABLE IF EXISTS ocr_cache;
CREATE TABLE ocr_cache (
    cache_id VARCHAR(250) PRIMARY KEY,
    file_hash VARCHAR(64) UNIQUE NOT NULL, -- SHA-256 hash dari file
    extracted_data JSONB NOT NULL,
    confidence_score DECIMAL(3,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);