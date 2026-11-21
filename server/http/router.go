package http

import (
	"io/fs"
	"log"
	"net/http"
	"strings"
	"unisched/server/http/routes"
	"unisched/server/internal/handler"
	"unisched/server/internal/services"
)

// NewRouter registers application routes, embeds Angular dist files,
// and injects the main Services container into Backend APIs.
func NewRouter(distFS fs.FS, svc *services.Services) http.Handler {

	// --------------------------------------------------
	// Root router
	// --------------------------------------------------
	mux := http.NewServeMux()

	// --------------------------------------------------
	// Backend APIs Sub-Router
	// --------------------------------------------------
	api := http.NewServeMux()

	// Register API Handlers using the Services Container
	registerAPIRoutes(api, svc)

	// Mount API sub-router under /api/
	//mux.Handle("/api/", http.StripPrefix("/api", api))
	mux.Handle("/api/", http.StripPrefix("/api", CORSMiddleware(api)))

	// --------------------------------------------------
	// Angular Frontend & Routing Handling
	// --------------------------------------------------
	fileServer := http.FileServer(http.FS(distFS))

	// Control the main root path "/" to serve Angular files and resolve the Refresh issue
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Trim the leading "/" to match the directory path inside the embed filesystem
		filePath := strings.TrimPrefix(r.URL.Path, "/")

		f, err := distFS.Open(filePath)
		if err != nil {
			// If the file is not found (SPA routing refresh like /professors),
			// fall back to index.html so Angular handles client-side routing.
			r.URL.Path = "/"
		} else {
			// If the file exists, close it immediately to prevent memory leaks
			if err := f.Close(); err != nil {
				log.Printf("Error closing file %s: %s\n", filePath, err)
			}
		}

		// Run the file server to deliver the correct static asset
		fileServer.ServeHTTP(w, r)
	})

	return mux
}

// registerAPIRoutes binds API endpoints to their respective Handlers.
func registerAPIRoutes(api *http.ServeMux, svc *services.Services) {
	// 1. Professors API
	profHandler := handler.NewProfessorHandler(svc)
	roomTypeHandler := handler.NewRoomTypeHandler(svc)
	roomHandler := handler.NewRoomHandler(svc)
	academicLevelHandler := handler.NewAcademicLevelHandler(svc)
	subjectHandler := handler.NewSubjectHandler(svc)
	academicYearHandler := handler.NewAcademicYearHandler(svc)
	academicTermHandler := handler.NewAcademicTermHandler(svc)
	studentGroupHandler := handler.NewStudentGroupHandler(svc)
	courseOfferingHandler := handler.NewCourseOfferingHandler(svc)
	timetableHandler := handler.NewTimetableHandler(svc)

	api.HandleFunc(routes.Professors, profHandler.HandleProfessors)
	api.HandleFunc(routes.RoomTypes, roomTypeHandler.HandleRoomType)
	api.HandleFunc(routes.Rooms, roomHandler.HandleRoom)
	api.HandleFunc(routes.AcademicLevels, academicLevelHandler.HandleAcademicLevel)
	api.HandleFunc(routes.Subjects, subjectHandler.HandleSubject)
	api.HandleFunc(routes.AcademicYears, academicYearHandler.HandleAcademicYear)
	api.HandleFunc(routes.AcademicTerms, academicTermHandler.HandleAcademicTerm)
	api.HandleFunc(routes.StudentGroups, studentGroupHandler.HandleStudentGroup)
	api.HandleFunc(routes.CourseOfferings, courseOfferingHandler.HandleCourseOfferings)
	api.HandleFunc(routes.CourseOfferingsBatch, courseOfferingHandler.HandleCourseOfferingBatch)
	api.HandleFunc(routes.Timetables, timetableHandler.HandleTimetables)
	api.HandleFunc(routes.TimetablesGenerate, timetableHandler.HandleGenerateTimetable)

	// 2. Future APIs (Rooms, Offerings, Timetables) can be seamlessly attached here:
	// roomHandler := handler.NewRoomHandler(svc)
	// api.HandleFunc("/rooms", roomHandler.HandleRooms)
}
