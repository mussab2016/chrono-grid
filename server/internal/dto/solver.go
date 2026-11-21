package dto

type TimetableSlotDTO struct {
	ID                int    `json:"id"`
	TermID            int    `json:"termId"`
	OfferingID        int    `json:"offeringId"`
	SessionUUID       string `json:"sessionUuid"`
	SubjectID         int    `json:"subjectId"`
	SubjectCode       string `json:"subjectCode"`
	SubjectName       string `json:"subjectName"`
	ProfessorID       int    `json:"professorId"`
	ProfessorName     string `json:"professorName"`
	RoomID            int    `json:"roomId"`
	RoomCode          string `json:"roomCode"`
	RoomTypeCode      string `json:"roomTypeCode"`
	AcademicLevelCode string `json:"academicLevelCode"`
	AcademicLevelID   int    `json:"academicLevelId"`
	DayNumber         int    `json:"dayNumber"`
	SlotNumber        int    `json:"slotNumber"`
	Groups            []int  `json:"groups"`
}
