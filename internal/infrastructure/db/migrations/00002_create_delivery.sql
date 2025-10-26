-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid),
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(200),
    region VARCHAR(100),
    email VARCHAR(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS delivery;
-- +goose StatementEnd
