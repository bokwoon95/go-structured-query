DROP MATERIALIZED VIEW IF EXISTS app.v_latest_cohort CASCADE;
CREATE MATERIALIZED VIEW app.v_latest_cohort AS
SELECT cohort
FROM cohort_enum
ORDER BY insertion_order DESC
LIMIT 1
;
-- Add unique index on materialized view so that concurrent refresh works
CREATE UNIQUE INDEX ON app.v_latest_cohort (cohort);

-- Function to refresh materialized view
CREATE OR REPLACE FUNCTION app.v_latest_cohort_refresh()
RETURNS TRIGGER AS $$ BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY app.v_latest_cohort;
    RETURN NULL;
END $$ LANGUAGE plpgsql;

-- Trigger to call refresh function
DROP TRIGGER IF EXISTS v_latest_cohort_refresh ON cohort_enum CASCADE;
CREATE TRIGGER v_latest_cohort_refresh
AFTER INSERT OR UPDATE OR DELETE OR TRUNCATE ON cohort_enum
FOR EACH STATEMENT EXECUTE PROCEDURE app.v_latest_cohort_refresh()
;
