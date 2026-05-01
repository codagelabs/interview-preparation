-- =============================================================================
-- V2__seed_customers.sql
-- Sample data for development / testing CDC events
-- =============================================================================

INSERT INTO public.customers (name, email) VALUES
    ('Alice Johnson',   'alice@example.com'),
    ('Bob Smith',       'bob@example.com'),
    ('Carol Williams',  'carol@example.com'),
    ('David Brown',     'david@example.com'),
    ('Eve Davis',       'eve@example.com')
ON CONFLICT (email) DO NOTHING;

SELECT COUNT(*) AS seeded_rows FROM public.customers;
