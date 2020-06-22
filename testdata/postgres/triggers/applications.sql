DROP FUNCTION IF EXISTS trg.applications() CASCADE;
CREATE OR REPLACE FUNCTION trg.applications()
RETURNS TRIGGER AS $$ DECLARE
    var_team_name TEXT;
    var_project_level TEXT;
    var_project_idea TEXT;
BEGIN
    -- If application_data column is not provided, don't bother doing anything and just return the row as is
    IF NEW.application_data IS NULL THEN
        RETURN NEW;
    END IF;

    -- data: {"team_name":["nusFoodHitch"]} => team_name: 'nusFoodHitch'
    IF NEW.application_data::JSONB->'team_name' IS NOT NULL THEN
        SELECT NEW.application_data::JSONB->'team_name'->>0 INTO var_team_name;
    END IF;
    NEW.team_name = var_team_name;

    -- data: {"project_level":["gemini"]} => project_level: 'gemini'
    IF NEW.application_data::JSONB->'project_level' IS NOT NULL THEN
        SELECT NEW.application_data::JSONB->'project_level'->>0 INTO var_project_level;
    END IF;
    NEW.project_level = COALESCE(var_project_level, 'gemini');

    -- data: {"project_idea":["lorem ipsum dolor sit amet"]} => project_idea: 'lorem ipsum dolor sit amet'
    IF NEW.application_data::JSONB->'project_idea' IS NOT NULL THEN
        SELECT NEW.application_data::JSONB->'project_idea'->>0 INTO var_project_idea;
    END IF;
    NEW.project_idea = COALESCE(var_project_idea, '');

    RETURN NEW;
END $$ LANGUAGE plpgsql;
CREATE TRIGGER applications_row BEFORE INSERT OR UPDATE ON applications FOR EACH ROW EXECUTE PROCEDURE trg.applications();
