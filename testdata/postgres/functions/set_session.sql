-- Set the session of a user associated with arg_email
DROP FUNCTION IF EXISTS app.set_session(TEXT, TEXT);
CREATE OR REPLACE FUNCTION app.set_session (arg_hash TEXT, arg_email TEXT)
RETURNS TABLE (_user_id INT) AS $$ DECLARE
    var_user_id INT;
BEGIN
    -- If user doesn't exist, raise exception
    SELECT u.user_id INTO var_user_id FROM users AS u WHERE u.email = arg_email;
    IF var_user_id IS NULL THEN
        RAISE EXCEPTION 'user user_id[%] does not exist', var_user_id USING ERRCODE = 'OLAMC';
    END IF;

    -- Create a new session
    INSERT INTO sessions (hash, user_id) VALUES (arg_hash, var_user_id);

    RETURN QUERY SELECT var_user_id AS user_id;
END $$ LANGUAGE plpgsql STRICT;
