-- +goose Up
CREATE TYPE "public"."event_action_enum" AS ENUM(
    'SAVE_USER',
    'LOGIN',
    'CALC_ACCRUAL'
);

CREATE TABLE events (
		id SERIAL PRIMARY KEY,
		root_id VARCHAR NOT NULL,
		action "public"."event_action_enum" NOT NULL,
		payload jsonb NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE events

DROP TYPE "public"."event_action_enum";
