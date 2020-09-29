DROP VIEW IF EXISTS v_team_evaluations;
CREATE VIEW v_team_evaluations AS
WITH pairs AS (
    SELECT
        t1.cohort
        ,tep.evaluatee_team_id
        ,t1.team_name AS evaluatee_team_name
        ,t1.project_level AS evaluatee_project_level
        ,tep.evaluator_team_id
        ,t2.team_name AS evaluator_team_name
        ,t2.project_level AS evaluator_project_level
    FROM
        team_evaluation_pairs AS tep
        JOIN teams AS t1 ON t1.team_id = tep.evaluatee_team_id
        JOIN teams AS t2 ON t2.team_id = tep.evaluator_team_id AND t2.cohort = t1.cohort
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
    ,s.submission_form_id
    ,sq.questions AS submission_questions
    ,s.submission_data AS submission_answers
    ,sq.start_at AS submission_start_at
    ,sq.end_at AS submission_end_at
    ,s.override_open AS submission_override_open
    ,s.submitted AS submission_submitted
    ,s.updated_at AS submission_updated_at

    -- Evaluation
    ,te.team_evaluation_id
    ,p.evaluator_team_id
    ,p.evaluator_team_name
    ,p.evaluator_project_level
    ,te.evaluation_form_id
    ,eq.questions AS evaluation_questions
    ,te.evaluation_data AS evaluation_answers
    ,eq.start_at AS evaluation_start_at
    ,eq.end_at AS evaluation_end_at
    ,te.override_open AS evaluation_override_open
    ,te.submitted AS evaluation_submitted
    ,te.updated_at AS evaluation_updated_at
FROM
    pairs AS p
    JOIN submission_questions AS sq ON sq.cohort = p.cohort
    LEFT JOIN submissions AS s ON s.submission_form_id = sq.form_id AND s.team_id = p.evaluatee_team_id
    JOIN evaluation_questions AS eq ON eq.cohort = sq.cohort AND eq.milestone = sq.milestone
    LEFT JOIN team_evaluations AS te ON te.evaluation_form_id = eq.form_id AND te.evaluator_team_id = p.evaluator_team_id AND te.evaluatee_submission_id = s.submission_id
;
