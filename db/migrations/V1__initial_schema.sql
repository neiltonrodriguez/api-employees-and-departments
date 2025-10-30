-- V1__initial_schema.sql
-- Initial database schema for employees and departments

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create departments table
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    manager_id UUID NOT NULL,
    parent_department_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Foreign key to itself for hierarchy
    CONSTRAINT fk_parent_department FOREIGN KEY (parent_department_id)
        REFERENCES departments(id) ON DELETE SET NULL
);

-- Create employees table
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cpf VARCHAR(11) NOT NULL,
    rg VARCHAR(20),
    department_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Foreign key to departments
    CONSTRAINT fk_department FOREIGN KEY (department_id)
        REFERENCES departments(id) ON DELETE RESTRICT,

    -- Unique constraints
    CONSTRAINT uk_cpf UNIQUE (cpf),
    CONSTRAINT uk_rg UNIQUE (rg)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_departments_deleted_at ON departments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_departments_manager_id ON departments(manager_id);
CREATE INDEX IF NOT EXISTS idx_departments_parent_id ON departments(parent_department_id);

CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees(deleted_at);
CREATE INDEX IF NOT EXISTS idx_employees_cpf ON employees(cpf);
CREATE INDEX IF NOT EXISTS idx_employees_rg ON employees(rg) WHERE rg IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees(department_id);

-- Comments for documentation
COMMENT ON TABLE departments IS 'Stores department information with hierarchical structure';
COMMENT ON TABLE employees IS 'Stores employee information linked to departments';
COMMENT ON COLUMN departments.parent_department_id IS 'Self-referencing FK for department hierarchy';
COMMENT ON COLUMN employees.cpf IS 'Brazilian CPF - must be unique and valid';
COMMENT ON COLUMN employees.rg IS 'Brazilian RG - optional but must be unique if provided';
