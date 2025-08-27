-- Create custom schema
CREATE SCHEMA IF NOT EXISTS cargo_tracker AUTHORIZATION cargo_tracker_role;

-- Set default search_path for the user
ALTER ROLE cargo_tracker_role SET search_path TO cargo_tracker;
