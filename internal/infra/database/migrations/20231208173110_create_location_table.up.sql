CREATE TABLE IF NOT EXISTS locations
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL,
    type         TEXT,
    address      TEXT,
    title        TEXT,
    description  TEXT,
    lat          NUMERIC(9, 6),
    lon          NUMERIC(9, 6),
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
