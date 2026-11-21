package services

import (
	"unisched/server/internal/repository"
)

// Services acts as a central container holding all application business logic services
type Services struct {
	Professor      *ProfessorService
	RoomType       *RoomTypeService
	Room           *RoomService
	AcademicLevel  *AcademicLevelService
	Subject        *SubjectService
	AcademicYear   *AcademicYearService
	AcademicTerm   *AcademicTermService
	StudentGroup   *StudentGroupService
	CourseOffering *CourseOfferingService
	Timetable      *TimetableService
}

// NewServices receives the Repositories container and injects each repository into its service
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Professor:      NewProfessorService(repos.Professor),
		RoomType:       NewRoomTypeService(repos.RoomType),
		Room:           NewRoomService(repos.Room),
		AcademicLevel:  NewAcademicLevelService(repos.AcademicLevel),
		Subject:        NewSubjectService(repos.Subject),
		AcademicYear:   NewAcademicYearService(repos.AcademicYear),
		AcademicTerm:   NewAcademicTermService(repos.AcademicTerm),
		StudentGroup:   NewStudentGroupService(repos.StudentGroup),
		CourseOffering: NewCourseOfferingService(repos.CourseOffering),
		Timetable:      NewTimetableService(repos.Timetable, "python", "./server/ai/solver.py"),
	}
}
