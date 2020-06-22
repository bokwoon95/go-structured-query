-- Create an application, given user details arg_displayname and arg_email
-- A user will be created first if he doesn't already exist
-- An applicant will be created first if he doesn't already exist
-- This function is idempotent, making it safe to repeatedly call it on the same user
DROP FUNCTION IF EXISTS app.idempotent_create_application(TEXT, TEXT);
CREATE OR REPLACE FUNCTION app.idempotent_create_application (arg_displayname TEXT, arg_email TEXT)
RETURNS TABLE (_user_id INT, _user_role_id INT, _application_id INT, _magicstring TEXT) AS $$ DECLARE
    var_cohort TEXT;
    var_user_id INT;
    var_user_role_id INT;
    var_application_id INT;
    var_form_id INT;

    var_number_of_applicants INT;
    var_magicstring TEXT;
BEGIN
    SELECT DATE_PART('year', CURRENT_DATE)::TEXT INTO var_cohort;

    -- Get user id
    SELECT users.user_id INTO var_user_id FROM users WHERE users.email = arg_email;
    -- If user doesn't exist, create new user
    IF var_user_id IS NULL THEN
        INSERT INTO users (displayname, email)
        VALUES (arg_displayname, arg_email)
        RETURNING users.user_id INTO var_user_id
        ;
    END IF;

    -- Get user role id for applicant
    SELECT ur.user_role_id
    INTO var_user_role_id
    FROM user_roles AS ur
    WHERE ur.cohort = var_cohort AND ur.user_id = var_user_id AND ur.role = 'applicant'
    ;
    -- If user role id doesn't exist, create new user role id
    IF var_user_role_id IS NULL THEN
        INSERT INTO user_roles (cohort, user_id, role)
        VALUES (var_cohort, var_user_id, 'applicant')
        RETURNING user_roles.user_role_id INTO var_user_role_id
        ;
    END IF;

    -- If application doesn't exist, create application and associate it with applicant
    SELECT application_id INTO var_application_id FROM user_roles_applicants AS ura WHERE ura.user_role_id = var_user_role_id;
    IF var_application_id IS NULL THEN
        -- Check if there are any deleted applications (of which the user is a creator of) to reuse first
        SELECT apn.application_id
        INTO var_application_id
        FROM applications AS apn
        WHERE
            apn.creator_user_role_id = var_user_role_id
            AND apn.status = 'deleted'
            AND apn.deleted_at IS NOT NULL
        LIMIT 1
        ;

        -- If no eligible applications to reuse, then create a new application
        IF var_application_id IS NULL THEN
            SELECT
                f.form_id
            INTO
                var_form_id
            FROM
                forms AS f
                JOIN periods AS p ON p.period_id = f.period_id
            WHERE
                p.cohort = var_cohort
                AND p.stage = 'application'
                AND p.milestone = ''
                AND f.name = ''
                AND f.subsection = 'application'
            ;
            IF var_form_id IS NULL THEN
                RAISE EXCEPTION 'Application form {cohort:%, stage:application, milestone:, name:, subsection:application} not yet created',
                var_cohort USING ERRCODE = 'OLAJX'
                ;
            END IF;
            INSERT INTO applications (application_form_id)
            VALUES (var_form_id)
            RETURNING applications.application_id INTO var_application_id
            ;
            UPDATE applications AS apn SET creator_user_role_id = var_user_role_id WHERE apn.application_id = var_application_id;
        END IF;

        -- Create a new entry in user_roles_applicants
        SELECT
            form_id
        INTO
            var_form_id
        FROM
            forms AS f
            JOIN periods AS p ON p.period_id = f.period_id
        WHERE
            p.cohort = var_cohort
            AND p.stage = 'application'
            AND p.milestone = ''
            AND f.name = ''
            AND f.subsection = 'applicant'
        ;
        INSERT INTO user_roles_applicants (user_role_id, application_id, applicant_form_id)
        VALUES (var_user_role_id, var_application_id, var_form_id)
        ON CONFLICT (user_role_id) DO UPDATE
        SET application_id = var_application_id
        ;
        UPDATE applications AS apn SET status = 'pending', deleted_at = NULL WHERE apn.application_id = var_application_id;
    END IF;

    -- If application has 1 applicant, ensure magicstring is present
    -- If application has 2 applicants, set magicstring to NULL as it is no longer needed
    SELECT COUNT(*) INTO var_number_of_applicants FROM user_roles_applicants AS ura WHERE ura.application_id = var_application_id;
    IF var_number_of_applicants = 1 THEN
        SELECT apn.magicstring INTO var_magicstring FROM applications AS apn WHERE apn.application_id = var_application_id;
        IF var_magicstring IS NULL THEN
            SELECT * INTO var_magicstring FROM translate(gen_random_uuid()::TEXT, '-', '');
            UPDATE applications AS apn SET magicstring = var_magicstring WHERE apn.application_id = var_application_id;
        END IF;
    ELSIF var_number_of_applicants = 2 THEN
        UPDATE applications AS apn SET magicstring = NULL WHERE apn.application_id = var_application_id;
    END IF;

    RETURN QUERY SELECT var_user_id, var_user_role_id, var_application_id, var_magicstring;
END $$ LANGUAGE plpgsql;
