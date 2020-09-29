DROP VIEW IF EXISTS v_teams;
CREATE VIEW v_teams AS
WITH students AS (
    SELECT u.user_id, u.displayname, u.email, ur.user_role_id, urs.team_id, urs.student_data
    FROM users AS u JOIN user_roles AS ur USING (user_id) LEFT JOIN user_roles_students AS urs USING (user_role_id)
    WHERE ur.role = 'student'
)
,advisers AS (
    SELECT u.user_id, u.displayname, u.email, ur.user_role_id
    FROM users AS u JOIN user_roles AS ur USING (user_id)
    WHERE ur.role = 'adviser'
)
,mentors AS (
    SELECT u.user_id, u.displayname, u.email, ur.user_role_id
    FROM users AS u JOIN user_roles AS ur USING (user_id)
    WHERE ur.role = 'mentor'
)
SELECT
    t.team_id
    ,t.cohort
    ,t.team_name
    ,t.project_level
    ,t.status
    -- Student 1
    ,stu1.user_id AS student1_user_id
    ,stu1.user_role_id AS student1_user_role_id
    ,stu1.displayname AS student1_displayname
    ,stu1.email AS student1_email
    -- Student 2
    ,stu2.user_id AS student2_user_id
    ,stu2.user_role_id AS student2_user_role_id
    ,stu2.displayname AS student2_displayname
    ,stu2.email AS student2_email
    -- Adviser
    ,adv.user_id AS adviser_user_id
    ,adv.user_role_id AS adviser_user_role_id
    ,adv.displayname AS adviser_displayname
    ,adv.email AS adviser_email
    -- Mentor
    ,mnt.user_id AS mentor_user_id
    ,mnt.user_role_id AS mentor_user_role_id
    ,mnt.displayname AS mentor_displayname
    ,mnt.email AS mentor_email
FROM
    teams AS t
    LEFT JOIN students AS stu1 ON stu1.team_id = t.team_id
    LEFT JOIN students AS stu2 ON stu2.team_id = t.team_id AND stu1.user_id < stu2.user_id
    LEFT JOIN advisers AS adv ON adv.user_role_id = t.adviser_user_role_id
    LEFT JOIN mentors AS mnt ON mnt.user_role_id = t.mentor_user_role_id
WHERE
    (
        (SELECT COUNT(*) FROM user_roles_students AS urs WHERE urs.team_id = t.team_id) = 2
        AND
        stu2.user_id IS NOT NULL
    )
    OR
    (
        (SELECT COUNT(*) FROM user_roles_students AS urs WHERE urs.team_id = t.team_id) = 1
        AND
        stu2.user_id IS NULL
    )
;
