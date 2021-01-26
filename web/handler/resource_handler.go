// Package handler in this package we will be handling the requests all almost 
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Resource"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// ResourceHandler struct
type ResourceHandler struct {
	ResourceService Resource.ResourceService
	SessionService  *session.Cookiehandler
	BranchService   Branch.BranchService
	AdminService    Admin.AdminService
	StudentService  Student.StudentService
	TeacherService  Teacher.TeacherService
	TrainerService  Trainer.TrainerService
	Templatehandler *TemplateHandler
}

// NewResourceHandler function
func NewResourceHandler(
	resourceser Resource.ResourceService,
	session *session.Cookiehandler,
	BranchService Branch.BranchService,
	AdminService Admin.AdminService,
	StudentService Student.StudentService,
	TeacherService Teacher.TeacherService,
	TrainerService Trainer.TrainerService,
	Templatehandler *TemplateHandler,
) *ResourceHandler {
	return &ResourceHandler{
		ResourceService: resourceser,
		SessionService:  session,
		BranchService:   BranchService,
		AdminService:    AdminService,
		StudentService:  StudentService,
		TeacherService:  TeacherService,
		TrainerService:  TrainerService,
		Templatehandler: Templatehandler,
	}
}

// UploadResource Rsource methdo
// Method POST Using Form Value
// Using  Form Submission
// Authority for Teachers  , TRAINERS  , SUPERADMIN  , SECRETARY
func (reshandler *ResourceHandler) UploadResource(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	request.ParseMultipartForm(999999999999999999)
	requestersession := reshandler.SessionService.GetSession(request)
	description := request.FormValue("description")
	file, fileHeader, newError := request.FormFile("file")
	title := request.FormValue("title")
	ds := struct {
		Success  bool
		Message  string
		Resource entity.Resource
	}{
		Success: false,
	}
	if newError != nil || fileHeader == nil {
		ds.Message = "Invalid Upload REquest"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	if title == "" {
		ds.Message = "The Resource Has to have a title "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	filename := ""
	filename = entity.RESOURCE_PATH + Helper.GenerateRandomString(10, Helper.CHARACTERS) + fileHeader.Filename
	newFile, erra := os.Create("../../web/templates/" + filename)
	if erra != nil {
		ds.Message = "Internal Server Error "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resource := &entity.Resource{
		Description:   description,
		Path:          filename,
		Uploadedby:    requestersession.Username,
		UploaderImage: requestersession.Imageurl,
		UploadDate:    *etc.NewDate(0),
		Title:         title,
	}
	// Setting Type for the Resource
	splitted := strings.Split(filename, ".")
	extension := strings.ToLower(splitted[len(splitted)-1])
	if resource.Type == 0 {
		for _, videoExtension := range entity.VIDEOS {
			if strings.Trim(extension, "") == videoExtension {
				resource.Type = entity.VIDEO
			}
		}
	}
	if resource.Type == 0 {
		for _, videoExtension := range entity.AUDIOS {
			if extension == videoExtension {
				resource.Type = entity.AUDIO
			}
		}
	}
	if resource.Type == 0 {
		for _, videoExtension := range entity.PICTURES {
			if extension == videoExtension {
				resource.Type = entity.IMAGES
			}
		}
	}
	if resource.Type == 0 && extension == "pdf" {
		resource.Type = entity.PDF
	}
	if resource.Type == 0 {
		resource.Type = entity.FILES
	}
	_, errors := io.Copy(newFile, file)
	if errors != nil {
		ds.Message = " Internal server Error while saving the file "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	// Closing the Opened File Before Getting into HLS File Generating
	file.Close()
	newFile.Close()
	//  if the file is audio or video i am gonna incode it to hls format
	// passing the file directory and the the file type
	if resource.Type == entity.AUDIO || resource.Type == entity.VIDEO {
		// Resource m3u8 directory Name Will Be Generated and Name
		// having that name Will Be Created
		// Then the file will be used and Populated in that page
		ResourceHLSDirectory := strconv.Itoa(int(etc.NewDate(0).Unix))
		m3u8Directory, message, errorOfHLS := Helper.EncodeToHLS((entity.TemplateDirectoryFromMain + filename), entity.TemplateDirectoryFromMain+entity.MediaDirectoryFromTemplates, ResourceHLSDirectory)
		if errorOfHLS != nil {
			log.Println(message)
		}
		resource.HLSDirectory = m3u8Directory
		if resource.Type == entity.VIDEO {
			input := entity.PathToTemplates + resource.Path
			OutPut := Helper.GetFirstFrameOfVideo(input)
			if OutPut != "" {
				resource.SnapShootImage = OutPut
			}
		}
	}

	resource = reshandler.ResourceService.CreateResource(resource)
	if resource == nil {
		ds.Message = "Error While Saving The Resource "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = " File Successfully Uploaded "
	ds.Success = true
	ds.Resource = *resource
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DeleteResource Get Method
// Authorized for the Uploaders
func (reshandler *ResourceHandler) DeleteResource(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	resourceidstring := request.FormValue("resource_id")
	resourceID, erra := strconv.Atoi(resourceidstring)
	requesterSession := reshandler.SessionService.GetSession(request)
	ds := struct {
		Success    bool
		Message    string
		ResourceID uint
	}{
		Success: false,
	}
	if erra != nil {
		ds.Message = "Invalid Resource ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resource := reshandler.ResourceService.GetResourceByID(uint(resourceID))
	if resource == nil || !((resource.Uploadedby == requesterSession.Username) || (requesterSession.Role == entity.SUPERADMIN)) {
		if resource == nil {
			ds.Message = " Record Not Found "
		} else {
			ds.Message = "You Are Not Authorized "
		}
		jsons, _ := json.Marshal(ds)
		response.Write(jsons)
		return
	}
	os.Remove(entity.PathToTemplates + resource.Path)
	if resource.HLSDirectory != "" {
		// Write A Function to FInd the Hls Folder Number
		val := Helper.GetHLSFolderNumberName(resource.HLSDirectory)
		os.RemoveAll(entity.PathToResources + "media/" + val)
	}

	if resource.FirstFrame != "" {
		os.RemoveAll(entity.PathToTemplates + resource.FirstFrame)
	}
	if resource.SnapShootImage != "" {
		fmt.Println(entity.PathToFirstFrames + resource.SnapShootImage)
		os.Remove(entity.PathToTemplates + resource.SnapShootImage)
	}
	success := reshandler.ResourceService.DeleteResource(uint(resourceID))
	if success {
		ds.Message = "SuccesFully deleted"
		ds.Success = true
		jsons, _ := json.Marshal(ds)
		response.Write(jsons)
		return
	}
	ds.Message = "Can't Delete The Resource!"
	ds.Success = true
	ds.ResourceID = uint(resourceID)
	jsons, _ := json.Marshal(ds)
	response.Write(jsons)
}

// GetResourceByID method
// Method GET
func (reshandler *ResourceHandler) GetResourceByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	resourceidstring := request.FormValue("resource_id")
	ds := struct {
		Success  bool
		Message  string
		Resource entity.Resource
	}{
		Success: false,
	}
	resourceID, erra := strconv.Atoi(strings.Trim(resourceidstring, " "))
	if erra != nil {
		ds.Message = "Invalid Resource Pointer "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	resource := reshandler.ResourceService.GetResourceByID(uint(resourceID))
	if resource == nil {
		ds.Message = "Resource Not Found "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = " Succesfull "
	ds.Success = true
	ds.Resource = *resource
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetResources (Offset uint, Limit uint) *[]entity.Resource
func (reshandler *ResourceHandler) GetResources(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requesterSession := reshandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		Resources []entity.Resource
	}{
		Success: false,
	}
	limitstring := strings.Trim(request.FormValue("limit"), " ")
	offsetstring := strings.Trim(request.FormValue("offset"), " ")
	offset, erra := strconv.Atoi(offsetstring)
	limit, erra := strconv.Atoi(limitstring)
	if erra != nil {
		ds.Message = "Invalid Request Niggoye"
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resources := reshandler.ResourceService.GetResources(uint(offset), uint(limit))
	if resources == nil {
		ds.Message = "Sorry Cant Fetch any Resource"
		ds.Success = false
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfull "
	ds.Success = true
	ds.Resources = *resources
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetResourceAudios method
// Method GET
// Authorized For all System Users
func (reshandler *ResourceHandler) GetResourceAudios(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requesterSession := reshandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		Resources []entity.Resource
	}{
		Success: false,
	}
	limitstring := strings.Trim(request.FormValue("limit"), " ")
	offsetstring := strings.Trim(request.FormValue("offset"), " ")
	offset, erra := strconv.Atoi(offsetstring)
	limit, erra := strconv.Atoi(limitstring)
	if erra != nil {
		ds.Message = "Invalid Request Niggoye"
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resources := reshandler.ResourceService.GetAudios(uint(offset), uint(limit))
	if resources == nil {
		ds.Message = "Sorry Cant Fetch any Audios"
		ds.Success = false
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfull "
	ds.Success = true
	ds.Resources = *resources
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetResourceVideos method
// Method GET
// Authorized For all System Users
func (reshandler *ResourceHandler) GetResourceVideos(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requesterSession := reshandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		Resources []entity.Resource
	}{
		Success: false,
	}
	limitstring := strings.Trim(request.FormValue("limit"), " ")
	offsetstring := strings.Trim(request.FormValue("offset"), " ")
	offset, erra := strconv.Atoi(offsetstring)
	limit, erra := strconv.Atoi(limitstring)
	if erra != nil {
		ds.Message = "Invalid Request Niggoye"
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resources := reshandler.ResourceService.GetVideos(uint(offset), uint(limit))
	if resources == nil {
		ds.Message = "Sorry Cant Fetch any Videos"
		ds.Success = false
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfull "
	ds.Success = true
	ds.Resources = *resources
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetResourceFiles method
// Method GET
// Authorized For all System Users
func (reshandler *ResourceHandler) GetResourceFiles(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requesterSession := reshandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		Resources []entity.Resource
	}{
		Success: false,
	}
	limitstring := strings.Trim(request.FormValue("limit"), " ")
	offsetstring := strings.Trim(request.FormValue("offset"), " ")
	offset, erra := strconv.Atoi(offsetstring)
	limit, erra := strconv.Atoi(limitstring)
	if erra != nil {
		ds.Message = "Invalid Request Niggoye"
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resources := reshandler.ResourceService.GetFiles(uint(offset), uint(limit))
	if resources == nil {
		ds.Message = "Sorry Cant Fetch any Files"
		ds.Success = false
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfull "
	ds.Success = true
	ds.Resources = *resources
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetResourcePicture method
// Method GET
// Authorized For all System Users
func (reshandler *ResourceHandler) GetResourcePicture(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requesterSession := reshandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		Resources []entity.Resource
	}{
		Success: false,
	}
	limitstring := strings.Trim(request.FormValue("limit"), " ")
	offsetstring := strings.Trim(request.FormValue("offset"), " ")
	offset, erra := strconv.Atoi(offsetstring)
	limit, erra := strconv.Atoi(limitstring)
	if erra != nil {
		ds.Message = "Invalid Request Niggoye"
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resources := reshandler.ResourceService.GetFiles(uint(offset), uint(limit))
	if resources == nil {
		ds.Message = "Sorry Cant Fetch any Files"
		ds.Success = false
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfull "
	ds.Success = true
	ds.Resources = *resources
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetResourcePDFS method
// Method GET
// Authorized For all System Users
func (reshandler *ResourceHandler) GetResourcePDFS(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// requesterSession := reshandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		Resources []entity.Resource
	}{
		Success: false,
	}
	limitstring := strings.Trim(request.FormValue("limit"), " ")
	offsetstring := strings.Trim(request.FormValue("offset"), " ")
	offset, erra := strconv.Atoi(offsetstring)
	limit, erra := strconv.Atoi(limitstring)
	if erra != nil {
		ds.Message = "Invalid Request"
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	resources := reshandler.ResourceService.GetFiles(uint(offset), uint(limit))
	if resources == nil {
		ds.Message = "Sorry Cant Fetch any PDFS"
		ds.Success = false
		ds.Resources = []entity.Resource{}
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = "Succesfull "
	ds.Success = true
	ds.Resources = *resources
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// SearchResource method returning resources matching the search text title
// Method Get
// AUTHORIZATION all Logged In Users
// variables q and t
func (reshandler *ResourceHandler) SearchResource(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	querytext := request.FormValue("q")
	if querytext == "" {
		jsonReturn, _ := json.Marshal([]entity.Resource{})
		response.Write(jsonReturn)
		return
	}
	val, era := strconv.Atoi(request.FormValue("t"))
	if era != nil || val < 0 || val > 5 {
		val = 0
	}
	resource := reshandler.ResourceService.SearchResource(querytext, uint(val))
	jsonReturn, _ := json.Marshal(resource)
	response.Write(jsonReturn)
	return
}

// TemplateLearningPage method
// Method GET
// AUTHORIZATION ALL LOGGED IN USERS
func (reshandler *ResourceHandler) TemplateLearningPage(response http.ResponseWriter, request *http.Request) {
	reqses := reshandler.SessionService.GetSession(request)
	if reqses == nil {
		return
	}
	branchID := reqses.BranchID
	branch := reshandler.BranchService.GetBranchByID(branchID)
	if branch == nil {
		response.Write([]byte("<h1>  UnAuthorized User   </h1>"))
		return
	}
	var user interface{}
	fmt.Println("The Id Of the User ", reqses.ID)
	switch strings.ToUpper(reqses.Role) {
	case entity.SUPERADMIN, entity.SECRETART:
		{
			something := reshandler.AdminService.GetAdminByID(reqses.ID, "")
			user = something
			break
		}
	case entity.TEACHER:
		{
			user = reshandler.TeacherService.GetTeacherByID(reqses.ID)
			break
		}
	case entity.STUDENT:
		{
			user = reshandler.StudentService.GetStudentByID(reqses.ID)
			break
		}
	case entity.FIELDMAN:
		{
			user = reshandler.TrainerService.GetTrainerByID(reqses.ID)
			break
		}
	}
//	 Lang := reshandler.SessionService.GetLang(request)

	 	Lang := branch.Lang
	learningPageDS := entity.LearningPageDS{
		HOST:     entity.PROTOCOL + entity.HOST,
		RouteMap: entity.StudentsNavigation,
		Branch:   *branch,
		User:     user,
		Lang:     Lang,
	}
	resources := reshandler.ResourceService.GetResources(0, 10)
	if resources == nil {
		learningPageDS.Success = false
		learningPageDS.Message = "No Record Found Please Create A Record "
		fmt.Println("No Record Found Please Create A Record ")
		reshandler.Templatehandler.Templates.ExecuteTemplate(response, "learning.html", learningPageDS)
		return
	}
	newActiveResource := reshandler.ResourceService.GetRandomActiveResource()
	if newActiveResource == nil {
		fmt.Println("The Resource is nil ")
		learningPageDS.Success = false
		learningPageDS.Message = "Internal Server Error "
		reshandler.Templatehandler.Templates.ExecuteTemplate(response, "learning.html", learningPageDS)
		return
	}
	learningPageDS.Success = true
	learningPageDS.Message = "Success "
	learningPageDS.Resources = *resources
	learningPageDS.ActiveResource = *newActiveResource
	reshandler.Templatehandler.Templates.ExecuteTemplate(response, "learning.html", learningPageDS)
}

// DownloadResource for downloading the Resources
// Method GET
// Authorization All Logged Is User of the System
func (reshandler *ResourceHandler) DownloadResource(response http.ResponseWriter, request *http.Request) {
	resourceid, era := strconv.Atoi(request.FormValue("resid"))
	if era != nil {
		return
	}
	// Getting the Resource path form the database
	resource := reshandler.ResourceService.GetResourceByID(uint(resourceid))
	if resource == nil {
		return
	}
	http.ServeFile(response, request, entity.PathToTemplates+resource.Path)
	file, erro := os.Open(entity.PathToTemplates + resource.Path)
	if erro != nil {
		return
	}
	reedSeeker := io.ReadSeeker(file)
	http.ServeContent(response, request, entity.PathToStaticTemplateFiles+entity.RoundTemplatingCSV, time.Now(), reedSeeker)
	file.Close()
}	
