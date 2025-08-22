CREATE TABLE IF EXISTS payment (
    transaction VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid),
    request_id VARCHAR(100),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount INTEGER,
    payment_dt BIGINT,
    bank VARCHAR(50),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

DROP TABLE IF EXISTS payment;