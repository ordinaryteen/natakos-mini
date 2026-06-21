-- 1. Create Properties Table (The high-level asset block)
CREATE TABLE IF NOT EXISTS properties (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create Rooms Table (Weak entity dependent on a property structural layout)
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    property_id INT REFERENCES properties(id) ON DELETE CASCADE,
    room_number VARCHAR(20) NOT NULL,
    price_per_month NUMERIC(12, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'VACANT',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create Tenants Table (Holds identity mappings and Role-Based context anchors)
CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES rooms(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    phone_jid VARCHAR(50) UNIQUE NOT NULL, -- The crucial unique identifier mapping incoming chats
    role VARCHAR(20) DEFAULT 'TENANT',     -- The RBAC structural authorization string
    balance_owed NUMERIC(12, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- BOOTSTRAP DATA INJECTION (SEEDING DATA)
-- ==========================================

-- Insert our core property asset
INSERT INTO properties (name, address) 
VALUES ('Kos Natakos flagship', 'Mulyosari Gang 3, Kota Surabaya')
ON CONFLICT DO NOTHING;

-- Insert operational rooms
INSERT INTO rooms (property_id, room_number, price_per_month, status)
VALUES 
(1, 'Kamar 101', 1500000.00, 'OCCUPIED'),
(1, 'Kamar 102', 1750000.00, 'VACANT')
ON CONFLICT DO NOTHING;

-- Seed our real relational records (Mapping Budi to Kamar 101)
INSERT INTO tenants (room_id, name, phone_jid, role, balance_owed)
VALUES 
(1, 'Budi', '628123456789@s.whatsapp.net', 'TENANT', 1500000.00)
ON CONFLICT DO NOTHING;

-- =======================================================
-- PHASE 3 VIDEO 15: PGVECTOR SEMANTIC SCHEMA EXTENSION
-- =======================================================

-- 1. Initialize the vector math extensions inside the database cluster
CREATE EXTENSION IF NOT EXISTS vector;

-- 2. Create Property Rules Table for our RAG Engine
CREATE TABLE IF NOT EXISTS property_rules (
    id SERIAL PRIMARY KEY,
    property_id INT REFERENCES properties(id) ON DELETE CASCADE,
    rule_text TEXT NOT NULL,
    -- 1536 matches OpenAI's 'text-embedding-3-small' coordinates footprint length
    embedding VECTOR(1536), 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Seed an official text rule row for local semantic tracking verification
-- (We leave the embedding column NULL for a brief moment; our Python extraction script will compute and inject the matrix next!)
INSERT INTO property_rules (property_id, rule_text, embedding)
VALUES (1, 'Semua penghuni Kos Natakos flagship dilarang keras membawa hewan peliharaan (anjing, kucing, burung, dll) ke dalam area kamar maupun koridor.', NULL)
ON CONFLICT DO NOTHING;