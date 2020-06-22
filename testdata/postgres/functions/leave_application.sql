-- Make an applicant associated with arg_user_id leave his current application
-- If applicant or his application doesn't exist, an exception will be raised
DROP FUNCTION IF EXISTS app.leave_application(INT);
CREATE OR REPLACE FUNCTION app.leave_application (arg_user_id INT)
RETURNS TABLE (_user_id INT, _user_role_id INT, _application_id INT, _magicstring TEXT) AS $$ DECLARE
    var_cohort TEXT;
    var_user_id INT;
    var_user_role_id INT;
    var_application_id INT;
    var_status TEXT;
    var_number_of_applicants INT;
    var_magicstring TEXT;

    var_rowcount INT;
BEGIN
    SELECT DATE_PART('year', CURRENT_DATE)::TEXT INTO var_cohort;

    -- Get applicant's current application id
    SELECT ur.user_role_id, ura.application_id
    INTO var_user_role_id, var_application_id
    FROM user_roles AS ur LEFT JOIN user_roles_applicants AS ura USING (user_role_id)
    WHERE ur.user_id = arg_user_id AND ur.cohort = var_cohort AND ur.role = 'applicant'
    ;
    RAISE DEBUG 'var_user_role_id[%], var_application_id[%]', var_user_role_id, var_application_id;

    -- If applicant doesn't exist, raise exception
    IF var_user_role_id IS NULL THEN
        RAISE EXCEPTION 'applicant doesnt exist'
        USING ERRCODE = 'OC8FY'
        ;
    END IF;

    -- If applicant doesn't have a current application, raise exception
    IF var_application_id IS NULL THEN
        RAISE EXCEPTION 'applicant doesnt have an application'
        USING ERRCODE = 'OC8EN'
        ;
    END IF;

    -- If application is already accepted, raise exception
    SELECT apn.status INTO var_status FROM applications AS apn WHERE apn.application_id = var_application_id;
    RAISE DEBUG 'var_status[%]', var_status;
    IF var_status = 'accepted' THEN
        RAISE EXCEPTION 'applicant user_role_id[%] tried leaving an application application_id[%] that was already accepted', var_user_role_id, var_application_id
        USING ERRCODE = 'OC8A4'
        ;
    END IF;

    -- Remove applicant from application
    UPDATE user_roles_applicants AS ura SET application_id = NULL WHERE ura.user_role_id = var_user_role_id;
    GET DIAGNOSTICS var_rowcount = ROW_COUNT;
    RAISE DEBUG 'updated % rows in user_roles_applicants', var_rowcount;

    -- Check if application is empty and set status and magicstring accordingly
    SELECT COUNT(*) INTO var_number_of_applicants FROM user_roles_applicants AS ura WHERE ura.application_id = var_application_id;
    RAISE DEBUG 'var_number_of_applicants[%]', var_number_of_applicants;
    IF var_number_of_applicants = 0 THEN
        UPDATE applications AS apn
        SET status = 'deleted', magicstring = NULL, deleted_at = NOW()
        WHERE apn.application_id = var_application_id
        ;
        GET DIAGNOSTICS var_rowcount = ROW_COUNT;
        RAISE DEBUG 'deleted % rows in applications', var_rowcount;
    ELSE
        -- Get application's current magic string. If NULL, insert it with a new magic string.
        SELECT apn.magicstring INTO var_magicstring FROM applications AS apn WHERE apn.application_id = var_application_id;
        IF var_magicstring IS NULL THEN
            SELECT * INTO var_magicstring FROM translate(gen_random_uuid()::TEXT, '-', '');
            UPDATE applications AS apn SET magicstring = var_magicstring WHERE apn.application_id = var_application_id;
            GET DIAGNOSTICS var_rowcount = ROW_COUNT;
            RAISE DEBUG 'updated % rows in applications', var_rowcount;
        END IF;
    END IF;

    RETURN QUERY SELECT var_user_id, var_user_role_id, var_application_id, var_magicstring;
END $$ LANGUAGE plpgsql;
