SET client_min_messages TO WARNING; -- Make this script a little more quiet
DROP EXTENSION IF EXISTS pgcrypto CASCADE;
DROP SCHEMA IF EXISTS app CASCADE;
DROP SCHEMA IF EXISTS trg CASCADE;

DROP TABLE IF EXISTS cohort_enum CASCADE;
DROP TABLE IF EXISTS stage_enum CASCADE;
DROP TABLE IF EXISTS milestone_enum CASCADE;
DROP TABLE IF EXISTS role_enum CASCADE;
DROP TABLE IF EXISTS project_level_enum CASCADE;
DROP TABLE IF EXISTS project_category_enum CASCADE;
DROP TABLE IF EXISTS applications_status_enum CASCADE;
DROP TABLE IF EXISTS teams_status_enum CASCADE;
DROP TABLE IF EXISTS mime_type_enum CASCADE;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS periods CASCADE;
DROP TABLE IF EXISTS forms CASCADE;
DROP TABLE IF EXISTS forms_authorized_roles CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS applications CASCADE;
DROP TABLE IF EXISTS team_evaluation_pairs CASCADE;
DROP TABLE IF EXISTS user_roles_applicants CASCADE;
DROP TABLE IF EXISTS user_roles_students CASCADE;
DROP TABLE IF EXISTS submissions CASCADE;
DROP TABLE IF EXISTS submissions_categories CASCADE;
DROP TABLE IF EXISTS team_evaluations CASCADE;
DROP TABLE IF EXISTS user_evaluations CASCADE;
DROP TABLE IF EXISTS feedback_on_teams CASCADE;
DROP TABLE IF EXISTS feedback_on_users CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS media CASCADE;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS app;
CREATE SCHEMA IF NOT EXISTS trg;

-----------
-- Enums --
-----------

-- cohort
CREATE TABLE cohort_enum (cohort TEXT PRIMARY KEY, insertion_order SERIAL UNIQUE);
DO $$ DECLARE
    var_year INT := 2016;
    var_current_year INT := DATE_PART('year', CURRENT_DATE);
BEGIN
    INSERT INTO cohort_enum (cohort) VALUES ('');
    WHILE var_year <= var_current_year LOOP
        INSERT INTO cohort_enum (cohort) VALUES (var_year::TEXT);
        var_year := var_year + 1;
    END LOOP;
END $$;
SELECT * FROM cohort_enum;

-- stage
CREATE TABLE stage_enum (stage TEXT PRIMARY KEY);
INSERT INTO stage_enum (stage) VALUES (''), ('application'), ('submission'), ('evaluation'), ('feedback') RETURNING *;

-- milestone
CREATE TABLE milestone_enum (milestone TEXT PRIMARY KEY);
INSERT INTO milestone_enum (milestone) VALUES (''), ('milestone1'), ('milestone2'), ('milestone3') RETURNING *;

-- role
CREATE TABLE role_enum (role TEXT PRIMARY KEY);
INSERT INTO role_enum (role) VALUES (''), ('applicant'), ('student'), ('adviser'), ('mentor'), ('admin') RETURNING *;

-- project_level
CREATE TABLE project_level_enum (project_level TEXT PRIMARY KEY);
INSERT INTO project_level_enum (project_level) VALUES ('vostok'), ('gemini'), ('apollo'), ('artemis') RETURNING *;

-- project_category
CREATE TABLE project_category_enum (project_category TEXT PRIMARY KEY);
INSERT INTO
    project_category_enum (project_category)
VALUES
    ('Website')
    ,('iOS')
    ,('Android')
    ,('VR')
    ,('Game')
RETURNING *
;

-- applications status
CREATE TABLE applications_status_enum (status TEXT PRIMARY KEY);
INSERT INTO applications_status_enum (status) VALUES ('pending'), ('accepted'), ('deleted') RETURNING *;

-- teams status
CREATE TABLE teams_status_enum (status TEXT PRIMARY KEY);
INSERT INTO teams_status_enum (status) VALUES ('good'), ('ok'), ('uncontactable') RETURNING *;

-- mime type
CREATE TABLE mime_type_enum (type TEXT PRIMARY KEY);
INSERT INTO
    mime_type_enum (type)
VALUES
    -- We only need to store the most common MIME types we are planning to
    -- serve to the user. 'application/octet-stream' (i.e. binary data) will be
    -- the fallback MIME type for everything else.
    ('application/octet-stream')

    -- images
    ,('image/jpeg')
    ,('image/png')
    ,('image/gif')
    ,('image/svg+xml')
    ,('image/apng')
    ,('image/bmp')
    ,('image/x-icon')
    ,('image/tiff')
    ,('image/webp')
RETURNING *
;

-- Trigger function that updates the 'updated_at' column of a table
CREATE OR REPLACE FUNCTION trg.updated_at()
RETURNS TRIGGER AS $$ BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END $$ LANGUAGE plpgsql;

------------
-- Tables --
------------

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY
    ,displayname TEXT NOT NULL DEFAULT ''
    ,email TEXT NOT NULL UNIQUE
    ,password TEXT

    ,UNIQUE(displayname, email)
);

CREATE TABLE periods (
    period_id SERIAL PRIMARY KEY
    ,cohort TEXT NOT NULL DEFAULT DATE_PART('year', CURRENT_DATE)
    ,stage TEXT NOT NULL DEFAULT ''
    ,milestone TEXT NOT NULL DEFAULT ''
    ,start_at TIMESTAMPTZ
    ,end_at TIMESTAMPTZ
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (cohort, stage, milestone)
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (stage) REFERENCES stage_enum (stage) ON UPDATE CASCADE
    ,FOREIGN KEY (milestone) REFERENCES milestone_enum (milestone) ON UPDATE CASCADE
);
INSERT INTO periods (period_id, cohort, stage, milestone) VALUES (0, '', '', ''); -- pseudo-null period (period_id: 0) for foreign key purposes
CREATE TRIGGER periods_updated_at BEFORE UPDATE ON periods FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE forms (
    form_id SERIAL PRIMARY KEY
    ,period_id INT NOT NULL DEFAULT 0
    ,name TEXT NOT NULL DEFAULT ''
    ,subsection TEXT NOT NULL DEFAULT ''
    ,questions JSONB
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (period_id, name, subsection)
    ,FOREIGN KEY (period_id) REFERENCES periods (period_id) ON UPDATE CASCADE
);
INSERT INTO forms (form_id, period_id) VALUES (0, 0); -- pseudo-null form schema (form_id: 0) for foreign key purposes
CREATE TRIGGER forms_updated_at BEFORE UPDATE ON forms FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE forms_authorized_roles (
    form_id INT NOT NULL
    ,role TEXT NOT NULL

    ,UNIQUE(form_id, role)
    ,FOREIGN KEY (form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (role) REFERENCES role_enum (role) ON UPDATE CASCADE
);

CREATE TABLE user_roles (
    user_role_id SERIAL PRIMARY KEY
    ,user_id INT NOT NULL
    ,cohort TEXT NOT NULL DEFAULT DATE_PART('year', CURRENT_DATE)
    ,role TEXT NOT NULL
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (user_id, cohort, role)
    ,FOREIGN KEY (user_id) REFERENCES users (user_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (role) REFERENCES role_enum (role) ON UPDATE CASCADE
);
CREATE TRIGGER user_roles_updated_at BEFORE UPDATE ON user_roles FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE teams (
    team_id SERIAL PRIMARY KEY
    ,project_level TEXT NOT NULL DEFAULT 'gemini'
    ,project_idea TEXT NOT NULL DEFAULT ''
    ,cohort TEXT NOT NULL DEFAULT DATE_PART('year', CURRENT_DATE)
    ,status TEXT NOT NULL DEFAULT 'ok'
    ,team_name TEXT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW())*1000)::BIGINT::TEXT
    ,mentor_user_role_id INT
    ,adviser_user_role_id INT
    ,team_data JSONB
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (cohort, team_name)
    ,FOREIGN KEY (mentor_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (adviser_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (project_level) REFERENCES project_level_enum (project_level) ON UPDATE CASCADE
    ,FOREIGN KEY (status) REFERENCES teams_status_enum (status) ON UPDATE CASCADE
);
CREATE TRIGGER teams_updated_at BEFORE UPDATE ON teams FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE applications (
    application_id SERIAL PRIMARY KEY
    ,creator_user_role_id INT
    ,team_id INT
    ,application_form_id INT NOT NULL
    ,application_data JSONB
    ,cohort TEXT NOT NULL DEFAULT DATE_PART('year', CURRENT_DATE)
    ,status TEXT NOT NULL DEFAULT 'pending'
    ,team_name TEXT
    ,project_level TEXT NOT NULL DEFAULT 'gemini'
    ,project_idea TEXT NOT NULL DEFAULT ''
    ,magicstring TEXT UNIQUE DEFAULT TRANSLATE(gen_random_uuid()::TEXT, '-', '')
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (cohort, team_name)
    ,FOREIGN KEY (team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (creator_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (application_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (project_level) REFERENCES project_level_enum (project_level) ON UPDATE CASCADE
    ,FOREIGN KEY (status) REFERENCES applications_status_enum (status) ON UPDATE CASCADE
);
CREATE TRIGGER applications_updated_at BEFORE UPDATE ON applications FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE team_evaluation_pairs (
    evaluatee_team_id INT NOT NULL
    ,evaluator_team_id INT NOT NULL

    ,UNIQUE (evaluatee_team_id, evaluator_team_id)
    ,FOREIGN KEY (evaluatee_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE ON DELETE CASCADE
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE user_roles_applicants (
    user_role_id INT PRIMARY KEY
    ,application_id INT
    ,applicant_form_id INT NOT NULL
    ,applicant_data JSONB

    ,FOREIGN KEY (user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE ON DELETE CASCADE
    ,FOREIGN KEY (application_id) REFERENCES applications (application_id) ON UPDATE CASCADE
    ,FOREIGN KEY (applicant_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);

CREATE TABLE user_roles_students (
    user_role_id INT PRIMARY KEY
    ,team_id INT
    ,student_data JSONB

    ,FOREIGN KEY (user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE ON DELETE CASCADE
    ,FOREIGN KEY (team_id) REFERENCES  teams (team_id) ON UPDATE CASCADE
);

CREATE TABLE submissions (
    submission_id SERIAL PRIMARY KEY
    ,team_id INT NOT NULL
    ,submission_form_id INT NOT NULL
    ,submission_data JSONB
    ,readme TEXT NOT NULL DEFAULT ''
    ,poster TEXT NOT NULL DEFAULT ''
    ,video TEXT NOT NULL DEFAULT ''
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (team_id, submission_form_id)
    ,FOREIGN KEY (team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (submission_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER team_submissions_updated_at BEFORE UPDATE ON submissions FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE submissions_categories (
    submission_id INT NOT NULL
    ,category TEXT NOT NULL

    ,UNIQUE(submission_id, category)
    ,FOREIGN KEY (submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
    ,FOREIGN KEY (category) REFERENCES project_category_enum (project_category) ON UPDATE CASCADE
);

CREATE TABLE team_evaluations (
    team_evaluation_id SERIAL PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_submission_id INT NOT NULL
    ,evaluation_form_id INT NOT NULL
    ,evaluation_data JSONB
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (evaluator_team_id, evaluatee_submission_id)
    ,FOREIGN KEY (evaluation_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
);
CREATE TRIGGER team_evaluations_updated_at BEFORE UPDATE ON team_evaluations FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE user_evaluations (
    user_evaluation_id SERIAL PRIMARY KEY
    ,evaluator_user_role_id INT NOT NULL
    ,evaluatee_submission_id INT NOT NULL
    ,evaluation_form_id INT NOT NULL
    ,evaluation_data JSONB
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE (evaluator_user_role_id, evaluatee_submission_id)
    ,FOREIGN KEY (evaluator_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluation_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER user_evaluations_updated_at BEFORE UPDATE ON user_evaluations FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE feedback_on_teams (
    feedback_id_on_team SERIAL PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_team_id INT NOT NULL
    ,feedback_form_id INT NOT NULL
    ,feedback_data JSONB
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE(evaluator_team_id, evaluatee_team_id)
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (feedback_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER feedback_on_teams_updated_at BEFORE UPDATE ON feedback_on_teams FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE feedback_on_users (
    feedback_id_on_user SERIAL PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_user_role_id INT NOT NULL
    ,feedback_form_id INT NOT NULL
    ,feedback_data JSONB
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,UNIQUE(evaluator_team_id, evaluatee_user_role_id)
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (feedback_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER feedback_on_users_updated_at BEFORE UPDATE ON feedback_on_users FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();

CREATE TABLE sessions (
    hash TEXT PRIMARY KEY
    ,user_id INT NOT NULL
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    ,FOREIGN KEY (user_id) REFERENCES users (user_id) ON UPDATE CASCADE
);

CREATE TABLE media (
    uuid UUID NOT NULL DEFAULT gen_random_uuid()
    ,name TEXT NOT NULL DEFAULT ''
    ,type TEXT NOT NULL DEFAULT 'application/octet-stream'
    ,description TEXT NOT NULL DEFAULT ''
    ,data BYTEA NOT NULL
    ,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    ,deleted_at TIMESTAMPTZ

    ,FOREIGN KEY (type) REFERENCES mime_type_enum (type) ON UPDATE CASCADE
);
CREATE TRIGGER media_updated_at BEFORE UPDATE ON media FOR EACH ROW EXECUTE PROCEDURE trg.updated_at();
