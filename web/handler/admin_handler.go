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

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/Projects/RidingTrainingSystem/pkg/form"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
)

// AdminHandler yes
type AdminHandler struct {
	AdminService    Admin.AdminService
	Templatehandler *TemplateHandler
	SessionHandler  *session.Cookiehandler
	BranchHandler   *BranchHandler
	StudentService  Student.StudentService
	TeacherService  Teacher.TeacherService
	TrainerService  Trainer.TrainerService
	CategoryService Category.CategoryService
	VehicleService  Vehicle.VehicleService
	RoundService    Round.RoundService
	RoomService     Room.RoomService
}

// NewAdminHandler returning admin handler for secretart and Super Admins
func NewAdminHandler(adminservice Admin.AdminService,
	temphan *TemplateHandler,
	sessHandler *session.Cookiehandler,
	branchh *BranchHandler,
	studentser *Student.StudentService,
	Teacher Teacher.TeacherService,
	trinerser Trainer.TrainerService,
	Cateser Category.CategoryService,
	vehicle Vehicle.VehicleService,
	roundSer Round.RoundService,
	RoomService Room.RoomService,
) *AdminHandler {
	return &AdminHandler{
		AdminService:    adminservice,
		Templatehandler: temphan,
		SessionHandler:  sessHandler,
		BranchHandler:   branchh,
		StudentService:  *studentser,
		TeacherService:  Teacher,
		TrainerService:  trinerser,
		CategoryService: Cateser,
		VehicleService:  vehicle,
		RoundService:    roundSer,
		RoomService:     RoomService,
	}
}

/*
Controll :  returning controll Page for admin
PrepareInfo : not Handler Func to prepare Info
SuperAdminRegistration : Handkler to handle Admin Registration returning the Page
AdminRegistration : Handler func to handle Admins
TeacherRegistration : Handle Func to Handle Teachers
TrainerRegistration : Handle rFunc to handle Trainers
*/
// Controll representing the Controll for SUPERADMIN and SECRETART
func (adminh *AdminHandler) Controll(response http.ResponseWriter, request *http.Request) {
	sessionOfAdmin := adminh.SessionHandler.GetSession(request)

	branch := adminh.BranchHandler.BranchService.GetBranchByID(sessionOfAdmin.BranchID)
	admin := adminh.AdminService.GetAdminByID(sessionOfAdmin.ID, sessionOfAdmin.Username)
	categories := adminh.CategoryService.GetCategories(sessionOfAdmin.BranchID)
	rooms := adminh.RoomService.GetRoomsOfABranch(sessionOfAdmin.BranchID)

	ControllPageStructure := entity.ControllPageStructure{
		Branch:     *branch,
		Rooms:      *rooms,
		Categories: categories,
		Admin:      *admin,
		Host:       entity.PROTOCOL + entity.HOST,
	}
	adminh.Templatehandler.Templates.ExecuteTemplate(response, "ControllPage.html", ControllPageStructure)
}

// PrepareInfo method
func (adminh *AdminHandler) PrepareInfo(request *http.Request, inputo form.Input) entity.Informing {
	session := adminh.SessionHandler.GetSession(request)
	branch := adminh.BranchHandler.BranchService.GetBranchByID(session.BranchID)
	admin := adminh.AdminService.GetAdminByID(session.ID, session.Username)
	info := entity.Informing{
		CSRF:     adminh.SessionHandler.RandomToken(),
		Branch:   *branch,
		Admin:    *admin,
		Host:     request.Host,
		Input:    inputo,
		Specific: adminh.CategoryService.GetCategories(session.BranchID),
	}
	return info
}

// SuperadminRegistration handling post Request and routing to only Admins Page
func (adminh *AdminHandler) SuperadminRegistration(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		session := adminh.SessionHandler.GetSession(request)
		branch := adminh.BranchHandler.BranchService.GetBranchByID(session.BranchID)
		admin := adminh.AdminService.GetAdminByID(session.ID, session.Username)

		adminregisterPageStructure := entity.AdminRegistrationsPageStructure{
			Host:       entity.PROTOCOL + entity.HOST,
			Branch:     *adminh.BranchHandler.BranchService.GetBranchByID(session.BranchID),
			Admin:      *admin,
			CSRF:       adminh.SessionHandler.RandomToken(),
			Categories: adminh.CategoryService.GetCategories(uint(branch.ID)),
		}
		// vehicles := adminh.VehicleService.GetFreeVehicles
		adminh.Templatehandler.Templates.ExecuteTemplate(response, "admin_register.html", adminregisterPageStructure)
		return
	}
	fmt.Println("Invalid Method ")
	// return something the Method is not Correct
}

// RegisterAdmin handling post Request and routing to only Admins Page
func (adminh *AdminHandler) RegisterAdmin(response http.ResponseWriter, request *http.Request) {
	reqses := adminh.SessionHandler.GetSession(request)
	if reqses == nil {
		return
	}
	if request.Method == http.MethodPost {
		request.ParseMultipartForm(3583453578987)
		singnUpForm := form.Input{Values: request.PostForm, VErrors: form.ValidationErrors{}}
		singnUpForm.Required("firstname", "lastname", "grandname", "email", "role", "phone")
		singnUpForm.MatchesPattern(request.FormValue("email"), form.EmailRX)
		singnUpForm.CSRF = adminh.SessionHandler.RandomToken()
		if !singnUpForm.Valid() {
			input := adminh.PrepareInfo(request, singnUpForm)
			fmt.Println(singnUpForm.VErrors)
			adminh.Templatehandler.Templates.ExecuteTemplate(response, "admin_register.html", input)
			return
		}
		session := adminh.SessionHandler.GetSession(request)
		branch := adminh.BranchHandler.BranchService.GetBranchByID(session.BranchID)
		admin := adminh.AdminService.GetAdminByID(session.ID, session.Username)
		imagedir := "img/adminImages/"
		imageurl := ""
		imageDirectory, header, erro := request.FormFile("profilepic")
		if erro != nil {
			fmt.Println("the Profile Picture is Null ")
		}
		imageurl = Helper.GenerateRandomString(10, Helper.CHARACTERS)
		filenew, fileError := os.Create("../../web/templates/" + imagedir + imageurl + header.Filename)
		imageurl = imagedir + imageurl
		if fileError != nil {
			log.Println("Error while saving the  file ")
			imageurl = ""
		}
		defer func() {
			imageDirectory.Close()
			filenew.Close()
		}()
		_, err := io.Copy(filenew, imageDirectory)
		if err != nil {
			imageurl = ""
		}
		Role := "SUPERADMIN"
		switch request.FormValue("role") {
		case "admin":
			{
				Role = entity.SUPERADMIN
				break
			}
		case "secretary":
			{
				Role = entity.SECRETART
				break
			}
		default:
			Role = entity.SECRETART
		}
		if Role == entity.SUPERADMIN {
			if reqses.Role != entity.OWNER {
				response.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		newAdmin := &entity.Admin{
			Branch:    *branch,
			Email:     request.FormValue("email"),
			Lang:      request.FormValue("lang"),
			Firstname: strings.Trim(request.FormValue("firstname"), " "),
			Lastname:  request.FormValue("lastname"),
			GrandName: request.FormValue("grandname"),
			Role:      Role,
			Phone:     request.FormValue("phone"),
			Createdby: admin.Username,
			Password:  Helper.GenerateRandomString(4, Helper.NUMBERS),
		}
		newAdmin.Imageurl = imageurl
		newAdmin.Username = "Admin/" + strings.Trim(newAdmin.Firstname, " ") + "/" + strconv.Itoa(int(adminh.AdminService.AdminsCount()))
		newAdmin = adminh.AdminService.RegisterAdmin(newAdmin)
		if newAdmin == nil {
			fmt.Println("Error  wile Saving the admin")
		}
		response.Header().Add("Content-Type", "application/json")
		jsonized, jsonerror := json.Marshal(newAdmin)
		if jsonerror != nil {
			fmt.Println("Json Error ")
		}
		response.Write(jsonized)
	} else {
		// return something the Method is not Correct
	}
}

// RegisterTeacher method POST
func (adminh *AdminHandler) RegisterTeacher(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		request.ParseMultipartForm(3583453578987)
		singnUpForm := form.Input{Values: request.PostForm, VErrors: form.ValidationErrors{}}
		singnUpForm.Required("firstname", "lastname", "grandname", "email", "phone")
		singnUpForm.MatchesPattern(request.FormValue("email"), form.EmailRX)
		singnUpForm.MinLength(request.FormValue("phone"), 10)
		singnUpForm.CSRF = adminh.SessionHandler.RandomToken()
		if !singnUpForm.Valid() {
			input := adminh.PrepareInfo(request, singnUpForm)
			response.Header().Add("Path", "/admin/registration/")
			adminh.Templatehandler.Templates.ExecuteTemplate(response, "admin_register.html", input)
			return
		}
		session := adminh.SessionHandler.GetSession(request)
		// branch := adminh.BranchHandler.BranchService.GetBranchByID(session.BranchID)
		admin := adminh.AdminService.GetAdminByID(session.ID, session.Username)
		imagedir := "img/Teachers/"
		imageurl := ""
		imageDirectory, header, erro := request.FormFile("profilepic")
		if erro == nil || header != nil {
			imageurl = Helper.GenerateRandomString(10, Helper.CHARACTERS)
			filenew, fileError := os.Create("../../web/templates/" + imagedir + imageurl + header.Filename)
			imageurl = imagedir + imageurl
			if fileError != nil {
				log.Println("Error while saving the  file ")
				imageurl = ""
			}
			defer func() {
				imageDirectory.Close()
				filenew.Close()
			}()
			_, err := io.Copy(filenew, imageDirectory)
			if err != nil {
				imageurl = ""
			}
		}
		newTeacher := &entity.Teacher{
			BranchNumber: session.BranchID,
			Email:        request.FormValue("email"),
			Lang:         request.FormValue("lang"),
			Firstname:    strings.Trim(request.FormValue("firstname"), " "),
			Lastname:     request.FormValue("lastname"),
			GrandName:    request.FormValue("grandname"),
			Phonenumber:  request.FormValue("phone"),
			Createdby:    admin.Username,
			Password:     Helper.GenerateRandomString(4, Helper.NUMBERS),
			BusyDates:    []etc.Date{},
		}
		newTeacher.Imageurl = imageurl
		newTeacher.Username = "Teacher/" + newTeacher.Firstname + "/" + strconv.Itoa(int(adminh.TeacherService.TeachersCount()))
		newTeacher = adminh.TeacherService.SaveTeacher(newTeacher)
		if newTeacher == nil {
			fmt.Println("Error  wile Saving the admin")
			http.Redirect(response, request, "/admin/registration/", 300)
			return
		}
		response.Header().Add("Content-Type", "application/json")
		jsonized, jsonerror := json.Marshal(newTeacher)
		if jsonerror != nil {
			fmt.Println("Json Error ")
		}
		response.Write(jsonized)
	} else {
		// return something the Method is not Correct
	}
}

// RegisterFieldMan method POST
// Method POST this creates a fieldman for category
func (adminh *AdminHandler) RegisterFieldMan(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		request.ParseMultipartForm(35834538987)
		singnUpForm := form.Input{Values: request.PostForm, VErrors: form.ValidationErrors{}}
		singnUpForm.Required("firstname", "lastname", "categoryid", "grandname", "phone")
		if request.FormValue("email") != "" {
			singnUpForm.MatchesPattern(request.FormValue("email"), form.EmailRX)
		}
		singnUpForm.MinLength(request.FormValue("phone"), 10)
		singnUpForm.CSRF = adminh.SessionHandler.RandomToken()
		vehicleid := uint(0)
		var erro error
		var number int
		if request.FormValue("vehicleid") != "" {
			number, erro = strconv.Atoi(request.FormValue("vehicle"))
			if erro != nil {
				log.Println("Error While Taking the Vehicle ID ")
				singnUpForm.VErrors.Add("error", "Vehicle Has To Be Selected Correctely ")
			}
		}
		vehicleid = uint(number)
		if !singnUpForm.Valid() {
			input := adminh.PrepareInfo(request, singnUpForm)
			response.Header().Add("path", "/admin/registration/")
			// http.Redirect(response, request, "/admin/registration/", http.StatusMovedPermanently)
			adminh.Templatehandler.Templates.ExecuteTemplate(response, "admin_register.html", input)
			return
		}
		session := adminh.SessionHandler.GetSession(request)
		BranchNumber := session.BranchID
		// branch := adminh.BranchHandler.BranchService.GetBranchByID(session.BranchID)
		admin := adminh.AdminService.GetAdminByID(session.ID, session.Username)
		imagedir := "img/Trainer/"
		imageurl := ""
		imageDirectory, header, erro := request.FormFile("profilepic")
		if erro != nil {
			imageurl = ""
		}
		number, errors := strconv.Atoi(request.FormValue("categoryid"))
		var CategoryID uint
		CategoryID = 0
		if errors == nil {
			CategoryID = uint(number)
		} else {
			singnUpForm.VErrors.Add("Category ID Error ", "Invalid Category ID ")
		}
		category := adminh.CategoryService.GetCategoryByID(CategoryID)
		if category != nil {
			singnUpForm.VErrors.Add("General ", "Invalid Category Selection")
		}
		vehicleid = uint(number)
		vehicle := adminh.VehicleService.GetVehicleByID(vehicleid)
		if vehicle == nil {
			singnUpForm.VErrors.Add("Vehicle ", "Vehicle Not Found ")
		}
		if !singnUpForm.Valid() {
			input := adminh.PrepareInfo(request, singnUpForm)
			response.Header().Add("path", "/admin/registration/")
			// http.Redirect(response, request, "/admin/registration/", http.StatusMovedPermanently)
			adminh.Templatehandler.Templates.ExecuteTemplate(response, "admin_register.html", input)
			return
		}
		imageurl = Helper.GenerateRandomString(10, Helper.CHARACTERS)
		if header != nil {
			if header.Filename != "" {
				filenew, fileError := os.Create("../../web/templates/" + imagedir + imageurl + header.Filename)
				fmt.Println("Ezihgar ... 1")
				imageurl = imagedir + imageurl
				if fileError != nil {
					log.Println("Error while saving the  file ")
					imageurl = ""
				}
				defer func() {
					imageDirectory.Close()
					filenew.Close()
				}()
				_, err := io.Copy(filenew, imageDirectory)
				if err != nil {
					imageurl = ""
				}
			}
		}
		newTrainer := &entity.FieldAssistant{
			Firstname:    strings.Trim(request.FormValue("firstname"), " "),
			Lastname:     strings.Trim(request.FormValue("lastname"), " "),
			GrandName:    request.FormValue("grandname"),
			Email:        request.FormValue("email"),
			BranchNumber: BranchNumber,
			Lang:         request.FormValue("lang"),
			Phonenumber:  request.FormValue("phone"),
			Createdby:    admin.Username,
			CategoryID:   CategoryID,
			Categoty:     *category,
			VehicleID:    vehicleid,
			Vehicle:      *vehicle,
			Password:     Helper.GenerateRandomString(4, Helper.NUMBERS),
			BusyDates:    []etc.Date{},
		}
		vehicle.Reserved = true
		newTrainer.Imageurl = imageurl
		newTrainer.Username = "Trainer/" + newTrainer.Firstname + "/" + strconv.Itoa(int(adminh.TrainerService.GetCount()))
		newTrainer = adminh.TrainerService.SaveTrainer(newTrainer)
		if newTrainer == nil {
			// http.RedirectHandler("/admin/registration"  ,  http.StatusMovedPermanently )
			http.Redirect(response, request, "/admin/registration/", http.StatusMovedPermanently)
			return
		}
		response.Header().Add("Content-Type", "application/json")
		jsonized, jsonerror := json.Marshal(newTrainer)
		if jsonerror != nil {
			fmt.Println("Json Error ")
		}
		response.Write(jsonized)
	} else {
		// return something the Method is not Correct
	}
}
