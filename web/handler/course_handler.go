package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Course"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
)

/*
*CreateCourse
*DeleteCourse
EditCourse
GetCoursesByCategoryID  using category_id
GetCoursesOfABranch
CreateCourseForRound
*/

// CourseHandler struct representing handler for course related requests
type CourseHandler struct {
	SessionHandler *session.Cookiehandler
	CourseService  Course.CourseService
	RoundService   Round.RoundService
}

// NewCourseHandler function
func NewCourseHandler(session *session.Cookiehandler, coureservice Course.CourseService, Roundservice Round.RoundService) *CourseHandler {
	return &CourseHandler{
		CourseService:  coureservice,
		SessionHandler: session,
		RoundService:   Roundservice,
	}
}

// CreateCourse method
//  this is gonna be a Form Input and Json Outout
func (coursehandler *CourseHandler) CreateCourse(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestSession := coursehandler.SessionHandler.GetSession(request)
	ds := struct {
		Success bool
		Message string
		Course  *entity.Course
		VErrors form.ValidationErrors
	}{
		Success: false,
		Message: "Can't Parse the Form Niggoye ",
	}
	// parse form
	erra := request.ParseForm()
	if erra != nil {
		ds.Message = "Parse Error "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	registrationForm := form.Input{
		VErrors: form.ValidationErrors{},
		Values:  request.PostForm,
	}
	registrationForm.Required("_csrf", "title", "duration", "categoryid")
	if request.FormValue("_csrf") == "" {
		registrationForm.VErrors.Add("_csrf_Erro", "InValid User ")
	}
	duration, newErra := strconv.Atoi(request.FormValue("duration"))
	if newErra != nil {
		registrationForm.VErrors.Add("durationError", "Duration Error Please Enter a Valid Number ")
	}
	categoryid, ewErra := strconv.Atoi(request.FormValue("categoryid"))
	if ewErra != nil {
		registrationForm.VErrors.Add("CategoryID", "Category Id Error  ")
	}
	if !registrationForm.Valid() {
		ds.VErrors = registrationForm.VErrors
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	course := &entity.Course{
		BranchID:   requestSession.BranchID,
		Duration:   uint(duration),
		Title:      request.FormValue("title"),
		Categoryid: uint(categoryid),
	}
	// Saving the Category
	course = coursehandler.CourseService.CreateCourse(course)
	if course == nil {
		ds.Success = false
		ds.Message = "Can't Able to Save the Course"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Course = course
	ds.Success = true
	ds.Message = "Succesfully Created "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DeleteCourse method  taking a json inpu
func (coursehandler *CourseHandler) DeleteCourse(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Courseid uint `json:"courseid"`
		BranchID uint
		Success  bool
		Message  string
	}{
		Success: false,
	}
	requesterSession := coursehandler.SessionHandler.GetSession(request)
	ds.BranchID = requesterSession.ID
	newDecoder := json.NewDecoder(request.Body)
	decodeErrro := newDecoder.Decode(&ds)
	if decodeErrro != nil || ds.Courseid <= 0 {
		ds.Message = "Invalid Request Body "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	course := &entity.Course{
		BranchID: requesterSession.BranchID,
	}
	course.ID = ds.Courseid

	result := coursehandler.CourseService.DeleteCourse(course)
	if result {
		ds.Success = true
		ds.Message = "Course Deleted Succesfully "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Can't delete the Course  "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// EditCourse method
//  Takes Json and Returns Json
func (coursehandler *CourseHandler) EditCourse(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requesterSession := coursehandler.SessionHandler.GetSession(request)
	ds := struct {
		Success      bool
		Message      string
		NewCourse    *entity.Course
		BranchNumber uint
	}{
		Success:      false,
		BranchNumber: requesterSession.BranchID,
	}
	course := &entity.Course{}
	newDecoder := json.NewDecoder(request.Body)
	decodeError := newDecoder.Decode(course)
	if decodeError != nil || course.ID <= 0 {
		ds.Message = "Please Fill The CourseId correctly "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	theCourse := coursehandler.CourseService.GetCourseByID(course.ID)
	if theCourse == nil {
		ds.Message = "No Course Present By This ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	theCourse.Title = course.Title
	theCourse.Duration = course.Duration
	course = coursehandler.CourseService.SaveCourse(theCourse)
	if course != nil {
		ds.Message = "Can't Update The Course "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Course Configured Correctly"
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// GetCoursesByCategoryID method method Get returning Jsons
// Method GET
// get variable category_id
func (coursehandler *CourseHandler) GetCoursesByCategoryID(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	categoryid := request.FormValue("category_id")
	ds := struct {
		Succesful bool
		Courses   []entity.Course
	}{
		Succesful: false,
		Courses:   []entity.Course{},
	}
	requesterSession := coursehandler.SessionHandler.GetSession(request)
	branchid := requesterSession.BranchID
	//  Converting
	categoryNum, erra := strconv.Atoi(categoryid)
	if erra != nil {
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	courses := coursehandler.CourseService.GetCourseofBranchAndCategory(branchid, uint(categoryNum))
	ds.Courses = *courses
	if courses != nil && len(*courses) > 0 {
		ds.Succesful = true
	}
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetCoursesOfABranch method GET  returning json
func (coursehandler *CourseHandler) GetCoursesOfABranch(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Succesful bool
		Courses   []entity.Course
	}{
		Succesful: false,
		Courses:   []entity.Course{},
	}
	requesterSession := coursehandler.SessionHandler.GetSession(request)
	branchid := requesterSession.BranchID
	//  Converting
	courses := coursehandler.CourseService.GetCourseofBranch(branchid)
	ds.Courses = *courses
	if courses != nil && len(*courses) > 0 {
		ds.Succesful = true
	}
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// CreateCourseForRound method for registering a courses for Specific Course
func (coursehandler *CourseHandler) CreateCourseForRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Title    string `json:"title"`
		Duration uint   `json:"duration"`
		RoundID  uint   `json:"roundID"`
		Message  string
		Course   entity.Course
		Round    entity.Round
	}{
		Success: false,
	}
	requesterSession := coursehandler.SessionHandler.GetSession(request)
	if requesterSession == nil {
		return
	}
	newDecoder := json.NewDecoder(request.Body)
	decodeError := newDecoder.Decode(&ds)
	if decodeError != nil {
		ds.Message = "Error While Decoding the message "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if ds.Title == "" || ds.Duration <= 0 || ds.RoundID <= 0 {
		ds.Message = "Please Write Appropriate Input type "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	round := coursehandler.RoundService.GetRoundByID(ds.RoundID)
	if round != nil && round.Branchnumber != requesterSession.BranchID {
		ds.Message = "Sorry You Are Not Authorized to Create A round In this ranch "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if round.Courses == nil {
		round.Courses = []entity.Course{}
	}
	// Creating a course
	course := &entity.Course{
		Title:      ds.Title,
		BranchID:   requesterSession.BranchID,
		Duration:   ds.Duration,
		Categoryid: round.CategoryRefer,
	}
	course = coursehandler.CourseService.CreateCourse(course)
	round.Courses = append(round.Courses, *course)
	if course != nil {
		ds.Message = "Course Succesfully Created "
		ds.Success = true
		ds.Course = *course
		ds.Round = *round
	}
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}
