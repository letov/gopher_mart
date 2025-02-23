-- +goose Up
CREATE TYPE "public"."order_statuses_enum" AS ENUM(
    'NEW',
    'INVALID',
    'PROCESSING',
    'PROCESSED'
);

CREATE TABLE orders (
		id SERIAL PRIMARY KEY,
		order_id VARCHAR NOT NULL,
		user_id integer NOT NULL,
        status "public"."order_statuses_enum" NOT NULL DEFAULT 'NEW',
        accrual integer DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (order_id),
		CONSTRAINT FK_orders__user_id___users__id FOREIGN KEY(user_id)
		    REFERENCES users(id) ON DELETE CASCADE ON UPDATE NO ACTION
);

-- +goose Down
DROP TABLE orders;

DROP TYPE "public"."order_statuses_enum";
