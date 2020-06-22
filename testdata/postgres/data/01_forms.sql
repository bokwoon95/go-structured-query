BEGIN;
DO $forms$ DECLARE
    var_cohort TEXT;
    var_period_id INT;
BEGIN
    SET TIME ZONE 'Asia/Singapore';

    SELECT cohort INTO var_cohort FROM cohort_enum ORDER BY cohort DESC LIMIT 1;

    -- Application
    INSERT INTO periods (cohort, stage, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'application'
            ,DATE(NOW() + '-1 week'::INTERVAL)
            ,DATE(NOW() + '1 week'::INTERVAL) + '23 hours 59 minutes 59 seconds'::INTERVAL
        )
    RETURNING
        period_id INTO var_period_id
    ;
    -- Application Form Subsection 'application'
    INSERT INTO forms (period_id, subsection, questions)
    VALUES
        (var_period_id, 'application' ,$$[
            {"Type":"radio","Text":"<b>Disclaimer: We will not be able to accept your application for Orbital should you not agree with the terms below.</b><p class=\"f6 mid-gray\">I hereby authorise, agree and consent to allow the Orbital organiser and its sponsors to: (a) collect, use, disclose and/or process personal data about me that I had previously provided, that I now provide, that I may in future provide with and/or that Orbital organisers possesses about me including but not limited to my name, my email address for contacting, tracking and marketing purposes. (b) disclose personal data about me to third party service providers in order to perform certain functions in connection with the abovementioned purposes so long as disclosure is necessary</p>","Name":"disclaimer","Options":[{"Value":"agree","Display":"I agree"}],"Subquestions":[]}
            ,{"Type":"short text","Name":"team_name","Text":"<b>Team Name</b>"}
            ,{"Type":"long text","Name":"project_idea","Text":"<b>Summary of your project idea</b>"}
            ,{"Type":"radio","Text":"<b>Which level of achievement are you interested in?</b>","Name":"project_level","Options":[{"Value":"vostok","Display":"Beginner (Восто́к; Vostok)"},{"Value":"gemini","Display":"Intermediate (Project Gemini)"},{"Value":"apollo","Display":"Advanced (Apollo 11)"},{"Value":"apollo_mentor","Display":"Advanced with mentorship (Apollo 11)"}],"Subquestions":[]}
            ,{"Type":"checkbox","Text":"<b>Your target audience(s)</b>","Name":"target_audience","Options":[{"Value":"nus","Display":"NUS"},{"Value":"universities","Display":"universities"},{"Value":"youths","Display":"youths"},{"Value":"adults","Display":"adults"},{"Value":"elderly","Display":"elderly"},{"Value":"children","Display":"children"},{"Value":"toddlers","Display":"toddlers"},{"Value":"families","Display":"families"},{"Value":"parents","Display":"parents"},{"Value":"others","Display":"Others"}],"Subquestions":[]}
        ]$$)
    ;
    -- Application Form Subsection 'applicant'
    INSERT INTO forms (period_id, subsection, questions)
    VALUES
        (var_period_id, 'applicant', $$[
            {"Type":"radio","Text":"<b>Your faculty?</b>","Name":"faculty","Options":[{"Value":"computing","Display":"School of Computing (BA/CS/CEG/InfoSys/InfoSec)"},{"Value":"business","Display":"Business"},{"Value":"engineering","Display":"Engineering"},{"Value":"fass","Display":"Arts and Social Science"},{"Value":"law","Display":"Law"},{"Value":"medicine","Display":"Medicine"},{"Value":"sde","Display":"School of Design and Environment"},{"Value":"science","Display":"Science"},{"Value":"others","Display":"Others"}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Which is your year of matriculation?</b>","Name":"matric_year","Options":[{"Value":"2020","Display":"2020"},{"Value":"2018","Display":"2019"},{"Value":"2017","Display":"2017"},{"Value":"others","Display":"others"}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Will you be overseas for any period of time during the summer?</b>","Name":"overseas","Options":[{"Value":"no","Display":"No, planning to stay in Singapore for the whole summer (minus a short vacation)"},{"Value":"some","Display":"Yes, for some portion (more than 2-3 days)"},{"Value":"most","Display":"Yes, for almost the whole summer"}],"Subquestions":[]}
            ,{"Type":"checkbox","Text":"<b>Your interests</b>","Name":"interests","Options":[{"Value":"raspberrypi","Display":"Raspberry Pi"},{"Value":"arduino","Display":"Arduino"},{"Value":"iot","Display":"IOT"}],"Subquestions":[]}
        ]$$)
    ;

    -- Milestones
    INSERT INTO periods (cohort, milestone, start_at, end_at)
    VALUES
        (var_cohort, 'milestone1', NOW(), DATE(NOW() + '1 month'::INTERVAL) + '23 hours 59 minutes'::INTERVAL)
        ,(var_cohort, 'milestone2', DATE(NOW() + '1 month'::INTERVAL), DATE(NOW() + '2 months'::INTERVAL) + '23 hours 59 minutes'::INTERVAL)
        ,(var_cohort, 'milestone3', DATE(NOW() + '2 months'::INTERVAL), DATE(NOW() + '3 months'::INTERVAL) + '23 hours 59 minutes'::INTERVAL)
    ;

    -- Milestone 1 Submission
    INSERT INTO periods (cohort, stage, milestone, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'submission'
            ,'milestone1'
            ,NOW()
            ,DATE(NOW() + '2 weeks'::INTERVAL) + '23 hours 59 minutes'::INTERVAL
        )
    RETURNING period_id INTO var_period_id
    ;
    -- Milestone 1 Submission Form
    INSERT INTO forms (period_id, questions)
    VALUES
        (var_period_id ,$$[
            {"Type":"long text","Text":"<b>Project Readme</b>","Name":"readme","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Project Log</b>","Name":"log","Options":[],"Subquestions":[]}
            ,{"Type":"short text","Name":"poster","Text":"<b>Poster Link</b>"}
            ,{"Type":"short text","Name":"video","Text":"<b>Video Link</b>"}
        ]$$)
    ;
    -- Milestone 1 Evaluation
    INSERT INTO periods (cohort, stage, milestone, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'evaluation'
            ,'milestone1'
            ,DATE(NOW() + '2 weeks'::INTERVAL)
            ,DATE(NOW() + '4 weeks'::INTERVAL) + '23 hours 59 minutes'::INTERVAL
        )
    RETURNING period_id INTO var_period_id
    ;
    -- Milestone 1 Evaluation Form
    INSERT INTO forms (period_id, questions)
    VALUES
        (var_period_id ,$$[
            {"Type":"paragraph","Text":"<h3>Project Ideation (Section 1 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>The project should serve some user needs. Is it clear who would benefit from the project?</b>","Name":"clear target audience","Options":[{"Value":"1","Display":"The project appears to be trying to serve some user needs, but I am not sure who would benefit from it."},{"Value":"2","Display":"I have a rough idea about who would benefit from the project. It would be better if more details were given."},{"Value":"3","Display":"I am very clear about who would benefit from the project."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Is the problem well defined? Is it clear what need the app is trying to satisfy? Is there a real need to solve the problem?</b>\n<p class=\"f6 mid-gray\">Please grade this with respect to the group of users specified. Be broad in deciding what a problem is. It can be to do something not done before, to do something better, faster, cheaper, or to provide entertainment, fun, etc.</p>","Name":"clear and relevant problem","Options":[{"Value":"1","Display":"The project appears to be trying to do something, but I am not sure what it is."},{"Value":"2","Display":"I can roughly see what problem is."},{"Value":"3","Display":"There is clearly a problem, but it is unclear that it really needs to be solved."},{"Value":"4","Display":"There is clearly a problem and it is likely that solving it would be useful."},{"Value":"5","Display":"The problem is real. There is a real need for such an project."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Does the project solve the most important aspects of the problem?</b>","Name":"solution works","Options":[{"Value":"1","Display":"The project does not really solve the problem."},{"Value":"2","Display":"The project would provide a minimal solution to the problem."},{"Value":"3","Display":"The project would provide an average solution to the problem."},{"Value":"4","Display":"The project would a good solution to the most important aspects of the problem."},{"Value":"5","Display":"The project would solve the problem, and creatively too"}],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please provide feedback to the project team on the problem they are trying to solve based on the ratings you gave to the three questions.</b>","Name":"explain problem+solution evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Features for the next phase (Section 2 of 5)</h3>\n<p class=\"f6 mid-gray\">Note: the team may have already started implementing and may have accomplished the necessary work to implement features. These are not being evaluated in this phase, as the focus is on ideation. In this section, we are evaluating the team's planning for features for the next phase (Jun, up to Milestone 2.)</p>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>\nFor each feature that is being proposed for this phase (minimum of 2-3):\n<ul>\n<li>Is the user role (e.g. public, member, admin) well specified?</li>\n<li>Is the desired outcome (user goal) clear?</li>\n<li>(Optional) Is the benefit clear?</li>\n</ul>\n</b>\n<p class=\"f6 mid-gray\">\nPlease choose one option per feature on whether you agree with this statement : The features are clearly specified and I will have no trouble deciding whether to accept them at the end of the phase.\n<br>\nStrongly Disagree, Disagree, Neutral, Agree, Strongly Agree.\n<p>\n<p class=\"f6 mid-gray\">\ne.g., Feature 1 (CRUD for recording Haze Pledges): Agree \n<br> Feature 2 (Facebook Login): Agree\n<br> Feature 3 (Leaderboard): Disagree\n</p>\n<p class=\"f6 mid-gray\">\n(You can/should be more fine grained in your evaluation than the sample).\n</p>","Name":"features evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please give written feedback to explain your rating for whether the features have been clearly specified.</b>\n<p class=\"f6 mid-gray\">Please write a minimum of 1-2 sentences and a maximum of 500 words (approximately 1 page).</p>","Name":"features clear","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>If the selected features are successfully implemented, there will be a minimum viable product to demonstrate by the end of this phase.</b>","Name":"features complete mvp","Options":[{"Value":"1","Display":"Strongly Disagree"},{"Value":"2","Display":"Disagree"},{"Value":"3","Display":"Neutral"},{"Value":"4","Display":"Agree"},{"Value":"5","Display":"Strongly Agree"}],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please explain your rating on whether the set of features selected for implementation is appropriate. Can the system be demonstrated after the feature set is completed in this phase?</b>","Name":"explain features evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Project README, poster, video, and log (Section 3 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The project README should explain the purpose of the project in detail, describe each of the feature to be implemented in the next phase clearly, and outline a draft plan for the remaining phases.<p>\n<b>After viewing the README...</b>","Name":"clear project direction","Options":[{"Value":"1","Display":"I still have no idea what the project is supposed to do."},{"Value":"2","Display":"I have a reasonable idea of what the project does, but not of the features."},{"Value":"3","Display":"I have a reasonable idea of what the project does and a rough idea of the features."},{"Value":"4","Display":"I have a good idea of what the project does and a reasonable idea of the features to be implemented."},{"Value":"5","Display":"I have a clear idea of what the project does and the features to be implemented."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The poster and the video (referenced in the README) should clearly showcase the project's purpose, scope and features implemented so far. They should assist you as an evaluating team in the peer evaluation process.</p>\n<b>After reading the poster and watching the video...</b>","Name":"poster and video quality","Options":[{"Value":"0","Display":"I couldn't view the poster and the video, or there were too many problems with the submitted poster and video. (0 out of 3)"},{"Value":"1","Display":"The poster and the video make minimal effort to document their project's work and features, and/or just repeats information from the README and log. (1 out of 3)"},{"Value":"2","Display":"The poster and the video are mostly complete. I have some sense of the purpose, scope and features of the project, but some aspects are not clear. (2 out of 3)"},{"Value":"3","Display":"The poster and the video are complete and I have a good sense of the goals, scope and features of the project. (3 out of 3)"},{"Value":"4","Display":"Excellent! The poster and the video are not only complete but also well-prepared. I have learnt from them how to improve my own work. (bonus point)"}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The project log should clearly document how much time the team (and its individual students, where applicable) have spent on their Orbital work so far.</p>\n<b>After viewing the project log...</b>","Name":"project log quality","Options":[{"Value":"1","Display":"I still have no idea of how much time the team members have invested in their project."},{"Value":"2","Display":"I have a reasonable idea of how much time the team members have invested in their project and some vague notion of what they have spent it on."},{"Value":"3","Display":"I have a reasonable idea of both how much time the team members have invested in their project and what they have spent it on."},{"Value":"4","Display":"I have a good idea of how much time the team members have invested in their project and what they have spent it on."},{"Value":"5","Display":"I learned from this group's log and what I can do in my own project for logging, it was excellent!"}],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Critical Feedback (Section 4 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>If there is any critical feedback that you do not want associated with your team, but feel would be helpful to the team to know, please provide it here. It will appear in the private section of the evaluation forms for the group but it will be anonymized.</b>","Name":"critical feedback","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Overall Evaluation (Section 5 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Please give your overall rating for the submission. This will be used to help us eventually decide whether the team passes and what level of achievement is obtained. This section is also private feedback and your team's name is anonymized when giving this feedback.</b>","Name":"overall rating","Options":[{"Value":"1","Display":"1 out of 4 stars. Likely to fail Orbital"},{"Value":"2","Display":"2 out of 4 stars. Sufficient to pass the beginner level (Vostok), maybe good enough for intermediate level (Project Gemini)"},{"Value":"3","Display":"3 out of 4 stars. Definitely intermediate level (Project Gemini). Maybe good enough for advanced level (Apollo 11)"},{"Value":"4","Display":"4 out of 4 stars. Definitely good enough for advanced level (Apollo 11)."},{"Value":"5","Display":"5 out of 4 stars. Wow! (Bonus point)"}],"Subquestions":[]}
        ]$$)
    ;

    -- Milestone 2 Submission
    INSERT INTO periods (cohort, stage, milestone, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'submission'
            ,'milestone2'
            ,DATE(NOW() + '4 weeks'::INTERVAL)
            ,DATE(NOW() + '6 weeks'::INTERVAL) + '23 hours 59 minutes'::INTERVAL
        )
    RETURNING period_id INTO var_period_id
    ;
    -- Milestone 2 Submission Form
    INSERT INTO forms (period_id, questions)
    VALUES
        (var_period_id ,$$[
            {"Type":"long text","Text":"<b>Project Readme</b>","Name":"readme","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Project Log</b>","Name":"log","Options":[],"Subquestions":[]}
            ,{"Type":"short text","Name":"poster","Text":"<b>Poster Link</b>"}
            ,{"Type":"short text","Name":"video","Text":"<b>Video Link</b>"}
        ]$$)
    ;
    -- Milestone 2 Evaluation
    INSERT INTO periods (cohort, stage, milestone, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'evaluation'
            ,'milestone2'
            ,DATE(NOW() + '6 weeks'::INTERVAL)
            ,DATE(NOW() + '8 weeks'::INTERVAL) + '23 hours 59 minutes'::INTERVAL
        )
    RETURNING period_id INTO var_period_id
    ;
    -- Milestone 2 Evaluation Form
    INSERT INTO forms (period_id, questions)
    VALUES
        (var_period_id ,$$[
            {"Type":"paragraph","Text":"<h3>Project Ideation (Section 1 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>The project should serve some user needs. Is it clear who would benefit from the project?</b>","Name":"clear target audience","Options":[{"Value":"1","Display":"The project appears to be trying to serve some user needs, but I am not sure who would benefit from it."},{"Value":"2","Display":"I have a rough idea about who would benefit from the project. It would be better if more details were given."},{"Value":"3","Display":"I am very clear about who would benefit from the project."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Is the problem well defined? Is it clear what need the app is trying to satisfy? Is there a real need to solve the problem?</b>\n<p class=\"f6 mid-gray\">Please grade this with respect to the group of users specified. Be broad in deciding what a problem is. It can be to do something not done before, to do something better, faster, cheaper, or to provide entertainment, fun, etc.</p>","Name":"clear and relevant problem","Options":[{"Value":"1","Display":"The project appears to be trying to do something, but I am not sure what it is."},{"Value":"2","Display":"I can roughly see what problem is."},{"Value":"3","Display":"There is clearly a problem, but it is unclear that it really needs to be solved."},{"Value":"4","Display":"There is clearly a problem and it is likely that solving it would be useful."},{"Value":"5","Display":"The problem is real. There is a real need for such an project."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Does the project solve the most important aspects of the problem?</b>","Name":"solution works","Options":[{"Value":"1","Display":"The project does not really solve the problem."},{"Value":"2","Display":"The project would provide a minimal solution to the problem."},{"Value":"3","Display":"The project would provide an average solution to the problem."},{"Value":"4","Display":"The project would a good solution to the most important aspects of the problem."},{"Value":"5","Display":"The project would solve the problem, and creatively too"}],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please provide feedback to the project team on the problem they are trying to solve based on the ratings you gave to the three questions.</b>","Name":"explain problem+solution evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Features for the next phase (Section 2 of 5)</h3>\n<p class=\"f6 mid-gray\">Note: the team may have already started implementing and may have accomplished the necessary work to implement features. These are not being evaluated in this phase, as the focus is on ideation. In this section, we are evaluating the team's planning for features for the next phase (Jun, up to Milestone 2.)</p>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>\nFor each feature that is being proposed for this phase (minimum of 2-3):\n<ul>\n<li>Is the user role (e.g. public, member, admin) well specified?</li>\n<li>Is the desired outcome (user goal) clear?</li>\n<li>(Optional) Is the benefit clear?</li>\n</ul>\n</b>\n<p class=\"f6 mid-gray\">\nPlease choose one option per feature on whether you agree with this statement : The features are clearly specified and I will have no trouble deciding whether to accept them at the end of the phase.\n<br>\nStrongly Disagree, Disagree, Neutral, Agree, Strongly Agree.\n<p>\n<p class=\"f6 mid-gray\">\ne.g., Feature 1 (CRUD for recording Haze Pledges): Agree \n<br> Feature 2 (Facebook Login): Agree\n<br> Feature 3 (Leaderboard): Disagree\n</p>\n<p class=\"f6 mid-gray\">\n(You can/should be more fine grained in your evaluation than the sample).\n</p>","Name":"features evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please give written feedback to explain your rating for whether the features have been clearly specified.</b>\n<p class=\"f6 mid-gray\">Please write a minimum of 1-2 sentences and a maximum of 500 words (approximately 1 page).</p>","Name":"features clear","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>If the selected features are successfully implemented, there will be a minimum viable product to demonstrate by the end of this phase.</b>","Name":"features complete mvp","Options":[{"Value":"1","Display":"Strongly Disagree"},{"Value":"2","Display":"Disagree"},{"Value":"3","Display":"Neutral"},{"Value":"4","Display":"Agree"},{"Value":"5","Display":"Strongly Agree"}],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please explain your rating on whether the set of features selected for implementation is appropriate. Can the system be demonstrated after the feature set is completed in this phase?</b>","Name":"explain features evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Project README, poster, video, and log (Section 3 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The project README should explain the purpose of the project in detail, describe each of the feature to be implemented in the next phase clearly, and outline a draft plan for the remaining phases.<p>\n<b>After viewing the README...</b>","Name":"clear project direction","Options":[{"Value":"1","Display":"I still have no idea what the project is supposed to do."},{"Value":"2","Display":"I have a reasonable idea of what the project does, but not of the features."},{"Value":"3","Display":"I have a reasonable idea of what the project does and a rough idea of the features."},{"Value":"4","Display":"I have a good idea of what the project does and a reasonable idea of the features to be implemented."},{"Value":"5","Display":"I have a clear idea of what the project does and the features to be implemented."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The poster and the video (referenced in the README) should clearly showcase the project's purpose, scope and features implemented so far. They should assist you as an evaluating team in the peer evaluation process.</p>\n<b>After reading the poster and watching the video...</b>","Name":"poster and video quality","Options":[{"Value":"0","Display":"I couldn't view the poster and the video, or there were too many problems with the submitted poster and video. (0 out of 3)"},{"Value":"1","Display":"The poster and the video make minimal effort to document their project's work and features, and/or just repeats information from the README and log. (1 out of 3)"},{"Value":"2","Display":"The poster and the video are mostly complete. I have some sense of the purpose, scope and features of the project, but some aspects are not clear. (2 out of 3)"},{"Value":"3","Display":"The poster and the video are complete and I have a good sense of the goals, scope and features of the project. (3 out of 3)"},{"Value":"4","Display":"Excellent! The poster and the video are not only complete but also well-prepared. I have learnt from them how to improve my own work. (bonus point)"}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The project log should clearly document how much time the team (and its individual students, where applicable) have spent on their Orbital work so far.</p>\n<b>After viewing the project log...</b>","Name":"project log quality","Options":[{"Value":"1","Display":"I still have no idea of how much time the team members have invested in their project."},{"Value":"2","Display":"I have a reasonable idea of how much time the team members have invested in their project and some vague notion of what they have spent it on."},{"Value":"3","Display":"I have a reasonable idea of both how much time the team members have invested in their project and what they have spent it on."},{"Value":"4","Display":"I have a good idea of how much time the team members have invested in their project and what they have spent it on."},{"Value":"5","Display":"I learned from this group's log and what I can do in my own project for logging, it was excellent!"}],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Critical Feedback (Section 4 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>If there is any critical feedback that you do not want associated with your team, but feel would be helpful to the team to know, please provide it here. It will appear in the private section of the evaluation forms for the group but it will be anonymized.</b>","Name":"critical feedback","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Overall Evaluation (Section 5 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Please give your overall rating for the submission. This will be used to help us eventually decide whether the team passes and what level of achievement is obtained. This section is also private feedback and your team's name is anonymized when giving this feedback.</b>","Name":"overall rating","Options":[{"Value":"1","Display":"1 out of 4 stars. Likely to fail Orbital"},{"Value":"2","Display":"2 out of 4 stars. Sufficient to pass the beginner level (Vostok), maybe good enough for intermediate level (Project Gemini)"},{"Value":"3","Display":"3 out of 4 stars. Definitely intermediate level (Project Gemini). Maybe good enough for advanced level (Apollo 11)"},{"Value":"4","Display":"4 out of 4 stars. Definitely good enough for advanced level (Apollo 11)."},{"Value":"5","Display":"5 out of 4 stars. Wow! (Bonus point)"}],"Subquestions":[]}
        ]$$)
    ;

    -- Milestone 3 Submission
    INSERT INTO periods (cohort, stage, milestone, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'submission'
            ,'milestone3'
            ,DATE(NOW() + '8 weeks'::INTERVAL)
            ,DATE(NOW() + '10 weeks'::INTERVAL) + '23 hours 59 minutes'::INTERVAL
        )
    RETURNING period_id INTO var_period_id
    ;
    -- Milestone 3 Submission Form
    INSERT INTO forms (period_id, questions)
    VALUES
        (var_period_id ,$$[
            {"Type":"long text","Text":"<b>Project Readme</b>","Name":"readme","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Project Log</b>","Name":"log","Options":[],"Subquestions":[]}
            ,{"Type":"short text","Name":"poster","Text":"<b>Poster Link</b>"}
            ,{"Type":"short text","Name":"video","Text":"<b>Video Link</b>"}
        ]$$)
    ;
    -- Milestone 3 Evaluation
    INSERT INTO periods (cohort, stage, milestone, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'evaluation'
            ,'milestone3'
            ,DATE(NOW() + '10 weeks'::INTERVAL)
            ,DATE(NOW() + '12 weeks'::INTERVAL) + '23 hours 59 minutes'::INTERVAL
        )
    RETURNING period_id INTO var_period_id
    ;
    -- Milestone 3 Evaluation Form
    INSERT INTO forms (period_id, questions)
    VALUES
        (var_period_id ,$$[
            {"Type":"paragraph","Text":"<h3>Project Ideation (Section 1 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>The project should serve some user needs. Is it clear who would benefit from the project?</b>","Name":"clear target audience","Options":[{"Value":"1","Display":"The project appears to be trying to serve some user needs, but I am not sure who would benefit from it."},{"Value":"2","Display":"I have a rough idea about who would benefit from the project. It would be better if more details were given."},{"Value":"3","Display":"I am very clear about who would benefit from the project."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Is the problem well defined? Is it clear what need the app is trying to satisfy? Is there a real need to solve the problem?</b>\n<p class=\"f6 mid-gray\">Please grade this with respect to the group of users specified. Be broad in deciding what a problem is. It can be to do something not done before, to do something better, faster, cheaper, or to provide entertainment, fun, etc.</p>","Name":"clear and relevant problem","Options":[{"Value":"1","Display":"The project appears to be trying to do something, but I am not sure what it is."},{"Value":"2","Display":"I can roughly see what problem is."},{"Value":"3","Display":"There is clearly a problem, but it is unclear that it really needs to be solved."},{"Value":"4","Display":"There is clearly a problem and it is likely that solving it would be useful."},{"Value":"5","Display":"The problem is real. There is a real need for such an project."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Does the project solve the most important aspects of the problem?</b>","Name":"solution works","Options":[{"Value":"1","Display":"The project does not really solve the problem."},{"Value":"2","Display":"The project would provide a minimal solution to the problem."},{"Value":"3","Display":"The project would provide an average solution to the problem."},{"Value":"4","Display":"The project would a good solution to the most important aspects of the problem."},{"Value":"5","Display":"The project would solve the problem, and creatively too"}],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please provide feedback to the project team on the problem they are trying to solve based on the ratings you gave to the three questions.</b>","Name":"explain problem+solution evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Features for the next phase (Section 2 of 5)</h3>\n<p class=\"f6 mid-gray\">Note: the team may have already started implementing and may have accomplished the necessary work to implement features. These are not being evaluated in this phase, as the focus is on ideation. In this section, we are evaluating the team's planning for features for the next phase (Jun, up to Milestone 2.)</p>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>\nFor each feature that is being proposed for this phase (minimum of 2-3):\n<ul>\n<li>Is the user role (e.g. public, member, admin) well specified?</li>\n<li>Is the desired outcome (user goal) clear?</li>\n<li>(Optional) Is the benefit clear?</li>\n</ul>\n</b>\n<p class=\"f6 mid-gray\">\nPlease choose one option per feature on whether you agree with this statement : The features are clearly specified and I will have no trouble deciding whether to accept them at the end of the phase.\n<br>\nStrongly Disagree, Disagree, Neutral, Agree, Strongly Agree.\n<p>\n<p class=\"f6 mid-gray\">\ne.g., Feature 1 (CRUD for recording Haze Pledges): Agree \n<br> Feature 2 (Facebook Login): Agree\n<br> Feature 3 (Leaderboard): Disagree\n</p>\n<p class=\"f6 mid-gray\">\n(You can/should be more fine grained in your evaluation than the sample).\n</p>","Name":"features evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please give written feedback to explain your rating for whether the features have been clearly specified.</b>\n<p class=\"f6 mid-gray\">Please write a minimum of 1-2 sentences and a maximum of 500 words (approximately 1 page).</p>","Name":"features clear","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>If the selected features are successfully implemented, there will be a minimum viable product to demonstrate by the end of this phase.</b>","Name":"features complete mvp","Options":[{"Value":"1","Display":"Strongly Disagree"},{"Value":"2","Display":"Disagree"},{"Value":"3","Display":"Neutral"},{"Value":"4","Display":"Agree"},{"Value":"5","Display":"Strongly Agree"}],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Please explain your rating on whether the set of features selected for implementation is appropriate. Can the system be demonstrated after the feature set is completed in this phase?</b>","Name":"explain features evaluation","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Project README, poster, video, and log (Section 3 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The project README should explain the purpose of the project in detail, describe each of the feature to be implemented in the next phase clearly, and outline a draft plan for the remaining phases.<p>\n<b>After viewing the README...</b>","Name":"clear project direction","Options":[{"Value":"1","Display":"I still have no idea what the project is supposed to do."},{"Value":"2","Display":"I have a reasonable idea of what the project does, but not of the features."},{"Value":"3","Display":"I have a reasonable idea of what the project does and a rough idea of the features."},{"Value":"4","Display":"I have a good idea of what the project does and a reasonable idea of the features to be implemented."},{"Value":"5","Display":"I have a clear idea of what the project does and the features to be implemented."}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The poster and the video (referenced in the README) should clearly showcase the project's purpose, scope and features implemented so far. They should assist you as an evaluating team in the peer evaluation process.</p>\n<b>After reading the poster and watching the video...</b>","Name":"poster and video quality","Options":[{"Value":"0","Display":"I couldn't view the poster and the video, or there were too many problems with the submitted poster and video. (0 out of 3)"},{"Value":"1","Display":"The poster and the video make minimal effort to document their project's work and features, and/or just repeats information from the README and log. (1 out of 3)"},{"Value":"2","Display":"The poster and the video are mostly complete. I have some sense of the purpose, scope and features of the project, but some aspects are not clear. (2 out of 3)"},{"Value":"3","Display":"The poster and the video are complete and I have a good sense of the goals, scope and features of the project. (3 out of 3)"},{"Value":"4","Display":"Excellent! The poster and the video are not only complete but also well-prepared. I have learnt from them how to improve my own work. (bonus point)"}],"Subquestions":[]}
            ,{"Type":"radio","Text":"<p>The project log should clearly document how much time the team (and its individual students, where applicable) have spent on their Orbital work so far.</p>\n<b>After viewing the project log...</b>","Name":"project log quality","Options":[{"Value":"1","Display":"I still have no idea of how much time the team members have invested in their project."},{"Value":"2","Display":"I have a reasonable idea of how much time the team members have invested in their project and some vague notion of what they have spent it on."},{"Value":"3","Display":"I have a reasonable idea of both how much time the team members have invested in their project and what they have spent it on."},{"Value":"4","Display":"I have a good idea of how much time the team members have invested in their project and what they have spent it on."},{"Value":"5","Display":"I learned from this group's log and what I can do in my own project for logging, it was excellent!"}],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Critical Feedback (Section 4 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>If there is any critical feedback that you do not want associated with your team, but feel would be helpful to the team to know, please provide it here. It will appear in the private section of the evaluation forms for the group but it will be anonymized.</b>","Name":"critical feedback","Options":[],"Subquestions":[]}
            ,{"Type":"paragraph","Text":"<h3>Overall Evaluation (Section 5 of 5)</h3>","Name":"","Options":[],"Subquestions":[]}
            ,{"Type":"radio","Text":"<b>Please give your overall rating for the submission. This will be used to help us eventually decide whether the team passes and what level of achievement is obtained. This section is also private feedback and your team's name is anonymized when giving this feedback.</b>","Name":"overall rating","Options":[{"Value":"1","Display":"1 out of 4 stars. Likely to fail Orbital"},{"Value":"2","Display":"2 out of 4 stars. Sufficient to pass the beginner level (Vostok), maybe good enough for intermediate level (Project Gemini)"},{"Value":"3","Display":"3 out of 4 stars. Definitely intermediate level (Project Gemini). Maybe good enough for advanced level (Apollo 11)"},{"Value":"4","Display":"4 out of 4 stars. Definitely good enough for advanced level (Apollo 11)."},{"Value":"5","Display":"5 out of 4 stars. Wow! (Bonus point)"}],"Subquestions":[]}
        ]$$)
    ;

    -- Feedback
    INSERT INTO periods (cohort, stage, start_at, end_at)
    VALUES
        (
            var_cohort
            ,'feedback'
            ,DATE(NOW() + '12 weeks'::INTERVAL)
            ,DATE(NOW() + '14 weeks'::INTERVAL) + '23 hours 59 minutes'::INTERVAL
        )
    RETURNING period_id INTO var_period_id
    ;
    -- Feedback Form
    INSERT INTO forms (period_id, questions)
    VALUES
        (var_period_id, $$[
            {"Type":"radio","Text":"<b>Your rating of this team/adviser's evaluations on your project, averaged all Evaluation Milestones that they completed for you.</b>\n<div class=\"f6 gray\">This only evaluates your peers' evaluation and comments on your project -- not their project's video, README or log. This will be used to help us to decide whether those teams achieve a higher than beginner (Vostok) level of achievement. This section will only be anonymously viewable to each target group.</div>","Name":"adviser rating","Options":[{"Value":"1","Display":"1 out of 4 stars. Did not meet peer evaluation requirements -- missing evaluations or did the bare minimum."},{"Value":"2","Display":"2 out of 4 stars. Met the minimum requirements for evaluations. Sufficient to pass the beginner level (Vostok)."},{"Value":"3","Display":"3 out of 4 stars. Exceeded the minimal requirements for feedback, with occasionally useful feedback given for my project. Meets intermediate level (Gemini), maybe good enough for the advanced level (Apollo 11)."},{"Value":"4","Display":"4 out of 4 stars. Often comprehensive and useful feedback on my project. Good enough for the advanced level (Apollo 11)."},{"Value":"5","Display":"5 out of 4 stars. Wow! (Bonus point) Consistently comprehensive and constructive feedback for all three Evaluations. We valued their contributions to my project. Definitely good enough for the advanced level (Apollo 11)."}],"Subquestions":[]}
            ,{"Type":"long text","Text":"<b>Feedback and comments on this team / adviser (if any).</b>","Name":"log","Options":[],"SubquestionAnswers":null,"Answer":null}
        ]$$)
    ;

END $forms$;
COMMIT;
