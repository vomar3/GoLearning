CREATE TABLE IF NOT EXISTS options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    votes_count BIGINT DEFAULT 0
);

CREATE INDEX idx_options_poll_id ON options(poll_id);