-- Make an applicant join an application associated with arg_magicstring
-- A user and applicant will be created first if he doesn't already exist
-- If no application with arg_magicstring exists, an exception will be raised
DROP FUNCTION IF EXISTS app.join_application(TEXT, TEXT, TEXT);
CREATE OR REPLACE FUNCTION app.join_application (arg_displayname TEXT, arg_email TEXT, arg_magicstring TEXT)
RETURNS TABLE (_user_id INT, _user_role_id INT, _application_id INT) AS $$ DECLARE
    var_cohort TEXT;
    var_user_role_id INT;
    var_user_id INT;
    var_form_id INT;
    var_old_application_id INT;
    var_new_application_id INT;
    var_magicstring TEXT;

    var_number_of_applicants INT;
    var_rowcount INT;
BEGIN
    SELECT DATE_PART('year', CURRENT_DATE)::TEXT INTO var_cohort;

    -- Get user
    SELECT user_id INTO var_user_id FROM users WHERE email = arg_email;
    -- If user doesn't exist, create new user
    IF var_user_id IS NULL THEN
        INSERT INTO users (displayname, email) VALUES (arg_displayname, arg_email) RETURNING users.user_id INTO var_user_id;
    END IF;
    RAISE DEBUG '{user_id:%}', var_user_id;

    -- Get applicant
    SELECT user_role_id INTO var_user_role_id FROM user_roles WHERE cohort = var_cohort AND user_id = var_user_id AND role = 'applicant';
    -- If applicant doesn't exist, create new applicant
    IF var_user_role_id IS NULL THEN
        INSERT INTO user_roles (cohort, user_id, role) VALUES (var_cohort, var_user_id, 'applicant') RETURNING user_roles.user_role_id INTO var_user_role_id;
    END IF;
    RAISE DEBUG '{user_role_id:%}', var_user_role_id;

    -- Get applicant's current application id if any
    SELECT application_id INTO var_old_application_id FROM user_roles_applicants WHERE user_role_id = var_user_role_id;
    RAISE DEBUG '{var_old_application_id:%}', var_old_application_id;

    -- Get application id for given arg_magicstring. If application doesn't exist, raise exception.
    SELECT application_id INTO var_new_application_id FROM applications WHERE magicstring = arg_magicstring;
    IF var_new_application_id IS NULL THEN
        RAISE EXCEPTION '{arg_magicstring:%} is not associated with any application.', arg_magicstring
        USING ERRCODE = 'OC8UM'
        ;
    END IF;
    RAISE DEBUG '{var_new_application_id:%}', var_new_application_id;

    -- Ensure applicant isn't trying to join his own application, else raise exception
    IF var_old_application_id = var_new_application_id THEN
        RAISE EXCEPTION '{applicant user_role_id:%} tried joining his own application application_id[%] magicstring[%]',
        var_user_role_id, var_new_application_id, arg_magicstring
        USING ERRCODE = 'OC8JK'
        ;
    END IF;

    -- Ensure var_new_application_id isn't already full else raise exception
    SELECT COUNT(*) INTO var_number_of_applicants FROM user_roles_applicants WHERE application_id = var_new_application_id;
    IF var_number_of_applicants >= 2 THEN
        RAISE EXCEPTION 'application {application_id:%} is already full', var_new_application_id
        USING ERRCODE = 'OC8FB'
        ;
    END IF;

    -- Upsert application for applicant
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
    VALUES (var_user_role_id, var_new_application_id, var_form_id)
    ON CONFLICT (user_role_id) DO UPDATE
    SET application_id = var_new_application_id
    ;
    GET DIAGNOSTICS var_rowcount = ROW_COUNT;
    RAISE DEBUG 'upserted % rows in user_roles_applicants', var_rowcount;

    -- Delete magicstring from application
    UPDATE applications SET magicstring = NULL WHERE application_id = var_new_application_id;
    GET DIAGNOSTICS var_rowcount = ROW_COUNT;
    RAISE DEBUG 'deleted magicstring from % rows in applications', var_rowcount;

    -- Delete all applications that the current applicant is the creator of, provided no one else is in the application
    UPDATE
        applications AS apn
    SET
        status = 'deleted'
        ,magicstring = NULL
        ,deleted_at = NOW()
    WHERE
        apn.creator_user_role_id = var_user_role_id
        AND apn.application_id <> var_new_application_id
        AND (SELECT COUNT(*) FROM user_roles_applicants AS ura WHERE ura.application_id = apn.application_id) = 0
    ;
    GET DIAGNOSTICS var_rowcount = ROW_COUNT;
    RAISE DEBUG 'deleted % rows in applications of which applicant {user_role_id:%} is the creator of', var_rowcount, var_user_role_id;

    -- If applicant's var_old_application_id still has another applicant in it, generate a new magicstring for that application
    SELECT COUNT(*) INTO var_number_of_applicants FROM user_roles_applicants WHERE application_id = var_old_application_id;
    IF var_number_of_applicants > 0 THEN
        -- First check if magicstring is NULL before generating a new one
        SELECT magicstring INTO var_magicstring FROM applications WHERE application_id = var_old_application_id;
        IF var_magicstring IS NULL THEN
            SELECT * INTO var_magicstring FROM translate(gen_random_uuid()::TEXT, '-', '');
            UPDATE applications AS apn SET magicstring = var_magicstring WHERE apn.application_id = var_old_application_id;
            GET DIAGNOSTICS var_rowcount = ROW_COUNT;
            RAISE DEBUG 'updated % rows in applications', var_rowcount;
        END IF;
    END IF;

    RETURN QUERY SELECT var_user_id AS _user_id, var_user_role_id AS _user_role_id, var_new_application_id AS _application_id;
END $$ LANGUAGE plpgsql;
