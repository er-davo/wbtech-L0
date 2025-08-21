CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email VARCHAR(255) NOT NULL
);