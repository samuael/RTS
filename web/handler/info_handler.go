package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Information"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// InfoHandler struct
type InfoHandler struct {
	InfoService    Information.InformationService
	SessionService *session.Cookiehandler
	BranchService  Branch.BranchService
}

// NewInfoHandler function
func NewInfoHandler(
	infoservice Information.InformationService,
	session *session.Cookiehandler,
	BranchService Branch.BranchService,
) *InfoHandler {
	return &InfoHandler{
		InfoService:    infoservice,
		SessionService: session,
		BranchService:  BranchService,
	}
}

// CreateInformation handler Function Using Form Post Input and Json Response
func (infohandler *InfoHandler) CreateInformation(response http.ResponseWriter, request *http.Request) {
	// Authorized only for Admins And Secretories
	response.Header().Add("Content-Type", "application/json")
	requesterSession := infohandler.SessionService.GetSession(request)
	ds := struct {
		Success bool
		Message string
		Info    entity.Information
	}{
		Success: false,
	}
	request.ParseForm()
	csrf := request.FormValue("_csrf")

	valid := infohandler.SessionService.ValidateForm(csrf)
	if !valid {
		ds.Message = " UnKnow User or Form  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	title := request.FormValue("title")
	description := request.FormValue("description")
	if title == "" || description == "" {
		ds.Message = " Invalid Form "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	branch := infohandler.BranchService.GetBranchByID(requesterSession.BranchID)
	info := &entity.Information{
		Title:       title,
		Description: description,
		Username:    requesterSession.Username,
		BranchID:    requesterSession.BranchID,
		BranchName:  branch.Name,
		Active:      true,
	}
	info = infohandler.InfoService.CreateInfo(info)
	if info == nil {
		ds.Message = " Internal Server Error "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Succesfully Created an Information"
	ds.Info = *info
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// UpdateInfo  method to update an Info Using the POstForm
func (infohandler *InfoHandler) UpdateInfo(response http.ResponseWriter, request *http.Request) {
	// This method is only visible for the person who wrote this Information or superadmn
	response.Header().Add("Content-Type", "application/json")
	requesterSession := infohandler.SessionService.GetSession(request)
	ds := struct {
		Success bool
		Message string
		Info    entity.Information
	}{
		Success: false,
	}

	// Taking the informationof the Update Using form request
	request.ParseForm()
	csrf := request.FormValue("_csrf")
	title := request.FormValue("title")
	description := request.FormValue("description")
	infoID := request.FormValue("id")

	infoid, erra := strconv.Atoi(infoID)
	if erra != nil || infoid <= 0 {
		ds.Message = "Invalid Info ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if title == "" || description == "" || csrf == "" {
		ds.Message = "Invalid Info Values "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if !infohandler.SessionService.ValidateForm(csrf) {
		ds.Message = " Invalid Form "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	information := entity.Information{}
	information.ID = uint(infoid)
	info := infohandler.InfoService.GetInfoByID(uint(infoid))
	if info == nil {
		ds.Message = " No Information By This Id "
		ds.Info = information
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if info.Username == requesterSession.Username || requesterSession.Role == entity.SUPERADMIN {
		info.Title = title
		info.Description = description
		info = infohandler.InfoService.SaveInfo(info)
		if info == nil {
			ds.Message = " Error While Updating the Info "
			ds.Info = information
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		ds.Info = *info
		ds.Message = "Succesfully Updated "
		ds.Success = true
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = " You Are Not Allowed To Make Change in this Information "
	ds.Info = information
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DeleteInformation method to delete Specific Information The Method is get
func (infohandler *InfoHandler) DeleteInformation(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success bool
		Message string
		InfoID  uint
	}{
		Success: false,
	}
	infoIdstring := request.FormValue("info_id")
	infoID, erra := strconv.Atoi(infoIdstring)
	if erra != nil {
		ds.Message = " Invalid Info ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	info := &entity.Information{}
	info.ID = uint(infoID)
	Success := infohandler.InfoService.DeleteInfo(info)
	if !Success {
		ds.Message = " Can't Delete the Information! "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfully Deleted The Info With ID " + infoIdstring
	ds.Success = false
	ds.InfoID = uint(infoID)
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetActiveInformations (BranchID uint) *[]entity.Information  the method is Get
// and Permitted to all Users this thing there doesnt ask both authorization or authentication
func (infohandler *InfoHandler) GetActiveInformations(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// Returning slice of informations
	ds := struct {
		Success      bool
		Message      string
		Informations []entity.Information
	}{
		Success: false,
	}
	requesterSession := infohandler.SessionService.GetSession(request)
	informations := infohandler.InfoService.GetActiveInfos(uint(requesterSession.BranchID))
	if informations == nil {
		ds.Message = "No Informations By This Branch "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Informations = *informations
	ds.Success = true
	ds.Message = "Succesfull "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetAllActiveInformations (BranchID uint) *[]entity.Information  the method is Get
// and Permitted to all Users this thing there doesnt ask both authorization or authentication
func (infohandler *InfoHandler) GetAllActiveInformations(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// Returning slice of informations
	ds := struct {
		Success      bool
		Message      string
		Informations []entity.Information
	}{
		Success: false,
	}
	informations := infohandler.InfoService.GetAllActiveInfos()
	if informations == nil {
		ds.Message = "No Informations By This Branch "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Informations = *informations
	ds.Success = true
	ds.Message = "Succesfull "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetAllInfos (BranchID uint) *[]entity.Information  the method is Get  returning a list of
func (infohandler *InfoHandler) GetAllInfos(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requesterSession := infohandler.SessionService.GetSession(request)
	ds := struct {
		Success      bool
		Message      string
		Informations []entity.Information
		Length       uint
	}{
		Success: false,
	}
	informations := infohandler.InfoService.GetAllInfos(requesterSession.BranchID)
	if informations == nil {
		ds.Message = "Informations Are Fetched Succesfully "
		ds.Length = 0
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Informations = *informations
	ds.Message = " Succesfully Fetched The Infos "
	ds.Success = true
	ds.Length = uint(len(ds.Informations))
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetInfoByID method GET variable info_id returning json
func (infohandler *InfoHandler) GetInfoByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success bool
		Message string
		Info    entity.Information
		InfoID  uint
	}{
		Success: false,
	}
	infoid := request.FormValue("info_id")
	InfoID, erra := strconv.Atoi(infoid)
	if erra != nil {
		ds.Message = "Invalid Get Request"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	info := infohandler.InfoService.GetInfoByID(uint(InfoID))
	if info == nil {
		ds.Message = "No Information Record Is Found"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Info = *info
	ds.Message = "Fetch Succesfull"
	ds.Success = true
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// ActivateInformation method rturning Success Datastructure METHOD Get   response JSON  METHOD POST
func (infohandler *InfoHandler) ActivateInformation(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// rsession := infohandler.SessionService.GetSession(request)
	ds := struct {
		Success bool
		Message string
		InfoID  uint
	}{
		Success: false,
	}
	request.ParseForm()
	infoIDString := request.FormValue("info_id")
	InfoID, erra := strconv.Atoi(infoIDString)
	if erra != nil {
		ds.Message = "Invalid Info ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	info := infohandler.InfoService.ActivateInformation(uint(InfoID))
	if info == nil {
		ds.Message = "Error While Activating the Info"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.InfoID = info.ID
	ds.Success = true
	ds.Message = "Succcesfull Activated "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DeactivateInformation method rturning Success Datastructure METHOD Get   response JSON  METHOD POST
func (infohandler *InfoHandler) DeactivateInformation(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// rsession := infohandler.SessionService.GetSession(request)
	ds := struct {
		Success bool
		Message string
		InfoID  uint
	}{
		Success: false,
	}
	request.ParseForm()
	infoIDString := request.FormValue("info_id")
	InfoID, erra := strconv.Atoi(infoIDString)
	if erra != nil {
		ds.Message = "Invalid Info ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	info := infohandler.InfoService.DeactivateInformation(uint(InfoID))
	if info == nil {
		ds.Message = "Error While Deactivating the Info"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.InfoID = info.ID
	ds.Success = true
	ds.Message = "Succcesfull Deactivated "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}
