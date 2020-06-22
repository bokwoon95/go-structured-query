DROP FUNCTION IF EXISTS trg.submissions() CASCADE;
CREATE OR REPLACE FUNCTION trg.submissions()
RETURNS TRIGGER AS $$ DECLARE
    var_readme TEXT;
    var_poster TEXT;
    var_video TEXT;
BEGIN
    -- If submission_data column is not provided, don't bother doing anything and just return the row as is
    IF NEW.submission_data IS NULL THEN
        RETURN NEW;
    END IF;

    -- data: {"readme":["<h1>README</h1>"]} => readme: '<h1>README</h1>'
    IF NEW.submission_data::JSONB->'readme' IS NOT NULL THEN
        SELECT NEW.submission_data::JSONB->'readme'->>0 INTO var_readme;
    END IF;
    NEW.readme = COALESCE(var_readme, '');

    -- data: {"poster":["https://i.imgur.com/RnlVETv.jpg"]} => poster: 'https://i.imgur.com/RnlVETv.jpg'
    IF NEW.submission_data::JSONB->'poster' IS NOT NULL THEN
        SELECT NEW.submission_data::JSONB->'poster'->>0 INTO var_poster;
    END IF;
    NEW.poster = COALESCE(var_poster, 'gemini');

    -- data: {"video":["https://www.youtube.com/watch?v=WTJSt4wP2ME"]} => video: 'https://www.youtube.com/watch?v=WTJSt4wP2ME'
    IF NEW.submission_data::JSONB->'video' IS NOT NULL THEN
        SELECT NEW.submission_data::JSONB->'video'->>0 INTO var_video;
    END IF;
    NEW.video = COALESCE(var_video, '');

    RETURN NEW;
END $$ LANGUAGE plpgsql;
CREATE TRIGGER submissions_row BEFORE INSERT OR UPDATE ON submissions FOR EACH ROW EXECUTE PROCEDURE trg.submissions();
