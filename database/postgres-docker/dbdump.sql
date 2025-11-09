DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS yards CASCADE;
DROP TABLE IF EXISTS blocks CASCADE;
DROP TABLE IF EXISTS yard_plans CASCADE;
DROP TABLE IF EXISTS container_positions CASCADE;

--users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO users (name, email, password)
VALUES (
        'Alice Johnson',
        'alice.j@example.com',
        '$2a$08$4O62PKNDy1kQvJeY9R5VjOAkSpjGt64M5UCc7UGGuQA52MKRqgdPC'
    ),
    (
        'Bob Smith',
        'bob.s@example.com',
        '$2a$08$4O62PKNDy1kQvJeY9R5VjOAkSpjGt64M5UCc7UGGuQA52MKRqgdPC'
    ),
    (
        'Charlie Brown',
        'charlie.b@example.com',
        '$2a$08$4O62PKNDy1kQvJeY9R5VjOAkSpjGt64M5UCc7UGGuQA52MKRqgdPC'
    ),
    (
        'Diana Prince',
        'diana.p@example.com',
        '$2a$08$4O62PKNDy1kQvJeY9R5VjOAkSpjGt64M5UCc7UGGuQA52MKRqgdPC'
    ),
    (
        'Evan Daniels',
        'evan.d@example.com',
        '$2a$08$4O62PKNDy1kQvJeY9R5VjOAkSpjGt64M5UCc7UGGuQA52MKRqgdPC'
    ),
    (
        'Fiona Glenn',
        'fiona.g@example.com',
        '$2a$08$4O62PKNDy1kQvJeY9R5VjOAkSpjGt64M5UCc7UGGuQA52MKRqgdPC'
    ),
    (
        'George King',
        'george.k@example.com',
        '$2a$08$4O62PKNDy1kQvJeY9R5VjOAkSpjGt64M5UCc7UGGuQA52MKRqgdPC'
    );
-- yards
CREATE TABLE yards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    location VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- blocks
CREATE TABLE blocks (
    id SERIAL PRIMARY KEY,
    yard_id INTEGER NOT NULL REFERENCES yards(id) ON DELETE RESTRICT,
    name VARCHAR(50) NOT NULL,
    slots INTEGER NOT NULL CHECK (slots > 0),
    rows INTEGER NOT NULL CHECK (rows > 0),
    tiers INTEGER NOT NULL CHECK (tiers > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (yard_id, name)
);
-- yard_plans
CREATE TABLE yard_plans (
    id SERIAL PRIMARY KEY,
    block_id INTEGER NOT NULL REFERENCES blocks(id) ON DELETE RESTRICT,
    plan_name VARCHAR(255) NOT NULL,
    slot_start INTEGER NOT NULL,
    slot_end INTEGER NOT NULL,
    row_start INTEGER NOT NULL,
    row_end INTEGER NOT NULL,
    container_size VARCHAR(5) NOT NULL CHECK (container_size IN ('20ft', '40ft')),
    container_height VARCHAR(5) NOT NULL CHECK (container_height IN ('8.6ft', '9.6ft')),
    container_type VARCHAR(50) NOT NULL,
    priority_stacking_direction VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (slot_start <= slot_end),
    CHECK (row_start <= row_end)
);
-- container_positions
CREATE TABLE container_positions (
    id SERIAL PRIMARY KEY,
    container_number VARCHAR(20) NOT NULL UNIQUE,
    block_id INTEGER NOT NULL REFERENCES blocks(id) ON DELETE RESTRICT,
    slot_number INTEGER NOT NULL,
    row_number INTEGER NOT NULL,
    tier_number INTEGER NOT NULL,
    container_size VARCHAR(5) NOT NULL,
    container_height VARCHAR(5) NOT NULL,
    container_type VARCHAR(50) NOT NULL,
    container_status VARCHAR(20) NOT NULL DEFAULT 'INBOUND',
    arrival_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    yard_plan_id INTEGER REFERENCES yard_plans(id) ON DELETE
    SET NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (block_id, slot_number, row_number, tier_number)
);

INSERT INTO yards (id, name, location) VALUES
(1, 'YRD-UTAMA', 'Terminal Kontainer Utama'),
(2, 'YRD-CADANGAN', 'Terminal Kapasitas Rendah'),
(3, 'YRD-REEFER', 'Area Khusus Pendingin'),
(4, 'YRD-IMPORT', 'Area Kedatangan Import'),
(5, 'YRD-EXPORT', 'Area Keberangkatan Export'),
(6, 'YRD-EMPTY', 'Gudang Kontainer Kosong'),
(7, 'YRD-HAZMAT', 'Area Material Berbahaya'),
(8, 'YRD-OOG', 'Area Out of Gauge'),
(9, 'YRD-BC', 'Area Breakdown Cargo'),
(10, 'YRD-LONGTERM', 'Penyimpanan Jangka Panjang')
ON CONFLICT (id) DO NOTHING;

INSERT INTO blocks (id, yard_id, name, slots, rows, tiers) VALUES
(1, 1, 'LC01', 10, 5, 5),    -- Yard 1, Dry General
(2, 1, 'LC02', 10, 5, 5),    -- Yard 1, Dry General
(3, 2, 'C01', 8, 4, 4),      -- Yard 2, Cadangan
(4, 3, 'RF01', 6, 3, 4),     -- Yard 3, Reefer (Cold)
(5, 3, 'RF02', 6, 3, 4),     -- Yard 3, Reefer (Warm)
(6, 4, 'IM01', 12, 6, 3),    -- Yard 4, Import Area
(7, 5, 'EX01', 12, 6, 3),    -- Yard 5, Export Area
(8, 6, 'E01', 20, 10, 3),    -- Yard 6, Empty Area
(9, 7, 'HAZ1', 4, 2, 2),     -- Yard 7, Hazmat (Small Capacity)
(10, 8, 'OOG1', 4, 4, 1)     -- Yard 8, OOG (Low Tier)
ON CONFLICT (id) DO NOTHING;


INSERT INTO yard_plans (id, block_id, plan_name, slot_start, slot_end, row_start, row_end, container_size, container_height, container_type, priority_stacking_direction, is_active) VALUES
(1, 1, '20ft DRY LC01 Low Row', 1, 4, 1, 2, '20ft', '8.6ft', 'DRY', 'BOTTOM_UP', TRUE),
(2, 1, '40ft DRY LC01 Mid Row', 5, 10, 1, 3, '40ft', '9.6ft', 'DRY', 'LEFT_RIGHT', TRUE),
(3, 1, '20ft OT LC01 High Row', 1, 10, 4, 5, '20ft', '8.6ft', 'Open Top', 'BOTTOM_UP', TRUE),
(4, 4, '40ft REEFER RF01 All', 1, 6, 1, 3, '40ft', '9.6ft', 'Reefer', 'BOTTOM_UP', TRUE),
(5, 2, '20ft DRY LC02 A', 1, 6, 1, 3, '20ft', '8.6ft', 'DRY', 'LEFT_RIGHT', TRUE),
(6, 2, '40ft DRY LC02 B', 7, 10, 1, 5, '40ft', '9.6ft', 'DRY', 'BOTTOM_UP', TRUE),
(7, 3, '20ft C01 Cadangan', 1, 8, 1, 4, '20ft', '8.6ft', 'DRY', 'BOTTOM_UP', TRUE),
(8, 6, '40ft IM01 Import', 1, 12, 1, 6, '40ft', '9.6ft', 'DRY', 'LEFT_RIGHT', TRUE),
(9, 8, '20ft E01 Empty', 1, 20, 1, 10, '20ft', '8.6ft', 'EMPTY', 'BOTTOM_UP', TRUE),
(10, 9, '20ft HAZ1 Hazmat', 1, 4, 1, 2, '20ft', '8.6ft', 'HAZMAT', 'BOTTOM_UP', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO container_positions (id, container_number, block_id, slot_number, row_number, tier_number, container_size, container_height, container_type, container_status, yard_plan_id) VALUES
(1, 'ALFI000001', 1, 1, 1, 1, '20ft', '8.6ft', 'DRY', 'STORAGE', 1),
(2, 'ALFI000002', 1, 1, 1, 2, '20ft', '8.6ft', 'DRY', 'STORAGE', 1),
(3, 'ALFI000003', 1, 5, 1, 1, '40ft', '9.6ft', 'DRY', 'STORAGE', 2), 
(4, 'ALFI000004', 1, 7, 2, 1, '40ft', '9.6ft', 'DRY', 'STORAGE', 2), 
(5, 'ALFI000005', 1, 10, 4, 1, '20ft', '8.6ft', 'Open Top', 'STORAGE', 3),
(6, 'ALFI000006', 1, 9, 1, 1, '20ft', '8.6ft', 'DRY', 'STORAGE', 2),
(7, 'ALFI000007', 4, 1, 1, 1, '40ft', '9.6ft', 'Reefer', 'STORAGE', 4), 
(8, 'ALFI000008', 4, 1, 1, 2, '40ft', '9.6ft', 'Reefer', 'STORAGE', 4), 
(9, 'ALFI000009', 2, 2, 2, 1, '20ft', '8.6ft', 'DRY', 'STORAGE', 5), 
(10, 'ALFI000010', 1, 1, 1, 3, '20ft', '8.6ft', 'DRY', 'STORAGE', 1)
ON CONFLICT (id) DO NOTHING;

SELECT setval('container_positions_id_seq', (SELECT MAX(id) FROM container_positions) + 1, false);