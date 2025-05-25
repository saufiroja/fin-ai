\c finaidb;

DROP TABLE IF EXISTS model_registries;
CREATE TABLE model_registries (
    model_registry_id VARCHAR(250) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

INSERT INTO model_registries (model_registry_id, name) VALUES
('01JW320X5V1N7QQS2YAFC6FS71', 'gpt-4o-mini'),
('01JW32155CVG7J5SEJE8MWXQW4', 'o3-mini');
