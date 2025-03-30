CREATE TABLE message_seen_by
(
    message_id INTEGER REFERENCES messages (id) ON DELETE CASCADE,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (message_id, user_id)
);
