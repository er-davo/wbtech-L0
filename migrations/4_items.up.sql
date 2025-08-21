CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    chrt_id INT NOT NULL,
    track_number TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    rid TEXT NOT NULL,
    name TEXT NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(100) NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INT NOT NULL
);