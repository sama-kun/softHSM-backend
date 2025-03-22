-- Удаляем индекс, если он есть
DROP INDEX IF EXISTS idx_blockchain_keys_blockchain;
DROP TABLE IF EXISTS blockchain_keys;

-- Удаляем ENUM, если он существует
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'blockchain_enum') THEN
        DROP TYPE blockchain_enum;
    END IF;
END $$;
