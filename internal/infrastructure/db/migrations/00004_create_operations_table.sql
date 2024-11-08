-- +goose Up
CREATE TYPE "public"."operation_statuses_enum" AS ENUM(
    'ADDED',
    'DEDUCTED'
);

CREATE TABLE operations (
		id SERIAL PRIMARY KEY,
		order_id VARCHAR NOT NULL,
		user_id integer NOT NULL,
        status "public"."operation_statuses_enum" NOT NULL,
        sum integer NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT FK_operations__order_id___orders__order_id FOREIGN KEY(order_id)
            REFERENCES orders(order_id) ON DELETE CASCADE ON UPDATE NO ACTION,
		CONSTRAINT FK_operations__user_id___users__id FOREIGN KEY(user_id)
		    REFERENCES users(id) ON DELETE CASCADE ON UPDATE NO ACTION
);

-- +goose Down
DROP TABLE operations;

DROP TYPE "public"."operation_statuses_enum";
