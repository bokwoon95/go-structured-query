DROP FUNCTION IF EXISTS app.upsert_form(INT, TEXT, TEXT, TEXT, TEXT, TEXT);
CREATE OR REPLACE FUNCTION app.upsert_form(arg_period_id INT, arg_cohort TEXT, arg_stage TEXT, arg_milestone TEXT, arg_name TEXT, arg_subsection TEXT)
RETURNS TABLE (_form_id INT) AS $$ DECLARE
    var_period_id INT;
    var_form_id INT;
BEGIN
    IF arg_period_id <> 0 AND arg_period_id IS NOT NULL THEN
        SELECT arg_period_id INTO var_period_id;
    ELSE
        SELECT period_id INTO var_period_id FROM periods WHERE cohort = arg_cohort AND stage = arg_stage AND milestone = arg_milestone;
        IF var_period_id IS NULL THEN
            INSERT INTO periods (cohort, stage, milestone)
            VALUES (arg_cohort, arg_stage, arg_milestone)
            RETURNING period_id INTO var_period_id
            ;
        END IF;
    END IF;

    SELECT form_id INTO var_form_id FROM forms WHERE period_id = var_period_id AND name = arg_name AND subsection = arg_subsection;
    IF var_form_id IS NULL THEN
        INSERT INTO forms (period_id, name, subsection)
        VALUES (var_period_id, arg_name, arg_subsection)
        RETURNING form_id INTO var_form_id
        ;
    END IF;

    RETURN QUERY SELECT var_form_id;
END $$ LANGUAGE plpgsql;
