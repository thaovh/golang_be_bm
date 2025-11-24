-- Migration: Create provinces table
-- Created: 2025-11-24

CREATE TABLE IF NOT EXISTS provinces (
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
    
    -- Foreign key to Country
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    
    -- Province fields
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    type VARCHAR(50),
    area DECIMAL(15,2),
    population BIGINT,
    coordinates VARCHAR(100),
    capital VARCHAR(100),
    postal_code VARCHAR(20),
    phone_prefix VARCHAR(10),
    sort_order INTEGER DEFAULT 0,
    
    -- Unique constraint: code must be unique within a country
    CONSTRAINT unique_province_code_per_country UNIQUE (country_id, code)
);

-- Create indexes separately (PostgreSQL syntax)
CREATE INDEX IF NOT EXISTS idx_provinces_deleted_at ON provinces(deleted_at);
CREATE INDEX IF NOT EXISTS idx_provinces_status ON provinces(status);
CREATE INDEX IF NOT EXISTS idx_provinces_country_id ON provinces(country_id);
CREATE INDEX IF NOT EXISTS idx_provinces_code ON provinces(code);
CREATE INDEX IF NOT EXISTS idx_provinces_name ON provinces(name);
CREATE INDEX IF NOT EXISTS idx_provinces_name_en ON provinces(name_en);
CREATE INDEX IF NOT EXISTS idx_provinces_type ON provinces(type);
CREATE INDEX IF NOT EXISTS idx_provinces_sort_order ON provinces(sort_order);

-- Trigger for updated_at
CREATE TRIGGER update_provinces_updated_at BEFORE UPDATE ON provinces
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE provinces IS 'Stores province/city information';
COMMENT ON COLUMN provinces.country_id IS 'Foreign key to countries table';
COMMENT ON COLUMN provinces.code IS 'Province/city code (unique within country)';
COMMENT ON COLUMN provinces.name IS 'Province/city name in local language';
COMMENT ON COLUMN provinces.name_en IS 'Province/city name in English';
COMMENT ON COLUMN provinces.type IS 'Administrative type: province, city, municipality';
COMMENT ON COLUMN provinces.area IS 'Area in square kilometers';
COMMENT ON COLUMN provinces.population IS 'Population count';
COMMENT ON COLUMN provinces.sort_order IS 'Sort order for display';

