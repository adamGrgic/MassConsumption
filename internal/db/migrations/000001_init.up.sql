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
CREATE TABLE IF NOT EXISTS taskitem (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT,
    description TEXT,
    statusid INT,
    listid BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER set_updated_at_taskitem
BEFORE UPDATE ON taskitem
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create Table TaskList
CREATE TABLE IF NOT EXISTS tasklist (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER set_updated_at_tasklist
BEFORE UPDATE ON taskitem
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

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

-- Create Table Accounts
CREATE TABLE IF NOT EXISTS Accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER set_updated_at_accounts
BEFORE UPDATE ON Accounts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- Create Table Statuses
CREATE TABLE IF NOT EXISTS TaskStatus (
    id SERIAL PRIMARY KEY,
    name TEXT
);