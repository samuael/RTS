// Package handler and this is Responsible to handle HTMLToPdfHandler
package handler

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Projects/RidingTrainingSystem/pkg/HtmlToPDF"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"

	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Section"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// HTMLToPdfHandler struct
type HTMLToPdfHandler struct {
	SectionService  Section.SectionService
	TemplateHandler *TemplateHandler
	SessionHandler  *session.Cookiehandler
	BranchService   Branch.BranchService
	RoundService    Round.RoundService
}

// NewHTMLToPdfHandler function
func NewHTMLToPdfHandler(
	sectionservice Section.SectionService,
	th *TemplateHandler,
	SessionHandler *session.Cookiehandler,
	BranchService Branch.BranchService,
	RoundService Round.RoundService,
) *HTMLToPdfHandler {
	return &HTMLToPdfHandler{
		SectionService:  sectionservice,
		TemplateHandler: th,
		SessionHandler:  SessionHandler,
		BranchService:   BranchService,
		RoundService:    RoundService,
	}
}

// GenerateSchedulePDF for Generating Schedule For Specific Students
// Method GET
// Accessible For All Loged In Users
func (htmlpdfh *HTMLToPdfHandler) GenerateSchedulePDF(response http.ResponseWriter, request *http.Request) {
	requesterSession := htmlpdfh.SessionHandler.GetSession(request)
	if requesterSession == nil {
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "four04.html", entity.Four04{
			Branch:     htmlpdfh.BranchService.GetDefaultBranch(),
			StatusCode: http.StatusProxyAuthRequired,
			Statustext: http.StatusText(http.StatusProxyAuthRequired),
			Path:       request.Header.Get("path"),
			Admin:      &entity.Admin{},
		})

	}
	fmt.Println("One Request /..... ")
	roundNumberString := request.FormValue("round_id")
	roundID, eerr := strconv.Atoi(roundNumberString)
	branch := htmlpdfh.BranchService.GetBranchByID(requesterSession.BranchID)
	round := htmlpdfh.RoundService.GetRoundByID(uint(roundID))
	if branch == nil {
		branch = &entity.Branch{}
	}
	ds := entity.ScheduleDataStructure{
		Success: false,
		Branch:  *branch,
		Host:    entity.PROTOCOL + entity.HOST,
	}
	if eerr != nil {
		// randomString := entity.PathToPdfs + Helper.GenerateRandomString(5, Helper.NUMBERS) + ".html"
		// pdfNewFile, er := os.Create(randomString)
		// if er != nil {
		// 	return
		// }
		ds.Message = " Invalid Round Number Value "
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "CantGetSchedule.html", ds)
		return
	}
	if round == nil {
		ds.Message = " The Round Doesn't Exist ...."
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "CantGetSchedule.html", ds)
		return
	}
	sectionso := (htmlpdfh.SectionService.GetSectionsOfRound(uint(roundID)))
	var sections []entity.Section
	if sectionso != nil {
		sections = *sectionso
	}
	if sectionso == nil {
		ds.Message = " The Round Has No Sections"
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "CantGetSchedule.html", ds)
		return
	}
	// Generation Of Schedule for Round Begins
	randomString := entity.PathToPdfs + Helper.GenerateRandomString(5, Helper.NUMBERS) + ".html"
	pdfNewFile, er := os.Create(randomString)
	if er != nil {
		ds.Message = " Internal  Server Errro Please Try Again "
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "CantGetSchedule.html", ds)
		return
	}
	round.Sections = sections
	ds.Branch = *branch
	ds.Round = *round
	ds.Sections = sections
	// fmt.Println(sections[0].Trainings[0].Students)
	htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(pdfNewFile, "schedule.html", ds)
	pdfNewFile.Close()
	thePdfFile := HtmlToPDF.GetThePdf(randomString)
	os.Remove(randomString)
	if thePdfFile == "" {
		ds.Message = " Error Generating Pdf File "
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "CantGetSchedule.html", ds)
		return
	}
	pdfScheduleFile, eras := os.Open(thePdfFile)
	if eras != nil {
		ds.Message = " Error Opening Schedule"
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "CantGetSchedule.html", ds)
		return
	}
	info, _ := pdfScheduleFile.Stat()

	response.Header().Set("Content-Type", "application/pdf")
	response.Header().Set("Content-Deposition", "attachment ; filename="+info.Name())
	response.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(response, pdfScheduleFile)
	pdfScheduleFile.Close()
	os.Remove(thePdfFile)
}

// RoundRegisteredStudentsInfoPDF method for generating a pdf for students data to be produced which means the studetns
// Username and Password
// Authorization SECRETARY and SUPERADMIN
// request value round_id
func (htmlpdfh *HTMLToPdfHandler) RoundRegisteredStudentsInfoPDF(response http.ResponseWriter, request *http.Request) {
	requestersession := htmlpdfh.SessionHandler.GetSession(request)
	if requestersession == nil {
		return
	} else if requestersession.BranchID <= 0 {
		return
	}
	ds := struct {
		Success  bool
		Message  string
		Host     string
		Branch   entity.Branch
		Round    entity.Round
		Students []entity.Student
	}{
		Success: false,
		Host:    entity.PROTOCOL + entity.HOST,
	}
	branch := htmlpdfh.BranchService.GetBranchByID(requestersession.BranchID)
	if branch == nil {
		return
	}
	ds.Branch = *branch
	roundID, era := strconv.Atoi(request.FormValue("round_id"))
	if era != nil {
		ds.Message = "There is No round By this ID "
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "roundStudentsInfo.html", ds)
		return
	}
	round := htmlpdfh.RoundService.GetRoundByID(uint(roundID))
	if round == nil {
		ds.Message = "Record Not Found "
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "roundStudentsInfo.html", ds)
		return
	}
	ds.Success = true
	ds.Round = *round
	ds.Branch = *branch
	ds.Students = round.Students
	randomString := entity.PathToPdfs + Helper.GenerateRandomString(5, Helper.NUMBERS) + ".html"
	theHTMLToBeUsed, er := os.Create(randomString)
	if er != nil {
		ds.Message = " Internal  Server Errro Please Try Again "
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "roundStudentsInfo.html", ds)
		return
	}
	htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(theHTMLToBeUsed, "roundStudentsInfo.html", ds)
	theHTMLToBeUsed.Close()
	thePdfFileName := HtmlToPDF.GetThePdf(randomString)
	os.Remove(randomString)
	if thePdfFileName == "" {
		ds.Message = " Error Generating Pdf File "
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "roundStudentsInfo.html", ds)
		return
	}
	pdfReturnableStudentsListFile, eras := os.Open(thePdfFileName)
	if eras != nil {
		ds.Message = " Error Opening Schedule"
		ds.Success = false
		htmlpdfh.TemplateHandler.Templates.ExecuteTemplate(response, "roundStudentsInfo.html", ds)
		return
	}
	info, _ := pdfReturnableStudentsListFile.Stat()
	response.Header().Set("Content-Type", "application/pdf")
	response.Header().Set("Content-Deposition", "attachment ; filename="+info.Name())
	response.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	io.Copy(response, pdfReturnableStudentsListFile)
	pdfReturnableStudentsListFile.Close()
	os.Remove(thePdfFileName)
}
