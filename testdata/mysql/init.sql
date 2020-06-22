SET FOREIGN_KEY_CHECKS=0;

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

SET FOREIGN_KEY_CHECKS=1;

-- ---------
-- Enums --
-- ---------

-- cohort
CREATE TABLE cohort_enum (cohort VARCHAR(255) PRIMARY KEY);
SELECT * FROM cohort_enum;
DROP PROCEDURE IF EXISTS insert_cohorts;
DELIMITER $$
CREATE PROCEDURE insert_cohorts() BEGIN
    DECLARE var_year INT DEFAULT 2016;
    DECLARE var_current_year INT DEFAULT CAST(YEAR(NOW()) AS SIGNED);
    INSERT INTO cohort_enum(cohort) VALUES ('');
    WHILE var_year <= var_current_year DO
        INSERT INTO cohort_enum(cohort) VALUES (CAST(var_year AS CHAR));
        SET var_year = var_year + 1;
    END WHILE;
END $$
DELIMITER ;
CALL insert_cohorts();
DROP PROCEDURE IF EXISTS insert_cohorts;

-- stage
CREATE TABLE stage_enum (stage VARCHAR(255) PRIMARY KEY);
INSERT INTO stage_enum (stage) VALUES (''), ('application'), ('submission'), ('evaluation'), ('feedback');

-- milestone
CREATE TABLE milestone_enum (milestone VARCHAR(255) PRIMARY KEY);
INSERT INTO milestone_enum (milestone) VALUES (''), ('milestone1'), ('milestone2'), ('milestone3');

-- role
CREATE TABLE role_enum (role VARCHAR(255) PRIMARY KEY);
INSERT INTO role_enum (role) VALUES (''), ('applicant'), ('student'), ('adviser'), ('mentor'), ('admin');

-- project_level
CREATE TABLE project_level_enum (project_level VARCHAR(255) PRIMARY KEY);
INSERT INTO project_level_enum (project_level) VALUES ('vostok'), ('gemini'), ('apollo'), ('artemis');

-- project_category
CREATE TABLE project_category_enum (project_category VARCHAR(255) PRIMARY KEY);
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
CREATE TABLE applications_status_enum (status VARCHAR(255) PRIMARY KEY);
INSERT INTO applications_status_enum (status) VALUES ('pending'), ('accepted'), ('deleted');

-- teams status
CREATE TABLE teams_status_enum (status VARCHAR(255) PRIMARY KEY);
INSERT INTO teams_status_enum (status) VALUES ('good'), ('ok'), ('uncontactable');

-- mime type
CREATE TABLE mime_type_enum (type VARCHAR(255) PRIMARY KEY);
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

-- MySQL 8 BIN_TO_UUID/ UUID_TO_BIN polyfill for MySQL 5.7
-- https://stackoverflow.com/a/58015720
DROP FUNCTION IF EXISTS BIN_TO_UUID_POLYFILL;
DROP FUNCTION IF EXISTS UUID_TO_BIN_POLYFILL;
DELIMITER $$
CREATE FUNCTION BIN_TO_UUID_POLYFILL(b BINARY(16), f BOOLEAN)
RETURNS CHAR(36) DETERMINISTIC
BEGIN
   DECLARE hexStr CHAR(32);
   SET hexStr = HEX(b);
   RETURN LOWER(CONCAT(
        IF(f,SUBSTR(hexStr, 9, 8),SUBSTR(hexStr, 1, 8)), '-',
        IF(f,SUBSTR(hexStr, 5, 4),SUBSTR(hexStr, 9, 4)), '-',
        IF(f,SUBSTR(hexStr, 1, 4),SUBSTR(hexStr, 13, 4)), '-',
        SUBSTR(hexStr, 17, 4), '-',
        SUBSTR(hexStr, 21)
    ));
END$$
CREATE FUNCTION UUID_TO_BIN_POLYFILL(uuid CHAR(36), f BOOLEAN)
RETURNS BINARY(16)
DETERMINISTIC
BEGIN
    RETURN UNHEX(CONCAT(
        IF(f,SUBSTRING(uuid, 15, 4),SUBSTRING(uuid, 1, 8)),
        SUBSTRING(uuid, 10, 4),
        IF(f,SUBSTRING(uuid, 1, 8),SUBSTRING(uuid, 15, 4)),
        SUBSTRING(uuid, 20, 4),
        SUBSTRING(uuid, 25)
    ));
END$$
DELIMITER ;

-- ----------
-- Tables --
-- ----------

CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY
    ,displayname VARCHAR(255) NOT NULL DEFAULT ''
    ,email VARCHAR(255) NOT NULL UNIQUE
    ,password VARCHAR(255)

    ,UNIQUE(displayname, email)
);

CREATE TABLE periods (
    period_id INT AUTO_INCREMENT PRIMARY KEY
    ,cohort VARCHAR(255) NOT NULL
    ,stage VARCHAR(255) NOT NULL DEFAULT ''
    ,milestone VARCHAR(255) NOT NULL DEFAULT ''
    ,start_at DATETIME
    ,end_at DATETIME
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (cohort, stage, milestone)
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (stage) REFERENCES stage_enum (stage) ON UPDATE CASCADE
    ,FOREIGN KEY (milestone) REFERENCES milestone_enum (milestone) ON UPDATE CASCADE
);
DELIMITER $$
CREATE TRIGGER before_insert_periods BEFORE INSERT ON periods FOR EACH ROW
BEGIN
    IF NEW.cohort IS NULL THEN
        SET NEW.cohort = YEAR(NOW());
    END IF;
END$$
DELIMITER ;
INSERT INTO periods (cohort, stage, milestone) VALUES ('', '', ''); -- pseudo-null period (period_id: 0) for foreign key purposes

CREATE TABLE forms (
    form_id INT AUTO_INCREMENT PRIMARY KEY
    ,period_id INT NOT NULL DEFAULT 0
    ,name VARCHAR(255) NOT NULL DEFAULT ''
    ,subsection VARCHAR(255) NOT NULL DEFAULT ''
    ,questions JSON
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (period_id, name, subsection)
    ,FOREIGN KEY (period_id) REFERENCES periods (period_id) ON UPDATE CASCADE
);
-- pseudo-null form schema (form_id: 0) for foreign key purposes
INSERT INTO forms (form_id, period_id)
SELECT 0, period_id FROM periods ORDER BY period_id ASC LIMIT 1
;

CREATE TABLE forms_authorized_roles (
    form_id INT NOT NULL
    ,role VARCHAR(255) NOT NULL

    ,UNIQUE(form_id, role)
    ,FOREIGN KEY (form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (role) REFERENCES role_enum (role) ON UPDATE CASCADE
);

CREATE TABLE user_roles (
    user_role_id INT AUTO_INCREMENT PRIMARY KEY
    ,user_id INT NOT NULL
    ,cohort VARCHAR(255) NOT NULL
    ,role VARCHAR(255) NOT NULL
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (user_id, cohort, role)
    ,FOREIGN KEY (user_id) REFERENCES users (user_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (role) REFERENCES role_enum (role) ON UPDATE CASCADE
);
DELIMITER $$
CREATE TRIGGER before_insert_user_roles BEFORE INSERT ON user_roles FOR EACH ROW
BEGIN
    IF NEW.cohort IS NULL THEN
        SET NEW.cohort = YEAR(NOW());
    END IF;
END$$
DELIMITER ;

CREATE TABLE teams (
    team_id INT AUTO_INCREMENT PRIMARY KEY
    ,project_level VARCHAR(255) NOT NULL DEFAULT 'gemini'
    ,project_idea VARCHAR(255) NOT NULL DEFAULT ''
    ,cohort VARCHAR(255) NOT NULL
    ,status VARCHAR(255) NOT NULL DEFAULT 'ok'
    ,team_name VARCHAR(255) NOT NULL
    ,mentor_user_role_id INT
    ,adviser_user_role_id INT
    ,team_data JSON
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (cohort, team_name)
    ,FOREIGN KEY (mentor_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (adviser_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (project_level) REFERENCES project_level_enum (project_level) ON UPDATE CASCADE
    ,FOREIGN KEY (status) REFERENCES teams_status_enum (status) ON UPDATE CASCADE
);
DELIMITER $$
CREATE TRIGGER before_insert_teams BEFORE INSERT ON teams FOR EACH ROW
BEGIN
    IF NEW.cohort IS NULL THEN
        SET NEW.cohort = YEAR(NOW());
    END IF;
END$$
DELIMITER ;

CREATE TABLE applications (
    application_id INT AUTO_INCREMENT PRIMARY KEY
    ,creator_user_role_id INT
    ,team_id INT
    ,application_form_id INT NOT NULL
    ,application_data JSON
    ,cohort VARCHAR(255) NOT NULL
    ,status VARCHAR(255) NOT NULL DEFAULT 'pending'
    ,team_name VARCHAR(255)
    ,project_level VARCHAR(255) NOT NULL DEFAULT 'gemini'
    ,project_idea MEDIUMTEXT NOT NULL
    ,magicstring VARCHAR(255) UNIQUE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (cohort, team_name)
    ,FOREIGN KEY (team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (creator_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (application_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (cohort) REFERENCES cohort_enum (cohort) ON UPDATE CASCADE
    ,FOREIGN KEY (project_level) REFERENCES project_level_enum (project_level) ON UPDATE CASCADE
    ,FOREIGN KEY (status) REFERENCES applications_status_enum (status) ON UPDATE CASCADE
);
DELIMITER $$
CREATE TRIGGER before_insert_applications BEFORE INSERT ON applications FOR EACH ROW
BEGIN
    IF NEW.cohort IS NULL THEN
        SET NEW.cohort = YEAR(NOW());
    END IF;
    IF NEW.project_idea IS NULL THEN
        SET NEW.project_idea = '';
    END IF;
END$$
DELIMITER ;

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
    submission_id INT AUTO_INCREMENT PRIMARY KEY
    ,team_id INT NOT NULL
    ,submission_form_id INT NOT NULL
    ,submission_data JSON
    ,readme VARCHAR(255) NOT NULL DEFAULT ''
    ,poster VARCHAR(255) NOT NULL DEFAULT ''
    ,video VARCHAR(255) NOT NULL DEFAULT ''
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (team_id, submission_form_id)
    ,FOREIGN KEY (team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (submission_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);

CREATE TABLE submissions_categories (
    submission_id INT NOT NULL
    ,category VARCHAR(255) NOT NULL

    ,UNIQUE(submission_id, category)
    ,FOREIGN KEY (submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
    ,FOREIGN KEY (category) REFERENCES project_category_enum (project_category) ON UPDATE CASCADE
);

CREATE TABLE team_evaluations (
    team_evaluation_id INT AUTO_INCREMENT PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_submission_id INT NOT NULL
    ,evaluation_form_id INT NOT NULL
    ,evaluation_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (evaluator_team_id, evaluatee_submission_id)
    ,FOREIGN KEY (evaluation_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
);

CREATE TABLE user_evaluations (
    user_evaluation_id INT AUTO_INCREMENT PRIMARY KEY
    ,evaluator_user_role_id INT NOT NULL
    ,evaluatee_submission_id INT NOT NULL
    ,evaluation_form_id INT NOT NULL
    ,evaluation_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE (evaluator_user_role_id, evaluatee_submission_id)
    ,FOREIGN KEY (evaluator_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_submission_id) REFERENCES submissions (submission_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluation_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);

CREATE TABLE feedback_on_teams (
    feedback_id_on_team INT AUTO_INCREMENT PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_team_id INT NOT NULL
    ,feedback_form_id INT NOT NULL
    ,feedback_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE(evaluator_team_id, evaluatee_team_id)
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (feedback_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);

CREATE TABLE feedback_on_users (
    feedback_id_on_user INT AUTO_INCREMENT PRIMARY KEY
    ,evaluator_team_id INT NOT NULL
    ,evaluatee_user_role_id INT NOT NULL
    ,feedback_form_id INT NOT NULL
    ,feedback_data JSON
    ,override_open BOOLEAN NOT NULL DEFAULT FALSE
    ,submitted BOOLEAN NOT NULL DEFAULT FALSE
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,UNIQUE(evaluator_team_id, evaluatee_user_role_id)
    ,FOREIGN KEY (evaluator_team_id) REFERENCES teams (team_id) ON UPDATE CASCADE
    ,FOREIGN KEY (evaluatee_user_role_id) REFERENCES user_roles (user_role_id) ON UPDATE CASCADE
    ,FOREIGN KEY (feedback_form_id) REFERENCES forms (form_id) ON UPDATE CASCADE
);

CREATE TABLE sessions (
    hash VARCHAR(255) PRIMARY KEY
    ,user_id INT NOT NULL
    ,created_at DATETIME NOT NULL DEFAULT NOW()

    ,FOREIGN KEY (user_id) REFERENCES users (user_id) ON UPDATE CASCADE
);

CREATE TABLE media (
    uuid BINARY(16) NOT NULL PRIMARY KEY
    ,name VARCHAR(255) NOT NULL DEFAULT ''
    ,type VARCHAR(255) NOT NULL DEFAULT 'application/octet-stream'
    ,description VARCHAR(255) NOT NULL DEFAULT ''
    ,data BLOB NOT NULL
    ,created_at DATETIME NOT NULL DEFAULT NOW()
    ,updated_at DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
    ,deleted_at DATETIME

    ,FOREIGN KEY (type) REFERENCES mime_type_enum (type) ON UPDATE CASCADE
);
DELIMITER $$
CREATE TRIGGER before_insert_media BEFORE INSERT ON media FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL THEN
        SET NEW.uuid = UUID_TO_BIN_POLYFILL(UUID());
    END IF;
END$$
DELIMITER ;
