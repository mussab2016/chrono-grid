package models

type ProfPrefInput struct {
	ProfessorID int `json:"professorId"`
	DayNumber   int `json:"dayNumber"`
	SlotNumber  int `json:"slotNumber"` // 0 = entire day
	Score       int `json:"score"`      // -1 = forbidden, >0 = preferred
}

type GenerateRequest struct {
	TermID      int             `json:"termId"`
	DaysCount   int             `json:"daysCount"`
	SlotsPerDay int             `json:"slotsPerDay"`
	Preferences []ProfPrefInput `json:"preferences"`
}

type RoomInput struct {
	ID         int `json:"id"`
	RoomTypeID int `json:"roomTypeId"`
}

type BlockInput struct {
	SessionUUID    string `json:"sessionUuid"`
	OfferingID     int    `json:"offeringId"`
	ProfessorID    int    `json:"professorId"`
	RoomTypeID     int    `json:"roomTypeId"`
	WeeklySessions int    `json:"weeklySessions"`
	GroupIDs       []int  `json:"groupIds"`
}

type SolverRequest struct {
	DaysCount   int             `json:"daysCount"`
	SlotsPerDay int             `json:"slotsPerDay"`
	Rooms       []RoomInput     `json:"rooms"`
	Blocks      []BlockInput    `json:"blocks"`
	Preferences []ProfPrefInput `json:"preferences"`
}

type AssignmentOutput struct {
	OfferingID  int    `json:"offeringId"`
	SessionUUID string `json:"sessionUuid"`
	RoomID      int    `json:"roomId"`
	DayNumber   int    `json:"dayNumber"`
	SlotNumber  int    `json:"slotNumber"`
}

type SolverResponse struct {
	Status      string             `json:"status"`
	Assignments []AssignmentOutput `json:"assignments"`
}
