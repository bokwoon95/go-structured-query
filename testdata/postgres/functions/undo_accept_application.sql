-- Un-accept the team by deleting the students and teams
DROP FUNCTION IF EXISTS app.undo_accept_application(INT);
CREATE OR REPLACE FUNCTION app.undo_accept_application (arg_application_id INT)
RETURNS VOID AS $$ DECLARE
    var_team_id INT;
    var_affected_rows INT;
BEGIN
    -- Get team of accepted application
    SELECT team_id INTO var_team_id FROM applications WHERE application_id = arg_application_id;
    IF var_team_id IS NULL THEN
        RAISE EXCEPTION 'tried unaccepting an application application_id[%] without a team', arg_application_id USING ERRCODE = 'OC8R1';
    END IF;
    RAISE DEBUG 'team to be deleted is team_id[%]', var_team_id;

    -- Set team to deleted
    UPDATE teams
    SET deleted_at = NOW()
    WHERE teams.team_id = var_team_id
    ;
    GET DIAGNOSTICS var_affected_rows = ROW_COUNT;
    RAISE DEBUG 'number of teams deleted is [%]', var_affected_rows;

    -- Set student users roles to deleted
    UPDATE user_roles AS ur
    SET deleted_at = NOW()
    FROM user_roles_students AS urs
    WHERE urs.user_role_id = ur.user_role_id AND ur.role = 'student' AND urs.team_id = var_team_id
    ;
    GET DIAGNOSTICS var_affected_rows = ROW_COUNT;
    RAISE DEBUG 'number of students deleted is [%]', var_affected_rows;

    -- Set application to not-deleted
    UPDATE applications SET status = 'pending', deleted_at = NULL WHERE application_id = arg_application_id;
END $$ LANGUAGE plpgsql;
