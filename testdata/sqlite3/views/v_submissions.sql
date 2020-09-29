DROP VIEW IF EXISTS v_submissions;
CREATE VIEW v_submissions AS
SELECT
    p.cohort
    ,p.milestone
    ,t.team_id
    ,t.team_name
    ,t.project_level
    ,s.submission_form_id
    ,s.submission_id
    ,f.questions
    ,s.submission_data AS answers
    ,p.start_at
    ,p.end_at
    ,s.submitted
    ,s.updated_at
    ,s.override_open
FROM
    teams AS t
    JOIN periods AS p ON p.cohort = t.cohort
    JOIN forms AS f ON f.period_id = p.period_id
    LEFT JOIN submissions AS s ON s.submission_form_id = f.form_id AND s.team_id = t.team_id
;
