-- SELECT name FROM sqlite_master WHERE type = 'table';
-- SELECT name, type FROM pragma_table_info('users');
DROP TABLE IF EXISTS cohort_enum;
DROP TABLE IF EXISTS stage_enum;
DROP TABLE IF EXISTS milestone_enum;
DROP TABLE IF EXISTS role_enum;
DROP TABLE IF EXISTS project_level_enum;
DROP TABLE IF EXISTS project_category_enum;
DROP TABLE IF EXISTS applications_status_enum;
DROP TABLE IF EXISTS teams_status_enum;
DROP TABLE IF EXISTS mime_type_enum;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS periods;
DROP TABLE IF EXISTS forms;
DROP TABLE IF EXISTS forms_authorized_roles;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS team_evaluation_pairs;
DROP TABLE IF EXISTS user_roles_applicants;
DROP TABLE IF EXISTS user_roles_students;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS submissions_categories;
DROP TABLE IF EXISTS team_evaluations;
DROP TABLE IF EXISTS user_evaluations;
DROP TABLE IF EXISTS feedback_on_teams;
DROP TABLE IF EXISTS feedback_on_users;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS media;

-----------
-- Enums --
-----------

-- cohort
CREATE TABLE cohort_enum (cohort TEXT PRIMARY KEY, insertion_order INT);
CREATE TRIGGER cohort_enum_insertion_order AFTER INSERT ON cohort_enum WHEN NEW.insertion_order IS NULL
BEGIN
    UPDATE cohort_enum
    SET insertion_order = IFNULL((SELECT MAX(insertion_order) FROM cohort_enum) + 1, 1)
    WHERE cohort = NEW.cohort;
END;
WITH RECURSIVE cohorts AS (
    SELECT 2016 AS cohort
    UNION
    SELECT cohort + 1 FROM cohorts WHERE cohort + 1 <= CAST(strftime('%Y', CURRENT_TIMESTAMP) AS INT)
)
INSERT OR IGNORE INTO cohort_enum (cohort) SELECT cohort FROM cohorts;

-- stage
CREATE TABLE stage_enum (stage TEXT PRIMARY KEY);
INSERT INTO stage_enum (stage) VALUES (''), ('application'), ('submission'), ('evaluation'), ('feedback');

-- milestone
CREATE TABLE milestone_enum (milestone TEXT PRIMARY KEY);
INSERT INTO milestone_enum (milestone) VALUES (''), ('milestone1'), ('milestone2'), ('milestone3');

-- role
CREATE TABLE role_enum (role TEXT PRIMARY KEY);
INSERT INTO role_enum (role) VALUES (''), ('applicant'), ('student'), ('adviser'), ('mentor'), ('admin');

-- project_level
CREATE TABLE project_level_enum (project_level TEXT PRIMARY KEY);
INSERT INTO project_level_enum (project_level) VALUES ('vostok'), ('gemini'), ('apollo'), ('artemis');

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
;

-- applications status
CREATE TABLE applications_status_enum (status TEXT PRIMARY KEY);
INSERT INTO applications_status_enum (status) VALUES ('pending'), ('accepted'), ('deleted');

-- teams status
CREATE TABLE teams_status_enum (status TEXT PRIMARY KEY);
INSERT INTO teams_status_enum (status) VALUES ('good'), ('ok'), ('uncontactable');

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
;

------------
-- Tables --
------------

CREATE TABLE users (
    user_id INT PRIMARY KEY
    ,displayname TEXT NOT NULL DEFAULT ''
    ,email TEXT NOT NULL UNIQUE
    ,password TEXT

    ,UNIQUE(displayname, email)
);

CREATE TABLE periods (
    period_id INT PRIMARY KEY
    ,cohort TEXT NOT NULL DEFAULT (strftime('%Y', CURRENT_TIMESTAMP))
    ,stage TEXT NOT NULL DEFAULT ''
    ,milestone TEXT NOT NULL DEFAULT ''
    ,start_at TIMESTAMP
    ,end_at TIMESTAMP
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (cohort, stage, milestone)
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (stage) REFERENCES stage_enum (stage) ON UPDATE CASCADE
    ,FOREIGN KEY (milestone) REFERENCES milestone_enum (milestone) ON UPDATE CASCADE
);
INSERT INTO periods (period_id, cohort, stage, milestone) VALUES (0, '', '', ''); -- pseudo-null period (period_id: 0) for foreign key purposes
CREATE TRIGGER periods_updated_at AFTER UPDATE ON periods WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE periods SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE forms (
    form_id INT PRIMARY KEY
    ,period_id INT NOT NULL DEFAULT 0
    ,name TEXT NOT NULL DEFAULT ''
    ,subsection TEXT NOT NULL DEFAULT ''
    ,questions JSON
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (period_id, name, subsection)
    ,FOREIGN KEY (period_id) REFERENCES periods (period_id) ON UPDATE CASCADE
);
INSERT INTO forms (form_id, period_id) VALUES (0, 0); -- pseudo-null form schema (form_id: 0) for foreign key purposes
CREATE TRIGGER forms_updated_at AFTER UPDATE ON forms WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE forms SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE forms_authorized_roles (
    form_id INT NOT NULL
    ,role TEXT NOT NULL

    ,UNIQUE(form_id, role)
    ,FOREIGN KEY (form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (role) REFERENCES role_enum (role) ON UPDATE CASCADE
);

CREATE TABLE user_roles (
    user_role_id INT PRIMARY KEY
    ,user_id INT NOT NULL
    ,cohort TEXT NOT NULL DEFAULT (strftime('%Y', CURRENT_TIMESTAMP))
    ,role TEXT NOT NULL
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (user_id, cohort, role)
    ,FOREIGN KEY (user_id) REFERENCES users (user_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (role) REFERENCES role_enum (role) ON UPDATE CASCADE
);
CREATE TRIGGER user_roles_updated_at AFTER UPDATE ON user_roles WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE user_roles SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE teams (
    team_id INT PRIMARY KEY
    ,project_level TEXT NOT NULL DEFAULT 'gemini'
    ,project_idea TEXT NOT NULL DEFAULT ''
    ,cohort TEXT NOT NULL DEFAULT (strftime('%Y', CURRENT_TIMESTAMP))
    ,status TEXT NOT NULL DEFAULT 'ok'
    ,team_name TEXT NOT NULL
    ,mentor_user_role_id INT
    ,adviser_user_role_id INT
    ,team_data JSON
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (cohort, team_name)
    ,FOREIGN KEY (mentor_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (adviser_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (project_level) REFERENCES project_level_enum (project_level) ON UPDATE CASCADE
    ,FOREIGN KEY (status) REFERENCES teams_status_enum (status) ON UPDATE CASCADE
);
CREATE TRIGGER teams_updated_at AFTER UPDATE ON teams WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE teams SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE applications (
    application_id INT PRIMARY KEY
    ,creator_user_role_id INT
    ,team_id INT
    ,application_form_id INT NOT NULL
    ,application_data JSON
    ,cohort TEXT NOT NULL DEFAULT (strftime('%Y', CURRENT_TIMESTAMP))
    ,status TEXT NOT NULL DEFAULT 'pending'
    ,team_name TEXT
    ,project_level TEXT NOT NULL DEFAULT 'gemini'
    ,project_idea TEXT NOT NULL DEFAULT ''
    ,magicstring TEXT UNIQUE
    ,submitted BOOLEAN NOT NULL DEFAULT 0
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (cohort, team_name)
    ,FOREIGN KEY (team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (creator_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (application_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (project_level) REFERENCES project_level_enum (project_level) ON UPDATE CASCADE
    ,FOREIGN KEY (status) REFERENCES applications_status_enum (status) ON UPDATE CASCADE
);
CREATE TRIGGER applications_updated_at AFTER UPDATE ON applications WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE applications SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

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
    ,applicant_data JSON

    ,FOREIGN KEY (user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE ON DELETE CASCADE
    ,FOREIGN KEY (application_id) REFERENCES applications (application_id) ON UPDATE CASCADE
    ,FOREIGN KEY (applicant_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);

CREATE TABLE user_roles_students (
    user_role_id INT PRIMARY KEY
    ,team_id INT
    ,student_data JSON

    ,FOREIGN KEY (user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE ON DELETE CASCADE
    ,FOREIGN KEY (team_id) REFERENCES  teams (team_id) ON UPDATE CASCADE
);

CREATE TABLE submissions (
    submission_id INT PRIMARY KEY
    ,team_id INT NOT NULL
    ,submission_form_id INT NOT NULL
    ,submission_data JSON
    ,readme TEXT NOT NULL DEFAULT ''
    ,poster TEXT NOT NULL DEFAULT ''
    ,video TEXT NOT NULL DEFAULT ''
    ,override_open BOOLEAN NOT NULL DEFAULT 0
    ,submitted BOOLEAN NOT NULL DEFAULT 0
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (team_id, submission_form_id)
    ,FOREIGN KEY (team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (submission_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER submissions_updated_at AFTER UPDATE ON submissions WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE periods SET submissions = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE submissions_categories (
    submission_id INT NOT NULL
    ,category TEXT NOT NULL

    ,UNIQUE(submission_id, category)
    ,FOREIGN KEY (submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
    ,FOREIGN KEY (category) REFERENCES project_category_enum (project_category) ON UPDATE CASCADE
);

CREATE TABLE team_evaluations (
    team_evaluation_id INT PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_submission_id INT NOT NULL
    ,evaluation_form_id INT NOT NULL
    ,evaluation_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT 0
    ,submitted BOOLEAN NOT NULL DEFAULT 0
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (evaluator_team_id, evaluatee_submission_id)
    ,FOREIGN KEY (evaluation_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
);
CREATE TRIGGER team_evaluations_updated_at_updated_at AFTER UPDATE ON team_evaluations WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE team_evaluations SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE user_evaluations (
    user_evaluation_id INT PRIMARY KEY
    ,evaluator_user_role_id INT NOT NULL
    ,evaluatee_submission_id INT NOT NULL
    ,evaluation_form_id INT NOT NULL
    ,evaluation_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT 0
    ,submitted BOOLEAN NOT NULL DEFAULT 0
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE (evaluator_user_role_id, evaluatee_submission_id)
    ,FOREIGN KEY (evaluator_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluation_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER user_evaluations_updated_at_updated_at AFTER UPDATE ON user_evaluations WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE user_evaluations SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE feedback_on_teams (
    feedback_id_on_team INT PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_team_id INT NOT NULL
    ,feedback_form_id INT NOT NULL
    ,feedback_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT 0
    ,submitted BOOLEAN NOT NULL DEFAULT 0
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE(evaluator_team_id, evaluatee_team_id)
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (feedback_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER feedback_on_teams_updated_at AFTER UPDATE ON feedback_on_teams WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE feedback_on_teams SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE feedback_on_users (
    feedback_id_on_user INT PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_user_role_id INT NOT NULL
    ,feedback_form_id INT NOT NULL
    ,feedback_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT 0
    ,submitted BOOLEAN NOT NULL DEFAULT 0
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,UNIQUE(evaluator_team_id, evaluatee_user_role_id)
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (feedback_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);
CREATE TRIGGER feedback_on_users_updated_at AFTER UPDATE ON feedback_on_users WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE feedback_on_users SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;

CREATE TABLE sessions (
    hash TEXT PRIMARY KEY
    ,user_id INT NOT NULL
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,FOREIGN KEY (user_id) REFERENCES users (user_id) ON UPDATE CASCADE
);

CREATE TABLE media (
    uuid UUID NOT NULL
    ,name TEXT NOT NULL DEFAULT ''
    ,type TEXT NOT NULL DEFAULT 'application/octet-stream'
    ,description TEXT NOT NULL DEFAULT ''
    ,data BYTEA NOT NULL
    ,created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,deleted_at TIMESTAMP

    ,FOREIGN KEY (type) REFERENCES mime_type_enum (type) ON UPDATE CASCADE
);
CREATE TRIGGER media_updated_at AFTER UPDATE ON media WHEN NEW.updated_at < CURRENT_TIMESTAMP
BEGIN UPDATE media SET updated_at = CURRENT_TIMESTAMP WHERE period_id = NEW.period_id; END;
