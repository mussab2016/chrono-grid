BEGIN;

-- 1. Reset all tables cleanly from zero
TRUNCATE TABLE
    timetables,
    course_offerings,
    student_groups,
    academic_terms,
    academic_years,
    subjects,
    academic_levels,
    rooms,
    room_types,
    professors
    RESTART IDENTITY CASCADE;

-- 2. Room Types (code, description)
INSERT INTO room_types (id, code, description) VALUES
                                                   (1, 'AMPHI', 'Amphithéâtre pour cours magistraux'),
                                                   (2, 'TD', 'Salle de travaux dirigés'),
                                                   (3, 'TP', 'Laboratoire informatique / TP');

-- 3. Physical Rooms & Auditoriums (room_type_id, code, description)
INSERT INTO rooms (room_type_id, code, description) VALUES
                                                        (1, 'A3', 'Amphithéâtre 3'),
                                                        (1, 'A4', 'Amphithéâtre 4'),
                                                        (1, 'A5', 'Amphithéâtre 5'),
                                                        (3, 'D9', 'Laboratoire Informatique D9'),
                                                        (3, 'D10', 'Laboratoire Informatique D10'),
                                                        (2, 'D6', 'Salle TD D6'),
                                                        (2, 'D7', 'Salle TD D7'),
                                                        (2, 'D8', 'Salle TD D8'),
                                                        (2, 'D11', 'Salle TD D11'),
                                                        (2, 'D12', 'Salle TD D12'),
                                                        (2, 'D13', 'Salle TD D13'),
                                                        (2, 'D14', 'Salle TD D14'),
                                                        (2, 'D15', 'Salle TD D15'),
                                                        (2, 'D16', 'Salle TD D16'),
                                                        (2, 'D17', 'Salle TD D17'),
                                                        (2, 'D18', 'Salle TD D18'),
                                                        (2, 'D19', 'Salle TD D19'),
                                                        (2, 'D20', 'Salle TD D20'),
                                                        (2, 'D21', 'Salle TD D21'),
                                                        (2, 'F11', 'Salle TD F11'),
                                                        (2, 'F16', 'Salle TD F16'),
                                                        (2, 'F17', 'Salle TD F17'),
                                                        (2, 'F18', 'Salle TD F18');

-- 4. Academic Levels (id, name, code)
INSERT INTO academic_levels (id, name, code) VALUES
                                                 (1, '1ère Année Licence', 'L1'),
                                                 (2, '2ème Année Licence', 'L2'),
                                                 (3, '3ème Année Licence', 'L3'),
                                                 (4, '1ère Année Master', 'M1'),
                                                 (5, '2ème Année Master', 'M2');

-- 5. Academic Year & Academic Terms
INSERT INTO academic_years (id, year_label) VALUES (1, '2026-2027');

INSERT INTO academic_terms (id, academic_year_id, semester, label) VALUES
                                                                       (1, 1, 1, 'S1'),
                                                                       (2, 1, 2, 'S2');

-- 6. Student Groups (42 groups bound to academic_year_id = 1)
-- L1: 15 groups (G1..G15)
INSERT INTO student_groups (academic_year_id, level_id, code)
SELECT 1, 1, 'G' || g FROM generate_series(1, 15) AS g;

-- L2: 9 groups (G1..G9)
INSERT INTO student_groups (academic_year_id, level_id, code)
SELECT 1, 2, 'G' || g FROM generate_series(1, 9) AS g;

-- L3: 6 specialty groups (Analyse: A1-A2, Proba: P1-P2, Algèbre: G1-G2)
INSERT INTO student_groups (academic_year_id, level_id, code) VALUES
                                                                  (1, 3, 'A1'), (1, 3, 'A2'), (1, 3, 'P1'), (1, 3, 'P2'), (1, 3, 'G1'), (1, 3, 'G2');

-- M1: 6 specialty cohorts (modeled as groups to support shared cross-track courses)
INSERT INTO student_groups (academic_year_id, level_id, code) VALUES
                                                                  (1, 4, 'EDP'), (1, 4, 'AMA'), (1, 4, 'PSA'), (1, 4, 'TASD'), (1, 4, 'MAEF'), (1, 4, 'SA');

-- M2: 6 specialty cohorts
INSERT INTO student_groups (academic_year_id, level_id, code) VALUES
                                                                  (1, 5, 'EDP'), (1, 5, 'MA'), (1, 5, 'SPA'), (1, 5, 'TASD'), (1, 5, 'MAEF'), (1, 5, 'SA');

-- 7. Subjects (name, code)
INSERT INTO subjects (name, code) VALUES
-- L1
('Analyse 1', 'L1-ANA'),
('Algèbre 1', 'L1-ALG'),
('Algorithmique et Structures de Données', 'L1-ALGO'),
('Statistique Descriptive et Probabilités', 'L1-STAT'),
('Economie Générale', 'L1-ECO'),
('Terminologie Scientifique', 'L1-TERM'),
('Histoire des Mathématiques', 'L1-HIST'),
-- L2
('Analyse 3', 'L2-ANA'),
('Algèbre 3', 'L2-ALG'),
('Topologie', 'L2-TOP'),
('Analyse Numérique', 'L2-ANUM'),
('Programmation MATLAB', 'L2-MAT'),
('Langue Anglaise', 'L2-ANG'),
('Histoire des Mathématiques L2', 'L2-HIST'),
-- L3
('Analyse Complexe', 'L3-ACOMP'),
('Equations Différentielles et Systèmes Dynamiques', 'L3-ESD'),
('Mesure et Intégration', 'L3-MINT'),
('Statistique Inférentielle', 'L3-SI'),
('Géométrie Affine et Euclidienne', 'L3-GAE'),
('Analyse Convexe', 'L3-AC'),
('Mathématiques de l''Ingénieur', 'L3-MI'),
('Initiation à la Didactique des Mathématiques', 'L3-IDM'),
-- M1
('Analyse Fonctionnelle 1', 'AF1'),
('Équations Différentielles Ordinaires', 'EDO'),
('Distributions et Analyse de Fourier 1', 'DAF'),
('Analyse Numérique Matricielle', 'ANM'),
('Éthique de l''Intelligence Artificielle', 'EIA'),
('Modélisation Mathématique', 'MM'),
('Calcul Différentiel', 'CD'),
('Analyse Convexe M1', 'AC'),
('Probabilités Avancées', 'PA'),
('Complément de la Théorie de la Mesure', 'CTM'),
('Statistique Inférentielle 1', 'SI1'),
('Analyse des Données', 'AD'),
('Programmation sous Logiciel R/Python', 'PLRP'),
('Statistiques des Valeurs Extrêmes', 'SVE'),
('Topologie Algébrique 1', 'TA1'),
('Théorie des Probabilités 1', 'TP1'),
('Méthodes Numériques 1', 'MN1'),
('Initiation à la Didactique de Mathématique', 'IDM'),
('Equations aux Dérivées Partielles', 'EDP'),
('Economie Générale M1', 'EG'),
('Statistique Mathématique', 'SM'),
('Statistique Appliquée au Domaine de la Santé', 'SADS'),
('Processus et leurs Statistiques', 'PS'),
('Techniques et Méthodes Descriptives de Datamining', 'TMDD'),
('Expression Écrite et Orale', 'EEO'),
-- M2
('Problèmes Mal-Posés', 'PMP'),
('Inégalités Énergétiques', 'IE'),
('Problèmes d''Évolution', 'PE'),
('Théorie de l''Approximation', 'TA'),
('Analyse et Contrôle des Systèmes Localisés', 'ACSL'),
('Analyse et Contrôle des Systèmes Distribués', 'ACSD'),
('Calcul Stochastique', 'CST'),
('Séries Temporelles', 'STE'),
('Données Censurées', 'DC'),
('Chaînes de Markov et Marches Aléatoires', 'CMMA'),
('Calcul Scientifique', 'CS'),
('Aspects Topologiques et Algébriques des Morphismes Multivoques', 'ATAMM'),
('Stabilité et Instabilité des Systèmes Dynamiques', 'SISD'),
('Théorie des Treillis', 'TT'),
('Modèles Stochastiques', 'MS'),
('Bases de Données et Modélisation avec R', 'BDMR'),
('Modélisation et Étude de Cas', 'MEC'),
('Plans d''Expérience', 'PEX'),
('Réseaux de Neurones 2', 'RN2'),
('Théorie de l''Apprentissage', 'TAP'),
('Techniques et Méthodes Prédictives de Datamining', 'TMPD'),
('Outils de Calcul MATLAB', 'MAT-M2'),
('Rédaction Scientifique LaTeX', 'LAT-M2');

-- 8. Professors (55 Professors with department = 'Mathématiques')
INSERT INTO professors (fullname, department) VALUES
                                                  ('Dr. Benali Ahmed (Analyse 1)', 'Mathématiques'),
                                                  ('Dr. Mansouri Karima (Analyse 2)', 'Mathématiques'),
                                                  ('Dr. Zerrouki Samira (Analyse 3)', 'Mathématiques'),
                                                  ('Dr. Khelifi Rachid (Analyse 4)', 'Mathématiques'),
                                                  ('Dr. Bouzid Nabila (Analyse 5)', 'Mathématiques'),
                                                  ('Dr. Hamidi Tarek (Analyse 6)', 'Mathématiques'),
                                                  ('Dr. Merah Sofiane (Analyse 7)', 'Mathématiques'),
                                                  ('Dr. Saidi Yasmine (Analyse 8)', 'Mathématiques'),
                                                  ('Dr. Belkacem Mourad (Analyse 9)', 'Mathématiques'),
                                                  ('Dr. Amrani Salima (Analyse 10)', 'Mathématiques'),
                                                  ('Dr. Boudiaf Farid (Algèbre 1)', 'Mathématiques'),
                                                  ('Dr. Chaib Nadia (Algèbre 2)', 'Mathématiques'),
                                                  ('Dr. Djebbar Omar (Algèbre 3)', 'Mathématiques'),
                                                  ('Dr. Ferhat Leila (Algèbre 4)', 'Mathématiques'),
                                                  ('Dr. Guendouz Adel (Algèbre 5)', 'Mathématiques'),
                                                  ('Dr. Haddad Meriem (Algèbre 6)', 'Mathématiques'),
                                                  ('Dr. Slimani Walid (Algèbre 7)', 'Mathématiques'),
                                                  ('Dr. Touati Amel (Algèbre 8)', 'Mathématiques'),
                                                  ('Dr. Meziane Bilal (Algo-Info 1)', 'Informatique & Math'),
                                                  ('Dr. Rahmani Houda (Algo-Info 2)', 'Informatique & Math'),
                                                  ('Dr. Cherif Youcef (Algo-Info 3)', 'Informatique & Math'),
                                                  ('Dr. Larbi Sonia (Algo-Info 4)', 'Informatique & Math'),
                                                  ('Dr. Messaoudi Karim (Algo-Info 5)', 'Informatique & Math'),
                                                  ('Dr. Benmalek Rania (Algo-Info 6)', 'Informatique & Math'),
                                                  ('Dr. Ziani Hamza (Algo-Info 7)', 'Informatique & Math'),
                                                  ('Dr. Kadri Nesrine (Algo-Info 8)', 'Informatique & Math'),
                                                  ('Dr. Yahiaoui Samir (Proba-Stat 1)', 'Mathématiques'),
                                                  ('Dr. Bouazza Imene (Proba-Stat 2)', 'Mathématiques'),
                                                  ('Dr. Chekroun Nadir (Proba-Stat 3)', 'Mathématiques'),
                                                  ('Dr. Derbal Asma (Proba-Stat 4)', 'Mathématiques'),
                                                  ('Dr. Fodil Mehdi (Proba-Stat 5)', 'Mathématiques'),
                                                  ('Dr. Gharbi Souad (Proba-Stat 6)', 'Mathématiques'),
                                                  ('Dr. Henni Lotfi (Proba-Stat 7)', 'Mathématiques'),
                                                  ('Dr. Lounis Djamila (Proba-Stat 8)', 'Mathématiques'),
                                                  ('Dr. Mokhtari Salim (Topologie 1)', 'Mathématiques'),
                                                  ('Dr. Naceri Ryma (Topologie 2)', 'Mathématiques'),
                                                  ('Dr. Ouali Hakim (Topologie 3)', 'Mathématiques'),
                                                  ('Dr. Rahal Sabrina (Topologie 4)', 'Mathématiques'),
                                                  ('Dr. Sellami Fouad (Topologie 5)', 'Mathématiques'),
                                                  ('Dr. Tali Wafa (EDP-Num 1)', 'Mathématiques'),
                                                  ('Dr. Zitouni Kamel (EDP-Num 2)', 'Mathématiques'),
                                                  ('Dr. Aissaoui Hind (EDP-Num 3)', 'Mathématiques'),
                                                  ('Dr. Brahimi Reda (EDP-Num 4)', 'Mathématiques'),
                                                  ('Dr. Chergui Nawel (EDP-Num 5)', 'Mathématiques'),
                                                  ('Dr. Dahmani Ilyes (EDP-Num 6)', 'Mathématiques'),
                                                  ('Dr. Fergani Lamia (EDP-Num 7)', 'Mathématiques'),
                                                  ('Dr. Gasmi Tariq (EDP-Num 8)', 'Mathématiques'),
                                                  ('Dr. Belhadj Salim (Economie 1)', 'Sciences Économiques'),
                                                  ('Dr. Kaci Malika (Economie 2)', 'Sciences Économiques'),
                                                  ('Dr. Madani Anis (Economie 3)', 'Sciences Économiques'),
                                                  ('Dr. Oukil Sarah (Economie 4)', 'Sciences Économiques'),
                                                  ('Dr. Benahmed Nadia (Langues-Hist 1)', 'Langues & Culture'),
                                                  ('Dr. Chikhi Mourad (Langues-Hist 2)', 'Langues & Culture'),
                                                  ('Dr. Dridi Selma (Langues-Hist 3)', 'Langues & Culture'),
                                                  ('Dr. Hachemi Sami (Langues-Hist 4)', 'Langues & Culture');

-- 9. Generate Course Offerings for Term 1 (Max 2 TD groups/prof, Max 4 TP groups/prof)
DO $$
    DECLARE
        v_term_id INT := 1;
        v_uuid UUID;
        v_prof_id INT;
        v_subj_id INT;
        v_grp_id INT;
        i INT;
        sec INT;
        grp_code TEXT;
    BEGIN
        -- =========================================================================
        -- A. L1 (15 Groups = 5 Sections of 3 groups each)
        -- =========================================================================
        FOR sec IN 1..5 LOOP
                -- Analyse Lecture (2/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ANA';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Analyse ' || sec || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 2
                FROM student_groups WHERE level_id = 1 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Algèbre Lecture (2/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ALG';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algèbre ' || sec || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 2
                FROM student_groups WHERE level_id = 1 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Algo Lecture (1/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ALGO';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algo-Info ' || sec || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
                FROM student_groups WHERE level_id = 1 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Stat Lecture (1/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-STAT';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Proba-Stat ' || sec || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
                FROM student_groups WHERE level_id = 1 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Economie Lecture (1/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ECO';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Economie ' || ((sec - 1) % 4 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
                FROM student_groups WHERE level_id = 1 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Terminologie Lecture (1/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-TERM';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Langues-Hist 1)';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
                FROM student_groups WHERE level_id = 1 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));
            END LOOP;

        -- L1 Individual TDs (Max 2 groups/prof) & TPs (Max 4 groups/prof)
        FOR i IN 1..15 LOOP
                SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 1 AND code = 'G' || i;

                -- Analyse TD (2/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ANA';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Analyse ' || ((i - 1) / 2 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 2);

                -- Algèbre TD (1/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ALG';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algèbre ' || ((i - 1) / 2 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                -- Algo TD (1/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ALGO';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algo-Info ' || ((i - 1) / 2 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                -- Algo TP in Lab (1/wk, max 4 groups/prof)
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algo-Info ' || ((i - 1) / 4 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 3, 1);

                -- Stat TD (1/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-STAT';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Proba-Stat ' || ((i - 1) / 2 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                -- Economie TD (1/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L1-ECO';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Economie ' || ((i - 1) / 4 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);
            END LOOP;

        -- =========================================================================
        -- B. L2 (9 Groups = 3 Sections of 3 groups each)
        -- =========================================================================
        FOR sec IN 1..3 LOOP
                -- Analyse L2 Lecture (2/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-ANA';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Analyse ' || (sec + 7) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 2
                FROM student_groups WHERE level_id = 2 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Algèbre L2 Lecture (1/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-ALG';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algèbre ' || (sec + 5) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
                FROM student_groups WHERE level_id = 2 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Topologie Lecture (2/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-TOP';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Topologie ' || sec || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 2
                FROM student_groups WHERE level_id = 2 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Analyse Numérique Lecture (1/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-ANUM';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(EDP-Num ' || sec || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
                FROM student_groups WHERE level_id = 2 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));

                -- Anglais L2 Lecture (1/wk)
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-ANG';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Langues-Hist 2)';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
                FROM student_groups WHERE level_id = 2 AND code IN ('G'||((sec-1)*3+1), 'G'||((sec-1)*3+2), 'G'||((sec-1)*3+3));
            END LOOP;

        -- L2 Group TDs (Max 2 groups/prof) & TPs (Max 4 groups/prof)
        FOR i IN 1..9 LOOP
                SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 2 AND code = 'G' || i;

                -- Analyse L2 TD (2/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-ANA';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Analyse ' || ((i - 1) / 2 + 6) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 2);

                -- Algèbre L2 TD (1/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-ALG';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algèbre ' || ((i - 1) / 2 + 4) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                -- Topologie TD (1/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-TOP';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Topologie ' || ((i - 1) / 2 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                -- Analyse Numérique TD (1/wk)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-ANUM';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(EDP-Num ' || ((i - 1) / 2 + 1) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                -- Analyse Numérique TP in Lab (1/wk, max 4 groups/prof)
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(EDP-Num ' || ((i - 1) / 4 + 6) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 3, 1);

                -- MATLAB TP in Lab (1/wk, max 4 groups/prof)
                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L2-MAT';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Algo-Info ' || ((i - 1) / 4 + 6) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 3, 1);
            END LOOP;

        -- =========================================================================
        -- C. L3 (6 Groups: A1, A2, P1, P2, G1, G2)
        -- =========================================================================
        FOR i IN 1..5 LOOP
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = (ARRAY['L3-ACOMP','L3-MINT','L3-ESD','L3-SI','L3-GAE'])[i];
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE (ARRAY['%(Analyse 9)','%(Analyse 10)','%(EDP-Num 7)','%(Proba-Stat 8)','%(Algèbre 8)'])[i];
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 2
                FROM student_groups WHERE level_id = 3;
            END LOOP;

        FOR grp_code IN SELECT unnest(ARRAY['A1','A2','P1','P2','G1','G2']) LOOP
                SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 3 AND code = grp_code;

                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L3-ESD';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(EDP-Num ' || (CASE WHEN grp_code IN ('A1','A2') THEN 1 WHEN grp_code IN ('P1','P2') THEN 2 ELSE 3 END) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L3-AC';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Analyse ' || (CASE WHEN grp_code IN ('A1','A2') THEN 8 WHEN grp_code IN ('P1','P2') THEN 9 ELSE 10 END) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);

                SELECT id INTO v_subj_id FROM subjects WHERE code = 'L3-SI';
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Proba-Stat ' || (CASE WHEN grp_code IN ('A1','A2') THEN 6 WHEN grp_code IN ('P1','P2') THEN 7 ELSE 8 END) || ')';
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                VALUES (v_term_id, gen_random_uuid(), v_subj_id, v_prof_id, v_grp_id, 2, 1);
            END LOOP;

        -- =========================================================================
        -- D. M1 Specialty Cohorts (With Shared Cross-Track Lectures)
        -- =========================================================================
        -- Shared AF1 (EDP + AMA)
        v_uuid := gen_random_uuid();
        SELECT id INTO v_subj_id FROM subjects WHERE code = 'AF1';
        SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Analyse 1)';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 2
        FROM student_groups WHERE level_id = 4 AND code IN ('EDP', 'AMA');

        -- Shared DAF (EDP + AMA)
        v_uuid := gen_random_uuid();
        SELECT id INTO v_subj_id FROM subjects WHERE code = 'DAF';
        SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(Analyse 2)';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 2, 2
        FROM student_groups WHERE level_id = 4 AND code IN ('EDP', 'AMA');

        -- Shared MM (EDP + AMA + PSA)
        v_uuid := gen_random_uuid();
        SELECT id INTO v_subj_id FROM subjects WHERE code = 'MM';
        SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE '%(EDP-Num 3)';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 1
        FROM student_groups WHERE level_id = 4 AND code IN ('EDP', 'AMA', 'PSA');

        -- Shared AF1, MN1, TP1 (TASD + MAEF)
        FOR i IN 1..3 LOOP
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = (ARRAY['AF1','MN1','TP1'])[i];
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE (ARRAY['%(Analyse 8)','%(EDP-Num 5)','%(Proba-Stat 4)'])[i];
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 1, 2
                FROM student_groups WHERE level_id = 4 AND code IN ('TASD', 'MAEF');
            END LOOP;

        -- M1 Track-Specific Modules
        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 4 AND code = 'EDP';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('ANM','%(EDP-Num 1)'), ('EDO','%(EDP-Num 2)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 4 AND code = 'AMA';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('CD','%(Analyse 4)'), ('AC','%(Analyse 5)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 4 AND code = 'PSA';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, r_type, 2
        FROM (VALUES ('PA','%(Proba-Stat 1)',2), ('CTM','%(Analyse 7)',2), ('SI1','%(Proba-Stat 2)',2), ('AD','%(Proba-Stat 3)',2), ('PLRP','%(Algo-Info 1)',3)) AS m(s_code, p_pat, r_type) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 4 AND code = 'TASD';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('TA1','%(Topologie 1)'), ('AC','%(Topologie 2)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 4 AND code = 'MAEF';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('EDP','%(EDP-Num 6)'), ('EG','%(Economie 1)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 4 AND code = 'SA';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('SM','%(Proba-Stat 6)'), ('SADS','%(Proba-Stat 7)'), ('AD','%(Proba-Stat 1)'), ('TMDD','%(Algo-Info 2)'), ('PS','%(Proba-Stat 2)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;

        -- =========================================================================
        -- E. M2 Specialty Cohorts (With Shared Cross-Track Lectures)
        -- =========================================================================
        FOR i IN 1..2 LOOP
                v_uuid := gen_random_uuid();
                SELECT id INTO v_subj_id FROM subjects WHERE code = (ARRAY['PE','TA'])[i];
                SELECT id INTO v_prof_id FROM professors WHERE fullname LIKE (ARRAY['%(EDP-Num 1)','%(Analyse 2)'])[i];
                INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
                SELECT v_term_id, v_uuid, v_subj_id, v_prof_id, id, 2, 2
                FROM student_groups WHERE level_id = 5 AND code IN ('EDP', 'MA');
            END LOOP;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 5 AND code = 'EDP';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, r_type, 2
        FROM (VALUES ('PMP','%(EDP-Num 3)',2), ('IE','%(EDP-Num 4)',2), ('MAT-M2','%(Algo-Info 3)',3)) AS m(s_code, p_pat, r_type) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 5 AND code = 'MA';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, r_type, 2
        FROM (VALUES ('ACSL','%(EDP-Num 5)',2), ('ACSD','%(EDP-Num 6)',2), ('LAT-M2','%(Algo-Info 4)',3)) AS m(s_code, p_pat, r_type) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 5 AND code = 'SPA';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('DC','%(Proba-Stat 3)'), ('CST','%(Proba-Stat 4)'), ('STE','%(Proba-Stat 5)'), ('CS','%(EDP-Num 1)'), ('CMMA','%(Proba-Stat 6)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 5 AND code = 'TASD';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('TT','%(Algèbre 7)'), ('TA','%(Topologie 3)'), ('SISD','%(Topologie 4)'), ('ATAMM','%(Topologie 5)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 5 AND code = 'MAEF';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, r_type, 2
        FROM (VALUES ('BDMR','%(Algo-Info 5)',3), ('SVE','%(Proba-Stat 7)',2), ('CS','%(EDP-Num 2)',2), ('MEC','%(Economie 2)',2), ('MS','%(Proba-Stat 1)',2)) AS m(s_code, p_pat, r_type) JOIN subjects s ON s.code = m.s_code;

        SELECT id INTO v_grp_id FROM student_groups WHERE level_id = 5 AND code = 'SA';
        INSERT INTO course_offerings (term_id, session_uuid, subject_id, professor_id, student_group_id, required_room_type_id, weekly_sessions)
        SELECT v_term_id, gen_random_uuid(), s.id, (SELECT id FROM professors WHERE fullname LIKE p_pat), v_grp_id, 2, 2
        FROM (VALUES ('TAP','%(Algo-Info 6)'), ('TMPD','%(Proba-Stat 2)'), ('PEX','%(Proba-Stat 3)'), ('AD','%(Proba-Stat 4)'), ('RN2','%(Algo-Info 1)')) AS m(s_code, p_pat) JOIN subjects s ON s.code = m.s_code;
    END $$;

COMMIT;