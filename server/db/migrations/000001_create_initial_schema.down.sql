-- Migration Rollback: Tear Down Initial Schema
-- Description: Safely drops all tables in reverse order of dependency.

BEGIN;

-- Drop Indexes
DROP INDEX IF EXISTS idx_timetables_offering;
DROP INDEX IF EXISTS idx_timetables_grid;
DROP INDEX IF EXISTS idx_timetables_term;
DROP INDEX IF EXISTS idx_course_offerings_prof;
DROP INDEX IF EXISTS idx_course_offerings_term;
DROP INDEX IF EXISTS idx_course_offerings_session_uuid;
DROP INDEX IF EXISTS idx_student_groups_year;
DROP INDEX IF EXISTS idx_student_groups_year_level;
DROP INDEX IF EXISTS idx_academic_terms_year;

-- Drop Tables
DROP TABLE IF EXISTS timetables CASCADE;
DROP TABLE IF EXISTS course_offering_targets CASCADE;
DROP TABLE IF EXISTS course_offerings CASCADE;
DROP TABLE IF EXISTS student_groups CASCADE;
DROP TABLE IF EXISTS academic_terms CASCADE;
DROP TABLE IF EXISTS academic_years CASCADE;
DROP TABLE IF EXISTS academic_levels CASCADE;
DROP TABLE IF EXISTS subjects CASCADE;
DROP TABLE IF EXISTS rooms CASCADE;
DROP TABLE IF EXISTS room_types CASCADE;
DROP TABLE IF EXISTS professors CASCADE;
DROP TABLE IF EXISTS departments CASCADE;

COMMIT;
