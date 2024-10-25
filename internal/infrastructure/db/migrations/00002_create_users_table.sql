-- +goose Up
CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		login VARCHAR NOT NULL,
		password_hash VARCHAR NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (login)
);

-- +goose Down
DROP TABLE users;
