DROP FUNCTION IF EXISTS app.create_user_role(TEXT, TEXT, TEXT, TEXT);
CREATE OR REPLACE FUNCTION app.create_user_role (arg_cohort TEXT, arg_role TEXT, arg_displayname TEXT, arg_email TEXT)
RETURNS TABLE (_user_id INT, _displayname TEXT, _email TEXT, _roles TEXT[], _already_created BOOLEAN) AS $$ DECLARE
    var_cohort TEXT;
    var_user_id INT;
    var_displayname TEXT;
    var_email TEXT;
    var_roles TEXT[];
    var_already_created BOOLEAN := FALSE;

    var_affected_rows INT;
BEGIN
    -- If arg_cohort is blank or NULL, use the current year as the cohort
    SELECT arg_cohort INTO var_cohort;
    IF var_cohort = '' OR var_cohort IS NULL THEN
        SELECT DATE_PART('year', CURRENT_DATE)::TEXT INTO var_cohort;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cohort_enum where cohort = var_cohort) THEN
        RAISE EXCEPTION 'cohort "%" is not a valid cohort', var_cohort USING ERRCODE = 'OLALE';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM role_enum WHERE role = arg_role) THEN
        RAISE EXCEPTION 'role "%" is not a valid role', arg_role USING ERRCODE = 'OLAZN';
    END IF;

    IF arg_email = '' OR arg_email IS NULL THEN
        RAISE EXCEPTION 'arg_email "%" cannot be empty', arg_email;
    END IF;

    -- Check if user already exists for the given arg_email. If not, create that user.
    SELECT user_id INTO var_user_id FROM users WHERE email = arg_email;
    IF var_user_id IS NULL THEN
        IF arg_displayname = '' OR arg_displayname IS NULL THEN
            RAISE EXCEPTION 'arg_displayname "%" cannot be empty', arg_displayname;
        END IF;
        INSERT INTO users (displayname, email) VALUES (arg_displayname, arg_email) RETURNING user_id INTO var_user_id;
    END IF;

    -- Create the role for the user
    INSERT INTO user_roles (user_id, cohort, role) VALUES (var_user_id, var_cohort, arg_role) ON CONFLICT DO NOTHING;
    GET DIAGNOSTICS var_affected_rows = ROW_COUNT;
    IF var_affected_rows = 0 THEN
        SELECT TRUE INTO var_already_created;
    END IF;

    SELECT ARRAY(SELECT role FROM user_roles WHERE user_id = var_user_id) INTO var_roles;

    RETURN QUERY SELECT var_user_id, var_displayname, var_email, var_roles, var_already_created;
END $$ LANGUAGE plpgsql;
