DO $data$ DECLARE
BEGIN
    INSERT INTO submissions
        (team_id, submission_form_id)
    VALUES
        (1, 3)
        ,(2, 3)
        ,(3, 3)
        ,(4, 3)
        ,(5, 3)
        ,(6, 3)
        ,(7, 3)
        ,(8, 3)
        ,(9, 3)
        ,(10, 3)
        ,(11, 3)
        ,(12, 3)
    ;
END $data$;
