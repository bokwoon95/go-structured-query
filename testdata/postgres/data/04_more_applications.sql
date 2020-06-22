-- pg_temp.make_application is a temporary function (will not be persisted in
-- the database) that combines the app.idempotent_create_application
-- function and the app.join_application function to create a full application
-- with two applicants.
CREATE FUNCTION pg_temp.make_application(
    arg_user1_name TEXT, arg_user1_email TEXT,
    arg_user2_name TEXT, arg_user2_email TEXT,
    arg_team_name TEXT
) RETURNS TABLE (application_id INT, applicant1_user_id INT, applicant2_user_id INT) AS $$ DECLARE
    var_magicstring TEXT;
    var_application_id INT;
    var_user_id_1 INT;
    var_user_id_2 INT;
BEGIN
    SELECT _magicstring, _user_id
    INTO var_magicstring, var_user_id_1
    FROM app.idempotent_create_application(arg_user1_name, arg_user1_email)
    ;
    SELECT _application_id, _user_id
    INTO var_application_id, var_user_id_2
    FROM app.join_application(arg_user2_name, arg_user2_email, var_magicstring)
    ;
    RETURN QUERY SELECT var_application_id, var_user_id_1, var_user_id_2;
END $$ LANGUAGE 'plpgsql';

SELECT COUNT(*) AS "teams created" FROM (
    SELECT * FROM pg_temp.make_application(
        'Bonzi Buddy', 'bonzibuddy@email.com',
        'Rover The Dog', 'roverthedog@email.com',
        'Rover and Bonzi'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Harlen Atkins', 'harlenatkins@email.com2',
        'Rumaysa Estrada', 'rumaysaestrada@email.com2',
        'Harlen and Rumaysa'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Linda Mclellan', 'lindamclellan@email.com3',
        'Idris Parker', 'idrisparker@email.com3',
        'Linda and Idris'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Dylon Mccall', 'dylonmccall@email.com',
        'Jamal Lynch', 'jamallynch@email.com',
        'Dylon and Jamal'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Mara Lowery', 'maralowery@email.com',
        'Cayson Howells', 'caysonhowells@email.com',
        'Mara and Cayson'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Aurelia Forbes', 'aureliaforbes@email.com',
        'Rabia Mckenna', 'rabiamckenna@email.com',
        'Aurelia and Rabia'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Tommy Stubbs', 'tommystubbs@email.com',
        'Gilbert Clemons', 'gilbertclemons@email.com',
        'Tommy and Gilbert'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Enya Abbott', 'enyaabbott@email.com',
        'Cathal Carpenter', 'cathalcarpenter@email.com',
        'Enya and Cathal'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Lillie Olsen', 'lillieolsen@email.com',
        'Luciano Henry', 'lucianohenry@email.com',
        'Lillie and Luciano'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Elize Newman', 'elizenewman@email.com',
        'Justin Townsend', 'justintownsend@email.com',
        'Elize and Justin'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Whitney Tierney', 'whitneytierney@email.com',
        'Danica Dotson', 'danicadotson@email.com',
        'Whitney and Danica'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Bentley Downes', 'bentleydownes@email.com',
        'Fahim Sullivan', 'fahimsullivan@email.com',
        'Bentley and Fahim'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Asha Bentley', 'ashabentley@email.com',
        'Nadeem Senior', 'nadeemsenior@email.com',
        'Asha and Nadeem'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Isla-Mae Mcghee', 'isla-maemcghee@email.com',
        'Dora Frank', 'dorafrank@email.com',
        'Isla and Dora'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Ayden Ferreira', 'aydenferreira@email.com',
        'Dante Dunne', 'dantedunne@email.com',
        'Ayden and Dante'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Emanuel Cassidy', 'emanuelcassidy@email.com',
        'Mea Mcintyre', 'meamcintyre@email.com',
        'Emanuel and Mea'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Wilbur Phelps', 'wilburphelps@email.com',
        'Samiha Frederick', 'samihafrederick@email.com',
        'Wilbur and Samiha'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Kabir Bloggs', 'kabirbloggs@email.com',
        'Paige Patton', 'paigepatton@email.com',
        'Kabir and Paige'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Priya Knight', 'priyaknight@email.com',
        'Melody Bailey', 'melodybailey@email.com',
        'Priya and Melody'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Ravi West', 'raviwest@email.com',
        'Declan Hampton', 'declanhampton@email.com',
        'Ravi and Declan'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Amy-Leigh Woodley', 'Amy-leighwoodley@email.com',
        'Maegan Blackwell', 'maeganblackwell@email.com',
        'Leigh and Maegan'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Lianne Wilks', 'liannewilks@email.com',
        'Ibrar Maynard', 'ibrarmaynard@email.com',
        'Lianne and Ibrar'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Jayde Morales', 'jaydemorales@email.com',
        'Madihah Benjamin', 'madihahbenjamin@email.com',
        'Jayde and Madihah'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Farrell Yates', 'farrellyates@email.com',
        'Yousef Ewing', 'yousefewing@email.com',
        'Farrell and Yousef'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Arabella Camacho', 'arabellacamacho@email.com',
        'Izzy Keenan', 'izzykeenan@email.com',
        'Arabella and Izzy'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Korben Solis', 'korbensolis@email.com',
        'Viola Ventura', 'violaventura@email.com',
        'Korben and Viola'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Adam Driscoll', 'adamdriscoll@email.com',
        'Judy Goddard', 'judygoddard@email.com',
        'Adam and Judy'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Larissa Porter', 'larissaporter@email.com',
        'Ronnie Petty', 'ronniepetty@email.com',
        'Larissa and Ronnie'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Dianne Day', 'dianneday@email.com',
        'Star Ingram', 'staringram@email.com',
        'Dianne and Star'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Ceara Jensen', 'cearajensen@email.com',
        'Bobby Beaumont', 'bobbybeaumont@email.com',
        'Ceara and Bobby'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Rea Chan', 'reachan@email.com',
        'Helena Cooley', 'helenacooley@email.com',
        'Rea and Helena'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Ayush Torres', 'ayushtorres@email.com',
        'Audrey Peel', 'audreypeel@email.com',
        'Ayush and Audrey'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Macy Roman', 'macyroman@email.com',
        'Leja Talbot', 'lejatalbot@email.com',
        'Macy and Leja'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Devonte Wicks', 'devontewicks@email.com',
        'Shereen Miles', 'shereenmiles@email.com',
        'Devonte and Shereen'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Anika Lake', 'anikalake@email.com',
        'Shabaz Emerson', 'shabazemerson@email.com',
        'Anika and Shabaz'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Amelia-Rose Curtis', 'amelia-rosecurtis@email.com',
        'Brady Becker', 'bradybecker@email.com',
        'Ameliaand Brady'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Lily-Anne Macdonald', 'lily-annemacdonald@email.com',
        'Philip Fuller', 'philipfuller@email.com',
        'Lilyand Philip'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Harleigh Sims', 'harleighsims@email.com',
        'Tashan Lowry', 'tashanlowry@email.com',
        'Harleigh and Tashan'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Jaylen Wormald', 'jaylenwormald@email.com',
        'Hawwa Huffman', 'hawwahuffman@email.com',
        'Jaylen and Hawwa'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Luciana Coleman', 'lucianacoleman@email.com',
        'Sarina Cowan', 'sarinacowan@email.com',
        'Luciana and Sarina'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Charles Hendricks', 'charleshendricks@email.com',
        'Wallace Whittle', 'wallacewhittle@email.com',
        'Charles and Wallace'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Ashlee Powell', 'ashleepowell@email.com',
        'Leigha Poole', 'leighapoole@email.com',
        'Ashlee and Leigha'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Nikodem Goodwin', 'nikodemgoodwin@email.com',
        'Jordi Trujillo', 'jorditrujillo@email.com',
        'Nikodem and Jordi'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Connie Marshall', 'conniemarshall@email.com',
        'Suhail Dixon', 'suhaildixon@email.com',
        'Connie and Suhail'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Aj Wynn', 'ajwynn@email.com',
        'Abdur-Rahman Lin', 'Abdur-rahmanlin@email.com',
        'Aj and Abdur-Rahman'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Bethan O''Quinn', 'quinnemail.com',
        'Denny Johns', 'dennyjohns@email.com',
        'Bethan and Denny'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Zahra Mansell', 'zahramansell@email.com',
        'Sanaya Sheridan', 'sanayasheridan@email.com',
        'Zahra and Sanaya'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Neive Golden', 'neivegolden@email.com',
        'Alana Whyte', 'alanawhyte@email.com',
        'Neive and Alana'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Kwame Neville', 'kwameneville@email.com',
        'Lexi Perry', 'lexiperry@email.com',
        'Kwame and Lexi'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Jibril Clay', 'jibrilclay@email.com',
        'Gia Edmonds', 'giaedmonds@email.com',
        'Jibril and Gia'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Emily Kennedy', 'emilykennedy@email.com',
        'Hendrix Brady', 'hendrixbrady@email.com',
        'Emily and Hendrix'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Farzana Figueroa', 'farzanafigueroa@email.com',
        'Lia Riggs', 'liariggs@email.com',
        'Farzana and Lia'
    )
    UNION SELECT * FROM pg_temp.make_application(
        'Lyle Fernandez', 'lylefernandez@email.com',
        'Khushi Bray', 'khushibray@email.com',
        'Lyle and Khushi'
    )
    ORDER BY application_id ASC
) AS temp;
