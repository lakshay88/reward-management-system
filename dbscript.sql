-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    user_password VARCHAR(255) NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Points Balance Table
CREATE TABLE points_balance (
    user_id INT PRIMARY KEY REFERENCES users(id),
    total_points INT DEFAULT 0,
    points_redeemed INT DEFAULT 0
);

-- Points History Table
CREATE TABLE points_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    transaction_id VARCHAR(50),
    points INT NOT NULL,
    points_type VARCHAR(10) CHECK (points_type IN ('earn', 'redeem', 'expired')),
    reason VARCHAR(255),
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

