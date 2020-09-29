DROP VIEW IF EXISTS v_user_evaluations;
CREATE VIEW v_user_evaluations AS
WITH pairs AS (
    SELECT
        t.cohort
        ,t.team_id AS evaluatee_team_id
        ,t.team_name AS evaluatee_team_name
        ,t.project_level AS evaluatee_project_level
        ,u.user_id AS evaluator_user_id
        ,ur.user_role_id AS evaluator_user_role_id
        ,u.displayname AS evaluator_displayname
        ,CASE ur.user_role_id
            WHEN t.adviser_user_role_id THEN 'adviser'
            WHEN t.mentor_user_role_id THEN 'mentor'
            ELSE ''
        END AS evaluator_role
    FROM
        teams AS t
        JOIN user_roles AS ur ON ur.user_role_id IN (t.adviser_user_role_id, t.mentor_user_role_id)
        JOIN users AS u ON u.user_id = ur.user_id
)
,submission_questions AS (
    SELECT p.cohort, p.stage, p.milestone, p.start_at, p.end_at, f.questions, f.form_id
    FROM periods AS p JOIN forms AS f ON f.period_id = p.period_id
    WHERE p.cohort <> '' AND p.stage = 'submission' AND p.milestone <> '' AND f.name = '' AND f.subsection = ''
)
,evaluation_questions AS (
    SELECT p.cohort, p.stage, p.milestone, p.start_at, p.end_at, f.questions, f.form_id
    FROM periods AS p JOIN forms AS f ON f.period_id = p.period_id
    WHERE p.cohort <> '' AND p.stage = 'evaluation' AND p.milestone <> '' AND f.name = '' AND f.subsection = ''
)
SELECT
    eq.cohort
    ,eq.stage
    ,eq.milestone

    -- Submission
    ,s.submission_id
    ,p.evaluatee_team_id
    ,p.evaluatee_team_name
    ,p.evaluatee_project_level
    ,sq.form_id AS submission_form_id
    ,sq.questions AS submission_questions
    ,s.submission_data AS submission_answers
    ,sq.start_at AS submission_start_at
    ,sq.end_at AS submission_end_at
    ,s.override_open AS submission_override_open
    ,s.submitted AS submission_submitted
    ,s.updated_at AS submission_updated_at

    -- Evaluation
    ,ue.user_evaluation_id
    ,p.evaluator_user_id
    ,p.evaluator_user_role_id
    ,p.evaluator_displayname
    ,p.evaluator_role
    ,eq.form_id AS evaluation_form_id
    ,eq.questions AS evaluation_questions
    ,ue.evaluation_data AS evaluation_answers
    ,eq.start_at AS evaluation_start_at
    ,eq.end_at AS evaluation_end_at
    ,ue.override_open AS evaluation_override_open
    ,ue.submitted AS evaluation_submitted
    ,ue.updated_at AS evaluation_updated_at
FROM
    pairs AS p
    JOIN submission_questions AS sq ON sq.cohort = p.cohort
    LEFT JOIN submissions AS s ON s.submission_form_id = sq.form_id AND s.team_id = p.evaluatee_team_id
    JOIN evaluation_questions AS eq ON eq.cohort = sq.cohort AND eq.milestone = sq.milestone
    LEFT JOIN user_evaluations AS ue
        ON ue.evaluation_form_id = eq.form_id
        AND ue.evaluator_user_role_id = p.evaluator_user_role_id
        AND ue.evaluatee_submission_id = s.submission_id
;
