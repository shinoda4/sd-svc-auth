ALTER TABLE users
ADD COLUMN reset_token TEXT,
ADD COLUMN reset_token_expire TIMESTAMP;
