DROP FUNCTION IF EXISTS app.update_team(INT, TEXT, TEXT, TEXT, INT, INT, INT, INT);
CREATE OR REPLACE FUNCTION app.update_team(
    arg_team_id INT
    ,arg_status TEXT
    ,arg_team_name TEXT
    ,arg_project_level TEXT
    ,arg_student_user_id_1 INT
    ,arg_student_user_id_2 INT
    ,arg_adviser_user_id INT
    ,arg_mentor_user_id INT
) RETURNS VOID AS $$ DECLARE
    var_adviser_user_role_id INT;
    var_mentor_user_role_id INT;
BEGIN
    UPDATE teams SET status = arg_status, team_name = arg_team_name, project_level = arg_project_level WHERE team_id = arg_team_id;

    UPDATE user_roles_students AS urs SET team_id = NULL WHERE team_id = arg_team_id;
    UPDATE
        user_roles_students AS urs
    SET
        team_id = arg_team_id
    FROM
        user_roles AS ur
        ,users AS u
    WHERE
        urs.user_role_id = ur.user_role_id
        AND ur.user_id = u.user_id
        AND u.user_id IN (arg_student_user_id_1, arg_student_user_id_2)
    ;

    SELECT user_role_id
    INTO var_adviser_user_role_id
    FROM users AS u JOIN user_roles AS ur USING (user_id)
    WHERE u.user_id = arg_adviser_user_id AND ur.role = 'adviser'
    ;
    UPDATE teams SET adviser_user_role_id = var_adviser_user_role_id WHERE team_id = arg_team_id;

    SELECT user_role_id
    INTO var_mentor_user_role_id
    FROM users AS u JOIN user_roles AS ur USING (user_id)
    WHERE u.user_id = arg_mentor_user_id AND ur.role = 'mentor'
    ;
    UPDATE teams SET mentor_user_role_id = var_mentor_user_role_id WHERE team_id = arg_team_id;
END $$ LANGUAGE plpgsql;
