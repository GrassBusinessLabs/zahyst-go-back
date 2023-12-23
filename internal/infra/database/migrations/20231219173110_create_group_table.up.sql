CREATE TABLE groups
(
    id           SERIAL PRIMARY KEY,
    title        TEXT,
    description  TEXT,
    user_id      INTEGER NOT NULL,
    access_code  TEXT,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE group_members (
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL,
    group_id     INTEGER NOT NULL,
    access_level TEXT,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);