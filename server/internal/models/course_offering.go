package models

import "time"

// CourseOffering represents a single database row linking a physical session block to one student group.
type CourseOffering struct {
	ID                 int       `json:"id"`
	SessionUUID        string    `json:"sessionUuid"`
	TermID             int       `json:"termId"`
	SubjectID          int       `json:"subjectId"`
	ProfessorID        int       `json:"professorId"`
	RequiredRoomTypeID int       `json:"requiredRoomTypeId"`
	StudentGroupID     int       `json:"studentGroupId"`
	WeeklySessions     int       `json:"weeklySessions"`
	CreatedAt          time.Time `json:"createdAt"`

	// Joined fields as pointers so unset values serialize as omitted rather than empty strings
	TermLabel        *string `json:"termLabel,omitempty"`
	SubjectName      *string `json:"subjectName,omitempty"`
	SubjectCode      *string `json:"subjectCode,omitempty"`
	ProfessorName    *string `json:"professorName,omitempty"`
	RoomTypeName     *string `json:"roomTypeName,omitempty"`
	StudentGroupCode *string `json:"studentGroupCode,omitempty"`
	LevelCode        *string `json:"levelCode,omitempty"`
}
