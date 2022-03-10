-- migrate:up
CREATE TABLE IF NOT EXISTS leaderboard (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    score INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_id ON leaderboard(user_id);

-- migrate:down
DROP TABLE leaderboard;
