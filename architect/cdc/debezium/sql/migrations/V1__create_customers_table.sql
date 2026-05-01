-- =============================================================================
-- V1__create_customers_table.sql
-- Creates the customers table tracked by Debezium CDC
-- =============================================================================

-- Drop if re-running migration in dev
DROP TABLE IF EXISTS public.customers;

-- Create customers table
CREATE TABLE public.customers (
    id         SERIAL          PRIMARY KEY,
    name       VARCHAR(255)    NOT NULL,
    email      VARCHAR(255)    NOT NULL UNIQUE,
    created_at TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Index on email for fast lookups
CREATE INDEX idx_customers_email ON public.customers (email);

-- -----------------------------------------------------------------------------
-- REPLICA IDENTITY FULL: required by Debezium to capture the full "before" row
-- image on UPDATE and DELETE events (default REPLICA IDENTITY only exposes the
-- primary key in the before image).
-- -----------------------------------------------------------------------------
ALTER TABLE public.customers REPLICA IDENTITY FULL;

-- Auto-update updated_at on every UPDATE
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_customers_updated_at
    BEFORE UPDATE ON public.customers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Confirm
SELECT 'customers table created successfully' AS status;
