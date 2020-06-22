DROP FUNCTION IF EXISTS app.upsert_submission(TEXT, TEXT, INT, JSONB);
CREATE OR REPLACE FUNCTION app.upsert_submission (arg_cohort TEXT, arg_milestone TEXT, arg_user_id INT, arg_submission_data JSONB)
RETURNS TABLE (_submission_id INT) AS $$ DECLARE
    var_cohort TEXT;
    var_team_id INT;
    var_submission_form_id INT;
    var_submission_id INT;

    var_rowcount INT;
BEGIN
    IF arg_cohort IS NULL OR arg_cohort = '' THEN
        SELECT DATE_PART('year', CURRENT_DATE)::TEXT INTO var_cohort;
    ELSIF NOT EXISTS(SELECT 1 FROM cohort_enum WHERE cohort = arg_cohort) THEN
        RAISE EXCEPTION '{cohort:%} is invalid', arg_cohort USING ERRCODE = 'OLALE';
    END IF;

    IF NOT EXISTS(SELECT 1 FROM milestone_enum WHERE milestone = arg_milestone) THEN
        RAISE EXCEPTION '{milestone:%} is invalid', arg_milestone USING ERRCODE = 'OLADN';
    END IF;

    SELECT urs.team_id
    INTO var_team_id
    FROM user_roles AS ur JOIN user_roles_students AS urs ON urs.user_role_id = ur.user_role_id
    WHERE ur.user_id = arg_user_id AND ur.role = 'student'
    ;
    GET DIAGNOSTICS var_rowcount = ROW_COUNT;
    IF var_rowcount = 0 THEN
        IF NOT EXISTS(SELECT 1 FROM users WHERE user_id = arg_user_id) THEN
            RAISE EXCEPTION 'User{user_id:%} does not exist', arg_user_id USING ERRCODE = 'OLAMC';
        ELSE
            RAISE EXCEPTION 'User{user_id:%} is not a student', arg_user_id USING ERRCODE = 'ONXIU';
        END IF;
    END IF;
    IF var_team_id IS NULL THEN
        RAISE EXCEPTION 'Student{user_id:%} does not have a team', arg_user_id USING ERRCODE = 'ONXDI';
    END IF;

    SELECT form_id
    INTO var_submission_form_id
    FROM forms AS f JOIN periods AS p ON p.period_id = f.period_id
    WHERE p.cohort = var_cohort AND p.stage = 'submission' AND p.milestone = arg_milestone AND f.name = '' AND f.subsection = ''
    ;
    IF var_submission_form_id IS NULL THEN
        RAISE EXCEPTION 'Submission form {cohort:%, stage:submission, milestone:%, name:, subsection:} not found',
        var_cohort, arg_milestone USING ERRCODE = 'OQ7WT'
        ;
    END IF;


    SELECT submission_id
    INTO var_submission_id
    FROM submissions
    WHERE team_id = var_team_id AND submission_form_id = var_submission_form_id
    ;
    IF var_submission_id IS NULL THEN
        INSERT INTO submissions (team_id, submission_data, submission_form_id) VALUES (var_team_id, arg_submission_data, var_submission_form_id) RETURNING submission_id INTO var_submission_id;
    ELSE
        UPDATE submissions SET submission_data = arg_submission_data, submission_form_id = var_submission_form_id WHERE submission_id = var_submission_id;
    END IF;

    RETURN QUERY SELECT var_submission_id;
END $$ LANGUAGE plpgsql;
