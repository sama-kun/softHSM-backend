CREATE TABLE
  users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    master_password TEXT,
    login TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active_master BOOLEAN DEFAULT FALSE,
    -- is_active_faceid BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP
  );

CREATE INDEX idx_users_email ON users(email)