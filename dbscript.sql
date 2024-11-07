-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    createdOn TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions Table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(50) UNIQUE NOT NULL,
    user_id INT REFERENCES users(id),
    transaction_amount DECIMAL(10, 2) NOT NULL,
    category VARCHAR(50) NOT NULL,
    transaction_date TIMESTAMP NOT NULL,
    product_code VARCHAR(50),
    points_earned INT NOT NULL,
    createdOn TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Points Balance Table
CREATE TABLE points_balance (
    user_id INT PRIMARY KEY REFERENCES users(id),
    total_points INT DEFAULT 0
);

-- Points History Table
CREATE TABLE points_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    points INT NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    reason VARCHAR(255),
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
