DROP FUNCTION IF EXISTS app.upsert_application_data(INT, JSONB, JSONB);
CREATE OR REPLACE FUNCTION app.upsert_application_data (arg_user_role_id INT, arg_applicant_data JSONB, arg_application_data JSONB)
RETURNS TABLE (_application_id INT) AS $$ DECLARE
    var_cohort TEXT;
    var_application_id INT;
    var_applicant_form_id INT;
    var_application_form_id INT;
BEGIN
    -- Ensure arg_user_role_id is a valid applicant
    IF NOT EXISTS(SELECT 1 FROM user_roles WHERE user_role_id = arg_user_role_id AND role = 'applicant') THEN
        RAISE EXCEPTION 'User{user_role_id:%} is not an applicant', arg_user_role_id USING ERRCODE = 'OC8FY';
    END IF;

    SELECT cohort INTO var_cohort FROM user_roles WHERE user_role_id = arg_user_role_id;

    IF EXISTS(SELECT 1 FROM user_roles_applicants WHERE user_role_id = arg_user_role_id) THEN
        UPDATE user_roles_applicants SET applicant_data = arg_applicant_data WHERE user_role_id = arg_user_role_id;
    ELSE
        -- Get form_id for the current cohort's applicant form
        SELECT form_id
        INTO var_applicant_form_id
        FROM forms AS f JOIN periods AS p ON p.period_id = f.period_id
        WHERE p.cohort = var_cohort AND p.stage = 'application' AND f.subsection = 'applicant'
        ;
        RAISE DEBUG 'var_applicant_form_id[%]', var_applicant_form_id;
        -- If form_id for the current cohort's applicant form is not found, raise exception
        IF var_applicant_form_id IS NULL THEN
            RAISE EXCEPTION 'applicant form_id not found' USING ERRCODE = 'OC8BK';
        END IF;
        -- Create new user_roles_applicants entry
        INSERT INTO user_roles_applicants(user_role_id, applicant_form_id, applicant_data)
        VALUES (arg_user_role_id, var_applicant_form_id, arg_applicant_data)
        ;
    END IF;

    SELECT application_id INTO var_application_id FROM user_roles_applicants WHERE user_role_id = arg_user_role_id;
    IF var_application_id IS NOT NULL THEN
        UPDATE applications SET application_data = arg_application_data WHERE application_id = var_application_id;
        RAISE NOTICE 'UPDATE applications SET application_data = arg_application_data WHERE application_id = %;', var_application_id;
    ELSE
        -- Get form_id for the current cohort's application form
        SELECT form_id
        INTO var_application_form_id
        FROM forms AS f JOIN periods AS p ON p.period_id = f.period_id
        WHERE p.cohort = var_cohort AND p.stage = 'application' AND f.subsection = 'application'
        ;
        RAISE DEBUG 'var_application_form_id[%]', var_application_form_id;
        -- If form_id for the current cohort's application form is not found, raise exception
        IF var_application_form_id IS NULL THEN
            RAISE EXCEPTION 'application form_id not found' USING ERRCODE = 'OC8BK';
        END IF;
        -- Create new application entry
        INSERT INTO applications(creator_user_role_id, application_form_id, application_data)
        VALUES (arg_user_role_id, var_application_form_id, arg_application_data)
        RETURNING application_id INTO var_application_id
        ;
        RAISE NOTICE 'INSERT {var_application_id:%}', var_application_id;
    END IF;

    RETURN QUERY SELECT var_application_id;
END $$ LANGUAGE plpgsql;
