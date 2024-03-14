CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    bank_account_id VARCHAR(255) REFERENCES accounts(id) NOT NULL,
    product_id INTEGER REFERENCES products(id) NOT NULL,
    payment_proof_image_url VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 1)
);
