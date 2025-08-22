CREATE TABLE IF EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid),
    chrt_id BIGINT,
    track_number VARCHAR(50),
    price INTEGER,
    rid VARCHAR(50),
    name VARCHAR(100),
    sale INTEGER,
    size VARCHAR(10),
    total_price INTEGER,
    nm_id BIGINT,
    brand VARCHAR(100),
    status INTEGER
);

DROP TABLE IF EXISTS items;