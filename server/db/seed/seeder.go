package seed

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Run populates the embedded PostgreSQL database with realistic university infrastructure,
// 5 academic levels (L1, L2, L3, M1, M2), specialization groups, professors, subjects,
// and full course offerings for testing.
func Run(ctx context.Context, db *sql.DB) error {
	var existingCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM academic_years WHERE year_label = $1`, "2026-2027").Scan(&existingCount); err != nil {
		return fmt.Errorf("seed check failed: %w", err)
	}
	if existingCount > 0 {
		return nil // Already seeded
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed tx: %w", err)
	}
	defer tx.Rollback()

	// ==========================================
	// 1. ROOM TYPES & PHYSICAL ROOMS
	// ==========================================
	roomTypes := []struct {
		Code string
		Desc string
	}{
		{"AMPHI", "Amphitheater for combined lectures"},
		{"TD", "Standard tutorial classroom (24 rooms)"},
		{"LAB", "Specialized applied laboratory (5 labs)"},
		{"INFO", "Computer science lab room (4 rooms)"},
	}

	roomTypeIDs := make(map[string]int)
	for _, rt := range roomTypes {
		var id int
		err := tx.QueryRowContext(ctx, `
			INSERT INTO room_types (code, description)
			VALUES ($1, $2)
			ON CONFLICT (code) DO UPDATE SET description = EXCLUDED.description
			RETURNING id
		`, rt.Code, rt.Desc).Scan(&id)
		if err != nil {
			return fmt.Errorf("insert room_type %s: %w", rt.Code, err)
		}
		roomTypeIDs[rt.Code] = id
	}

	insertRooms := func(typeCode, prefix, label string, count int) error {
		rtID := roomTypeIDs[typeCode]
		for i := 1; i <= count; i++ {
			code := fmt.Sprintf("%s-%02d", prefix, i)
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO rooms (room_type_id, code, description)
				VALUES ($1, $2, $3)
			`, rtID, code, fmt.Sprintf("%s %d", label, i)); err != nil {
				return err
			}
		}
		return nil
	}

	if err := insertRooms("AMPHI", "AMPHI", "Amphitheater", 4); err != nil {
		return err
	}
	if err := insertRooms("TD", "TD", "Tutorial Room", 24); err != nil {
		return err
	}
	if err := insertRooms("LAB", "LAB", "Applied Lab", 5); err != nil {
		return err
	}
	if err := insertRooms("INFO", "INFO", "Computer Lab", 4); err != nil {
		return err
	}

	// ==========================================
	// 2. ACADEMIC LEVELS (5 LEVELS ONLY: L1, L2, L3, M1, M2)
	// ==========================================
	levels := []struct {
		Code string
		Name string
	}{
		{"L1", "1st Year License"},
		{"L2", "2nd Year License"},
		{"L3", "3rd Year License"},
		{"M1", "Master 1"},
		{"M2", "Master 2"},
	}

	levelIDs := make(map[string]int)
	for _, lvl := range levels {
		var id int
		err := tx.QueryRowContext(ctx, `
			INSERT INTO academic_levels (name, code)
			VALUES ($1, $2)
			RETURNING id
		`, lvl.Name, lvl.Code).Scan(&id)
		if err != nil {
			return fmt.Errorf("insert level %s: %w", lvl.Code, err)
		}
		levelIDs[lvl.Code] = id
	}

	// ==========================================
	// 3. ACADEMIC YEAR & TERMS (S1, S2)
	// ==========================================
	var yearID int
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO academic_years (year_label)
		VALUES ($1)
		RETURNING id
	`, "2026-2027").Scan(&yearID); err != nil {
		return fmt.Errorf("insert academic_year: %w", err)
	}

	var termS1ID, termS2ID int
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO academic_terms (academic_year_id, semester, label)
		VALUES ($1, 1, 'S1') RETURNING id
	`, yearID).Scan(&termS1ID); err != nil {
		return err
	}
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO academic_terms (academic_year_id, semester, label)
		VALUES ($1, 2, 'S2') RETURNING id
	`, yearID).Scan(&termS2ID); err != nil {
		return err
	}

	// ==========================================
	// 4. STUDENT GROUPS INSIDE EACH LEVEL (Specializations & Groups)
	// ==========================================
	groupIDsByLevel := make(map[string][]int)
	groupIDByCode := make(map[string]int)

	insertGroup := func(levelCode, groupCode string) error {
		var gID int
		err := tx.QueryRowContext(ctx, `
			INSERT INTO student_groups (academic_year_id, level_id, code)
			VALUES ($1, $2, $3)
			RETURNING id
		`, yearID, levelIDs[levelCode], groupCode).Scan(&gID)
		if err != nil {
			return fmt.Errorf("insert group %s/%s: %w", levelCode, groupCode, err)
		}
		groupIDsByLevel[levelCode] = append(groupIDsByLevel[levelCode], gID)
		groupIDByCode[levelCode+":"+groupCode] = gID
		return nil
	}

	// L1: 12 Groups (G01 .. G12)
	for i := 1; i <= 12; i++ {
		if err := insertGroup("L1", fmt.Sprintf("G%02d", i)); err != nil {
			return err
		}
	}

	// L2: 9 Groups (G01 .. G09)
	for i := 1; i <= 9; i++ {
		if err := insertGroup("L2", fmt.Sprintf("G%02d", i)); err != nil {
			return err
		}
	}

	// L3: 3 Specializations x 2 groups each = 6 groups
	l3Groups := []string{"AM-G1", "AM-G2", "PROBA-G1", "PROBA-G2", "RO-G1", "RO-G2"}
	for _, g := range l3Groups {
		if err := insertGroup("L3", g); err != nil {
			return err
		}
	}

	// M1: 4 Specialization Groups
	masterSpecs := []string{"EDP", "AN", "PROBA", "RO"}
	for _, spec := range masterSpecs {
		if err := insertGroup("M1", spec); err != nil {
			return err
		}
	}

	// M2: Same 4 Specialization Groups
	for _, spec := range masterSpecs {
		if err := insertGroup("M2", spec); err != nil {
			return err
		}
	}

	// ==========================================
	// 5. PROFESSORS (38 FACULTY MEMBERS)
	// ==========================================
	professors := []struct {
		Name string
		Dept string
	}{
		{"Prof. Abdelkrim Benali", "Applied Mathematics"},
		{"Prof. Samira Khelifi", "Applied Mathematics"},
		{"Dr. Rachid Boudiaf", "Numerical Analysis"},
		{"Dr. Nadia Mansouri", "Numerical Analysis"},
		{"Prof. Mohamed Amine Zerrouki", "Partial Differential Equations"},
		{"Dr. Yasmine Belkacem", "Partial Differential Equations"},
		{"Dr. Tariq Hamzaoui", "Computer Science"},
		{"Dr. Meriem Ouali", "Computer Science"},
		{"Prof. Kamel Meziane", "Probability & Statistics"},
		{"Dr. Leila Saidi", "Probability & Statistics"},
		{"Dr. Sofiane Rahmani", "Operations Research"},
		{"Dr. Amel Guerfi", "Operations Research"},
		{"Prof. Farid Bensalah", "Pure Mathematics"},
		{"Dr. Houda Chergui", "Pure Mathematics"},
		{"Dr. Walid Touati", "Numerical Analysis"},
		{"Dr. Nesrine Djebbar", "Computer Science"},
		{"Prof. Mourad Haddad", "Applied Mathematics"},
		{"Dr. Imene Slimani", "Probability & Statistics"},
		{"Dr. Bilal Messaoudi", "Operations Research"},
		{"Dr. Rania Bouzid", "Partial Differential Equations"},
		{"Dr. Adel Kaci", "Computer Science"},
		{"Dr. Salima Ferhat", "Pure Mathematics"},
		{"Prof. Nabil Benmohamed", "Numerical Analysis"},
		{"Dr. Kenza Ait-Ali", "Applied Mathematics"},
		{"Dr. Hakim Larbi", "Operations Research"},
		{"Dr. Souad Brahimi", "Probability & Statistics"},
		{"Dr. Zakaria Taleb", "Computer Science"},
		{"Dr. Asma Ziani", "Partial Differential Equations"},
		{"Prof. Youcef Merabet", "Pure Mathematics"},
		{"Dr. Hanane Belhadj", "Numerical Analysis"},
		{"Dr. Mehdi Chikhi", "Applied Mathematics"},
		{"Dr. Wissam Kadri", "Computer Science"},
		{"Dr. Samir Boukhalfa", "Operations Research"},
		{"Dr. Dalila Rezig", "Probability & Statistics"},
		{"Dr. Redouane Amrani", "Partial Differential Equations"},
		{"Dr. Sabrina Yahiaoui", "Pure Mathematics"},
		{"Dr. Omar Benyahia", "Computer Science"},
		{"Dr. Lamia Mokhtari", "Numerical Analysis"},
	}

	profIDs := make([]int, 0, len(professors))
	for _, p := range professors {
		var pID int
		err := tx.QueryRowContext(ctx, `
			INSERT INTO professors (fullname, department)
			VALUES ($1, $2)
			ON CONFLICT (fullname) DO UPDATE SET department = EXCLUDED.department
			RETURNING id
		`, p.Name, p.Dept).Scan(&pID)
		if err != nil {
			return fmt.Errorf("insert professor %s: %w", p.Name, err)
		}
		profIDs = append(profIDs, pID)
	}

	// ==========================================
	// 6. SUBJECTS CATALOG
	// ==========================================
	subjects := []struct {
		Code string
		Name string
	}{
		{"MATH101", "Calculus & Real Analysis I"},
		{"ALGE101", "Linear Algebra I"},
		{"INFO101", "Algorithms & Programming (C/Python)"},
		{"PHYS101", "Mechanics & General Physics"},
		{"MATH201", "Multivariable Analysis & Series"},
		{"ALGE201", "Matrix Reduction & Euclidean Spaces"},
		{"NUM201", "Introduction to Numerical Methods"},
		{"PROB201", "Probability & Combinatorics"},
		{"PDE301", "Partial Differential Equations"},
		{"NA301", "Numerical Linear Algebra & MATLAB"},
		{"STAT301", "Inferential Statistics"},
		{"OR301", "Linear Programming & Simplex"},
		{"PDE401", "Finite Difference & Finite Element Methods"},
		{"NA401", "Scientific Computing & Optimization"},
		{"STAT401", "Stochastic Processes & Martingales"},
		{"OR401", "Combinatorial & Graph Optimization"},
		{"PDE501", "Evolution Equations & Stability Analysis"},
		{"NA501", "Interior-Point & Convex Optimization"},
		{"STAT501", "Time Series & Applied Econometrics"},
		{"OR501", "Optimal Control & Metaheuristics"},
	}

	subjectIDs := make(map[string]int)
	for _, s := range subjects {
		var sID int
		err := tx.QueryRowContext(ctx, `
			INSERT INTO subjects (name, code)
			VALUES ($1, $2)
			ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, s.Name, s.Code).Scan(&sID)
		if err != nil {
			return fmt.Errorf("insert subject %s: %w", s.Code, err)
		}
		subjectIDs[s.Code] = sID
	}

	// ==========================================
	// 7. COURSE OFFERINGS (FULL SEMESTER S1)
	// ==========================================
	addBlock := func(termID int, subjectCode string, profIdx int, roomTypeCode string, weeklySessions int, gIDs []int) error {
		sessionUUID := uuid.New().String()
		subID := subjectIDs[subjectCode]
		profID := profIDs[profIdx]
		rtID := roomTypeIDs[roomTypeCode]

		for _, gID := range gIDs {
			_, err := tx.ExecContext(ctx, `
				INSERT INTO course_offerings (
					session_uuid, term_id, subject_id, professor_id,
					required_room_type_id, student_group_id, weekly_sessions
				) VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, sessionUUID, termID, subID, profID, rtID, gID, weeklySessions)
			if err != nil {
				return fmt.Errorf("insert course_offering (%s): %w", subjectCode, err)
			}
		}
		return nil
	}

	// --- L1 Offerings (12 Groups: Section A = G01..G06, Section B = G07..G12) ---
	l1 := groupIDsByLevel["L1"]
	if err := addBlock(termS1ID, "MATH101", 0, "AMPHI", 2, l1[0:6]); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "MATH101", 1, "AMPHI", 2, l1[6:12]); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "ALGE101", 12, "AMPHI", 1, l1[0:6]); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "ALGE101", 13, "AMPHI", 1, l1[6:12]); err != nil {
		return err
	}
	for i, gID := range l1 {
		if err := addBlock(termS1ID, "MATH101", 2+(i%4), "TD", 1, []int{gID}); err != nil {
			return err
		}
		if err := addBlock(termS1ID, "ALGE101", 14+(i%4), "TD", 1, []int{gID}); err != nil {
			return err
		}
		if err := addBlock(termS1ID, "INFO101", 6+(i%4), "INFO", 1, []int{gID}); err != nil {
			return err
		}
	}

	// --- L2 Offerings (9 Groups: G01..G09) ---
	l2 := groupIDsByLevel["L2"]
	if err := addBlock(termS1ID, "MATH201", 16, "AMPHI", 2, l2); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "NUM201", 2, "AMPHI", 1, l2); err != nil {
		return err
	}
	for i, gID := range l2 {
		if err := addBlock(termS1ID, "MATH201", 20+(i%3), "TD", 1, []int{gID}); err != nil {
			return err
		}
		if err := addBlock(termS1ID, "NUM201", 23+(i%3), "INFO", 1, []int{gID}); err != nil {
			return err
		}
	}

	// --- L3 Offerings (3 Specializations x 2 groups each) ---
	l3AM := []int{groupIDByCode["L3:AM-G1"], groupIDByCode["L3:AM-G2"]}
	l3PROBA := []int{groupIDByCode["L3:PROBA-G1"], groupIDByCode["L3:PROBA-G2"]}
	l3RO := []int{groupIDByCode["L3:RO-G1"], groupIDByCode["L3:RO-G2"]}

	if err := addBlock(termS1ID, "PDE301", 4, "AMPHI", 1, l3AM); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "NA301", 22, "LAB", 1, l3AM); err != nil {
		return err
	}
	for i, gID := range l3AM {
		if err := addBlock(termS1ID, "PDE301", 5+i, "TD", 1, []int{gID}); err != nil {
			return err
		}
	}

	if err := addBlock(termS1ID, "STAT301", 8, "AMPHI", 1, l3PROBA); err != nil {
		return err
	}
	for i, gID := range l3PROBA {
		if err := addBlock(termS1ID, "STAT301", 9+i, "TD", 1, []int{gID}); err != nil {
			return err
		}
	}

	if err := addBlock(termS1ID, "OR301", 10, "AMPHI", 1, l3RO); err != nil {
		return err
	}
	for i, gID := range l3RO {
		if err := addBlock(termS1ID, "OR301", 11+i, "TD", 1, []int{gID}); err != nil {
			return err
		}
	}

	// --- M1 Offerings (4 Specialization Groups: EDP, AN, PROBA, RO) ---
	if err := addBlock(termS1ID, "PDE401", 4, "LAB", 2, []int{groupIDByCode["M1:EDP"]}); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "NA401", 22, "INFO", 2, []int{groupIDByCode["M1:AN"]}); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "STAT401", 8, "TD", 2, []int{groupIDByCode["M1:PROBA"]}); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "OR401", 10, "LAB", 2, []int{groupIDByCode["M1:RO"]}); err != nil {
		return err
	}

	// --- M2 Offerings (Same 4 Specialization Groups: EDP, AN, PROBA, RO) ---
	if err := addBlock(termS1ID, "PDE501", 5, "LAB", 2, []int{groupIDByCode["M2:EDP"]}); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "NA501", 29, "INFO", 2, []int{groupIDByCode["M2:AN"]}); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "STAT501", 17, "TD", 2, []int{groupIDByCode["M2:PROBA"]}); err != nil {
		return err
	}
	if err := addBlock(termS1ID, "OR501", 18, "LAB", 2, []int{groupIDByCode["M2:RO"]}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed tx: %w", err)
	}
	return nil
}
