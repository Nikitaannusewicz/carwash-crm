-- +goose Up

-- Drop the foreign key constraint from bookings first
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_bay_id_fkey;

-- Drop the incorrect bays table
DROP TABLE IF EXISTS bays CASCADE;

-- Create the correct bay table (singular, matching repository)
CREATE TABLE bay (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES locations(id),
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Re-add the foreign key constraint from bookings to the new bay table
ALTER TABLE bookings
    ADD CONSTRAINT bookings_bay_id_fkey
    FOREIGN KEY (bay_id) REFERENCES bay(id);

-- +goose Down

-- Drop the foreign key constraint from bookings
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_bay_id_fkey;

-- Drop the correct bay table
DROP TABLE IF EXISTS bay CASCADE;

-- Recreate the old incorrect bays table for rollback
CREATE TABLE bays (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    location_id BIGINT REFERENCES locations(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Re-add the foreign key constraint to the old table
ALTER TABLE bookings
    ADD CONSTRAINT bookings_bay_id_fkey
    FOREIGN KEY (bay_id) REFERENCES bays(id);
