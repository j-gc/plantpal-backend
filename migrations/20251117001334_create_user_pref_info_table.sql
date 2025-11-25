-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE gender_options (
    option VARCHAR(20) PRIMARY KEY
);

INSERT INTO gender_options (option) VALUES
('male'),
('female'),
('non-binary'),
('other'),
('prefer_not_to_say');


CREATE TABLE user_preferences (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id),
    gender VARCHAR(20) NOT NULL REFERENCES gender_options(option),
    max_age INT NULL,
    min_age INT NULL,
    bio TEXT NOT NULL,
    profile_visibility BOOLEAN DEFAULT TRUE,
    profile_picture_url VARCHAR(255) NOT NULL,
    notification_settings JSONB NULL,
    looking_for VARCHAR(100) NULL,
    preferred_gender VARCHAR(20) REFERENCES gender_options(option),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS gender_options;

-- +goose StatementEnd
