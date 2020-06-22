-- Accept the application referenced by arg_application_id
-- A new team will always be created
-- Students 1 and 2 will be upserted (they might have been accpeted before, but deleted by admin)
-- If the application does not have a team name to create the team with, use arg_default_name
DROP FUNCTION IF EXISTS app.accept_application(INT, TEXT);
CREATE OR REPLACE FUNCTION app.accept_application (arg_application_id INT, arg_default_name TEXT)
RETURNS TABLE (_team_id INT, _student_user_id_1 INT, _student_user_id_2 INT) AS $$ DECLARE
    var_deleted_at TIMESTAMPTZ;

    var_cohort TEXT;
    var_status TEXT;
    var_project_level TEXT;
    var_application_data JSONB;

    var_user_id_1 INT;
    var_user_id_2 INT;
    var_applicant_data_1 JSONB;
    var_applicant_data_2 JSONB;

    var_team_id INT;
    var_team_name TEXT;
    var_student_user_role_id_1 INT;
    var_student_user_role_id_2 INT;
BEGIN
    -- Get application details
    SELECT
        team_id
        ,team_name
        ,cohort
        ,status
        ,project_level
        ,deleted_at
        ,application_data
    INTO
        var_team_id
        ,var_team_name
        ,var_cohort
        ,var_status
        ,var_project_level
        ,var_deleted_at
        ,var_application_data
    FROM
        applications
    WHERE
        application_id = arg_application_id
    ;
    RAISE DEBUG 'Application {application_id:%, tid:%, cohort:%, status:%, project_level:%, data:%}',
    arg_application_id, var_team_id, var_cohort, var_status, var_project_level, var_application_data::TEXT
    ;

    -- If application doesn't exist, raise exception
    IF var_cohort IS NULL THEN
        RAISE EXCEPTION 'Tried accepting a non existent application{application_id:%}', arg_application_id
        USING ERRCODE = 'OC8U9'
        ;
    END IF;

    -- If application is deleted, raise exception
    IF var_deleted_at IS NOT NULL THEN
        RAISE EXCEPTION 'Tried accepting an already deleted application{application_id:%}', arg_application_id
        USING ERRCODE = 'OC8W6'
        ;
    END IF;

    -- Get applicant 1 details
    SELECT
        ur.user_id
        ,ura.applicant_data
    INTO
        var_user_id_1
        ,var_applicant_data_1
    FROM
        user_roles AS ur
        JOIN user_roles_applicants AS ura USING (user_role_id)
    WHERE
        ura.application_id = arg_application_id
    ORDER BY
        ura.user_role_id
    LIMIT 1
    ;
    RAISE DEBUG 'Applicant1 {uid:%} data:%', var_user_id_1, var_applicant_data_1::TEXT;

    -- If applicant1 doesn't exist, raise exception
    IF var_user_id_1 IS NULL THEN
        RAISE EXCEPTION 'Tried accepting incomplete application{application_id:%}, missing applicant1', arg_application_id
        USING ERRCODE = 'OC8KH'
        ;
    END IF;

    -- Get applicant 2 details
    SELECT
        ur.user_id
        ,ura.applicant_data
    INTO
        var_user_id_2
        ,var_applicant_data_2
    FROM
        user_roles AS ur
        JOIN user_roles_applicants AS ura USING (user_role_id)
    WHERE
        ura.application_id = arg_application_id
    ORDER BY
        ura.user_role_id
    LIMIT 1
    OFFSET 1
    ;
    RAISE DEBUG 'Applicant2 {uid:%} data:%', var_user_id_2, var_applicant_data_2::TEXT;

    -- If applicant2 doesn't exist, raise exception
    IF var_user_id_2 IS NULL THEN
        RAISE EXCEPTION 'Tried accepting incomplete application{application_id:%}, missing applicant2', arg_application_id
        USING ERRCODE = 'OC8KH'
        ;
    END IF;

    -- Upsert student1
    INSERT INTO user_roles (user_id, cohort, role)
    VALUES (var_user_id_1, var_cohort, 'student')
    ON CONFLICT (user_id, cohort, role) DO UPDATE
    SET updated_at = NOW(), deleted_at = NULL
    RETURNING user_roles.user_role_id INTO var_student_user_role_id_1
    ;
    INSERT INTO user_roles_students (user_role_id, student_data)
    VALUES (var_student_user_role_id_1, var_applicant_data_1)
    ON CONFLICT (user_role_id) DO UPDATE
    SET student_data = var_applicant_data_1
    ;

    -- Upsert student2
    INSERT INTO user_roles (user_id, cohort, role)
    VALUES (var_user_id_2, var_cohort, 'student')
    ON CONFLICT (user_id, cohort, role) DO UPDATE
    SET updated_at = NOW(), deleted_at = NULL
    RETURNING user_roles.user_role_id INTO var_student_user_role_id_2
    ;
    INSERT INTO user_roles_students (user_role_id, student_data)
    VALUES (var_student_user_role_id_2, var_applicant_data_2)
    ON CONFLICT (user_role_id) DO UPDATE
    SET student_data = var_applicant_data_2
    ;

    -- Create new team if not exists, and update application's team accordingly
    IF var_team_id IS NULL THEN
        INSERT INTO teams (cohort, team_name, project_level, team_data)
        VALUES (
            var_cohort
            ,COALESCE(var_team_name, arg_default_name, TRANSLATE(gen_random_uuid()::TEXT, '-', ''))
            ,var_project_level
            ,var_application_data
        )
        RETURNING teams.team_id INTO var_team_id
        ;
        UPDATE applications SET team_id = var_team_id WHERE application_id = arg_application_id;
    END IF;
    RAISE DEBUG '{var_team_id:%}', var_team_id;

    -- Update team for student1 and student2
    UPDATE teams SET deleted_at = NULL WHERE team_id = var_team_id;
    UPDATE user_roles_students SET team_id = var_team_id WHERE user_role_id IN (var_student_user_role_id_1, var_student_user_role_id_2);

    -- Update application for arg_application_id
    UPDATE applications SET status = 'accepted', deleted_at = NOW() WHERE application_id = arg_application_id;

    RETURN QUERY SELECT var_team_id, var_user_id_1, var_user_id_2;
END $$ LANGUAGE plpgsql;
