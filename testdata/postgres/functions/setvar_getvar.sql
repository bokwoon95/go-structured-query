-- setvar and getvar are complementary functions that emulate variables in pure
-- SQL without needing to drop into a procedural language like PL/pgSQL.
-- http://okbob.blogspot.com/2019/05/how-to-dont-emulate-schema-global.html
-- https://wiki.postgresql.org/wiki/Variable_Design#With_Secure_Session_Variables

-- setvar sets the result of a query into a variable. Example usage:
-- SELECT setvar('cohort', (SELECT cohort FROM cohort_enum ORDER BY cohort DESC LIMIT 1));
DROP FUNCTION IF EXISTS setvar(TEXT, TEXT);
CREATE OR REPLACE FUNCTION setvar(arg_name TEXT, arg_value TEXT)
RETURNS VOID AS $$ DECLARE
BEGIN
    PERFORM set_config('var.' || arg_name, arg_value, FALSE);
END $$ LANGUAGE plpgsql;

-- getvar gets the value of a variable set by setvar. Example usage:
-- SELECT * FROM user_roles WHERE cohort = getvar('cohort');
DROP FUNCTION IF EXISTS getvar(TEXT);
CREATE OR REPLACE FUNCTION getvar(arg_name TEXT)
RETURNS TEXT AS $$ DECLARE
BEGIN
    RETURN current_setting('var.' || arg_name, TRUE);
END $$ LANGUAGE plpgsql;
