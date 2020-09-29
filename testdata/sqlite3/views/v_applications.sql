DROP VIEW IF EXISTS v_applications;
CREATE VIEW v_applications AS
WITH applicants AS (
    SELECT u.user_id, ur.user_role_id, u.displayname, u.email, ura.application_id, ura.applicant_data
    FROM users AS u JOIN user_roles AS ur USING (user_id) LEFT JOIN user_roles_applicants AS ura USING (user_role_id)
    WHERE ur.role = 'applicant'
)
,application_questions AS (
    SELECT p.cohort, p.milestone, p.start_at, p.end_at, f.questions, f.form_id
    FROM periods AS p JOIN forms AS f ON f.period_id = p.period_id
    WHERE p.cohort <> '' AND p.stage = 'application' AND p.milestone = '' AND f.name = '' AND f.subsection = 'application'
)
,applicant_questions AS (
    SELECT p.cohort, p.milestone, p.start_at, p.end_at, f.questions, f.form_id
    FROM periods AS p JOIN forms AS f ON f.period_id = p.period_id
    WHERE p.cohort <> '' AND p.stage = 'application' AND p.milestone = '' AND f.name = '' AND f.subsection = 'applicant'
)
SELECT
    -- Application
    applications.application_id
    ,applications.cohort
    ,applications.status
    ,applications.creator_user_role_id
    ,applications.project_level
    ,applications.magicstring
    ,applications.submitted

    -- Applicant 1
    ,applicant1.user_id AS applicant1_user_id
    ,applicant1.user_role_id AS applicant1_user_role_id
    ,applicant1.displayname AS applicant1_displayname
    ,applicant1.email AS applicant1_email

    -- Applicant 2
    ,applicant2.user_id AS applicant2_user_id
    ,applicant2.user_role_id AS applicant2_user_role_id
    ,applicant2.displayname AS applicant2_displayname
    ,applicant2.email AS applicant2_email

    -- Questions and Answers
    ,application_questions.form_id AS application_form_id
    ,applicant_questions.form_id AS applicant_form_id
    ,application_questions.questions AS application_questions
    ,applications.application_data AS application_answers
    ,applicant_questions.questions AS applicant_questions
    ,applicant1.applicant_data AS applicant1_answers
    ,applicant2.applicant_data AS applicant2_answers

    ,applications.created_at
    ,applications.updated_at
    ,applications.deleted_at
FROM
    applications
    JOIN application_questions ON application_questions.cohort = applications.cohort
    JOIN applicant_questions ON applicant_questions.cohort = applications.cohort
    JOIN applicants AS applicant1 ON applicant1.application_id = applications.application_id
    LEFT JOIN applicants AS applicant2 ON applicant2.application_id = applications.application_id AND applicant1.user_id < applicant2.user_id
WHERE
    (
        (SELECT COUNT(*) FROM user_roles_applicants AS ura WHERE ura.application_id = applications.application_id) = 2
        AND
        applicant2.user_id IS NOT NULL
    )
    OR
    (
        (SELECT COUNT(*) FROM user_roles_applicants AS ura WHERE ura.application_id = applications.application_id) = 1
        AND
        applicant2.user_id IS NULL
    )
;
