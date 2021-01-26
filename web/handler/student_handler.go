package handler

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	"github.com/Projects/RidingTrainingSystem/pkg/HtmlToPDF"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"

	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

/*
* GetRoundStudentsImageZip
*TheoryTestingForm
*  StudentBiographyPDF
* GenerateTreatyPage
*  Login
*  ChangeProfilePicture
* StudentsIDCardGenerate
**/

// StudentHandler y
type StudentHandler struct {
	StudentService  Student.StudentService
	TempHandler     *TemplateHandler
	SessionSeervice *session.Cookiehandler
	BranchService   Branch.BranchService
	AdminService    Admin.AdminService
	RoundService    Round.RoundService
}

// NewStudentHandler returning new Instance of Student Handler
func NewStudentHandler(
	StudentSer Student.StudentService,
	th *TemplateHandler,
	SessionSeervice *session.Cookiehandler,
	BranchService Branch.BranchService,
	AdminService Admin.AdminService,
	RoundService Round.RoundService,
) *StudentHandler {
	return &StudentHandler{
		StudentService:  StudentSer,
		TempHandler:     th,
		SessionSeervice: SessionSeervice,
		BranchService:   BranchService,
		AdminService:    AdminService,
		RoundService:    RoundService,
	}
}

// Login student Login Controll
func (studenthandler *StudentHandler) Login(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		fmt.Println("POSTMEThod in Student/login")
		return
	}
	http.Redirect(response, request, "/", http.StatusPermanentRedirect)
}

// GenerateTreatyPage   method
// Method GEt
// Authorized  SUPERADMIN  , SECRETARY
// Variables student_id
func (studenthandler *StudentHandler) GenerateTreatyPage(response http.ResponseWriter, request *http.Request) {
	requestersession := studenthandler.SessionSeervice.GetSession(request)
	if requestersession == nil {
		response.Write([]byte("<h1>   Unknown User    </h1>"))
		return
	}
	ds := struct {
		Success bool
		Message string
		Branch  entity.Branch
		Student entity.Student
		Host    string
	}{
		Success: false,
		Host:    entity.PROTOCOL + entity.HOST,
	}
	branch := studenthandler.BranchService.GetBranchByID(requestersession.BranchID)
	studentID, era := strconv.Atoi(request.FormValue("student_id"))
	if era != nil || studentID <= 0 {
		ds.Message = " Invalid Student ID "
		ds.Branch = *branch
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_treaty_page.html", ds)
		return
	}
	student := studenthandler.StudentService.GetStudentByID(uint(studentID))
	if student == nil {
		ds.Message = " Record Not Found "
		ds.Branch = *branch
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_treaty_page.html", ds)
		return
	}
	ds.Success = true
	ds.Student = *student
	ds.Branch = *branch
	ds.Message = "Succesful"
	randomstring := entity.PathToPdfs + Helper.GenerateRandomString(4, Helper.CHARACTERS) + ".html"
	htmlfile, era := os.Create(randomstring)
	if era != nil {
		ds.Message = " Internal Server Error "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_treaty_page.html", ds)
		return
	}
	studenthandler.TempHandler.Templates.ExecuteTemplate(htmlfile, "student_treaty_page.html", ds)
	htmlfile.Close()
	thefileDirectory := HtmlToPDF.GetThePdf(randomstring)
	if thefileDirectory == "" {
		ds.Message = "  Error While Generating the File "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_treaty_page.html", ds)
		return
	}
	thePDfFileDiractoryForStudentProfile, er := os.Open(thefileDirectory)
	if er != nil || thePDfFileDiractoryForStudentProfile == nil {
		ds.Message = "  Error While Generating the File "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_treaty_page.html", ds)
		return
	}
	info, _ := thePDfFileDiractoryForStudentProfile.Stat()
	response.Header().Set("Content-Type", "application/pdf")
	response.Header().Set("Content-Deposition", "attachment ; filename="+info.Name())
	response.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(response, thePDfFileDiractoryForStudentProfile)
	thePDfFileDiractoryForStudentProfile.Close()
	os.Remove(thefileDirectory)

}

// StudentBiographyPDF method to gnerate pdf for student Biography
// Method GET
// Variale student_id
func (studenthandler *StudentHandler) StudentBiographyPDF(response http.ResponseWriter, request *http.Request) {
	requestersession := studenthandler.SessionSeervice.GetSession(request)
	if requestersession == nil {
		response.Write([]byte("Un-Authorized "))
		return
	}
	branch := studenthandler.BranchService.GetBranchByID(requestersession.BranchID)
	if branch == nil {
		response.Write([]byte("Not-Authorized "))
	}
	studentidstring := request.FormValue("studentid")
	studentID, era := strconv.Atoi(studentidstring)
	ds := struct {
		Success bool
		Message string
		Host    string
		Student entity.Student
		Branch  entity.Branch
	}{
		Success: false,
		Branch:  *branch,
		Host:    entity.PROTOCOL + entity.HOST,
	}
	if era != nil || studentID <= 0 {
		ds.Message = " Invalid Input Value Please Enter appropriate Value  "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_biography_template.html", ds)
		return
	}
	student := studenthandler.StudentService.GetStudentByID(uint(studentID))
	if student == nil {
		ds.Message = "  Record Not Found No Student By This ID"
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_biography_template.html", ds)
		return
	}
	randomstring := entity.PathToPdfs + Helper.GenerateRandomString(4, Helper.CHARACTERS) + ".html"
	htmlfile, era := os.Create(randomstring)
	if era != nil {
		ds.Message = " Internal Server Error "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_biography_template.html", ds)
		return
	}
	ds.Student = *student
	ds.Success = true
	studenthandler.TempHandler.Templates.ExecuteTemplate(htmlfile, "student_biography_template.html", ds)
	htmlfile.Close()
	thefileDirectory := HtmlToPDF.GetThePdf(randomstring)
	if thefileDirectory == "" {
		ds.Message = "  Error While Generating the File "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_biography_template.html", ds)
		return
	}
	thePDfFileDiractoryForStudentProfile, er := os.Open(thefileDirectory)
	if er != nil || thePDfFileDiractoryForStudentProfile == nil {
		ds.Message = "  Error While Generating the File "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_biography_template.html", ds)
		return
	}

	info, _ := thePDfFileDiractoryForStudentProfile.Stat()
	response.Header().Set("Content-Type", "application/pdf")
	response.Header().Set("Content-Deposition", "attachment ; filename="+info.Name())
	response.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(response, thePDfFileDiractoryForStudentProfile)
	thePDfFileDiractoryForStudentProfile.Close()
	os.Remove(thefileDirectory)
}

// TheoryTestingForm  method to Generate Theory Testing Form For Each Students
// Satisfying the Given Criteria
// Method Get
// Parameters round_id  , Payment_limit
func (studenthandler *StudentHandler) TheoryTestingForm(response http.ResponseWriter, request *http.Request) {
	requesterSession := studenthandler.SessionSeervice.GetSession(request)
	if requesterSession == nil {
		response.Write([]byte("<h1>Not Authorized </h1>"))
		return
	}

	admin := studenthandler.AdminService.GetAdminByID(requesterSession.ID, "")
	branch := studenthandler.BranchService.GetBranchByID(requesterSession.BranchID)
	if admin == nil || branch == nil {
		response.Write([]byte("<h1>Not Authorized </h1>"))
		return
	}
	roundidstring := request.FormValue("round_id")
	paymentlimitstring := request.FormValue("payment_limit")
	ds := struct {
		Success  bool
		Message  string
		Admin    entity.Admin
		Branch   entity.Branch
		Students []entity.Student
		Host     string
	}{
		Success: false,
		Branch:  *branch,
		Admin:   *admin,
		Host:    entity.PROTOCOL + entity.HOST,
	}
	roundID, era := strconv.Atoi(roundidstring)
	PaymentLimit, era := strconv.ParseFloat(paymentlimitstring, 64)
	fmt.Println(roundID, PaymentLimit)
	if era != nil || roundID <= 0 || PaymentLimit < 0 {
		ds.Message = "Invalid Request Values  "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_theory_testing_form.html", ds)
		return
	}
	students := studenthandler.StudentService.GetStudentPaidMoreThanPaymentLimit(uint(roundID), float64(PaymentLimit))
	if students == nil {
		ds.Message = fmt.Sprintf(" Internal Server Error ")
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_theory_testing_form.html", ds)
		return
	}
	randomstring := entity.PathToPdfs + Helper.GenerateRandomString(4, Helper.CHARACTERS) + ".html"
	htmlfile, era := os.Create(randomstring)
	if era != nil {
		ds.Message = " Internal Server Error "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_theory_testing_form.html", ds)
		return
	}
	ds.Students = *students
	ds.Success = true
	studenthandler.TempHandler.Templates.ExecuteTemplate(htmlfile, "student_theory_testing_form.html", ds)
	htmlfile.Close()
	thefileDirectory := HtmlToPDF.GetThePdf(randomstring)
	if thefileDirectory == "" {
		ds.Message = "  Error While Generating the File "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_theory_testing_form.html", ds)
		return
	}
	thePDfFileDiractoryForStudentProfile, er := os.Open(thefileDirectory)
	if er != nil || thePDfFileDiractoryForStudentProfile == nil {
		ds.Message = "  Error While Generating the File "
		studenthandler.TempHandler.Templates.ExecuteTemplate(response, "student_theory_testing_form.html", ds)
		return
	}
	info, _ := thePDfFileDiractoryForStudentProfile.Stat()
	response.Header().Set("Content-Type", "application/pdf")
	response.Header().Set("Content-Deposition", "attachment ; filename="+info.Name())
	response.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(response, thePDfFileDiractoryForStudentProfile)
	thePDfFileDiractoryForStudentProfile.Close()
	os.Remove(thefileDirectory)
}

// GetRoundStudentsImageZip  method returning the List of Image of students ordered By their Name
// Method GET
// variables round_id
func (studenthandler *StudentHandler) GetRoundStudentsImageZip(response http.ResponseWriter, request *http.Request) {
	requestersession := studenthandler.SessionSeervice.GetSession(request)
	if requestersession == nil {
		response.WriteHeader(http.StatusProxyAuthRequired)
		response.Write([]byte("<h1>  invalid Request </h1>"))
		return
	}
	roundID, erra := strconv.Atoi(request.FormValue("round_id"))
	if erra != nil || roundID <= 0 {
		fmt.Println(" Invalid Input ")
		response.Write([]byte("<h1  style=\"color:#006699;  \"> Invalid Round ID  </h1>"))
		return
	}
	round := studenthandler.RoundService.GetRoundByID(uint(roundID))
	if round == nil {
		fmt.Println(" No Round By This ID  ")
		response.Write([]byte("<h1  style=\"color:#006699;text-align:center;vertical-align:center; \"> No Round By This ID </h1>"))
		return
	}
	if len(round.Students) <= 0 {
		fmt.Println(" No Students To Deal With Input ")
		response.Write([]byte("<h1  style=\"color:#006699;text-align:center;vertical-align:center; \"> No Students In This Round </h1>"))
		return
	}

	var WasZip bool // telling whether there was a for the Round Or Not
	Counter := 1
	// var ZipWriter *zip.Writer
	zipFileName := entity.PATHToZipFiles + strconv.Itoa(int(round.Roundnumber)) + "_" + "Students_images" + ".zip"
	zipfile, erra := os.Open(zipFileName)
	WasZip = true
	if erra != nil {
		erra = nil
		zipfile, erra = os.Create(zipFileName)
		WasZip = false
	}
	if erra != nil {
		fmt.Println(" Error While Creating File  ")
		response.Write([]byte("<h1  style=\"color:#006699;text-align:center;vertical-align:center; \"> Error While Creating File  </h1>"))
		return
	}
	defer zipfile.Close()

	if !WasZip {
		ZipWriter := zip.NewWriter(zipfile)
		defer ZipWriter.Close()
		for j := 0; j < len(round.Students); j++ {
			student := &round.Students[j]
			imagename := strconv.Itoa(Counter) + "_" + student.Firstname + "_" + student.Lastname + "_" + student.GrandFatherName + "."
			imageDirectoryOfTheStudents := student.Imageurl
			extension := Helper.GetExtension(imageDirectoryOfTheStudents)
			imagename += extension
			errors := Helper.AddFileToZip(ZipWriter, entity.PathToTemplates+imageDirectoryOfTheStudents, imagename)
			if errors != nil {
				fmt.Println(errors.Error())
				fmt.Println(" Error While Adding the File to the Zip ")
				response.Write([]byte("<h1  style=\"color:#006699;text-align:center;vertical-align:center; \"> Internal Server Error AdFileToZip </h1>"))
				return
			}
			student.ImageNumber = string(Counter)
			Counter++
		}
	}
	round = studenthandler.RoundService.SaveRound(round)
	if round == nil {
		fmt.Println(" Error While Saving The Round  ")
		response.Write([]byte("<h1  style=\"color:#006699;text-align:center;vertical-align:center; \"> Internal Server Error  </h1>"))
		return
	}
	zipinfo, eras := zipfile.Stat()
	if eras != nil {
		fmt.Println(erra.Error())
		return
	}
	response.Header().Set("Content-Type", "application/zip")
	response.Header().Set("Content-Deposition", "attachment;filename="+zipinfo.Name())
	response.Header().Set("Content-Length", strconv.FormatInt(zipinfo.Size(), 10))
	io.Copy(response, zipfile)
}

// StudentsOfCategory method
// Method GEt
// Authorization for all Logged In Users
// Varaible category_id  , offset  , limit , active
func (studenthandler *StudentHandler) StudentsOfCategory(response http.ResponseWriter, request *http.Request) {
	categoryID, era := strconv.Atoi(request.FormValue("category_id"))
	Offset, era := strconv.Atoi(request.FormValue("offset"))
	Limit, era := strconv.Atoi(request.FormValue("limit"))
	BranchID := studenthandler.SessionSeervice.GetSession(request).BranchID
	Active, era := strconv.ParseBool(request.FormValue("active"))
	ds := struct {
		Success  bool
		Message  string
		Students []entity.Student
	}{
		Success: false,
	}
	if era != nil || categoryID <= 0 || Offset < 0 || Limit == 0 {
		ds.Message = " Invalid Request Values "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	students := studenthandler.StudentService.GetStudentsOfCategory(BranchID, uint(categoryID), uint(Offset), uint(Limit), Active)
	if students == nil {
		ds.Message = " Record ot Found "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
}

// ChangeProfilePicture  method
// AUTHORIZATION  Only For The Student
// Input FormValue Output  Json
// Method POST
func (adminh *AdminHandler) ChangeProfilePicture(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestersession := adminh.SessionHandler.GetSession(request)
	if requestersession == nil || requestersession.ID <= 0 {
		fmt.Println("Bla Ble blu ")
		return
	}
	// fmt.Printf("User %s \n ID : %d\n Username : %s\n", requestersession.Role, requestersession.ID, requestersession.Username)
	ds := struct {
		Success  bool
		Message  string
		ImageURL string
	}{
		Success: false,
	}
	var user entity.Users
	erra := request.ParseMultipartForm(9999999999999999)
	if erra != nil {
		ds.Message = " Parse Error Please Select a file with lesser Size "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	switch strings.ToUpper(requestersession.Role) {
	case entity.TEACHER:
		{
			user = adminh.TeacherService.GetTeacherByID(requestersession.ID)
			break
		}
	case entity.STUDENT:
		{
			user = adminh.StudentService.GetStudentByID(requestersession.ID)
			break
		}
	case entity.FIELDMAN:
		{
			user = adminh.TrainerService.GetTrainerByID(requestersession.ID)
			break
		}
	case entity.SUPERADMIN, entity.SECRETART:
		{
			user = adminh.AdminService.GetAdminByID(requestersession.ID, "")
			break
		}
	}
	// user := adminh.StudentService.GetStudentByID(requestersession.ID)
	if user == nil {
		ds.Message = "There Is No User By This ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	// var ImageDirectory string
	image, header, erra := request.FormFile("image")
	if erra != nil || header == nil {
		ds.Message = "No Image Found In the Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	defer image.Close()
	filename := header.Filename
	if !(Helper.IsImage(filename)) {
		ds.Message = "Unsupported File Format"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	randomname := Helper.GenerateRandomString(6, Helper.CHARACTERS) + ".jpg"
	previousPath := user.GetImageURL()
	newFile, erro := os.Create(entity.PathToTemplates + entity.PathToStudentsFromTemplates + randomname)
	if erro != nil {
		ds.Message = "Error While Saving the new file "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	defer newFile.Close()
	user.SetImageURL(entity.PathToStudentsFromTemplates + randomname)
	var success bool
	success = adminh.HelpSaveTheProfilePicture(requestersession.ID, user.GetImageURL(), strings.ToUpper(requestersession.Role))
	if !success {
		ds.Message = "Internale Server Error a"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	_, erer := io.Copy(newFile, image)
	if erer != nil {
		user.SetImageURL(previousPath)
		success = adminh.HelpSaveTheProfilePicture(requestersession.ID, user.GetImageURL(), strings.ToUpper(requestersession.Role))
		ds.Message = " Internal Server Error"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	os.Remove(entity.PathToTemplates + previousPath)
	ds.Success = true
	ds.Message = "succesful"
	ds.ImageURL = user.GetImageURL()
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// HelpSaveTheProfilePicture function returning bool
func (adminh *AdminHandler) HelpSaveTheProfilePicture(ID uint, Imageurl, role string) bool {
	var success bool
	switch role {
	case entity.TEACHER:
		{
			success = adminh.TeacherService.ChangeImageURL(ID, Imageurl)
			break
		}
	case entity.STUDENT:
		{
			success = adminh.StudentService.ChangeImageURL(ID, Imageurl)
			break
		}
	case entity.FIELDMAN:
		{
			success = adminh.TrainerService.ChangeImageURL(ID, Imageurl)
			break
		}
	case entity.SUPERADMIN, entity.SECRETART:
		{
			success = adminh.AdminService.ChangeImageURL(ID, Imageurl)
			break
		}
	}
	return success
}

// StudentsIDCardGenerate method to Generate an ID Card For Studeents
// Method GET
// Variables round_id  , paid_amount , year  , month , days
func (studenthandler *StudentHandler) StudentsIDCardGenerate(response http.ResponseWriter, request *http.Request) {
	session := studenthandler.SessionSeervice.GetSession(request)
	if session == nil {
		return
	}
	fmt.Println("This is rinning")
	ds := struct {
		Success  bool
		Message  string
		Students []entity.Student
		HOST     string
		Year     int
		Month    int
		Day      int
		Lang     string
		Branch   *entity.Branch
		Round    entity.Round
		Date     etc.Date
	}{
		Success: false,
		HOST:    entity.PROTOCOL + entity.HOST,
	}
	var (
		PagesFront []string
		PagesBack  []string
		// Students   []entity.Student
	)
	branch := studenthandler.BranchService.GetBranchByID(session.BranchID)
	ds.Branch = branch
	roundID, er := strconv.Atoi(request.FormValue("round_id"))
	lang := request.FormValue("lang")
	ds.Lang = lang
	if lang == "" {
		ds.Lang = branch.Lang
	}
	paymentAmount, er := strconv.ParseFloat(request.FormValue("payment_amount"), 64)
	if er != nil {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(response, request)
		response.Write([]byte("<h1>   Invalid Request Value  </h1>"))
		return
	}
	round := studenthandler.RoundService.GetRoundByID(uint(roundID))
	if round == nil {
		fmt.Println()
		handler := http.NotFoundHandler()
		handler.ServeHTTP(response, request)
		response.Write([]byte("<h1> Round Not Found! </h1>"))
		return
	}
	ds.Round = *round
	students := studenthandler.StudentService.StudentsOfRoundWithPayment(uint(roundID), paymentAmount)
	if students == nil {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(response, request)
		response.Write([]byte("<h1> Record Not Found </h1>"))
		return
	}
	var PagesCount int
	PagesCount = len(*students) / 4
	if len(*students)%4 > 0 {
		PagesCount++
	}
	if PagesCount == 0 {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(response, request)
		response.Write([]byte("<h1> No Student Found!</h1>"))
		return
	}
	year, er := strconv.Atoi(request.FormValue("year"))
	month, er := strconv.Atoi(request.FormValue("month"))
	day, er := strconv.Atoi(request.FormValue("day"))
	if er == nil {
		ds.Year = year
		ds.Month = month
		ds.Day = day
	}
	var init, end int
	init = 0
	for k := 0; k < PagesCount; k++ {
		var frontFile, backFile *os.File
		var frontFilename, backFilename string
		var selectedStudents []entity.Student
		frontFilename = entity.PathToPdfs + strings.ToLower(Helper.GenerateRandomString(5, Helper.CHARACTERS)) + "a.html"
		backFilename = entity.PathToPdfs + strings.ToLower(Helper.GenerateRandomString(5, Helper.CHARACTERS)) + "a.html"
		frontFile, era := os.Create(frontFilename)
		backFile, era = os.Create(backFilename)
		if era != nil {
			os.Remove(frontFilename)
			os.Remove(backFilename)
			return
		}
		if k < PagesCount-1 {
			end += 4
		} else {
			end = -1
		}
		if end > 0 {
			selectedStudents = (*students)[init:end]
			init += 4
		} else {
			selectedStudents = (*students)[init:]
		}
		ds.Students = selectedStudents
		// fmt.Println(selectedStudents[0].BirthDate)
		studenthandler.TempHandler.Templates.ExecuteTemplate(frontFile, "frontStudentID.html", ds)
		studenthandler.TempHandler.Templates.ExecuteTemplate(backFile, "backStudentID.html", ds)
		PagesFront = append(PagesFront, frontFilename)
		PagesBack = append(PagesBack, backFilename)
		frontFile.Close()
		backFile.Close()
	}
	idsPDFFileDirectory := HtmlToPDF.GetIDSPdf(append(PagesFront, PagesBack...)[0])
	if idsPDFFileDirectory == "" {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(response, request)
		response.Write([]byte("<h1> Internal Server ERROR !</h1>"))
		return
	}
	// Sice the Pdf is Generated I Am Gonna Delete all the Html files that are used to generate teh Pdf files
	for l := 0; l < PagesCount; l++ {
		os.Remove(PagesFront[l])
		os.Remove(PagesBack[l])
	}
	file, eraaa := os.OpenFile(idsPDFFileDirectory, os.O_RDWR, os.ModePerm)
	if eraaa != nil {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(response, request)
		response.Write([]byte("<h1> Internal Server ERROR ! ..</h1>"))
		return
	}
	reedSeeker := io.ReadSeeker(file)
	http.ServeContent(response, request, idsPDFFileDirectory, time.Now(), reedSeeker)
	file.Close()
	os.Remove(idsPDFFileDirectory)
}
