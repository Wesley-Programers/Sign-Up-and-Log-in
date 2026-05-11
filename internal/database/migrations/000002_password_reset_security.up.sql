ALTER TABLE users
CHANGE password passwrod_hash VARCHAR(255) NOT NULL;

ALTER TABLE users
ADD COLUMN updated_at TIMESTAMP NOT NULL
DEFAULT CURRENT_TIMESTAMP
ON UPDATE CURRENT_TIMESTAMP;

ALTER TABLE reset_password
CHANGE token token_hash VARCHAR(255) NOT NULL;

ALTER TABLE reset_password
ADD COLUMN user_agent VARCHAR(255),
ADD COLUMN consumed_at TIMESTAMP NULL;

ALTER TABLE reset_password
MODIFY created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE reset_password
MODIFY used BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_user_id
ON reset_password(user_id);

CREATE INDEX idx_expires_at
ON reset_password(expires_at);

CREATE INDEX idx_token_used_expires
ON reset_password(token_hash, used, expires_at);

DROP TABLE attempts;