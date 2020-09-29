DROP VIEW IF EXISTS v_teams_and_students;
CREATE VIEW v_teams_and_students AS
WITH students AS (
    SELECT u.user_id, u.displayname, u.email, ur.user_role_id, urs.team_id, urs.student_data
    FROM users AS u JOIN user_roles AS ur USING (user_id) LEFT JOIN user_roles_students AS urs USING (user_role_id)
    WHERE ur.role = 'student'
)
SELECT
    t.team_id
    ,t.team_name
    ,t.project_level
    ,t.adviser_user_role_id
    ,t.mentor_user_role_id
    ,student1.displayname AS student1_displayname
    ,student2.displayname AS student2_displayname
FROM
    teams AS t
    LEFT JOIN students AS student1 ON student1.team_id = t.team_id
    LEFT JOIN students AS student2 ON student2.team_id = t.team_id AND student1.user_id < student2.user_id
WHERE
    (
        (SELECT COUNT(*) FROM user_roles_students AS urs WHERE urs.team_id = t.team_id) = 2
        AND
        student2.user_id IS NOT NULL
    )
    OR
    (SELECT COUNT(*) FROM user_roles_students AS urs WHERE urs.team_id = t.team_id) IN (0, 1)
;
