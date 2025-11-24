-- Migration: Create wards table
-- Created: 2025-11-24

CREATE TABLE IF NOT EXISTS wards (
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
    
    -- Foreign key to Province
    province_id UUID NOT NULL REFERENCES provinces(id) ON DELETE RESTRICT,
    
    -- Ward fields
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    type VARCHAR(50),
    area DECIMAL(15,2),
    population BIGINT,
    coordinates VARCHAR(100),
    postal_code VARCHAR(20),
    address VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    
    -- Unique constraint: code must be unique within a province
    CONSTRAINT unique_ward_code_per_province UNIQUE (province_id, code)
);

-- Create indexes separately (PostgreSQL syntax)
CREATE INDEX IF NOT EXISTS idx_wards_deleted_at ON wards(deleted_at);
CREATE INDEX IF NOT EXISTS idx_wards_status ON wards(status);
CREATE INDEX IF NOT EXISTS idx_wards_province_id ON wards(province_id);
CREATE INDEX IF NOT EXISTS idx_wards_code ON wards(code);
CREATE INDEX IF NOT EXISTS idx_wards_name ON wards(name);
CREATE INDEX IF NOT EXISTS idx_wards_name_en ON wards(name_en);
CREATE INDEX IF NOT EXISTS idx_wards_type ON wards(type);
CREATE INDEX IF NOT EXISTS idx_wards_sort_order ON wards(sort_order);

-- Trigger for updated_at
CREATE TRIGGER update_wards_updated_at BEFORE UPDATE ON wards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE wards IS 'Stores ward/commune information';
COMMENT ON COLUMN wards.province_id IS 'Foreign key to provinces table';
COMMENT ON COLUMN wards.code IS 'Ward/commune code (unique within province)';
COMMENT ON COLUMN wards.name IS 'Ward/commune name in local language';
COMMENT ON COLUMN wards.name_en IS 'Ward/commune name in English';
COMMENT ON COLUMN wards.type IS 'Administrative type: ward, commune, town';
COMMENT ON COLUMN wards.area IS 'Area in square kilometers';
COMMENT ON COLUMN wards.population IS 'Population count';
COMMENT ON COLUMN wards.sort_order IS 'Sort order for display';

