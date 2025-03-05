DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'root') THEN
        CREATE DATABASE root;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'vi_database') THEN
        CREATE DATABASE vi_database;
    END IF;
END $$;