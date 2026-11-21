package dto

type CreateStudentGroupRequest struct {
	AcademicYearID int    `json:"academicYearId"`
	LevelID        int    `json:"levelId"`
	Code           string `json:"code"`
}

type UpdateStudentGroupRequest struct {
	ID             int    `json:"id"`
	AcademicYearID int    `json:"academicYearId,omitempty"`
	LevelID        int    `json:"levelId,omitempty"`
	Code           string `json:"code,omitempty"`
}

type StudentGroupResponse struct {
	ID                int    `json:"id"`
	AcademicYearID    int    `json:"academicYearId,omitempty"`
	AcademicYearLabel string `json:"academicYearLabel,omitempty"`
	LevelID           int    `json:"levelId,omitempty"`
	LevelCode         string `json:"levelCode,omitempty"`
	LevelName         string `json:"levelName,omitempty"`
	Code              string `json:"code,omitempty"`
}
