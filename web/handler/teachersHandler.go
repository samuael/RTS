package handler

import (
	"encoding/json"
	"net/http"

	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher/TeacherService"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// TeacherHandler struct
type TeacherHandler struct {
	TeacherService *TeacherService.TeacherService
	SessionService *session.Cookiehandler
}

// NewTeacherHandler method
func NewTeacherHandler(
	service *TeacherService.TeacherService,
	SessionService *session.Cookiehandler,
) *TeacherHandler {
	return &TeacherHandler{
		TeacherService: service,
		SessionService: SessionService,
	}
}

// GetMyActiveLectures Lectures
// Method GEt
// Authority TEACHERS
func (teacherhandler *TeacherHandler) GetMyActiveLectures(response http.ResponseWriter, request *http.Request) {
	// This has to satisfy
	// Select All Active Rounds Lectures
	// Has to Match Branch
	// Match Teacher ID
	// Their Round Has To Be Populeted
	// Gruped By Round Number
	requestersession := teacherhandler.SessionService.GetSession(request)
	if requestersession == nil || requestersession.Role != entity.TEACHER {
		response.Write([]byte("<h2> Un Authorized User </h2>"))
		return
	}
	ActiveLecturesDS := teacherhandler.
		TeacherService.
		GetActiveLectures(requestersession.BranchID, requestersession.ID)
	jsonReturn, _ := json.Marshal(ActiveLecturesDS)
	response.Write(jsonReturn)
}

// GetTodaysLectures Lectures
// Method GEt
// Authority TEACHERS
func (teacherhandler *TeacherHandler) GetTodaysLectures(response http.ResponseWriter, request *http.Request) {
	// This has to satisfy
	// Select All Active Rounds Lectures
	// Has to Match Branch
	// Match Teacher ID
	// Their Round Has To Be Populeted
	// Gruped By Round Number
	requestersession := teacherhandler.SessionService.GetSession(request)
	if requestersession == nil || requestersession.Role != entity.TEACHER {
		response.Write([]byte("<h2> Un Authorized User </h2>"))
		return
	}
	ActiveLecturesDS := teacherhandler.
		TeacherService.
		GetTodaysLectures(requestersession.BranchID, requestersession.ID)
	jsonReturn, _ := json.Marshal(ActiveLecturesDS)
	response.Write(jsonReturn)
}

// PostponeLecture method
// Method GET
// Parameter lecture_id
func (teacherhandler *TeacherHandler) PostponeLecture(response http.ResponseWriter, request *http.Request) {
	// TheLecture Should Be In The Last Date where the Room Is Free
	// Loop Over the Reserved Dates and If You Find Make The Shift and The Date to be the Class Date
	// And if there is training in that day or time postpone to the Last
}
