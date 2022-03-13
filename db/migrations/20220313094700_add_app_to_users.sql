-- migrate:up
ALTER TABLE users ADD COLUMN app TEXT NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_app ON users(app);
DROP INDEX IF EXISTS idx_users_username;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_app ON users(username, app);

-- migrate:down
ALTER TABLE users DROP COLUMN app;
