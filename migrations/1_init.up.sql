CREATE TABLE users (
    id SERIAL PRIMARY KEY,   
    name VARCHAR(50) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES users(id),
    receiver_id INT NOT NULL REFERENCES users(id),
    idempotency_key UUID UNIQUE NOT NULL,
	amount BIGINT NOT NULL,
    date_operation TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES users(id),
    CONSTRAINT fk_receiver FOREIGN KEY (receiver_id) REFERENCES users(id),
	CONSTRAINT chk_sender_not_receiver CHECK (sender_id <> receiver_id)
);

CREATE INDEX idx_transactions_sender_id ON transactions (sender_id);
CREATE INDEX idx_transactions_receiver_id ON transactions (receiver_id);