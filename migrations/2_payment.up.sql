CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank TEXT NOT NULL,
    delivery_cost NUMERIC(10,2) NOT NULL,
    goods_total NUMERIC(10,2) NOT NULL,
    custom_fee NUMERIC(10,2) NOT NULL
);