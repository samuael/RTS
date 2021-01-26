package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	"github.com/Projects/RidingTrainingSystem/pkg/translation"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Course"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
	"github.com/gocarina/gocsv"
)

// Routes OF File
/*
	CreateRound
	GetRounds
	GetRoundsOfCategory
	UpdateRound
	AddTrainerToRound
	AddTeacherToRound
	RemoveStudentFromRound
	GetRoundByID
	NewCourseToRound
	AddCourseToRound

*/

// RoundHandler struct
type RoundHandler struct {
	RoundService    Round.RoundService
	SessionService  *session.Cookiehandler
	CategoryService Category.CategoryService
	AdminService    Admin.AdminService
	TrainerService  Trainer.TrainerService
	TeacherService  Teacher.TeacherService
	StudentService  Student.StudentService
	CourseService   Course.CourseService
}

// NewRoundHandler function
func NewRoundHandler(Roundservice Round.RoundService,
	sessions *session.Cookiehandler,
	CategoryService Category.CategoryService,
	adminservice Admin.AdminService,
	TrainerService Trainer.TrainerService,
	TeacherService Teacher.TeacherService,
	StudentService Student.StudentService,
	CourseService Course.CourseService,
) *RoundHandler {
	return &RoundHandler{
		RoundService:    Roundservice,
		SessionService:  sessions,
		CategoryService: CategoryService,
		AdminService:    adminservice,
		TrainerService:  TrainerService,
		TeacherService:  TeacherService,
		StudentService:  StudentService,
		CourseService:   CourseService,
	}
}

/*
	CreateRound : Aking Form and Returning an application/json
	GetRounds : get Request returning jsons of Rounds

*/

//CreateRound method for creating a Round using Form Input
// METHOD POST
// variables categoryid   , active  , _csrf  , training_shift_number  , round_number  ,year  ,
func (roundhandler RoundHandler) CreateRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// Return structure
	ds := struct {
		Success bool
		Message string
		Round   entity.Round
		VErrors form.ValidationErrors
	}{
		Success: false,
	}
	// First Checking to Parse the Form
	erra := request.ParseForm()
	if erra != nil {
		return
	}
	signupForm := form.Input{
		VErrors: form.ValidationErrors{},
		Values:  request.PostForm,
	}
	signupForm.Required("categoryid", "active", "_csrf")
	if csrf := request.FormValue("_csrf"); !roundhandler.SessionService.ValidateForm(csrf) {
		signupForm.VErrors.Add("CSRF", "Invalid CSRF ! You are not allowed to enter the System ")
	}
	if !signupForm.Valid() {
		ds.Message = "Sorrry the request is Invalid "
		ds.VErrors = signupForm.VErrors
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	requesterSession := roundhandler.SessionService.GetSession(request)
	if requesterSession == nil {
		ds.Message = "UnKnown User "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	newActive, era := strconv.ParseBool(request.FormValue("active"))
	if newActive == false && era == nil {
		signupForm.Required("year", "round_number")
	} else if newActive == true && era == nil {
		signupForm.Required("cost", "max_students", "training_duration")
	} else {
		ds.Message = "Round Status Info Active is Incorrect "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	fmt.Println("The Round IS ", newActive)
	trainingShiftADay, erro := strconv.Atoi(request.FormValue("training_shift_number"))
	if newActive && (erro != nil || trainingShiftADay > 2 || trainingShiftADay < 1) {
		signupForm.VErrors.Add("TrainingShiftNumber ", "Invalid TrainingShift Number the value Should Be Number and a value  1 or 2")
	}

	categoryI := request.FormValue("categoryid")
	categoryID, era := strconv.Atoi(categoryI)
	if era != nil {
		ds.Message = "Invalid Category Number"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	category := roundhandler.CategoryService.GetCategoryByID(uint(categoryID))
	if category == nil {
		ds.Message = "No Category By this ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	maxStudent := request.FormValue("max_students")
	maxStudents, era := strconv.Atoi(maxStudent)
	if era != nil && newActive {
		ds.Message = " Invalid Max Student Value! "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	trainingDurationstring := request.FormValue("training_duration")
	trainingDuration, erra := strconv.Atoi(trainingDurationstring)
	if erra != nil && newActive {
		ds.Message = " Invalid Training Duration Value  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	cost := request.FormValue("cost")
	Cost, ers := strconv.Atoi(cost)
	if (ers != nil || Cost <= 0) && newActive {
		ds.Message = " Invalid Cost Value  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	createdby := roundhandler.AdminService.GetAdminByID(requesterSession.ID, "")
	round := &entity.Round{
		Active:           newActive,
		AdminRefer:       requesterSession.ID,
		Branchnumber:     category.Branchid,
		CategoryRefer:    uint(categoryID),
		Category:         *category,
		Cost:             float64(Cost),
		Courses:          []entity.Course{},
		Sections:         []entity.Section{},
		Students:         []entity.Student{},
		Lectures:         []entity.Lecture{},
		Teachers:         []entity.Teacher{},
		Trainers:         []entity.FieldAssistant{},
		OnRegistration:   newActive,
		Year:             uint(etc.NewDate(0).Year),
		TotalPaid:        0,
		MaxStudents:      uint(maxStudents),
		TrainingDuration: uint(trainingDuration),
		// CreatedBY
	}
	round.Roundnumber = category.LastRoundNumber + 1
	if category.Branchid != requesterSession.BranchID {
		ds.Message = " No Such Category In Your Branch"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	if !newActive {
		round.Learning = false
		round.OnRegistration = false
		round.Active = false
		roundNumber, erra := strconv.Atoi(strings.Trim(request.FormValue("round_number"), " "))
		fmt.Println(roundNumber)
		if erra != nil || roundNumber <= 0 {
			// fmt.Println(" Error While Saving the Category ")
			ds.Message = " Invalid Round Number "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		// Checking whether the round  number is reserved or not
		reserved := roundhandler.RoundService.IsRoundNumberReseerved(category.Branchid, uint(categoryID), uint(roundNumber))
		if reserved {
			ds.Message = " Round number Is Reserved "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		round.Roundnumber = uint(roundNumber)
		if category.LastRoundNumber < round.Roundnumber {
			category.LastRoundNumber = round.Roundnumber
		}

		year, era := strconv.Atoi(request.FormValue("year"))
		if era != nil || year <= 1900 || year >= int(etc.ETYear())+2 {
			ds.Message = fmt.Sprintf(" Invalid Year Input ..%d.. for Passive Round", year)
			response.Write(Helper.MarshalThis(ds))
			return
		}
		round.Year = uint(year)
	} else {
		category.LastRoundNumber++
	}
	result := roundhandler.CategoryService.UpdateCategoryLasrRoundCount(category, category.ID, category.LastRoundNumber)
	if !result {
		fmt.Println(" Error While Saving the Category ")
		ds.Message = " Error While Saving the Category "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if createdby != nil {
		round.CreatedBY = *createdby
	}
	round = roundhandler.RoundService.CreateRound(round)
	if round != nil {
		ds.Round = *round
		ds.Success = true
		ds.Message = "Round Succesfully Created "
		ds.VErrors = nil
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Sorry Can't Create the Round "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// GetRounds method  returning a list of branchs   id: as a parameter
func (roundhandler *RoundHandler) GetRounds(response http.ResponseWriter, request *http.Request) {
	// the Method should Be GEt and the Return will be slice of branches
	response.Header().Add("Content-Type", "application/json")
	if request.Method == http.MethodGet {
		ds := struct {
			Success bool
			Message string
			Rounds  []entity.Round
		}{
			Success: false,
			Rounds:  []entity.Round{},
		}
		requestersession := roundhandler.SessionService.GetSession(request)
		if requestersession == nil {
			ds.Message = "Unknown User "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		branchID := requestersession.BranchID
		rounds := roundhandler.RoundService.GetRounds(branchID)
		if rounds == nil {
			ds.Message = "Can't Find Rounds In this Branch "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		ds.Success = true
		ds.Message = " Succesfull F,etched The Rounds "
		ds.Rounds = *rounds
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	} else {
		// This method is not allowed niggoye
	}
}

// GetRoundsOfCategory of Category
// Method GET
// Authority SECRETARY and SUPERADMIN
func (roundhandler *RoundHandler) GetRoundsOfCategory(response http.ResponseWriter, request *http.Request) {
	requesterSession := roundhandler.SessionService.GetSession(request)
	if requesterSession == nil {
		response.Write([]byte("<h1> Not Authorized </h1>"))
	}
	response.Header().Add("Content-Type", "application/json")
	categoryidstring := request.FormValue("category_id")
	isActive := request.FormValue("active")
	ds := struct {
		Success bool
		Message string
		Rounds  []entity.Round
	}{
		Success: false,
	}
	categoryID, era := strconv.Atoi(categoryidstring)
	Active, era := strconv.ParseBool(isActive)
	if era != nil {
		ds.Message = "Error While Parsing the Get Values "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if categoryID <= 0 {
		ds.Message = "Invalid Category ID Values "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	var rounds *[]entity.Round
	if Active {
		rounds = roundhandler.RoundService.GetActiveRoundsOfCategory(uint(categoryID))
	} else {
		rounds = roundhandler.RoundService.GetRoundsOfCategory(uint(categoryID))
	}
	if rounds == nil {
		ds.Message = " No Record Found By Thid Category ID  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Rounds = *rounds
	ds.Success = true
	ds.Message = "Successfully Fetched  " + strconv.Itoa(len(*rounds)) + " Rounds "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// UpdateRound method post authority SUPERADMIN
//  The form inpus is Json and the Out put also Json
func (roundhandler *RoundHandler) UpdateRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success bool
		Message string
		Round   entity.Round
	}{
		Success: false,
		Round:   entity.Round{},
	}
	newDecoder := json.NewDecoder(request.Body)
	round := &entity.Round{}
	decodeError := newDecoder.Decode(round)
	if decodeError != nil {
		ds.Message = "InValid Input Type "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if round.ID <= 0 {
		ds.Message = "Can't Update Round of ID " + strconv.Itoa(int(round.Roundnumber))
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if round.Sections != nil || round.Students != nil || round.Trainers != nil || round.Cost <= 0 {
		ds.Message = "InValid Request "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	// Getting the Round bt the Specified Round ID
	roundnew := roundhandler.RoundService.GetRoundByID(round.ID)
	if roundnew == nil {
		ds.Message = "No Round By This ID  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if round.Cost != 0 {
		roundnew.Cost = round.Cost
	}
	if round.Roundnumber != 0 {
		roundnew.Roundnumber = round.Roundnumber
	}
	if round.MaxStudents != 0 {
		roundnew.MaxStudents = round.MaxStudents
	}

	roundnew = roundhandler.RoundService.SaveRound(roundnew)
	if roundnew == nil {
		ds.Message = "Can't Update the round  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = " Succesfully Updated"
	ds.Round = *roundnew
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// AddTrainerToRound method
func (roundhandler *RoundHandler) AddTrainerToRound(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		// signupForm := form.Input{
		// 	VErrors: form.ValidationErrors{},
		// 	Values:  request.Form,
		// }
		// signupForm.Required("round_id", "trainer_id"  , )
		response.Header().Add("Content-Type", "application/jsons")
		ds := struct {
			Success   bool
			Message   string
			Roundid   uint `json:"round_id"`
			Trainerid uint `json:"trainer_id"`
		}{
			Success: false,
			Message: "Not Succesfull ",
		}
		requesterSession := roundhandler.SessionService.GetSession(request)
		if requesterSession == nil {
			ds.Message = "Unknown User "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		newDecoder := json.NewDecoder(request.Body)
		decodeError := newDecoder.Decode(&ds)
		if decodeError != nil {
			ds.Message = "Input Error Pleae Use Appropriate "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		if ds.Roundid <= 0 || ds.Trainerid <= 0 {
			ds.Message = "Invalid Trainer or Round Number "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}

		round := roundhandler.RoundService.GetRoundByID(ds.Roundid)
		fmt.Println(round)
		trainer := roundhandler.TrainerService.GetTrainerByID(ds.Trainerid)
		if trainer == nil || round == nil {
			ds.Message = "Can't Find The Trainer || Round  "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		round.Trainers = append(round.Trainers, *trainer)
		round = roundhandler.RoundService.SaveRound(round)
		if round == nil {
			ds.Message = "Can't Update the Round "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		ds.Success = true
		ds.Message = " Succesfuly Updated the Round "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	} else {
		//  Other method is not allowed
	}
}

// GetRoundByID method for Specific Round
func (roundhandler *RoundHandler) GetRoundByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success bool
		Message string
		Round   entity.Round
	}{
		Success: false,
	}
	roundNumber := request.FormValue("round_id")
	roundno, era := strconv.Atoi(roundNumber)
	if era != nil {
		ds.Message = "Incorrect Round Number"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	round := roundhandler.RoundService.GetRoundByID(uint(roundno))
	if round == nil {

		ds.Message = "Can't Find Round By this ID "

	} else {
		ds.Success = true
		ds.Message = "Round is Found "
		ds.Round = *round
	}
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// AddTeacherToRound  method to add a Teacher to a round
func (roundhandler *RoundHandler) AddTeacherToRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success   bool
		Message   string
		Round     entity.Round
		Roundid   uint `json:"round_id"`
		Teacherid uint `json:"teacher_id"`
	}{
		Success: false,
	}
	newDecoder := json.NewDecoder(request.Body)
	jsonError := newDecoder.Decode(&ds)
	if jsonError != nil {
		ds.Message = "Error Invalid Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if ds.Teacherid <= 0 || ds.Roundid <= 0 {
		ds.Message = "Error Invalid Data type"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	teacher := roundhandler.TeacherService.GetTeacherByID(ds.Teacherid)
	round := roundhandler.RoundService.GetRoundByID(ds.Roundid)
	if teacher == nil || round == nil {
		ds.Message = "Undefined Id Number"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	hasTeacher := false
	for _, teachero := range round.Teachers {
		if teachero.ID == teacher.ID {
			hasTeacher = true
		}
	}
	if hasTeacher {
		ds.Message = "The Teacher Does exist "
		ds.Round = *round
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	round.Teachers = append(round.Teachers, *teacher)
	round = roundhandler.RoundService.SaveRound(round)
	if round == nil {
		ds.Message = "Internal Server Error"
	} else {
		ds.Success = true
		ds.Message = "Teacher Succesfuly added to a round "
		ds.Round = *round
	}
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// RemoveStudentFromRound method to remove a student from a round
func (roundhandler *RoundHandler) RemoveStudentFromRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requesterSession := roundhandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		Studentid uint `json:"student_id"`
		Roundid   uint `json:"round_id"`
	}{
		Success: false,
	}
	newDecoder := json.NewDecoder(request.Body)
	jsonError := newDecoder.Decode(&ds)
	if jsonError != nil {
		ds.Message = "Invalid Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if ds.Studentid <= 0 || ds.Roundid <= 0 {
		ds.Message = "Invalid Student || Round Number "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	student := roundhandler.StudentService.GetStudentByID(ds.Studentid)
	round := roundhandler.RoundService.GetRoundByID(ds.Roundid)
	if student == nil || round == nil {
		ds.Message = "Unknown Record "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	removed := false
	students := []entity.Student{}
	// Deleting the Student from a round
	for _, stus := range round.Students {
		if student.ID == stus.ID {
			removed = true
			continue
		}
		students = append(students, stus)
	}
	if !removed {
		ds.Message = "The student don't exist in this round"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	round.Students = students
	round = roundhandler.RoundService.SaveRound(round)
	if round == nil {
		ds.Message = "Error While Updating the Round"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Succesfully Removed The Round"
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// NewCourseToRound method to add a course and to the Round also
// Method GET
// Input Post Form Value
//   title  , duration  , round_id
func (roundhandler *RoundHandler) NewCourseToRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requesterSession := roundhandler.SessionService.GetSession(request)
	if requesterSession == nil {
		response.Header().Add("Content-Type", "text/html")
		response.Write([]byte("<h1>Invalid User :</h1>"))
		return
	}

	ds := struct {
		Success bool
		Message string
		Course  entity.Course
	}{
		Success: false,
	}
	// newDecoder := json.NewDecoder(request.Body)
	request.ParseForm()
	title := request.PostFormValue("title")
	duration := request.PostFormValue("duration")
	roundIDString := request.PostFormValue("round_id")
	roundID, era := strconv.Atoi(roundIDString)
	Duration, era := strconv.Atoi(duration)
	if era != nil || title == "" || duration == "" || roundID <= 0 || Duration <= 1 {
		ds.Message = " Invalid Input Variables "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	round := roundhandler.RoundService.GetRoundByID(uint(roundID))
	if round == nil {
		ds.Message = "No Such Round "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	course := entity.Course{
		Duration:   uint(Duration),
		BranchID:   requesterSession.BranchID,
		Title:      title,
		Categoryid: round.CategoryRefer,
	}
	round.Courses = append(round.Courses, course)
	round = roundhandler.RoundService.SaveRound(round)
	if round == nil {
		ds.Message = "Error While Saving the Course to a round "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Course Succesfuly Created"
	ds.Course = course
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// AddCourseToRound method returning the Round Struct with the Round Info
func (roundhandler *RoundHandler) AddCourseToRound(response http.ResponseWriter, request *http.Request) {
	//  The Method is Get Oly
	// requesterSeddion := roundhandler.SessionService.GetSession(request)
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		Round    entity.Round
		Roundid  uint `json:"roundid"`
		Courseid uint `json:"courseid"`
	}{
		Success: false,
		Round:   entity.Round{},
	}
	newDecoder := json.NewDecoder(request.Body)
	decodeError := newDecoder.Decode(&ds)
	if decodeError != nil {
		ds.Message = "can't Decode ! Invalid Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if ds.Roundid <= 0 || ds.Courseid <= 0 {
		ds.Message = "Invalid Request"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	course := roundhandler.CourseService.GetCourseByID(uint(ds.Courseid))
	round := roundhandler.RoundService.GetRoundByID(uint(ds.Roundid))
	if course == nil || round == nil {
		ds.Message = "Can't Find Respective Record "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	for _, courso := range round.Courses {
		if courso.ID == course.ID {
			fmt.Println(courso.ID)
			ds.Message = " The Course Already Exit in the Round "
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
	}
	round.Courses = append(round.Courses, *course)
	round = roundhandler.RoundService.SaveRound(round)
	if round == nil {
		ds.Message = " Error While Updating the Round  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Round = *round
	ds.Message = " Course Succesfully Added "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetPassiveRoundPopulatingForm having a header first Entry Filled
// Method GET
// Authorization for SUPERADMIN  , and SECRETARIES
// input variable file_type for now will be only csv file but after modofication it will decide accordingly their machine
// machie : win , lin
// lang  : eng , amh
func (roundhandler *RoundHandler) GetPassiveRoundPopulatingForm(response http.ResponseWriter, request *http.Request) {
	headers := []string{}
	lang := request.FormValue("lang")
	if lang != "" {
		for _, value := range entity.DynamicRoundPopulatingHeadersEnglish {
			headers = append(headers, translation.Translate(lang, value))
		}
	} else {
		headers = entity.DynamicRoundPopulatingHeadersEnglish
	}
	file, era := os.OpenFile(entity.PathToStaticTemplateFiles+entity.RoundTemplatingCSV, os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if era != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("<small>Internal Server Error! Please Try Again!</small>"))
		return
	}
	defer file.Close()
	CsvEncoder := csv.NewWriter(file)
	// fmt.Println(headers)
	HeaderWritingError := CsvEncoder.Write(headers)
	CsvEncoder.Flush()
	HeaderWritingError = CsvEncoder.Error()
	if HeaderWritingError != nil {
		response.WriteHeader(http.StatusInternalServerError)
		// fmt.Println(HeaderWritingError.Error())
		response.Write([]byte("<small>Internal Server Error! Please Try Again! ..... </small>"))
		return
	}
	reedSeeker := io.ReadSeeker(file)
	http.ServeContent(response, request, entity.PathToStaticTemplateFiles+entity.RoundTemplatingCSV, time.Now(), reedSeeker)
	file.Close()
	ioutil.WriteFile(entity.PathToStaticTemplateFiles+entity.RoundTemplatingCSV, []byte{}, os.ModePerm)
}

// PopulateRoundUsingXlsx method to take input from xslx Input
// variables file , round_id
// Authorization SECREETARY , SUPERADMINS ,
func (roundhandler *RoundHandler) PopulateRoundUsingXlsx(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success          bool
		Message          string
		RoundID          uint
		RoundNumber      uint
		NumberOfStudents uint
		Objections       []string
	}{
		Success: false,
	}
	Objections := []string{}
	// inputs file.xslx   round id  or ods
	file, header, era := request.FormFile("file")
	if era != nil || header == nil {
		ds.Message = "No File Is Found Please Upload the file Containing the appropriate names of a students "
		response.Write(Helper.MarshalThis(ds))
		return
	}
	filename := header.Filename
	extension := Helper.GetExtension(filename)
	pass := false
	switch strings.ToLower(extension) {
	case "ods", "xlsx", "xls", "csv":
		pass = true
	}
	if !pass {
		ds.Message = "Invalid File Type"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	roundID, era := strconv.Atoi(request.FormValue("round_id"))
	if era != nil || roundID <= 0 {
		ds.Message = "Invalid Round ID!"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	ds.RoundID = uint(roundID)
	round := roundhandler.RoundService.GetRoundByID(uint(roundID))
	if round == nil {
		ds.Message = "No Record Found "
		response.Write(Helper.MarshalThis(ds))
		return
	}
	if round.Active {
		ds.Message = " You Cant Populate Active Round Using Spreed Sheet Files"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	ds.RoundNumber = round.Roundnumber
	// Starting reading the file
	var studentscsv *[]entity.StudentCsv
	studentscsv = &[]entity.StudentCsv{}
	newFilename := entity.PathToStaticTemplateFiles + Helper.GenerateRandomString(5, Helper.CHARACTERS) + "." + Helper.GetExtension(header.Filename)
	fileNew, era := os.Create(newFilename)
	if era != nil {
		ds.Message = "Internal Server Error .."
		response.Write(Helper.MarshalThis(ds))
		return
	}
	_, era = io.Copy(fileNew, file)
	if era != nil {
		ds.Message = "Internal Server Error .. .."
		response.Write(Helper.MarshalThis(ds))
		return
	}
	file.Close()
	// Closign the file and Openning a new file
	fileNew.Close()
	fileNew, era = os.OpenFile(newFilename, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm|os.ModeDir)
	if era != nil {
		ds.Message = "Internal Server Error ..open"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	// defer fileNew.Close()
	defer os.Remove(newFilename)
	// Replacing the First Line To Generate the Schedule And Unmarshal the file
	NewReader := csv.NewReader(fileNew)
	newSlice, er := NewReader.ReadAll()
	// fmt.Println(newSlice)
	if len(newSlice) <= 0 || er != nil {
		if len(newSlice) == 1 {
			ds.Message = "No Students Data Is Added"
		} else if len(newSlice) == 0 {
			ds.Message = " Invalid File. Please Use the Appropriate File Inputing CSV file "
		} else {
			ds.Message = " Invalid Input , System was not able to Read the File "
		}
		response.Write(Helper.MarshalThis(ds))
		return
	}
	newSlice[0] = []string{"fullname", "sex", "age", "academic_status", "lang", "city", "kebele", "phone"}
	NewReader = nil
	// Changing the Header name End
	NewWriter := csv.NewWriter(fileNew)
	ioutil.WriteFile(newFilename, []byte{}, os.ModePerm)
	NewWriter.WriteAll(newSlice)
	NewWriter.Flush()
	fileNew.Close()
	fileNew, _ = os.OpenFile(newFilename, os.O_RDWR, os.ModePerm)
	NewWriter = nil
	erro := gocsv.UnmarshalFile(fileNew, studentscsv)
	if erro != nil || len(*studentscsv) == 0 {
		ds.Message = "Invalid File Body ! Please Use Appropriate File With Safe File Headers"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	var students []entity.Student
	var exception error
	for index, studentcsv := range *studentscsv {
		student := entity.Student{}
		studentnames := strings.Split(strings.Trim(studentcsv.Fullname, " "), " ")
		if len(studentnames) == 3 {
			student.Firstname = studentnames[0]
			student.Lastname = studentnames[1]
			student.GrandFatherName = studentnames[2]
		} else if len(studentnames) == 2 {
			student.Firstname = studentnames[0]
			student.Lastname = studentnames[1]
			Objections = append(Objections, fmt.Sprintf("Grand Fathers name is Missing For %d Student ", index))
		} else if len(studentnames) == 1 {
			student.Firstname = studentnames[0]
			Objections = append(Objections, fmt.Sprintf("  Only First Name Is inputed In the  %d Student ", index))
		} else {
			Objections = append(Objections, fmt.Sprintf(" Student name is not Mentioned For %d Student ", index))
			exception = fmt.Errorf("Name is Not Mentioned For Student %d", index)
			break
		}
		student.Sex = studentcsv.Sex
		student.Username = fmt.Sprintf("%s/%d/%d/%d", round.Category.Title, round.Roundnumber, index+1, etc.ETYear())
		// student.Category = round.Category
		student.CategoryID = round.CategoryRefer
		student.AcademicStatus = studentcsv.AcademicStatus
		student.Address = &entity.Address{}
		student.Address.City = studentcsv.City
		student.Address.Kebele = studentcsv.Kebele
		student.Lang = studentcsv.Lang
		student.PhoneNumber = studentcsv.Phone
		student.RoundRefer = round.ID
		student.BranchID = round.Branchnumber
		student.Active = false
		student.BirthDate = &etc.Date{Day: int(etc.ETDay()), Month: int(etc.ETMonth()), Year: int(etc.ETYear()) - studentcsv.Age}
		students = append(students, student)

	}
	ds.Objections = Objections
	if exception != nil {
		// ERROR DO SOME THING
		ds.Message = exception.Error()
		response.Write(Helper.MarshalThis(ds))
		return
	}
	if len(students) == 0 {
		ds.Message = "No Student Record Is Found "
		response.Write(Helper.MarshalThis(ds))
		return
	}
	success := roundhandler.StudentService.SaveStudents(int(round.ID), &students)
	newStudentsCount := len(students)
	roundhandler.RoundService.UpdateStudentsCount(ds.RoundID, uint(newStudentsCount))
	if !success {
		ds.Message = "Internal Server Error"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	ds.Message = "Succesfuly Saved "
	ds.Success = true
	ds.NumberOfStudents = uint(len(students))
	ds.Objections = Objections
	response.Write(Helper.MarshalThis(ds))
}

// GetStudentsOfRound for Round
// Method GET
// Authority SUPERADMIN  , SECRETAART  , TEACHERS  , FIELDASSISTANTS
// VARIABELe round_id
func (roundhandler *RoundHandler) GetStudentsOfRound(response http.ResponseWriter, request *http.Request) {
	requestersession := roundhandler.SessionService.GetSession(request)
	BranchID := requestersession.BranchID
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		Students []entity.Student
	}{
		Success: false,
	}
	roundID, era := strconv.Atoi(request.FormValue("round_id"))
	if era != nil || roundID <= 0 {
		ds.Message = "Invalid Round Number Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	round := roundhandler.RoundService.GetRoundByID(uint(roundID))
	if round == nil {
		ds.Message = "No Record By This Round ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if round.Branchnumber != BranchID {
		ds.Message = " You Are not Allowed To Access this Page "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	students := roundhandler.StudentService.StudentsOfRound(round.ID)
	if students == nil {
		ds.Message = "No Students Found For this Round "
		ds.Students = []entity.Student{}
		response.Write(Helper.MarshalThis(ds))
		return
	}
	ds.Students = *students
	// ds.Students = round.Students
	ds.Success = true
	ds.Message = " Succesful "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}
