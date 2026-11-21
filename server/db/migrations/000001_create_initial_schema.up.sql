-- Migration: Initial Schema (Minimalist, Dynamic & High-Performance Structure)
-- Description: Creates master entities, academic hierarchy, course offerings, and scheduling grid.

BEGIN;

-- ==========================================
-- 1. MASTER INFRASTRUCTURE TABLES (STATIC DATA)
-- ==========================================

-- Academic Departments
-- CREATE TABLE IF NOT EXISTS departments (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     code VARCHAR(50) UNIQUE NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
--    );

-- Professors / Faculty Members
CREATE TABLE IF NOT EXISTS professors (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) UNIQUE NOT NULL,
    department VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- Room Types Catalog (Free-Text Lookup: Amphi, Computer Lab, TD Room, etc.)
CREATE TABLE IF NOT EXISTS room_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- Physical Rooms & Auditoriums (Simplified without capacity)
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    room_type_id INT NOT NULL REFERENCES room_types(id),
    code VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- Static Academic Specializations / Levels (e.g., L1 Applied Math, M1 PDE)
CREATE TABLE IF NOT EXISTS academic_levels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- Global Subject Catalog (Linked to Academic Level)
CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- ==========================================
-- 2. ACADEMIC PERIOD & DYNAMIC HIERARCHY
-- ==========================================

-- Academic Years Master Table (e.g., "2026-2027")
CREATE TABLE IF NOT EXISTS academic_years (
    id SERIAL PRIMARY KEY,
    year_label VARCHAR(20) UNIQUE NOT NULL, -- e.g., "2026-2027"
--     is_current BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- Academic Terms (Linked to Academic Year ID + Semester 1 or 2)
CREATE TABLE IF NOT EXISTS academic_terms (
    id SERIAL PRIMARY KEY,
    academic_year_id INT NOT NULL REFERENCES academic_years(id),
    semester INT NOT NULL,               -- 1 or 2
    label VARCHAR(20) NOT NULL,                             -- e.g., 'S1', 'S2', 'S3', 'S4'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_academic_term UNIQUE (academic_year_id, semester)
    );

-- Student Groups per Level (Bound directly to Academic Year ID, stable across S1 & S2)
CREATE TABLE IF NOT EXISTS student_groups (
    id SERIAL PRIMARY KEY,
    academic_year_id INT NOT NULL REFERENCES academic_years(id),
    level_id INT NOT NULL REFERENCES academic_levels(id),
    code VARCHAR(10) NOT NULL,          -- e.g., "Group 01"
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_group_per_year_level UNIQUE (academic_year_id, level_id, code)
    );

-- ==========================================
-- 3. COURSE OFFERINGS & TARGET RELATIONSHIPS
-- ==========================================

-- Academic Course Offerings (Inputs tied to specific Term S1/S2)
CREATE TABLE IF NOT EXISTS course_offerings (
    id SERIAL PRIMARY KEY,
    session_uuid UUID NOT NULL, -- Shared across rows that belong to the SAME physical session
    term_id INT NOT NULL REFERENCES academic_terms(id),
    subject_id INT NOT NULL REFERENCES subjects(id),
    professor_id INT NOT NULL REFERENCES professors(id),
    required_room_type_id INT NOT NULL REFERENCES room_types(id),
    student_group_id INT NOT NULL REFERENCES student_groups(id),
    weekly_sessions INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Junction Table mapping Offerings to Targets (Level for Lectures, Group for TD/TP)
-- CREATE TABLE IF NOT EXISTS course_offering_targets (
--     id SERIAL PRIMARY KEY,
--     offering_id INT NOT NULL REFERENCES course_offerings(id),
--     group_id INT REFERENCES student_groups(id),
--     CONSTRAINT unique_group_per_offering UNIQUE (group_id, offering_id)
--     );

-- ==========================================
-- 4. TIMETABLE GRID ASSIGNMENTS (OUTPUTS)
-- ==========================================

-- Final Timetable Grid (Uses abstract integer slots 1..6 and days 1..5)
CREATE TABLE IF NOT EXISTS timetables (
    id SERIAL PRIMARY KEY,
    term_id INT NOT NULL REFERENCES academic_terms(id),
    offering_id INT NOT NULL REFERENCES course_offerings(id),
    room_id INT NOT NULL REFERENCES rooms(id),

    day_number INT NOT NULL,               -- Abstract day index (1 to 5)
    slot_number INT NOT NULL,              -- Abstract slot index (1 to 6)

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_room_schedule UNIQUE (term_id, room_id, day_number, slot_number)
    );


-- ==========================================
-- 5. FAST SEARCH INDEXES (INDEXING FOR HIGH PERFORMANCE)
-- ==========================================

-- Fast lookup for Terms by Academic Year
CREATE INDEX IF NOT EXISTS idx_academic_terms_year ON academic_terms(academic_year_id);

-- Fast lookup for Student Groups by Year and Specialization
CREATE INDEX IF NOT EXISTS idx_student_groups_year ON student_groups(academic_year_id);
CREATE INDEX IF NOT EXISTS idx_student_groups_year_level ON student_groups(academic_year_id, level_id);

-- Fast lookup for Course Offerings by Semester Term & Professor
CREATE INDEX IF NOT EXISTS idx_course_offerings_term ON course_offerings(term_id);
CREATE INDEX IF NOT EXISTS idx_course_offerings_prof ON course_offerings(professor_id);
CREATE INDEX IF NOT EXISTS idx_course_offerings_session_uuid ON course_offerings(session_uuid);

-- Fast lookup for Timetable grid queries & Conflict validation
CREATE INDEX IF NOT EXISTS idx_timetables_term ON timetables(term_id);
CREATE INDEX IF NOT EXISTS idx_timetables_grid ON timetables(term_id, day_number, slot_number);
CREATE INDEX IF NOT EXISTS idx_timetables_offering ON timetables(offering_id);

COMMIT;