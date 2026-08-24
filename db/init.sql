-- Schema initialization for the Devices API.
-- Mounted into /docker-entrypoint-initdb.d, Postgres runs this automatically
-- the FIRST time the container starts with an empty data volume.

CREATE TABLE IF NOT EXISTS devices (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    brand         VARCHAR(255) NOT NULL,
    state         VARCHAR(20)  NOT NULL DEFAULT 'AVAILABLE'
                  CHECK (state IN ('AVAILABLE', 'IN_USE', 'INACTIVE')),
    creation_time TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Speeds up the FindByBrand and FindByState endpoints.
CREATE INDEX IF NOT EXISTS idx_devices_brand ON devices (brand);
CREATE INDEX IF NOT EXISTS idx_devices_state ON devices (state);

-- Seed data for local development / manual testing.
-- creation_time is left to the column DEFAULT (NOW()) rather than hardcoded,
-- so timestamps always reflect when the container was actually initialized.
INSERT INTO devices (name, brand, state) VALUES
    ('iPhone 15 Pro',      'Apple',   'AVAILABLE'),
    ('MacBook Air M3',     'Apple',   'IN_USE'),
    ('Galaxy S24 Ultra',   'Samsung', 'AVAILABLE'),
    ('Galaxy Tab S9',      'Samsung', 'INACTIVE'),
    ('ThinkPad X1 Carbon', 'Lenovo',  'IN_USE'),
    ('Pixel 8',            'Google',  'AVAILABLE');

