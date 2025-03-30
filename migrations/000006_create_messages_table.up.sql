CREATE TABLE messages
(
    id         SERIAL PRIMARY KEY,
    chat_id    INTEGER REFERENCES chats (id) ON DELETE CASCADE,
    sender_id  INTEGER REFERENCES users (id) ON DELETE CASCADE,
    type       VARCHAR(10) CHECK (type IN ('text', 'image', 'video', 'audio', 'file')),
    content    TEXT,
    file_url   TEXT,
    created_at TIMESTAMP DEFAULT now()
);