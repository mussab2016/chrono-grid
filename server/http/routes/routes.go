package routes

const (
	Professors           = "/professors/{id...}"
	RoomTypes            = "/room-types/{id...}"
	Rooms                = "/rooms/{id...}"
	AcademicLevels       = "/academic-levels/{id...}"
	Subjects             = "/subjects/{id...}"
	AcademicYears        = "/academic-years/{id...}"
	AcademicTerms        = "/academic-terms/{id...}"
	StudentGroups        = "/student-groups/{id...}"
	CourseOfferings      = "/course-offerings/{id...}"
	CourseOfferingsBatch = "/course-offerings-batch/{sessionUuid...}"
	Timetables           = "/timetables"
	TimetablesGenerate   = "/timetables/generate"
)
