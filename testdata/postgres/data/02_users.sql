BEGIN;
DO $users$ DECLARE
    var_magicstring TEXT;
    var_application_id INT;
    var_adviser_magicstring TEXT;
    var_mentor_magicstring TEXT;
    var_admin_magicstring TEXT;
    var_cohort TEXT;

    var_adviser_user_role_id_1 INT;
    var_adviser_user_role_id_2 INT;
    var_adviser_user_role_id_3 INT;

    var_mentor_user_role_id_1 INT;
    var_mentor_user_role_id_2 INT;
    var_mentor_user_role_id_3 INT;

    var_team_id_01 INT;
    var_team_id_02 INT;
    var_team_id_03 INT;
    var_team_id_04 INT;
    var_team_id_05 INT;
    var_team_id_06 INT;
    var_team_id_07 INT;
    var_team_id_08 INT;
    var_team_id_09 INT;
    var_team_id_10 INT;
    var_team_id_11 INT;
    var_team_id_12 INT;
BEGIN
    SET TIME ZONE 'Asia/Singapore';

    SELECT cohort INTO var_cohort FROM cohort_enum ORDER BY cohort DESC LIMIT 1;

    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser01', 'adviser01@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser02', 'adviser02@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser03', 'adviser03@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser04', 'adviser04@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser05', 'adviser05@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser06', 'adviser06@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser07', 'adviser07@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser08', 'adviser08@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser09', 'adviser09@u.nus.edu');
    PERFORM app.create_user_role(var_cohort, 'adviser', 'Adviser10', 'adviser10@u.nus.edu');
    SELECT ur.user_role_id INTO var_adviser_user_role_id_1 FROM users AS u JOIN user_roles AS ur USING (user_id) WHERE u.email = 'adviser01@u.nus.edu';
    SELECT ur.user_role_id INTO var_adviser_user_role_id_2 FROM users AS u JOIN user_roles AS ur USING (user_id) WHERE u.email = 'adviser02@u.nus.edu';
    SELECT ur.user_role_id INTO var_adviser_user_role_id_3 FROM users AS u JOIN user_roles AS ur USING (user_id) WHERE u.email = 'adviser03@u.nus.edu';

    PERFORM app.create_user_role(var_cohort, 'mentor', 'Mentor01', 'mentor01@gmail.com');
    PERFORM app.create_user_role(var_cohort, 'mentor', 'Mentor02', 'mentor02@gmail.com');
    PERFORM app.create_user_role(var_cohort, 'mentor', 'Mentor03', 'mentor03@gmail.com');
    PERFORM app.create_user_role(var_cohort, 'mentor', 'Mentor04', 'mentor04@gmail.com');
    PERFORM app.create_user_role(var_cohort, 'mentor', 'Mentor05', 'mentor05@gmail.com');
    SELECT ur.user_role_id INTO var_mentor_user_role_id_1 FROM users AS u JOIN user_roles AS ur USING (user_id) WHERE u.email = 'mentor01@gmail.com';
    SELECT ur.user_role_id INTO var_mentor_user_role_id_2 FROM users AS u JOIN user_roles AS ur USING (user_id) WHERE u.email = 'mentor02@gmail.com';
    SELECT ur.user_role_id INTO var_mentor_user_role_id_3 FROM users AS u JOIN user_roles AS ur USING (user_id) WHERE u.email = 'mentor03@gmail.com';


    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Aqua', 'aqua@kono.suba')
    ;
    PERFORM app.join_application('Dustiness Ford Lalatina', 'darkness@kono.suba', var_magicstring);
    SELECT _team_id INTO var_team_id_01 FROM app.accept_application(var_application_id, 'Aqua x Darkness');
    UPDATE applications SET submitted = TRUE WHERE application_id = var_application_id;

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Satou Kazuma', 'kazuma@kono.suba')
    ;
    PERFORM app.join_application('Megumin', 'megumin@kono.suba', var_magicstring);
    SELECT _team_id INTO var_team_id_02 FROM app.accept_application(var_application_id, 'Kazuma x Megumin');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Adolph Schweigert ','adolphschweigert@gmail.com')
    ;
    PERFORM app.join_application('Adrianna Ballweg ','adriannaballweg@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_03 FROM app.accept_application(var_application_id, 'adolphschweigert x adriannaballweg');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Tamar Mauricio ','tamarmauricio@gmail.com')
    ;
    PERFORM app.join_application('Jacquetta Pitman ','jacquettapitman@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_04 FROM app.accept_application(var_application_id, 'tamarmauricio x jacquettapitman');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Yulanda Griest ','yulandagriest@gmail.com')
    ;
    PERFORM app.join_application('Sharika Brase ','sharikabrase@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_05 FROM app.accept_application(var_application_id, 'yulandagriest x sharikabrase');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Barbara Olinger ','barbaraolinger@gmail.com')
    ;
    PERFORM app.join_application('Kamilah Phoenix ','kamilahphoenix@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_06 FROM app.accept_application(var_application_id, 'barbaraolinger x kamilahphoenix');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Lowell Savarese ','lowellsavarese@gmail.com')
    ;
    PERFORM app.join_application('Cynthia Grand ','cynthiagrand@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_07 FROM app.accept_application(var_application_id, 'lowellsavarese x cynthiagrand');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Sharice Neubauer ','shariceneubauer@gmail.com')
    ;
    PERFORM app.join_application('Corrine Tarlton ','corrinetarlton@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_08 FROM app.accept_application(var_application_id, 'shariceneubauer x corrinetarlton');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Shane Calo ','shanecalo@gmail.com')
    ;
    PERFORM app.join_application('Roman Thronson ','romanthronson@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_09 FROM app.accept_application(var_application_id, 'shanecalo x romanthronson');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Celeste Town ','celestetown@gmail.com')
    ;
    PERFORM app.join_application('Wonda Mccluskey ','wondamccluskey@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_10 FROM app.accept_application(var_application_id, 'celestetown x wondamccluskey');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Reginald Dillenbeck ','reginalddillenbeck@gmail.com')
    ;
    PERFORM app.join_application('Deana Haupt ','deanahaupt@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_11 FROM app.accept_application(var_application_id, 'reginalddillenbeck x deanahaupt');

    SELECT _magicstring, _application_id
    INTO var_magicstring, var_application_id
    FROM app.idempotent_create_application('Kandace Lauer ','kandacelauer@gmail.com')
    ;
    PERFORM app.join_application('Shaquana Donley ','shaquanadonley@gmail.com', var_magicstring);
    SELECT _team_id INTO var_team_id_12 FROM app.accept_application(var_application_id, 'kandacelauer x shaquanadonley');

    -- team_01, team_02
    UPDATE
        teams
    SET
        adviser_user_role_id = var_adviser_user_role_id_1
        ,mentor_user_role_id = var_mentor_user_role_id_1
        ,project_level = 'apollo'
    WHERE
        team_id IN (var_team_id_01, var_team_id_02)
    ;
    -- team_03
    UPDATE
        teams
    SET
        adviser_user_role_id = var_adviser_user_role_id_1
        ,mentor_user_role_id = var_mentor_user_role_id_2
        ,project_level = 'apollo'
    WHERE
        team_id = var_team_id_03
    ;
    -- team_04
    UPDATE
        teams
    SET
        adviser_user_role_id = var_adviser_user_role_id_1
        ,mentor_user_role_id = var_mentor_user_role_id_3
        ,project_level = 'apollo'
    WHERE
        team_id = var_team_id_04
    ;
    -- team_05, team_06, team_07, team_08
    UPDATE
        teams
    SET
        adviser_user_role_id = var_adviser_user_role_id_2
        ,mentor_user_role_id = NULL
        ,project_level = 'gemini'
    WHERE
        team_id IN (var_team_id_05, var_team_id_06, var_team_id_07, var_team_id_08)
    ;
    -- team_09, team_10, team_11, team_12
    UPDATE
        teams
    SET
        adviser_user_role_id = var_adviser_user_role_id_3
        ,mentor_user_role_id = NULL
        ,project_level = 'vostok'
    WHERE
        team_id IN (var_team_id_09, var_team_id_10, var_team_id_11, var_team_id_12)
    ;

    INSERT INTO team_evaluation_pairs (evaluatee_team_id, evaluator_team_id)
    VALUES
        -- adviser_01: team_01, team_02, team_03, team_04
        -- 01
        (var_team_id_01, var_team_id_02)
        ,(var_team_id_01, var_team_id_03)
        -- 02
        ,(var_team_id_02, var_team_id_01)
        ,(var_team_id_02, var_team_id_03)
        -- 03
        ,(var_team_id_03, var_team_id_01)
        ,(var_team_id_03, var_team_id_02)
        -- 04
        ,(var_team_id_04, var_team_id_01)
        ,(var_team_id_04, var_team_id_02)

        -- adviser_02: team_05, team_06, team_07, team_08
        -- 05
        ,(var_team_id_05, var_team_id_06)
        ,(var_team_id_05, var_team_id_07)
        -- 06
        ,(var_team_id_06, var_team_id_07)
        ,(var_team_id_06, var_team_id_08)
        -- 07
        ,(var_team_id_07, var_team_id_08)
        ,(var_team_id_07, var_team_id_05)
        -- 08
        ,(var_team_id_08, var_team_id_05)
        ,(var_team_id_08, var_team_id_06)

        -- adviser_03: team_09, team_10, team_11, team_12
        -- 09
        ,(var_team_id_09, var_team_id_10)
        ,(var_team_id_09, var_team_id_11)
        -- 10
        ,(var_team_id_10, var_team_id_11)
        ,(var_team_id_10, var_team_id_12)
        -- 11
        ,(var_team_id_11, var_team_id_12)
        ,(var_team_id_11, var_team_id_09)
        -- 12
        ,(var_team_id_12, var_team_id_09)
        ,(var_team_id_12, var_team_id_10)
    ON CONFLICT DO NOTHING
    ;

    PERFORM app.idempotent_create_application('John', 'e1119090@u.nus.edu');
    PERFORM app.idempotent_create_application('Shrek', 'shrek.2001@dreamworks.com');
    PERFORM app.idempotent_create_application('Donkey', 'donkey.2001@dreamworks.com');

END $users$;
COMMIT;
