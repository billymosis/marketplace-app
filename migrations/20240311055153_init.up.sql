CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username VARCHAR(15) NOT NULL,
    password VARCHAR(72) NOT NULL,
    name VARCHAR(50) NOT NULL,
    UNIQUE(username)
);
