package dto

// CourseOfferingRequest is used for creating or updating a single offering row (/course-offerings/0 or /course-offerings/{id}).
type CourseOfferingRequest struct {
	ID                 int     `json:"id"`
	SessionUUID        *string `json:"sessionUuid,omitempty"`
	TermID             int     `json:"termId"`
	SubjectID          int     `json:"subjectId"`
	ProfessorID        int     `json:"professorId"`
	RequiredRoomTypeID int     `json:"requiredRoomTypeId"`
	StudentGroupID     int     `json:"studentGroupId"`
	WeeklySessions     int     `json:"weeklySessions"`
}

// CourseOfferingBatchRequest is sent by Angular to create or update a physical session block across multiple groups.
// If IsSplitPerGroup is true, the service generates a separate session_uuid for each group in StudentGroupIDs.
// If false, all groups share a single session_uuid (combined lecture/session).
type CourseOfferingBatchRequest struct {
	SessionUUID        *string `json:"sessionUuid,omitempty"`
	TermID             int     `json:"termId"`
	SubjectID          int     `json:"subjectId"`
	ProfessorID        int     `json:"professorId"`
	RequiredRoomTypeID int     `json:"requiredRoomTypeId"`
	StudentGroupIDs    []int   `json:"studentGroupIds"`
	WeeklySessions     int     `json:"weeklySessions"`
	IsSplitPerGroup    bool    `json:"isSplitPerGroup"`
}

// GroupSummary represents a compact student group item inside an aggregated session block.
type GroupSummary struct {
	OfferingID int     `json:"offeringId"`
	GroupID    int     `json:"groupId"`
	Code       *string `json:"code,omitempty"`
	LevelCode  *string `json:"levelCode,omitempty"`
}

// CourseOfferingBlockResponse aggregates rows sharing the same session_uuid into one logical card for UI & Timetable Solver.
type CourseOfferingBlockResponse struct {
	SessionUUID        string         `json:"sessionUuid"`
	TermID             int            `json:"termId"`
	TermLabel          *string        `json:"termLabel,omitempty"`
	SubjectID          int            `json:"subjectId"`
	SubjectName        *string        `json:"subjectName,omitempty"`
	SubjectCode        *string        `json:"subjectCode,omitempty"`
	ProfessorID        int            `json:"professorId"`
	ProfessorName      *string        `json:"professorName,omitempty"`
	RequiredRoomTypeID int            `json:"requiredRoomTypeId"`
	RoomTypeName       *string        `json:"roomTypeName,omitempty"`
	WeeklySessions     int            `json:"weeklySessions"`
	Groups             []GroupSummary `json:"groups"`
}
