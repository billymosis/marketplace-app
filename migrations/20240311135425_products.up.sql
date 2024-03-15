CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(60) NOT NULL,
    price INTEGER NOT NULL, 
    image_url TEXT NOT NULL,
    stock INTEGER NOT NULL,
    condition VARCHAR(10) NOT NULL CHECK (condition IN ('new', 'second')),
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_purchasable BOOLEAN NOT NULL,
    purchase_count INTEGER NOT NULL DEFAULT 0, 
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL
);
