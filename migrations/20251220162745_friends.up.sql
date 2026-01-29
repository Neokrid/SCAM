CREATE TABLE IF NOT EXISTS friendships (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status        BOOL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_friendships PRIMARY KEY (user_id, friend_id),
    CONSTRAINT chk_not_self_friend CHECK (user_id <> friend_id)
);

CREATE TABLE IF NOT EXISTS friend_requests (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status        VARCHAR(20) NOT NULL,
    CONSTRAINT uq_friend_request UNIQUE (sender_id, receiver_id),
    CONSTRAINT chk_not_self_request CHECK (sender_id <> receiver_id)
);