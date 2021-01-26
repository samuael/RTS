// Package handler
package handler

import (
	"html/template"
	"net/http"

	// "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// TemplateHandler  handling template Related Tas
type TemplateHandler struct {
	Templates     *template.Template
	BranchService Branch.BranchService
	SessioHandler *session.Cookiehandler
	AdminService  Admin.AdminService
}

// NewTemplateHandler returning new instance of Template Handler
func NewTemplateHandler(
	temp *template.Template,
	branchSerivice Branch.BranchService,
	SessionHandler *session.Cookiehandler,
	AdminService Admin.AdminService,
) *TemplateHandler {
	return &TemplateHandler{
		Templates:     temp,
		BranchService: branchSerivice,
		AdminService:  AdminService,
		SessioHandler: SessionHandler,
	}
}

// IndexPage serving the firsthome page
func (th TemplateHandler) IndexPage(writer http.ResponseWriter, response *http.Request) {

	branch := th.BranchService.GetDefaultBranch()
	branchs := th.BranchService.GetBranchs()
	indexPageData := entity.MainPageInfo{
		Branch:  *branch,
		Branchs: branchs,
	}
	// fmt.Println(branchs)
	th.Templates.ExecuteTemplate(writer, "sam_index.html", indexPageData)
}

// ChartPage  represents the chart page
func (th TemplateHandler) ChartPage(writer http.ResponseWriter, response *http.Request) {
	th.Templates.ExecuteTemplate(writer, "sam_chart.html", nil)
}

// PageNotFound method
func (th *TemplateHandler) PageNotFound(response http.ResponseWriter, request *http.Request) {
	newSession := th.SessioHandler.GetSession(request)
	admin := th.AdminService.GetAdminByID(newSession.ID, newSession.Username)
	branch := th.BranchService.GetBranchByID(newSession.BranchID)
	path := request.URL.Path

	Four04Info := entity.Four04{
		Admin:      admin,
		Branch:     branch,
		Path:       path,
		Statustext: "Page Not Found",
		StatusCode: 404,
	}

	th.Templates.ExecuteTemplate(response, "four04.html", Four04Info)
	return
}
