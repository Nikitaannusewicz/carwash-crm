-- +goose Up

-- 1. Locations (Physical Places)
CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    adress TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 2. Users (identity)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL, -- 'admin', 'manager', 'worker'
    location_id BIGINT REFERENCES locations(id), -- Nullable (Admins might not have one)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()    
);

-- 3. BAYS (Resources at a location)
CREATE TABLE bay (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL, -- 'admin', 'manager', 'worker'
    location_id BIGINT REFERENCES locations(id),    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 4. Services (Products sold)
CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    duration_minutes INT NOT NULL,
    price_cents BIGINT NOT NULL
);

-- 5. Bookings (The core transaction)
CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES locations(id),
    bay_id BIGINT NOT NULL REFERENCES bays(id),
    service_id BIGINT NOT NULL REFERENCES locations(id),
    customer_id BIGINT NOT NULL REFERENCES users(id), -- With assumption that customers are Users for now

    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'confirmed', 'cancelled'

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for last lookup of "Is this bay free between X and Y?"
CREATE TABLE shifts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    location_id BIGINT NOT NULL REFERENCES locations(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

-- Index for "Who is working right now?"
CREATE INDEX idx_shifts_time ON shifts (location_id, start_time, end_time);

-- +goose Down
DROP TABLE IF EXISTS shifts;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS bays   ;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS locations;