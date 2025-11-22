-- +goose Up
-- 1. Fix the Service ID Foreign Key
ALTER TABLE bookings DROP CONSTRAINT bookings_service_id_fkey;
ALTER TABLE bookings 
    ADD CONSTRAINT bookings_service_id_fkey 
    FOREIGN KEY (service_id) REFERENCES services(id);

-- 2. Create a Join Table for Booking Workers (Many-to-Many relationship)
CREATE TABLE booking_workers (
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    PRIMARY KEY (booking_id, user_id)
);

-- +goose Down
DROP TABLE booking_workers;
ALTER TABLE bookings DROP CONSTRAINT bookings_service_id_fkey;
ALTER TABLE bookings 
    ADD CONSTRAINT bookings_service_id_fkey 
    FOREIGN KEY (service_id) REFERENCES locations(id);