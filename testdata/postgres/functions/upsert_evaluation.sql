DROP FUNCTION IF EXISTS app.upsert_evaluation(TEXT, INT, INT, JSONB);
CREATE OR REPLACE FUNCTION app.upsert_evaluation (arg_milestone TEXT, arg_evaluator_team_id INT, arg_evaluatee_submission_id INT, arg_evaluation_data JSONB)
RETURNS TABLE (_team_evaluation_id INT) AS $$ DECLARE
    var_cohort TEXT;
    var_form_id INT;
    var_submission_id INT;
    var_team_evaluation_id INT;
BEGIN
    SELECT DATE_PART('year', CURRENT_DATE)::TEXT INTO var_cohort;

    IF NOT EXISTS(SELECT 1 FROM milestone_enum WHERE milestone = arg_milestone) THEN
        RAISE EXCEPTION '{milestone:%} is invalid', arg_milestone USING ERRCODE = 'OLADN';
    END IF;

    IF NOT EXISTS(SELECT 1 FROM teams WHERE team_id = arg_evaluator_team_id) THEN
        RAISE EXCEPTION 'evaluator{team_id:%} does not exist', arg_evaluator_team_id USING ERRCODE = 'OWHZT';
    END IF;

    IF NOT EXISTS(SELECT 1 FROM teams WHERE team_id = arg_evaluatee_submission_id) THEN
        RAISE EXCEPTION 'evaluatee{team_id:%} does not exist', arg_evaluatee_submission_id USING ERRCODE = 'OWHZT';
    END IF;

    SELECT
        submission_id
    INTO
        var_submission_id
    FROM
        submissions AS s
        JOIN forms AS f ON f.form_id = s.submission_form_id
        JOIN periods AS p ON p.period_id = f.period_id
    WHERE
        p.cohort = var_cohort
        AND p.stage = 'submission'
        AND p.milestone = arg_milestone
        AND f.name = ''
        AND f.subsection = ''
        AND s.submission_id = arg_evaluatee_submission_id
    ;
    IF var_submission_id IS NULL THEN
        RAISE EXCEPTION 'evaluatee{team_id:%} does not have a valid submission for {cohort:%, stage:submission, milestone:%, name:, subsection:}',
        arg_evaluatee_submission_id, var_cohort, arg_milestone
        ;
    END IF;

    SELECT form_id
    INTO var_form_id
    FROM forms AS f JOIN periods AS p ON p.period_id = f.period_id
    WHERE p.cohort = var_cohort AND p.stage = 'evaluation' AND p.milestone = arg_milestone AND f.name = '' AND f.subsection = ''
    ;
    IF var_form_id IS NULL THEN
        RAISE EXCEPTION 'Evaluation form {cohort:%, stage:evaluation, milestone:%, name:, subsection:} not found', var_cohort, arg_milestone;
    END IF;

    SELECT team_evaluation_id
    INTO var_team_evaluation_id
    FROM team_evaluations
    WHERE evaluator_team_id = arg_evaluator_team_id AND evaluatee_submission_id = arg_evaluatee_submission_id
    ;
    RAISE NOTICE '{var_team_evaluation_id:%}', var_team_evaluation_id;
    IF var_team_evaluation_id IS NULL IS NULL THEN
        INSERT INTO team_evaluations (evaluator_team_id, evaluatee_submission_id, evaluation_data, evaluation_form_id)
        VALUES (arg_evaluator_team_id, var_submission_id, arg_evaluation_data, var_form_id)
        RETURNING team_evaluation_id INTO var_team_evaluation_id
        ;
    ELSE
        UPDATE team_evaluations SET evaluation_data = arg_evaluation_data WHERE team_evaluation_id = var_team_evaluation_id;
    END IF;

    RETURN QUERY SELECT var_team_evaluation_id;
END $$ LANGUAGE plpgsql;
