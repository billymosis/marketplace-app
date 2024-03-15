CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id) NOT NULL,
    product_id INTEGER REFERENCES products(id) NOT NULL,
    payment_proof_image_url VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 1)
);
