-- V2__seed_initial_data.sql
-- Seed initial data for development/testing

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Function to generate UUIDv7 (time-ordered UUID)
-- UUIDv7 format: 48-bit timestamp + version bits + random bits
-- Based on RFC 9562 specification
CREATE OR REPLACE FUNCTION uuid_generate_v7()
RETURNS UUID AS $$
DECLARE
    unix_ts_ms BIGINT;
    uuid_bytes BYTEA;
    random_bytes BYTEA;
BEGIN
    -- Get current Unix timestamp in milliseconds
    unix_ts_ms := (EXTRACT(EPOCH FROM CLOCK_TIMESTAMP()) * 1000)::BIGINT;

    -- Generate 10 random bytes
    random_bytes := gen_random_bytes(10);

    -- Build UUIDv7 bytes
    uuid_bytes :=
        -- 48 bits of timestamp (6 bytes)
        SUBSTRING(INT8SEND(unix_ts_ms) FROM 3 FOR 6) ||
        -- Set version to 7 (0111) in the most significant 4 bits of byte 7
        SET_BYTE(SUBSTRING(random_bytes FROM 1 FOR 2), 0,
                 (GET_BYTE(random_bytes, 0) & 15) | 112) ||
        -- Set variant to 2 (10) in the most significant 2 bits of byte 9
        SET_BYTE(SUBSTRING(random_bytes FROM 3 FOR 8), 0,
                 (GET_BYTE(random_bytes, 2) & 63) | 128);

    RETURN ENCODE(uuid_bytes, 'hex')::UUID;
END;
$$ LANGUAGE plpgsql VOLATILE;

-- ==============================================================================
-- SEED DATA
-- ==============================================================================

-- Department 1: TI (Tecnologia da Informação)
-- This will be the parent department
DO $$
DECLARE
    dept_ti_id UUID;
BEGIN
    dept_ti_id := uuid_generate_v7();

    INSERT INTO departments (id, name, manager_id, parent_department_id, created_at, updated_at)
    VALUES (
        dept_ti_id,
        'TI - Tecnologia da Informação',
        uuid_generate_v7(), -- Temporary manager_id, will be updated
        NULL, -- No parent department
        NOW(),
        NOW()
    );

    -- Store for later use
    PERFORM set_config('seed.dept_ti_id', dept_ti_id::TEXT, false);
END $$;

-- Department 2: Desenvolvimento (child of TI)
DO $$
DECLARE
    dept_dev_id UUID;
    dept_ti_id UUID;
    manager_dev_id UUID;
    emp_ids UUID[];
BEGIN
    dept_ti_id := current_setting('seed.dept_ti_id')::UUID;
    dept_dev_id := uuid_generate_v7();
    manager_dev_id := uuid_generate_v7();

    -- Create department with temporary manager
    INSERT INTO departments (id, name, manager_id, parent_department_id, created_at, updated_at)
    VALUES (
        dept_dev_id,
        'Desenvolvimento',
        manager_dev_id,
        dept_ti_id, -- Parent is TI
        NOW(),
        NOW()
    );

    -- Create 5 employees for Desenvolvimento
    -- Employee 1: Manager
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        manager_dev_id,
        'Carlos Silva',
        '12345678901',
        '123456789',
        dept_dev_id,
        NOW(),
        NOW()
    );

    -- Employee 2: Developer
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'Ana Santos',
        '23456789012',
        '234567890',
        dept_dev_id,
        NOW(),
        NOW()
    );

    -- Employee 3: Developer
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'Roberto Lima',
        '34567890123',
        '345678901',
        dept_dev_id,
        NOW(),
        NOW()
    );

    -- Employee 4: Developer
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'Juliana Costa',
        '45678901234',
        '456789012',
        dept_dev_id,
        NOW(),
        NOW()
    );

    -- Employee 5: Developer
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'Pedro Oliveira',
        '56789012345',
        '567890123',
        dept_dev_id,
        NOW(),
        NOW()
    );

    PERFORM set_config('seed.dept_dev_id', dept_dev_id::TEXT, false);
END $$;

-- Department 3: Infraestrutura (child of TI)
DO $$
DECLARE
    dept_infra_id UUID;
    dept_ti_id UUID;
    manager_infra_id UUID;
BEGIN
    dept_ti_id := current_setting('seed.dept_ti_id')::UUID;
    dept_infra_id := uuid_generate_v7();
    manager_infra_id := uuid_generate_v7();

    -- Create department with temporary manager
    INSERT INTO departments (id, name, manager_id, parent_department_id, created_at, updated_at)
    VALUES (
        dept_infra_id,
        'Infraestrutura',
        manager_infra_id,
        dept_ti_id, -- Parent is TI
        NOW(),
        NOW()
    );

    -- Create 5 employees for Infraestrutura
    -- Employee 1: Manager
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        manager_infra_id,
        'Maria Fernandes',
        '67890123456',
        '678901234',
        dept_infra_id,
        NOW(),
        NOW()
    );

    -- Employee 2: DevOps
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'João Pereira',
        '78901234567',
        '789012345',
        dept_infra_id,
        NOW(),
        NOW()
    );

    -- Employee 3: SysAdmin
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'Fernanda Alves',
        '89012345678',
        '890123456',
        dept_infra_id,
        NOW(),
        NOW()
    );

    -- Employee 4: Network Engineer
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'Ricardo Souza',
        '90123456789',
        '901234567',
        dept_infra_id,
        NOW(),
        NOW()
    );

    -- Employee 5: Security Analyst
    INSERT INTO employees (id, name, cpf, rg, department_id, created_at, updated_at)
    VALUES (
        uuid_generate_v7(),
        'Camila Rodrigues',
        '01234567890',
        '012345678',
        dept_infra_id,
        NOW(),
        NOW()
    );

    PERFORM set_config('seed.dept_infra_id', dept_infra_id::TEXT, false);
END $$;

-- Update TI department manager (must be from TI department)
-- We'll use Carlos Silva (from Desenvolvimento) as the TI manager
DO $$
DECLARE
    dept_ti_id UUID;
    manager_ti_id UUID;
BEGIN
    dept_ti_id := current_setting('seed.dept_ti_id')::UUID;

    -- Get Carlos Silva's ID (manager of Desenvolvimento)
    SELECT id INTO manager_ti_id
    FROM employees
    WHERE cpf = '12345678901';

    -- Update TI department with the actual manager
    UPDATE departments
    SET manager_id = manager_ti_id
    WHERE id = dept_ti_id;

    -- Update Carlos Silva's department to TI
    UPDATE employees
    SET department_id = dept_ti_id
    WHERE id = manager_ti_id;
END $$;

-- ==============================================================================
-- VERIFICATION QUERIES (commented out, uncomment to verify data)
-- ==============================================================================

-- SELECT 'Departments:' as info;
-- SELECT id, name, manager_id, parent_department_id FROM departments ORDER BY created_at;

-- SELECT 'Employees:' as info;
-- SELECT e.id, e.name, e.cpf, d.name as department
-- FROM employees e
-- JOIN departments d ON e.department_id = d.id
-- ORDER BY d.name, e.created_at;

-- SELECT 'Department Hierarchy:' as info;
-- WITH RECURSIVE dept_hierarchy AS (
--     SELECT id, name, manager_id, parent_department_id, 0 as level
--     FROM departments
--     WHERE parent_department_id IS NULL
--     UNION ALL
--     SELECT d.id, d.name, d.manager_id, d.parent_department_id, dh.level + 1
--     FROM departments d
--     JOIN dept_hierarchy dh ON d.parent_department_id = dh.id
-- )
-- SELECT REPEAT('  ', level) || name as department_tree, level
-- FROM dept_hierarchy
-- ORDER BY level, name;

COMMENT ON FUNCTION uuid_generate_v7() IS 'Generates time-ordered UUIDv7 following RFC 9562';
