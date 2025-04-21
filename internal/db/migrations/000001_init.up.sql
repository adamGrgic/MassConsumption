-- Enable UUID support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Function to auto-update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create Table TaskItem
CREATE TABLE IF NOT EXISTS pricerecord (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT,
    categoryid UUID,
    price TEXT,
    currencyid INT,
    processid UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create Table Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER set_updated_at_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create Table Categories
CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER set_updated_at_category
BEFORE UPDATE ON category
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS currencies (
    code CHAR(3) PRIMARY KEY,          -- ISO 4217 currency code (e.g., USD, EUR)
    name TEXT NOT NULL,                -- Currency name
    symbol TEXT NOT NULL          -- Symbol (e.g., $, €, ¥)
);


CREATE TABLE IF NOT EXISTS importprocess (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER set_updated_at_importprocess
BEFORE UPDATE ON importprocess
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS importentry (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    importlogid UUID,
    success BIT,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER set_updated_at_importentry
BEFORE UPDATE ON importentry
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();