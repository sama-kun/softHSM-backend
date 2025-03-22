DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'blockchain_enum') THEN
        CREATE TYPE blockchain_enum AS ENUM ('ethereum', 'bitcoin', 'solana');
    END IF;
END $$;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE blockchain_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id SERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT,
    description TEXT,

    blockchain blockchain_enum NOT NULL,         -- Тип блокчейна (ethereum, bitcoin, solana и т. д.)
    network TEXT NOT NULL,            -- Сеть (mainnet, testnet, goerli и т. д.)

    address TEXT UNIQUE NOT NULL,     -- Публичный адрес
    encrypted_key TEXT NOT NULL,      -- Зашифрованный приватный ключ
    public_key TEXT NOT NULL,
    mnemonic_hash TEXT NOT NULL,
    salt TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_blockchain_keys_user_id ON blockchain_keys(user_id);
CREATE INDEX idx_blockchain_keys_blockchain ON blockchain_keys(blockchain);
CREATE INDEX idx_blockchain_keys_address ON blockchain_keys(address);
