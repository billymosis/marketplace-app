CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    bank_name VARCHAR(15) NOT NULL ,
    bank_account_name VARCHAR(15) NOT NULL,
    bank_account_number VARCHAR(15) NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL
);
