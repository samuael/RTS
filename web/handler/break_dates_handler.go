package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/BreakDates"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// BreakDateHandler struct
type BreakDateHandler struct {
	SessionService   *session.Cookiehandler
	BreakDateService BreakDates.BreakDateService
}

// NewBreakDateHandler function
func NewBreakDateHandler(bdservice BreakDates.BreakDateService, session *session.Cookiehandler) *BreakDateHandler {
	return &BreakDateHandler{
		SessionService:   session,
		BreakDateService: bdservice,
	}
}

// CreateBreakDatesHandler method
// Method Post
// Input type Json
func (bdhandler *BreakDateHandler) CreateBreakDatesHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestersession := bdhandler.SessionService.GetSession(request)
	ds := struct {
		Success bool
		Message string
		Date    *etc.Date
	}{
		Success: false,
	}
	newDecoder := json.NewDecoder(request.Body)
	newDate := &etc.Date{}
	todaye := etc.NewDate(0)
	decodeError := newDecoder.Decode(newDate)
	if decodeError != nil || newDate.Day <= 0 || newDate.Month <= 0 || newDate.Year < todaye.Year {
		ds.Message = "Invalid Date Input The Date Be After today"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if !newDate.IsPassed(todaye) {
		ds.Message = "The Date You Have entered Has Passed "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	newDate.BranchID = requestersession.BranchID
	newDate.IsBreakDate = true
	newDate = bdhandler.BreakDateService.CreateBreakDate(newDate)
	if newDate == nil {
		ds.Message = " Error While Saving the Day Please Try Again  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "SuccessFully Saved the Date "
	ds.Success = true
	ds.Date = newDate
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DeleteBreakDatesHandler method
// Method GET
// Input for date_id uint
func (bdhandler *BreakDateHandler) DeleteBreakDatesHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requestersession := bdhandler.SessionService.GetSession(request)
	ds := struct {
		Success bool
		Message string
		DateID  uint
	}{
		Success: false,
	}
	dateidstring := request.FormValue("date_id")
	dateID, Erra := strconv.Atoi(dateidstring)
	if Erra != nil || dateID <= 0 {
		ds.Message = "Invalid Request Body"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	ds.DateID = uint(dateID)
	success := bdhandler.BreakDateService.DeleteBreakDate(uint(dateID))
	if success {
		ds.Message = "Successfully Deleted "
		ds.Success = true
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Error While Deleting the Reserved Date Please Try again "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetBreakDates method
// Method is Get
// Authorization ALL LOGGED IN USERS
// Input from session requesterSession BRANCHID
func (bdhandler *BreakDateHandler) GetBreakDates(response http.ResponseWriter, request *http.Request) {
	requestersession := bdhandler.SessionService.GetSession(request)
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success    bool
		Message    string
		BreakDates []etc.Date
	}{
		Success: false,
	}
	BranchID := requestersession.BranchID
	breakDates := bdhandler.BreakDateService.GetBreakDates(BranchID, etc.NewDate(0))
	if breakDates != nil {
		ds.Message = "Internal Server Error "
		ds.BreakDates = *breakDates
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Break Dates Succesfuly Fetched "
	ds.BreakDates = *breakDates
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DateInformationHandler Information
// Method GET
// Input Form Value
/*
	day
	month
	and
	Year */
func (bdhandler *BreakDateHandler) DateInformationHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	daystring := request.FormValue("day")
	monthstring := request.FormValue("month")
	yearstring := request.FormValue("year")

	year, erra := strconv.Atoi(yearstring)
	month, erra := strconv.Atoi(monthstring)
	day, erra := strconv.Atoi(daystring)

	ds := struct {
		Success bool
		Message string
		Date    *etc.Date
	}{
		Success: false,
	}
	if erra != nil || year <= 0 || month > 13 || month < 0 || day > 30 || day < 0 || (month == 13 && day > 6) {
		ds.Message = "Invalid Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	dates := &etc.Date{
		Day:   day,
		Month: month,
		Year:  year,
	}

	if !dates.IsPassed(etc.NewDate(0)) {
		ds.Message = "The Dates has Passed "
		ds.Success = false
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	dates = dates.FulFill()
	dates = dates.Modify()

	ds.Date = dates
	ds.Message = "The Fullfiled Date is :... "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}
