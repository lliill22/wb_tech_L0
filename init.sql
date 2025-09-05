CREATE TABLE orders (
    order_uid VARCHAR PRIMARY KEY,
    track_number VARCHAR,
    entry VARCHAR,
    locale VARCHAR,
    internal_signature VARCHAR,
    customer_id VARCHAR,
    delivery_service VARCHAR,
    shardkey VARCHAR,
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR
);

CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR,
    phone VARCHAR,
    zip VARCHAR,
    city VARCHAR,
    address VARCHAR,
    region VARCHAR,
    email VARCHAR,
    CONSTRAINT deliveries_order_uid_key UNIQUE(order_uid)
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR,
    request_id VARCHAR,
    currency VARCHAR,
    provider VARCHAR,
    amount INT,
    payment_dt BIGINT,
    bank VARCHAR,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    CONSTRAINT payments_order_uid_key UNIQUE(order_uid)
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INT,
    track_number VARCHAR,
    price INT,
    rid VARCHAR,
    name VARCHAR,
    sale INT,
    size VARCHAR,
    total_price INT,
    nm_id INT,
    brand VARCHAR,
    status INT,
    CONSTRAINT items_order_uid_chrt_id_key UNIQUE(order_uid, chrt_id)
);
