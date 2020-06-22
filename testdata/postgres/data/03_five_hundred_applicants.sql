BEGIN;
DO $applicants$ DECLARE
    var_magicstring TEXT;
    var_application_id INT;
    var_user_id_1 INT;
    var_user_id_2 INT;
    var_user_role_id_1 INT;
    var_user_role_id_2 INT;
BEGIN

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Hilde Manigault', 'hildemanigault@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Tyron Speaks ','tyronspeaks@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ashlyn Storer ','ashlynstorer@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Timothy Sitzes ','timothysitzes@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Bruce Dave ','brucedave@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Travis Vanwinkle ','travisvanwinkle@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Damion Snoddy ','damionsnoddy@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Shenna Entrekin ','shennaentrekin@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Gianna Jurgensen ','giannajurgensen@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rosalie Struck ','rosaliestruck@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lowell Betts ','lowellbetts@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Zina Devaney ','zinadevaney@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lynda Feely ','lyndafeely@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Viola Cottingham ','violacottingham@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Adah Lupton ','adahlupton@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Arleen Burgett ','arleenburgett@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Renea Mcelfresh ','reneamcelfresh@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Candi Dorgan ','candidorgan@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jacki Manross ','jackimanross@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Ellsworth Haydon ','ellsworthhaydon@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Deloise Lubinski ','deloiselubinski@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Nicolasa Bourne ','nicolasabourne@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Erich Trost ','erichtrost@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Britta Brand ','brittabrand@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Melisa Artiaga ','melisaartiaga@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Judson Fuentez ','judsonfuentez@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jazmine Timko ','jazminetimko@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Oliver Wurster ','oliverwurster@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dia Baden ','diabaden@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Laurence Waldrip ','laurencewaldrip@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ta Levay ','talevay@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Creola Scudder ','creolascudder@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Julie Navarro ','julienavarro@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Samella Hunsberger ','samellahunsberger@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Marguerita Turley ','margueritaturley@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Arline Tienda ','arlinetienda@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nicol Barak ','nicolbarak@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Corrinne Glueck ','corrinneglueck@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ranee Atwater ','raneeatwater@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Jillian Rathburn ','jillianrathburn@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Shiloh Baron ','shilohbaron@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Randolph Bridgeman ','randolphbridgeman@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Shakira Barrs ','shakirabarrs@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Nancie Meads ','nanciemeads@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jerilyn Greaver ','jerilyngreaver@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Zachery Buonocore ','zacherybuonocore@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Molly Haag ','mollyhaag@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Shane Clewis ','shaneclewis@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Muoi Doyel ','muoidoyel@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Moon Arizmendi ','moonarizmendi@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Penni Casimir ','pennicasimir@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Tyler Craft ','tylercraft@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jeremiah Clever ','jeremiahclever@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Tiny Hewey ','tinyhewey@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dani Biggers ','danibiggers@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Chery Leleux ','cheryleleux@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Blake Doss ','blakedoss@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rosalva Heater ','rosalvaheater@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Earnestine Olmos ','earnestineolmos@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Benjamin Bias ','benjaminbias@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Gwendolyn Mccants ','gwendolynmccants@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Wyatt Laudenslager ','wyattlaudenslager@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Raylene Fleck ','raylenefleck@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sari Drinnon ','saridrinnon@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Elinor Tilly ','elinortilly@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Augustine Mitchel ','augustinemitchel@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Raina Lokey ','rainalokey@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Maurice Nay ','mauricenay@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Shin Brzozowski ','shinbrzozowski@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Ana Noonkester ','ananoonkester@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Fredrick Zielke ','fredrickzielke@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Marcelina Gabriel ','marcelinagabriel@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Karmen Shawn ','karmenshawn@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Larae Sally ','laraesally@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Chieko Beauchemin ','chiekobeauchemin@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Shelba Bentley ','shelbabentley@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Logan Croley ','logancroley@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Candyce Wedeking ','candycewedeking@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Wenona Lindgren ','wenonalindgren@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Twana Eckstein ','twanaeckstein@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Cathryn Stellmacher ','cathrynstellmacher@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Manuela Topping ','manuelatopping@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Laurice Ochs ','lauriceochs@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Melissia Mauch ','melissiamauch@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jerri Bilderback ','jerribilderback@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Leonida Hilgefort ','leonidahilgefort@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Bambi Hazelwood ','bambihazelwood@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Janett Vasta ','janettvasta@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nereida Mattis ','nereidamattis@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Mana Sartor ','manasartor@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Missy Lopinto ','missylopinto@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Eartha Varnado ','earthavarnado@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Wilda Lanier ','wildalanier@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Kassandra Rybak ','kassandrarybak@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jacqulyn Hokanson ','jacqulynhokanson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lavera Sexton ','laverasexton@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lavern Fane ','lavernfane@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Michelina Egger ','michelinaegger@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Laquanda Pinkard ','laquandapinkard@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Malorie Muir ','maloriemuir@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Elias Hake ','eliashake@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Dawna Caison ','dawnacaison@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Aaron Helbig ','aaronhelbig@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lulu Appel ','luluappel@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Luann Cissell ','luanncissell@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Kathey Mickles ','katheymickles@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Rhea Araujo ','rheaaraujo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Dirk Hayner ','dirkhayner@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jamison Hambrick ','jamisonhambrick@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Deetta Cleek ','deettacleek@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sunny Engebretson ','sunnyengebretson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Phillis Poucher ','phillispoucher@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Yevette Barham ','yevettebarham@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Philomena Cambareri ','philomenacambareri@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nelida Frey ','nelidafrey@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Anh Nath ','anhnath@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Suellen Langevin ','suellenlangevin@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Melvin Margolin ','melvinmargolin@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Kim Grange ','kimgrange@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Odette Hiles ','odettehiles@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Charity Turgeon ','charityturgeon@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Bernardo Rhodes ','bernardorhodes@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Chang Hendrixson ','changhendrixson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Tomasa Chen ','tomasachen@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Niesha Fischetti ','nieshafischetti@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Antonetta Pasquale ','antonettapasquale@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sheridan Brassfield ','sheridanbrassfield@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Veronique Venezia ','veroniquevenezia@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ma Ghoston ','maghoston@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Teresa Pitzer ','teresapitzer@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Augusta Hannon ','augustahannon@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Denny Lessig ','dennylessig@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jacquelin Tilton ','jacquelintilton@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Keena Lebo ','keenalebo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Shaunta Maher ','shauntamaher@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Emelda Perfecto ','emeldaperfecto@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Emily Dee ','emilydee@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Anderson Salazar ','andersonsalazar@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Maximo Hougland ','maximohougland@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Joie Hollier ','joiehollier@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Barb Mcglothlen ','barbmcglothlen@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Madge Weed ','madgeweed@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Vennie Hertlein ','venniehertlein@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Ed Creighton ','edcreighton@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Moshe Daughtrey ','moshedaughtrey@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Florence Dehn ','florencedehn@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Mable Steinfeldt ','mablesteinfeldt@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Nancie Liverman ','nancieliverman@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sandee Spargo ','sandeespargo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Annamae Woodford ','annamaewoodford@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Willy Kuester ','willykuester@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Juana Gregorio ','juanagregorio@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jenette Mcintosh ','jenettemcintosh@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lenny Senegal ','lennysenegal@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Patience Tumlinson ','patiencetumlinson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Harlan Gunter ','harlangunter@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lizabeth Canova ','lizabethcanova@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Kala Ostler ','kalaostler@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Raymonde Fickett ','raymondefickett@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Elena Scheidler ','elenascheidler@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Salley Deer ','salleydeer@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Bennett Victorian ','bennettvictorian@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Tomasa Anselmo ','tomasaanselmo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Idalia Buda ','idaliabuda@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Marlene Marquez ','marlenemarquez@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rico Fausnaught ','ricofausnaught@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jeri Cyr ','jericyr@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Angele Ahearn ','angeleahearn@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Aja Nimmo ','ajanimmo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lee Villescas ','leevillescas@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Cordelia Depaz ','cordeliadepaz@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sharon Tirrell ','sharontirrell@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dorothea Kestner ','dorotheakestner@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Mabel Tedesco ','mabeltedesco@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Marie Nicastro ','marienicastro@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Eladia Edington ','eladiaedington@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Misti Lucius ','mistilucius@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Maudie Stallone ','maudiestallone@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lexie Ralston ','lexieralston@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sophia Sinegal ','sophiasinegal@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Caren Creek ','carencreek@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Nicholas Gendreau ','nicholasgendreau@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Torri Sumrall ','torrisumrall@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Georgie Ratliff ','georgieratliff@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Leola Goosby ','leolagoosby@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Priscilla Ennis ','priscillaennis@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Tyson Ryman ','tysonryman@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Hope Hathorn ','hopehathorn@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Adriana Mendieta ','adrianamendieta@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Chas Cowen ','chascowen@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Chester Labadie ','chesterlabadie@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Hugo Rainer ','hugorainer@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Carlita Providence ','carlitaprovidence@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Gracie Shomo ','gracieshomo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Kandice Olea ','kandiceolea@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Aliza Anker ','alizaanker@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Joana Gohr ','joanagohr@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Gayle Borden ','gayleborden@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Vita Phillip ','vitaphillip@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Dacia Rotondo ','daciarotondo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Beatris Epping ','beatrisepping@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Giovanni Axley ','giovanniaxley@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dayle Depew ','dayledepew@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Josue Kight ','josuekight@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Fidela Haar ','fidelahaar@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Marlon Hardison ','marlonhardison@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Rufus Toothaker ','rufustoothaker@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Jae Lappin ','jaelappin@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Angella Killeen ','angellakilleen@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Ninfa Grindle ','ninfagrindle@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Melda Dansereau ','meldadansereau@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Cheryl Dubose ','cheryldubose@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nadia Antonelli ','nadiaantonelli@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Julian Fester ','julianfester@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dian Jarrett ','dianjarrett@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Daphine Rhein ','daphinerhein@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nolan Vanatta ','nolanvanatta@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Roseanna Milburn ','roseannamilburn@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Gilma Pursell ','gilmapursell@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Monte Mayville ','montemayville@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Kelle Kilduff ','kellekilduff@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('America Zacharias ','americazacharias@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Britt Numbers ','brittnumbers@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Willa Michael ','willamichael@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Alyce Sinegal ','alycesinegal@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rebbecca Mcwilliam ','rebbeccamcwilliam@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sumiko Prevo ','sumikoprevo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Willia Frase ','williafrase@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Gianna Cadena ','giannacadena@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Alvina Schild ','alvinaschild@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Doug Gehl ','douggehl@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Hedwig Kluck ','hedwigkluck@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Beulah Downs ','beulahdowns@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Cassy Mitchem ','cassymitchem@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Suk Spadafora ','sukspadafora@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Daniell Seigel ','daniellseigel@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Annita Wendel ','annitawendel@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Odilia Martino ','odiliamartino@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dania Kroeger ','daniakroeger@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Darren Mitchener ','darrenmitchener@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jacquelyn Brice ','jacquelynbrice@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Robena Kollman ','robenakollman@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Providencia Hedgpeth ','providenciahedgpeth@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Kyong Senior ','kyongsenior@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Rona Marek ','ronamarek@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Clifford Lefever ','cliffordlefever@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Billy Cambra ','billycambra@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sana Mooneyham ','sanamooneyham@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Branda Calero ','brandacalero@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Raylene Hersey ','raylenehersey@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Karisa Waites ','karisawaites@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Betty Huston ','bettyhuston@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dell Minder ','dellminder@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Clara Fulp ','clarafulp@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Treena Presti ','treenapresti@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Nilsa Lamberton ','nilsalamberton@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Camille Grippo ','camillegrippo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Renetta Westlund ','renettawestlund@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ines Samples ','inessamples@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Stevie Greenfield ','steviegreenfield@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Emiko Racca ','emikoracca@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lachelle Brogdon ','lachellebrogdon@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nora Triplett ','noratriplett@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Yuki Helmuth ','yukihelmuth@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Krista Shellhammer ','kristashellhammer@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Dennise Weise ','denniseweise@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Spencer Lichtman ','spencerlichtman@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Deon Vivian ','deonvivian@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lovie Riedel ','lovieriedel@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Ingrid Spagnola ','ingridspagnola@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Hiroko Latson ','hirokolatson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Shonna Prejean ','shonnaprejean@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Demetrice Ambrosino ','demetriceambrosino@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Alyssa Halperin ','alyssahalperin@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Johnna Mushrush ','johnnamushrush@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Jeannine Beckert ','jeanninebeckert@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Matt Downard ','mattdownard@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Selina Valdovinos ','selinavaldovinos@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Faith Etherton ','faithetherton@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Genevive Locher ','genevivelocher@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Walton Grissett ','waltongrissett@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sarai Lagasse ','sarailagasse@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Santiago Depew ','santiagodepew@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Norman Steeves ','normansteeves@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jani Youngs ','janiyoungs@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rossie Ingwersen ','rossieingwersen@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Quiana Strobl ','quianastrobl@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Catrina Okamura ','catrinaokamura@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Racheal Mix ','rachealmix@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Celine Laird ','celinelaird@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Belva Patnaude ','belvapatnaude@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Janna Hugo ','jannahugo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jermaine Kostka ','jermainekostka@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Pablo Litt ','pablolitt@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Keeley Tignor ','keeleytignor@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Leora Mclauchlin ','leoramclauchlin@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Darleen Schlager ','darleenschlager@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Ashely Ober ','ashelyober@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Judith Gile ','judithgile@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rebecca Daum ','rebeccadaum@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Shaunte Parras ','shaunteparras@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Ted Flor ','tedflor@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jovan Nickson ','jovannickson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Adolph Schupp ','adolphschupp@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Arlene Martindale ','arlenemartindale@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Venessa Hudon ','venessahudon@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Mariella Rottman ','mariellarottman@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Thomas Brien ','thomasbrien@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sherri Barrera ','sherribarrera@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Carrol Counce ','carrolcounce@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jackeline Daughtridge ','jackelinedaughtridge@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Victorina Farfan ','victorinafarfan@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Reginald Preuss ','reginaldpreuss@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Alejandra Echevarria ','alejandraechevarria@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jerry Hubbs ','jerryhubbs@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Terrance Wilton ','terrancewilton@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ciara Randolph ','ciararandolph@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Trinidad Averill ','trinidadaverill@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Melodi Witherington ','melodiwitherington@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Karri Gerhart ','karrigerhart@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Margarita Peterson ','margaritapeterson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Mauricio Rusch ','mauriciorusch@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Andra Greb ','andragreb@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Pamula Blades ','pamulablades@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jenna Kroner ','jennakroner@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Melania Romig ','melaniaromig@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Cassey Bezio ','casseybezio@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Carlos Pflug ','carlospflug@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Kacie Gulbranson ','kaciegulbranson@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Cathrine Ibrahim ','cathrineibrahim@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Emily Bashir ','emilybashir@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lorraine Landers ','lorrainelanders@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Divina Straus ','divinastraus@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lauran Gagliardo ','laurangagliardo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Bev Jiron ','bevjiron@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Almeta Sipe ','almetasipe@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sherry Learned ','sherrylearned@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Mckinley Bracamonte ','mckinleybracamonte@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ricki Duhon ','rickiduhon@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Jewell Darbonne ','jewelldarbonne@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Mabelle Verdin ','mabelleverdin@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Georgina Ordway ','georginaordway@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Cherie Delima ','cheriedelima@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Merlin Tillinghast ','merlintillinghast@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Gala Danziger ','galadanziger@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Barrett Pedro ','barrettpedro@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dexter Liner ','dexterliner@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Teresia Boor ','teresiaboor@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Stefani Saxon ','stefanisaxon@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Regina Gramlich ','reginagramlich@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Elodia Llanos ','elodiallanos@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Arleen Searles ','arleensearles@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jc Kier ','jckier@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Alyse Chalfant ','alysechalfant@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Bambi Geibel ','bambigeibel@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Evita Buechner ','evitabuechner@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Maribel Zayas ','maribelzayas@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Clinton Carter ','clintoncarter@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Arlyne Plewa ','arlyneplewa@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Chong Fobbs ','chongfobbs@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Raguel Banas ','raguelbanas@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sherwood Losey ','sherwoodlosey@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Kiley Hambrick ','kileyhambrick@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Astrid Cantrell ','astridcantrell@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Taisha Merriam ','taishamerriam@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lakeshia Nasser ','lakeshianasser@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Kecia Pavone ','keciapavone@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lourdes Mcarthur ','lourdesmcarthur@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Gussie Hogge ','gussiehogge@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rosalva Harstad ','rosalvaharstad@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Maya Dively ','mayadively@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Christine Heer ','christineheer@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Li Riggan ','liriggan@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Yen Maniscalco ','yenmaniscalco@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Miguelina Rapozo ','miguelinarapozo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Tiffiny Kiernan ','tiffinykiernan@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Alisha Carrillo ','alishacarrillo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Charlena Schroth ','charlenaschroth@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Raelene Natera ','raelenenatera@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Imogene Geisel ','imogenegeisel@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Beau Pavao ','beaupavao@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sunday Carwile ','sundaycarwile@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jimmy Brandt ','jimmybrandt@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Althea Blessing ','altheablessing@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dennise Magoon ','dennisemagoon@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Min Alberico ','minalberico@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lorean Mckeown ','loreanmckeown@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Jackelyn Frausto ','jackelynfrausto@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Ellyn Whitting ','ellynwhitting@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Celestine Peasley ','celestinepeasley@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Bao Croley ','baocroley@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Wm Stratton ','wmstratton@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Dalton Hoppes ','daltonhoppes@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Johnson Marsh ','johnsonmarsh@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Delinda Beardsley ','delindabeardsley@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Eleonora Gruber ','eleonoragruber@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Chan Dale ','chandale@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Phyllis Resler ','phyllisresler@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Alleen Tourville ','alleentourville@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Glynis Catalan ','glyniscatalan@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Felisha Bertsch ','felishabertsch@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Tynisha Claunch ','tynishaclaunch@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Christi Wiebe ','christiwiebe@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Randall Coplin ','randallcoplin@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Marquita Ebarb ','marquitaebarb@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Antwan Natoli ','antwannatoli@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Clorinda Kell ','clorindakell@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Elin Podesta ','elinpodesta@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Era Coaxum ','eracoaxum@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Hong Duff ','hongduff@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Raguel Rux ','raguelrux@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Rochelle Derrick ','rochellederrick@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Olene Judd ','olenejudd@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Jocelyn Fesler ','jocelynfesler@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nelly Hazelton ','nellyhazelton@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Marilou Moll ','mariloumoll@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Alida Lazenby ','alidalazenby@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lizette Caraveo ','lizettecaraveo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Barrie Middlebrooks ','barriemiddlebrooks@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Cordia Barretta ','cordiabarretta@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Pei Bogen ','peibogen@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Elvin Loo ','elvinloo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Elanor Snedden ','elanorsnedden@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Reynaldo Sibrian ','reynaldosibrian@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sharan Wycoff ','sharanwycoff@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Everette Kai ','everettekai@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lauran Reames ','lauranreames@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Sixta Harry ','sixtaharry@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Thuy Corlett ','thuycorlett@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Trinidad Warfield ','trinidadwarfield@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Jesusita Malm ','jesusitamalm@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Euna Mcafee ','eunamcafee@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sabra Armagost ','sabraarmagost@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Grisel Mcintosh ','griselmcintosh@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Nell Manuelito ','nellmanuelito@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Dorris Eilers ','dorriseilers@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Reinaldo Aylward ','reinaldoaylward@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Mica Pedroza ','micapedroza@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Vanita Simoes ','vanitasimoes@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Clotilde Crawford ','clotildecrawford@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Brandee Kalman ','brandeekalman@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Shantae Kennemer ','shantaekennemer@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Shara Armistead ','sharaarmistead@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Easter Dane ','easterdane@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Kirk Knight ','kirkknight@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Alona Hasting ','alonahasting@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lona Cleary ','lonacleary@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Clark Cashion ','clarkcashion@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Joel Karney ','joelkarney@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Anika Romo ','anikaromo@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Lory Blom ','loryblom@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Gianna Acklin ','giannaacklin@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Joana Schomer ','joanaschomer@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Fernanda Holben ','fernandaholben@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Mickie Drown ','mickiedrown@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Zachary Roberti ','zacharyroberti@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Miyoko Izzo ','miyokoizzo@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Guillermo Belding ','guillermobelding@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Britt Dineen ','brittdineen@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Dawne Ranney ','dawneranney@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Almeda Beauford ','almedabeauford@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Altagracia Fucci ','altagraciafucci@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Zina Manzella ','zinamanzella@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Danny Romanik ','dannyromanik@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Apolonia Luthy ','apolonialuthy@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Danika Abee ','danikaabee@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Summer Dinwiddie ','summerdinwiddie@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Luanne Downing ','luannedowning@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Olga Wirt ','olgawirt@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Trena Prow ','trenaprow@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Veronica League ','veronicaleague@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Lona Wasmund ','lonawasmund@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Melodie Faith ','melodiefaith@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Dorothy Villafane ','dorothyvillafane@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Catina Orduna ','catinaorduna@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Wally Zurita ','wallyzurita@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Sheryl Mcvay ','sherylmcvay@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Loren Garand ','lorengarand@gmail.com', var_magicstring);

    SELECT _user_id, _user_role_id, _application_id, _magicstring INTO var_user_id_1, var_user_role_id_1, var_application_id, var_magicstring FROM app.idempotent_create_application('Percy Gover ','percygover@gmail.com');
    SELECT _user_id, _user_role_id INTO var_user_id_2, var_user_role_id_2 FROM app.join_application('Gerald Eifert ','geraldeifert@gmail.com', var_magicstring);

END $applicants$;
COMMIT;
