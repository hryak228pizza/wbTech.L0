CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50),
    entry VARCHAR(10),
    locale VARCHAR(10),
    internal_signature VARCHAR(100),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INTEGER,
    date_created TIMESTAMP WITH TIME ZONE,
    oof_shard VARCHAR(10)
);

CREATE TABLE delivery (
    order_uid VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid),
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(200),
    region VARCHAR(100),
    email VARCHAR(100)
);

CREATE TABLE payment (
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

CREATE TABLE items (
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


INSERT INTO orders (
    order_uid, track_number, entry, locale, internal_signature, 
    customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
) VALUES (
    'b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '', 
    'test', 'meest', '9', 99, '2021-11-26T06:22:19Z', '1'
);

INSERT INTO delivery (
    order_uid, name, phone, zip, city, address, region, email
) VALUES (
    'b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809', 
    'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com'
);

INSERT INTO payment (
    transaction, request_id, currency, provider, amount, payment_dt, 
    bank, delivery_cost, goods_total, custom_fee
) VALUES (
    'b563feb7b2b84b6test', '', 'USD', 'wbpay', 1817, 1637907727, 
    'alpha', 1500, 317, 0
);

INSERT INTO items (
    order_uid, chrt_id, track_number, price, rid, name, sale, size, 
    total_price, nm_id, brand, status
) VALUES (
    'b563feb7b2b84b6test', 9934930, 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest', 
    'Mascaras', 30, '0', 317, 2389212, 'Vivienne Sabo', 202
);