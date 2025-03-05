ALTER TABLE users ADD COLUMN role_id INTEGER NOT NULL DEFAULT 3;
ALTER TABLE users ADD CONSTRAINT fk_users_roles FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT;
