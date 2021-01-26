package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Projects/RidingTrainingSystem/pkg/Helper"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Question"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// QuestionsHandler striuct
type QuestionsHandler struct {
	QuestionService Question.QuestionService
	SessionHandler  *session.Cookiehandler
	StudentService  Student.StudentService
	BranchService   Branch.BranchService
	TemplateHandler *TemplateHandler
}

// NewQuestionHandler function
func NewQuestionHandler(
	QuestionService Question.QuestionService,
	Sessionhandler *session.Cookiehandler,
	Studentservice Student.StudentService,
	BranchService Branch.BranchService,
	TemplateHandler *TemplateHandler,
) *QuestionsHandler {

	return &QuestionsHandler{
		QuestionService: QuestionService,
		SessionHandler:  Sessionhandler,
		StudentService:  Studentservice,
		BranchService:   BranchService,
		TemplateHandler: TemplateHandler,
	}
}

// QuestionsTemplate method serving questions page Filling in the Data
// Method GEt
// Authorization for Students Only
func (questhandler *QuestionsHandler) QuestionsTemplate(response http.ResponseWriter, request *http.Request) {
	session := questhandler.SessionHandler.GetSession(request)
	branch := questhandler.BranchService.GetBranchByID(session.BranchID)
	if branch == nil {
		http.RedirectHandler("/", http.StatusNotFound)
		return
	}
	ds := struct {
		Success   bool
		Message   string
		HOST      string
		Branch    *entity.Branch
		Questions []entity.Question
		Student   *entity.Student
		Lang      string
		Count     int
	}{
		Success: false,
	}
	student := questhandler.StudentService.GetStudentByID(session.ID)
	if student == nil {
		http.RedirectHandler("/", http.StatusNotFound)
		return
	}
	ds.Branch = branch
	ds.Student = student
	ds.HOST = entity.PROTOCOL + entity.HOST
	if student.Lang != "" {
		ds.Lang = student.Lang
	} else {
		ds.Lang = branch.Lang
	}
	// Getting the Questions that the Student take
	questions := questhandler.QuestionService.GetQuestions(session.ID, 0, 10)
	if questions == nil {
		ds.Message = "Welcome To Our Question and Answer Page \n No Question Found In the System"
	} else {
		ds.Questions = *questions
		ds.Success = true
		ds.Count = len(*questions)
		ds.Message = "Welcome To Our Question And Answer Page We Have Offered You " + strconv.Itoa(ds.Count) + " Questins For More Questions Please Press the << See More >> Button"
	}
	questhandler.TemplateHandler.Templates.ExecuteTemplate(response, "question.html", ds)
}

// CreateQuestion method getting the
// Input Variables answer_index  , answers  , body and image
func (questhandler *QuestionsHandler) CreateQuestion(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// JsonValues
	/*
		answer_index
		answers
		body
	*/
	question := &entity.Question{}
	ds := struct {
		Success    bool
		Message    string
		QuestionID uint
	}{
		Success: false,
	}
	var erra error
	body := request.FormValue("body")
	answerIndex, erra := strconv.Atoi(request.FormValue("answer_index"))
	postForm := request.PostForm
	answers := postForm["answers"]
	if erra != nil || len(answers) == 0 {
		ds.Message = " Invalid Request Please Try Again"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if body == "" || answerIndex == 0 || len(answers) < 2 {
		ds.Message = " Missing Question Values  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	question.Body = body
	question.Answers = answers
	question.Answerindex = uint(answerIndex)
	image, header, erra := request.FormFile("image")
	var imageDir string
	if erra == nil && header != nil {
		if !(Helper.IsImage(header.Filename)) {
			ds.Message = "Please Upload Only Image Files"
			response.Write(Helper.MarshalThis(ds))
			return
		}
		RandomFilename := Helper.GenerateRandomString(5, Helper.NUMBERS) + "." + Helper.GetExtension(header.Filename)
		file, erra := os.Create(entity.PathToQuestionImages + RandomFilename)
		if erra != nil {
			ds.Message = "Internal Server Error Please Try Again"
			response.Write(Helper.MarshalThis(ds))
			return
		}
		defer file.Close()
		_, erra = io.Copy(file, image)
		if erra != nil {
			ds.Message = "Internal Server ERROR!"
			response.Write(Helper.MarshalThis(ds))
			return
		}
		imageDir = entity.PathToQuestionImagesFromTemplates + RandomFilename
	}
	if image != nil {
		image.Close()
	}
	question.ImageDirectory = imageDir
	question = questhandler.QuestionService.CreateQuestion(question)
	ds.QuestionID = question.ID
	if question == nil {
		ds.Message = " Error While Creating the Message "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfully Created "
	ds.Success = true
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// DeleteQuestion (ID uint) bool
// Request method POST AND input Question ID this thing is Authorized Only for the
//  ADMINS ( SECRETARY , SUPERADMIN  )  , TEACHERS , TRAINERS ,
func (questhandler *QuestionsHandler) DeleteQuestion(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	request.ParseForm()
	ds := struct {
		Success    bool
		Message    string
		QuestionID uint
	}{
		Success: false,
	}
	questionidstring := request.FormValue("question_id")
	questionID, erra := strconv.Atoi(questionidstring)
	if erra != nil {
		ds.Message = "Invalid Question Information"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	question := questhandler.QuestionService.GetQuestionByID(uint(questionID))
	if question == nil {
		ds.Message = "Record Not Found \nNo Question By this ID !"
		ds.QuestionID = uint(questionID)
		response.Write(Helper.MarshalThis(ds))
		return
	}
	success := questhandler.QuestionService.DeleteQuestion(uint(questionID))
	if success {
		os.Remove("../../web/templates/" + question.ImageDirectory)
		ds.Message = "Succesfully Deleted"
		ds.Success = success
		ds.QuestionID = uint(questionID)
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Can't Delete the Info"
	ds.QuestionID = uint(questionID)
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// GetQuestions method
// method GEt Returning Json Responnse and Success Status Message
//  ROLE STUDENS
// variables offset and limit
func (questhandler *QuestionsHandler) GetQuestions(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	response.Header().Add("Access-Control", "Allow-Origin")
	ds := struct {
		Success   bool
		Message   string
		Questions []entity.Question
	}{
		Success: false,
	}
	limits := request.FormValue("limit")
	limit, era := strconv.Atoi(limits)
	offset, era := strconv.Atoi(request.FormValue("offset"))
	if era != nil || limit <= 0 {
		ds.Message = "Invalid Request Value "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	requesterSession := questhandler.SessionHandler.GetSession(request)
	questions := questhandler.QuestionService.GetQuestions(requesterSession.ID, uint(offset), uint(limit))
	if questions == nil {
		ds.Message = "Sorry !\nCan't Get More Questions "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	response.WriteHeader(http.StatusOK)
	ds.Success = true
	ds.Message = "Succesfully Get Questions "
	ds.Questions = *questions
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetAnswer method for getting the Answer of the WQuestion and save the answered Question Count
// Method get
// Variable QuestionId AnswerId
// Update the Student asked Questions and Answered Questions Count
func (questhandler *QuestionsHandler) GetAnswer(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success    bool
		HasError   bool
		Message    string
		QuestionID uint
		AnswerID   uint
	}{
		HasError: true,
		Success:  false,
	}
	requesterSession := questhandler.SessionHandler.GetSession(request)
	questionidstring := request.FormValue("question_id")
	selectedIndexstring := request.FormValue("selected_index")
	questionID, era := strconv.Atoi(strings.Trim(questionidstring, " "))
	selectedIndex, era := strconv.Atoi(strings.Trim(selectedIndexstring, " "))
	if era != nil || questionID <= 0 || selectedIndex <= 0 {
		ds.Message = "Invalid Request Body "
		ds.HasError = true
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.QuestionID = uint(questionID)
	student := questhandler.StudentService.GetStudentByID(requesterSession.ID)
	if student == nil {
		ds.Message = "Unknown Student!"
		ds.HasError = true
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	answerIndex := questhandler.QuestionService.GetAnswer(uint(questionID))
	studentAskedQuestions := questhandler.QuestionService.GetStudentsAskedQuestionByID(requesterSession.ID)
	if studentAskedQuestions == nil {
		studentAskedQuestions = &entity.AskedQuetion{}
		studentAskedQuestions.Studentid = requesterSession.ID
	}
	studentAskedQuestions.Questionsid = append(studentAskedQuestions.Questionsid, int64(questionID))
	student.AskedQuestionsCount++
	if answerIndex == selectedIndex {
		ds.Success = true
		ds.HasError = false
		ds.Message = "The Answer is Correct"
		ds.AnswerID = uint(answerIndex)
		student.AnsweredQuestionCount++
	} else {
		ds.Success = false
		ds.HasError = false
		ds.Message = "The Answer is InCorrect"
		ds.AnswerID = uint(answerIndex)
	}
	studentAskedQuestions = questhandler.QuestionService.SaveAskedQuestion(studentAskedQuestions)
	student = questhandler.StudentService.SaveStudent(student)
	if student == nil || studentAskedQuestions == nil {
		ds.Message = " Error While Saving the Result "
		ds.Success = false
		ds.HasError = true
	}
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// ShowQuestionResult method
// Method GET
// Authority Student
func (questhandler *QuestionsHandler) ShowQuestionResult(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success       bool
		Message       string
		AnsweredCount uint
		AskedCount    uint
	}{
		Success: false,
	}
	requestersession := questhandler.SessionHandler.GetSession(request)
	student := questhandler.StudentService.GetStudentByID(requestersession.ID)
	if student != nil {
		ds.Success = true
		ds.Message = " You Have Scored " + strconv.Itoa(int(student.AnsweredQuestionCount)) + "/" + strconv.Itoa(int(student.AskedQuestionsCount))
		ds.AskedCount = student.AskedQuestionsCount
		ds.AnsweredCount = student.AnsweredQuestionCount
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = false
	ds.Message = "Can't Generate Result for this student"
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetGradeResultForStudent methdo to get the Score the Student have from the Gived Questions He Take
func (questhandler *QuestionsHandler) GetGradeResultForStudent(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	rsession := questhandler.SessionHandler.GetSession(request)
	if rsession == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	ds := struct {
		Success                bool
		Message                string
		AskedQuestionCount     uint `json:"askedCount"`
		AnsweredQuestionsCount uint `json:"answeredCount"`
		StudentID              uint
	}{
		Success: false,
	}
	asked, answered, success := questhandler.QuestionService.GetGradeResult(rsession.ID)
	if success {
		ds.Success = true
		ds.AskedQuestionCount = asked
		ds.AnsweredQuestionsCount = answered
		ds.StudentID = rsession.ID
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Grade Result Fetch Ws Not Succesful "
	ds.StudentID = rsession.ID
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// ResetResult method
// METHOD GET
// AYTHORIZATION STUDENT
func (questhandler *QuestionsHandler) ResetResult(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	session := questhandler.SessionHandler.GetSession(request)
	if session == nil {
		response.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	ds := struct {
		Success   bool
		Meesage   string
		StudentID uint
	}{
		Success: false,
	}
	studentID := session.ID
	askedquestion := questhandler.QuestionService.GetStudentsAskedQuestionByID(studentID)
	student := questhandler.StudentService.GetStudentByID(studentID)
	if askedquestion == nil {
		askedquestion = &entity.AskedQuetion{}
		askedquestion.Studentid = studentID
		ds.Meesage = "No Asked Questions Yet!"
		ds.StudentID = studentID
		response.WriteHeader(http.StatusOK)
		response.Write(Helper.MarshalThis(ds))
		return
	}
	askedquestion.Questionsid = []int64{}
	student.AskedQuestionsCount = 0
	student.AnsweredQuestionCount = 0
	student = questhandler.StudentService.SaveStudent(student)
	askedquestion = questhandler.QuestionService.SaveAskedQuestion(askedquestion)
	if askedquestion == nil || student == nil {
		ds.Meesage = "Error While Saving the Change"
		ds.StudentID = studentID
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(Helper.MarshalThis(ds))
		return
	}
	ds.StudentID = studentID
	ds.Success = true
	ds.Meesage = "Succescfuly Updated"
	response.WriteHeader(http.StatusOK)
	response.Write(Helper.MarshalThis(ds))
}
