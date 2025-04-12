CREATE TABLE vacancy_feed_views
(
    id             SERIAL PRIMARY KEY,
    user_id        INT       NOT NULL UNIQUE,
    last_viewed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
