package repository

import (
	"database/sql"
	lvlRepo "unisched/server/internal/repository/academiclevel"
	acdTrmRepo "unisched/server/internal/repository/academicterms"
	acdYearRepo "unisched/server/internal/repository/academicyears"
	coursOfferRepo "unisched/server/internal/repository/courseoffering"
	profRepo "unisched/server/internal/repository/professor"
	roomRepo "unisched/server/internal/repository/room"
	roomTypeRepo "unisched/server/internal/repository/roomType"
	stdntGrpRepo "unisched/server/internal/repository/studentgroups"
	subjRepo "unisched/server/internal/repository/subject"
	timetables "unisched/server/internal/repository/timetable"
)

// Repositories acts as a central container holding references to all entity repositories
type Repositories struct {
	Professor      profRepo.Repository
	RoomType       roomTypeRepo.Repository
	Room           roomRepo.Repository
	AcademicLevel  lvlRepo.Repository
	Subject        subjRepo.Repository
	AcademicYear   acdYearRepo.Repository
	AcademicTerm   acdTrmRepo.Repository
	StudentGroup   stdntGrpRepo.Repository
	CourseOffering coursOfferRepo.Repository
	Timetable      timetables.Repository
}

// NewRepositories initializes all entity repositories and packs them into a single container
func NewRepositories(db *sql.DB) (*Repositories, error) {
	profRepository, err := profRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	roomTypeRepository, err := roomTypeRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	roomRepository, err := roomRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	academicLevelRepository, err := lvlRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	subjectRepository, err := subjRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	academicYearRepository, err := acdYearRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	academicTermRepository, err := acdTrmRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	studentGroupRepository, err := stdntGrpRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	courseOfferingRepository, err := coursOfferRepo.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	timetableRepo, err := timetables.NewPostgresRepository(db)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Professor:      profRepository,
		RoomType:       roomTypeRepository,
		Room:           roomRepository,
		AcademicLevel:  academicLevelRepository,
		Subject:        subjectRepository,
		AcademicYear:   academicYearRepository,
		AcademicTerm:   academicTermRepository,
		StudentGroup:   studentGroupRepository,
		CourseOffering: courseOfferingRepository,
		Timetable:      timetableRepo,
	}, nil
}

// Close gracefully closes all prepared statements across all repositories
func (r *Repositories) Close() {
	if r.Professor != nil {
		_ = r.Professor.Close()
	}
	if r.RoomType != nil {
		_ = r.RoomType.Close()
	}
	if r.Room != nil {
		_ = r.Room.Close()
	}
	if r.AcademicLevel != nil {
		_ = r.AcademicLevel.Close()
	}
	if r.Subject != nil {
		_ = r.Subject.Close()
	}
	if r.Timetable != nil {
		_ = r.Timetable.Close()
	}
}
