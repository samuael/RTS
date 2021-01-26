package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Payment"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	"github.com/Projects/RidingTrainingSystem/pkg/HtmlToPDF"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// PaymentHandler struct
type PaymentHandler struct {
	PaymentService  Payment.PaymentService
	SessionService  *session.Cookiehandler
	AdminService    Admin.AdminService
	StudentService  Student.StudentService
	RoundService    Round.RoundService
	BranchService   Branch.BranchService
	TemplateHandler *TemplateHandler
}

// NewPaymentHandler function
func NewPaymentHandler(
	PaymentService Payment.PaymentService,
	SessionService *session.Cookiehandler,
	AdminService Admin.AdminService,
	StudentService Student.StudentService,
	RoundService Round.RoundService,
	BranchService Branch.BranchService,
	TemplateHandler *TemplateHandler,
) *PaymentHandler {
	return &PaymentHandler{
		PaymentService:  PaymentService,
		SessionService:  SessionService,
		AdminService:    AdminService,
		StudentService:  StudentService,
		RoundService:    RoundService,
		BranchService:   BranchService,
		TemplateHandler: TemplateHandler,
	}
}

// CreatePayment method for Creating new Payment
func (paymenthandler *PaymentHandler) CreatePayment(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success bool
		Message string
		// Student     entity.Student
		PaidAmount  float64
		Adminname   string
		Payment     entity.Payment
		Amount      float64 `json:"amount"`
		Roundid     uint    `json:"round_id"`
		Studentid   uint    `json:"student_id"`
		RepayAmount float64
	}{
		Success: false,
	}
	requesterSession := paymenthandler.SessionService.GetSession(request)
	if requesterSession == nil {
		ds.Message = "UnKnown User "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	newDecoder := json.NewDecoder(request.Body)
	decodeErrro := newDecoder.Decode(&ds)
	if decodeErrro != nil {
		ds.Message = "Error While Reading the Addresses "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if ds.Roundid <= 0 || ds.Studentid <= 0 || ds.Amount <= 0.0 {
		ds.Message = " Invalid Amount "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	admin := paymenthandler.AdminService.GetAdminByID(requesterSession.ID, "")
	student := paymenthandler.StudentService.GetStudentByID(ds.Studentid)
	round := paymenthandler.RoundService.GetRoundByID(ds.Roundid)
	if admin == nil || student == nil || round == nil {
		ds.Message = " Record Not Found "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	payment := &entity.Payment{
		Amount:       ds.Amount,
		BranchRefer:  requesterSession.BranchID,
		RoundRefer:   ds.Roundid,
		StudentRefer: ds.Studentid,
		AdminRefer:   admin.ID,
		Admin:        *admin,
		Student:      *student,
		Date:         etc.NewDate(0),
		Round:        round,
	}

	payment = paymenthandler.PaymentService.CreatePayment(payment)
	round.TotalPaid = round.TotalPaid + ds.Amount

	student.PaidAmount = student.PaidAmount + ds.Amount

	if student.PaidAmount >= round.Cost {
		repay := round.Cost - student.PaidAmount
		ds.RepayAmount = repay
		student.PaidAmount = round.Cost
	}

	student = paymenthandler.StudentService.SaveStudent(student)
	round = paymenthandler.RoundService.SaveRound(round)
	if round == nil || student == nil {
		ds.Message = " Error while Update "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if payment == nil {
		ds.Message = " Error whle Creating the Payment  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	// student.Round = nil
	student.Address = &entity.Address{}
	student.BirthDate = &etc.Date{}
	student.Category = entity.Category{}
	student.Section = &entity.Section{}
	student.GuarantorAddress = &entity.Address{}
	payment.Round = &entity.Round{}
	payment.Admin = entity.Admin{}
	student.PaidAmount += ds.Amount
	ds.Success = true
	ds.Message = " Payment Success fully Created "
	ds.Adminname = admin.Username
	ds.Payment = *payment
	ds.PaidAmount = student.PaidAmount
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// ReciptForPayment   method to generate a recipt for a payment
// Method GET
// PArameter payment_id
// Role Secretary (Admin )
func (paymenthandler *PaymentHandler) ReciptForPayment(response http.ResponseWriter, request *http.Request) {
	requestersessioin := paymenthandler.SessionService.GetSession(request)
	if requestersessioin == nil {
		response.Write([]byte("<h1> Un Authorized User </h1>"))
		return
	}
	paymentidstring := request.FormValue("payment_id")
	paymentID, erro := strconv.Atoi(paymentidstring)
	if erro != nil || paymentID <= 0 {
		// response.Header().Add("Content-Type", "application/json")
	}
	paymentsResponseds := paymenthandler.PaymentService.GetSinglePaymentReciptData(uint(paymentID))
	branch := paymenthandler.BranchService.GetBranchByID(uint(requestersessioin.BranchID))
	if branch == nil {
		response.Write([]byte("<h1> Un-Authorized User </h1>"))
		return
	}
	paymentsResponseds.Branch = *branch
	paymentsResponseds.Host = entity.PROTOCOL + entity.HOST
	randomString := entity.PathToPdfs + Helper.GenerateRandomString(5, Helper.NUMBERS) + ".html"
	pdfNewFile, _ := os.Create(randomString)
	paymenthandler.TemplateHandler.Templates.ExecuteTemplate(pdfNewFile, "singlePaymentRecipt.html", paymentsResponseds)
	pdfNewFile.Close()

	theNewPDFFileReciptDirectory := HtmlToPDF.GetReciptThePdf(randomString)
	if theNewPDFFileReciptDirectory == "" {

	}
	thereciptpdffile, _ := os.Open(theNewPDFFileReciptDirectory)
	info, _ := thereciptpdffile.Stat()
	response.Header().Set("Content-Type", "application/pdf")
	response.Header().Set("Content-Deposition", "attachment ; filename="+info.Name())
	response.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(response, thereciptpdffile)
	thereciptpdffile.Close()
	os.Remove(theNewPDFFileReciptDirectory)
}

// PaymentsOfRound (Roundid, BranchNumber uint) *[]entity.Payment
func (paymenthandler *PaymentHandler) PaymentsOfRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success   bool
		Message   string
		Roundid   uint
		Paymens   []entity.Payment
		TotalPaid float64
		Cost      float64
	}{
		Success: false,
	}
	requestersessioin := paymenthandler.SessionService.GetSession(request)

	roundidstring := request.FormValue("round_id")
	roundid, era := strconv.Atoi(roundidstring)
	if era != nil {
		ds.Message = "Invalid Request "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	round := paymenthandler.RoundService.GetRoundByID(uint(roundid))
	if round == nil {
		ds.Message = " Round Doesnt Exit "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if requestersessioin.BranchID != round.Branchnumber {
		ds.Message = " Not Authorized "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	payments := paymenthandler.PaymentService.PaymentsOfRound(round.ID, requestersessioin.BranchID)
	if payments == nil {
		ds.Message = "Sorry ! Not Payments of this Round  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	sum := 0.0
	for _, Pay := range *payments {
		sum += Pay.Amount
	}

	ds.Success = true
	ds.Paymens = *payments
	// ds.TotalPaid = paymenthandler.PaymentService.RoundPaidAmount(uint(roundid))
	ds.Message = "Succesfull"
	ds.TotalPaid = sum
	ds.Roundid = round.ID
	ds.Cost = round.Cost
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// PaymentsOfStudent (StudentID uint) *[]entity.Payment
func (paymenthandler *PaymentHandler) PaymentsOfStudent(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success      bool
		Message      string
		Student      entity.Student
		Payments     []entity.Payment
		Studentid    uint
		TotalPays    float64
		RemainingPay float64
	}{
		Success: true,
	}
	studenidstring := request.FormValue("student_id")
	studentid, srra := strconv.Atoi(studenidstring)
	if srra != nil || studentid <= 0 {
		ds.Message = " Invalid Input  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	payments := paymenthandler.PaymentService.PaymentsOfStudent(uint(studentid))
	if payments == nil {
		ds.Message = " Internal Server Errror "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	student := paymenthandler.StudentService.GetStudentByID(uint(studentid))
	if student == nil {
		ds.Message = "Unknown User"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Payments Of Student Succesfull Fetched "
	ds.Success = true
	ds.Payments = *payments
	ds.TotalPays = student.PaidAmount
	ds.RemainingPay = student.Round.Cost - student.PaidAmount
	ds.Student = *student
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// PaymentsOfSecretary (SecretaryID, limit uint) *[]entity.Payment
func (paymenthandler *PaymentHandler) PaymentsOfSecretary(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success   bool
		Message   string
		Payments  []entity.Payment
		Secretary entity.Admin
		// LastMonthTransaction float64
		RoundToCost []interface{}
		SecretaryID uint
	}{
		Success: false,
	}
	requesterSession := paymenthandler.SessionService.GetSession(request)
	secretaryidstring := request.FormValue("secretary_id")
	limitstr := request.FormValue("limit")
	limit, eraa := strconv.Atoi(limitstr)
	offsetstr := request.FormValue("offset")
	offset, eraa := strconv.Atoi(offsetstr)
	secretaryid, eraa := strconv.Atoi(secretaryidstring)
	if eraa != nil || secretaryid <= 0 {
		ds.Message = "Invalid Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.SecretaryID = uint(secretaryid)
	admin := paymenthandler.AdminService.GetAdminByID(uint(secretaryid), requesterSession.Username)
	admin.Password = ""
	payments := paymenthandler.PaymentService.PaymentsOfSecretary(uint(secretaryid), uint(limit), uint(offset))
	if admin == nil || payments == nil {
		ds.Message = "Record Not Found "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if admin.Role != entity.SECRETART {
		ds.Message = " Invalid Role "
		ds.Secretary = *admin
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Payments = *payments
	ds.Secretary = *admin
	ds.Success = true
	ds.Message = "Succesfull"
	mapo := map[uint]float64{}
	realMapo := []interface{}{}
	for _, payment := range *payments {
		mapo[payment.RoundRefer] += payment.Amount
	}
	for key, value := range mapo {
		val := struct {
			Round  entity.Round
			Amount float64
		}{
			Round:  *paymenthandler.RoundService.GetRoundByID(key),
			Amount: value,
		}
		realMapo = append(realMapo, val)
	}
	ds.RoundToCost = realMapo
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DailyPaymentReport method
// Method Get
// Return Page
func (paymenthandler *PaymentHandler) DailyPaymentReport(response http.ResponseWriter, request *http.Request) {
	// rsession := paymenthandler.SessionService.GetSession(request)
	// if rsession == nil {
	// 	response.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	daystring := request.FormValue("day")
	monthstring := request.FormValue("month")
	yearstring := request.FormValue("year")
	day, era := strconv.Atoi(daystring)
	month, era := strconv.Atoi(monthstring)
	year, era := strconv.Atoi(yearstring)

	branchID, era := strconv.Atoi(request.FormValue("branch_id"))

	now := etc.NewDate(0)
	ds := struct {
		Success bool
		Message string
		Branch  entity.Branch
	}{
		Success: false,
	}
	if era != nil || day <= 0 || day > 30 || month <= 0 || month >= 13 || year > now.Year {
		ds.Message = "Invalid Date Input.."
		paymenthandler.TemplateHandler.Templates.ExecuteTemplate(response, "CantGetSchedule.html", ds)
	}
	fmt.Printf(" %d /%d /%d %d \n\n", day, month, year, branchID)
	paymenthandler.PaymentService.DPR(etc.Date{Day: day, Month: month, Year: year}, branchID)
	// dailyPaymentReport := paymenthandler.PaymentService.GetDailyPaymentReport(etc.Date{Day: day, Month: month, Year: year}, rsession.BranchID)
	// if dailyPaymentReport == nil {

	// }

}

// GetMonthlyPaymentReport method Get
// Methdo GET
// AUTHORIZATION SUPERADMIN
// func (paymenthandler *PaymentHandler) GetMonthlyPaymentReport(response http.ResponseWriter, request *http.Request) {
// 	reqses := paymenthandler.SessionService.GetSession(request)
// 	if reqses == nil {
// 		return
// 	}
// 	ds := struct {
// 		Success              bool
// 		Message              string
// 		MonthlyPaymentResult entity.MonthlyPaymentResult
// 		Branch               entity.Branch
// 		Admin                entity.Admin
// 	}{
// 		Success: false,
// 	}
// 	year, era := strconv.Atoi(request.FormValue("year"))
// 	month, era := strconv.Atoi(request.FormValue("month"))
// 	if era != nil {
// 		return
// 	}

// 	admin := paymenthandler.AdminService.GetAdminByID(reqses.ID, "")
// 	branch := paymenthandler.BranchService.GetBranchByID(uint(reqses.BranchID))
// 	ds.Admin = *admin
// 	ds.Branch = *branch
// 	totalInComePaymentsOfEachSecretories := paymenthandler.PaymentService.PaymentsOfEachSecretariesForMonth(uint(branchID), year, month)
// 	roundToPaymentsOfThebranch := paymenthandler.PaymentService.PaymentsOfRoundForTheMonth(uint(branchID), year, month)
// 	categoryPaymentsOfTheMonth := paymenthandler.PaymentService.PaymentsPfCategoryForTheMonth(uint(branchID), year, month)
// 	totalPaidAmount := paymenthandler.PaymentService.TotalPaymentofMonth(uint(branchID), year, month)
// 	ds.MonthlyPaymentResult.TotalPaidAmount = totalPaidAmount
// }
