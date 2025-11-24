-- Migration: Create countries table
-- Created: 2025-11-24

CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Audit fields
    created_by UUID NULL,
    updated_by UUID NULL,
    
    -- Optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    
    -- Country fields
    code VARCHAR(2) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    region VARCHAR(100),
    sub_region VARCHAR(100),
    currency_code VARCHAR(3),
    currency_symbol VARCHAR(10),
    phone_code VARCHAR(10),
    time_zone VARCHAR(50),
    flag VARCHAR(10),
    capital VARCHAR(100),
    population BIGINT,
    iso3166_alpha3 VARCHAR(3) UNIQUE,
    iso3166_numeric VARCHAR(3),
    
    -- Indexes
    CONSTRAINT idx_countries_deleted_at ON countries(deleted_at),
    CONSTRAINT idx_countries_status ON countries(status),
    CONSTRAINT idx_countries_code ON countries(code),
    CONSTRAINT idx_countries_name ON countries(name),
    CONSTRAINT idx_countries_name_en ON countries(name_en),
    CONSTRAINT idx_countries_region ON countries(region)
);

-- Create indexes separately (PostgreSQL syntax)
CREATE INDEX IF NOT EXISTS idx_countries_deleted_at ON countries(deleted_at);
CREATE INDEX IF NOT EXISTS idx_countries_status ON countries(status);
CREATE INDEX IF NOT EXISTS idx_countries_code ON countries(code);
CREATE INDEX IF NOT EXISTS idx_countries_name ON countries(name);
CREATE INDEX IF NOT EXISTS idx_countries_name_en ON countries(name_en);
CREATE INDEX IF NOT EXISTS idx_countries_region ON countries(region);

-- Trigger for updated_at
CREATE TRIGGER update_countries_updated_at BEFORE UPDATE ON countries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE countries IS 'Stores country information';
COMMENT ON COLUMN countries.code IS 'ISO 3166-1 alpha-2 country code (e.g., VN, US)';
COMMENT ON COLUMN countries.name IS 'Country name in local language';
COMMENT ON COLUMN countries.name_en IS 'Country name in English';
COMMENT ON COLUMN countries.region IS 'Geographic region (e.g., Asia, Europe)';
COMMENT ON COLUMN countries.currency_code IS 'ISO 4217 currency code (e.g., VND, USD)';
COMMENT ON COLUMN countries.phone_code IS 'International phone code (e.g., +84)';
COMMENT ON COLUMN countries.iso3166_alpha3 IS 'ISO 3166-1 alpha-3 country code (e.g., VNM)';

