CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_uid TEXT UNIQUE NOT NULL,
    track_number TEXT NOT NULL,
    entry TEXT NOT NULL,
    delivery_id BIGINT REFERENCES delivery(id),
    payment_id BIGINT REFERENCES payment(id),
    locale VARCHAR(5) NOT NULL,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INT,
    date_created TIMESTAMP DEFAULT now(),
    oof_shard VARCHAR(10) NOT NULL
);
