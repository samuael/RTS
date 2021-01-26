package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Projects/RidingTrainingSystem/pkg/Helper"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"

	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
	"github.com/Projects/RidingTrainingSystem/pkg/permission"
)

// UserHandler y
type UserHandler struct {
	Studenthandler  *StudentHandler
	Adminhandler    *AdminHandler
	Templatehandler *TemplateHandler
	SessionHandler  *session.Cookiehandler
	BranchService   Branch.BranchService
	StudentService  Student.StudentService
	TeacherService  Teacher.TeacherService
	TrainerService  Trainer.TrainerService
}

// NewUserHandler representing Userhandler handling cases common to all Users of the System
func NewUserHandler(adminhandler *AdminHandler,
	studenthandler *StudentHandler,
	th *TemplateHandler,
	sessioHandler *session.Cookiehandler,
	branchsser Branch.BranchService,
	StudentService Student.StudentService,
	TeacherService Teacher.TeacherService,
	TrainerService Trainer.TrainerService) *UserHandler {
	return &UserHandler{
		Studenthandler:  studenthandler,
		Adminhandler:    adminhandler,
		Templatehandler: th,
		SessionHandler:  sessioHandler,
		BranchService:   branchsser,
		StudentService:  StudentService,
		TeacherService:  TeacherService,
		TrainerService:  TrainerService,
	}
}

// CreateSession function returning a session Struct for any User structs
func CreateSession(user entity.Users) *entity.Session {
	session := &entity.Session{
		ID:       user.GetID(),
		Username: user.GetUsername(),
		Imageurl: user.GetImageURL(),
		Lang:     user.GetLang(),
		BranchID: user.GetBranchID(),
		Role:     user.GetRole(),
	}
	return session
}

// LoginGetWithMessage returning html Page With A message
func (useh *UserHandler) LoginGetWithMessage(response http.ResponseWriter, message string) {
	useh.Templatehandler.Templates.ExecuteTemplate(response, "login.html",
		entity.Informing{Message: message,
			CSRF:     useh.SessionHandler.RandomToken(),
			HasError: true,
			Host:     entity.PROTOCOL + entity.HOST,
			Branch:   *useh.BranchService.GetDefaultBranch(),
			Branches: *useh.BranchService.GetBranchs(),
		})
	return
}

// // PageNotFound representing Page Not Found Page for UnAuthorized Acceses
// func (useh *UserHandler) PageNotFound(response http.ResponseWriter, message string) {
// 	useh.Templatehandler.Templates.ExecuteTemplate(response, "four04.html", nil)
// }

// Controll Page
func (useh *UserHandler) Controll(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {

	} else {

	}
}

// Login handling the Login request for All User
func (useh *UserHandler) Login(response http.ResponseWriter, request *http.Request) {
	ds := struct {
		Success bool
		Message string
	}{
		Success: false,
	}

	if request.Method == http.MethodGet {
		useh.Templatehandler.Templates.ExecuteTemplate(response, "login.html",
			entity.Informing{
				Message:  "",
				CSRF:     useh.SessionHandler.RandomToken(),
				HasError: false,
				Host:     entity.PROTOCOL + entity.HOST,
				Branch:   *useh.BranchService.GetDefaultBranch(),
				Branches: *useh.BranchService.GetBranchs(),
			})
		return
	} else if request.Method == http.MethodPost {
		erro := request.ParseForm()
		if erro != nil {
			return
		}
		api, _ := strconv.ParseBool(request.FormValue("api"))
		if api {
			response.Header().Set("Content-Type", "application/json")
		}
		ok := useh.SessionHandler.ValidateForm(request.FormValue("_csrf"))
		if !ok {
			if api {
				ds.Message = " You Are Not Authorized!"
				response.Header().Set("Content-Type", "application/json")
				response.Write(Helper.MarshalThis(ds))
				return
			}
			http.Error(response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		signUpForm := form.Input{
			VErrors: form.ValidationErrors{},
			Values:  request.Form,
		}

		signUpForm.Required("role", "username", "password")
		role := request.FormValue("role")
		username := strings.Trim(request.FormValue("username"), " ")
		password := request.FormValue("password")
		if !signUpForm.Valid() {
			// Message := "Please Fill the datas "
			mesa := ""
			if username == "" && password != "" {
				mesa = " Please Fill the Username ! "
			} else if password == "" && username != "" {
				mesa = " Please Fill the Password ! "
			} else {
				mesa = "Please Input Username and Passwords Correctly!"
			}
			if api {
				ds.Message = "Please Fill the Values Correctly "
				response.Header().Set("Content-Type", "application/json")
				response.Write(Helper.MarshalThis(ds))
				return
			}
			useh.LoginGetWithMessage(response, mesa)
			// useh.Templatehandler.Templates.ExecuteTemplate(response, "login.html", signUpForm.VErrors)
			return
		}

		var session *entity.Session
		// var user *entity.User
		switch strings.ToUpper(role) {
		case "ADMIN", entity.OWNER:
			{
				user := useh.Adminhandler.AdminService.GetAdmin(username, password)
				if user == nil {
					if api {
						ds.Message = "Invalid Username or Password ! Please Try Again "
						response.Header().Set("Content-Type", "application/json")
						response.Write(Helper.MarshalThis(ds))
						return
					}
					useh.LoginGetWithMessage(response, "Invalid Username or Password ! Please Try Again ")
					return
				}
				session = CreateSession(user)
				break
			}
		case entity.TEACHER:
			{
				user := useh.TeacherService.LogTeacher(username, password)
				if user == nil {
					if api {
						ds.Message = "Invalid Username or Password ! Please Try Again "
						response.Header().Set("Content-Type", "application/json")
						response.Write(Helper.MarshalThis(ds))
						return
					}
					useh.LoginGetWithMessage(response, "Invalid Username or Password ! Please Try Again ")
					return
				}
				session = CreateSession(user)
				break
			}
		case entity.STUDENT:
			{
				user := useh.StudentService.LogStudent(username, password)
				if user == nil {
					if api {
						ds.Message = "Invalid Username or Password ! Please Try Again "
						response.Header().Set("Content-Type", "application/json")
						response.Write(Helper.MarshalThis(ds))
						return
					}
					useh.LoginGetWithMessage(response, "Invalid Username or Password ! Please Try Again ")
					return
				}
				session = CreateSession(user)
				break
			}
		case entity.FIELDMAN:
			{
				user := useh.TrainerService.LogTrainer(username, password)
				if user == nil {
					if api {
						ds.Message = "Invalid Username or Password ! Please Try Again "
						response.Header().Set("Content-Type", "application/json")
						response.Write(Helper.MarshalThis(ds))
						return
					}
					useh.LoginGetWithMessage(response, "Invalid Username or Password ! Please Try Again ")
					return
				}
				session = CreateSession(user)
				break
			}
		default:
			session = nil
			if api {
				ds.Message = " Not Succesful Please Try Again !"
				response.Header().Set("Content-Type", "application/json")
				response.Write(Helper.MarshalThis(ds))
				return
			}
			http.Redirect(response, request, "/login", http.StatusPermanentRedirect)
			break
		}
		if session != nil {
			success := useh.SessionHandler.SaveSession(response, session, request.Host)
			if api && success {
				response.WriteHeader(http.StatusOK)
				ds.Success = true
				ds.Message = "Succesfuly Logged In "
				response.Header().Set("Content-Type", "application/json")
				response.Write(Helper.MarshalThis(ds))
				return
			} else if api && !success {
				ds.Success = false
				ds.Message = " Logging In Was Not Succesful "
				response.Header().Set("Content-Type", "application/json")
				response.Write(Helper.MarshalThis(ds))
				return
			}
			if !success {
				useh.LoginGetWithMessage(response, "Internal Server Error  ! Please Try Again ")
				return
			}
			request.Method = http.MethodGet
			switch session.Role {
			case entity.SUPERADMIN, entity.SECRETART:
				{
					http.Redirect(response, request, "/admin/controll/", http.StatusPermanentRedirect)
					break
				}
			case entity.TEACHER:
				{
					http.Redirect(response, request, "/teacher/controll/", http.StatusPermanentRedirect)
					break
				}
			case entity.FIELDMAN:
				{
					http.Redirect(response, request, "/trainer/controll/", http.StatusPermanentRedirect)
					break
				}
			case entity.STUDENT:
				{
					http.Redirect(response, request, "/student/controll/", http.StatusMovedPermanently)
					break
				}
			}
		}
	}
}

// Logout method to Log Out  this works for all the Users of the System
func (useh *UserHandler) Logout(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestSession := useh.SessionHandler.GetSession(request)
	lang := requestSession.Lang
	datastructure := struct {
		Success bool
		Message string
	}{
		Success: false,
	}
	// Saving Lang Session
	useh.SessionHandler.SaveLang(response, lang, request.Host)
	success := useh.SessionHandler.DeleteSession(response, request)
	if success {
		datastructure.Success = true
		datastructure.Message = " Succesfully Logged Out Niggory "
	}
	jsonReturn, _ := json.Marshal(datastructure)
	response.Write(jsonReturn)
}

// Authenticated checks if a user is authenticated to access a given route
func (useh *UserHandler) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := useh.LoggedIn(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		//  Here the Context Use is uncear
		// ctx := context.WithValue(r.Context(), ctxUserSessionKey, uh.userSess)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// PageNotFound method returning the Page Not Found Template embeded in the Body
func (useh *UserHandler) PageNotFound(response http.ResponseWriter, request *http.Request, statusText string, StatusCode int) {
	path := request.URL.Path
	Four04Info := struct {
		Path       string
		StatusCode int
		Statustext string
		Host       string
	}{
		Path:       path,
		Statustext: statusText,
		StatusCode: StatusCode,
		Host:       entity.PROTOCOL + entity.HOST,
	}
	useh.Templatehandler.Templates.ExecuteTemplate(response, "four04.html", Four04Info)
	return
}

// Authorized checks if a user has proper authority to access a give route
func (useh *UserHandler) Authorized(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !useh.LoggedIn(r) {
			// http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			useh.PageNotFound(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		session := useh.SessionHandler.GetSession(r)
		if session == nil {
			useh.PageNotFound(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			// http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// fmt.Println("Here ", session.Role)
		role := session.Role
		permitted := permission.HasPermission(r.URL.Path, role, r.Method)
		if !permitted {
			useh.PageNotFound(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			// http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// if r.Method == http.MethodPost {
		// 	// erro := r.ParseForm()
		// 	// if erro != nil {
		// 	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	// 	return
		// 	// }
		// }
		next.ServeHTTP(w, r)
	})
}

// LoggedIn checks whether the user is Authenticated or not
func (useh *UserHandler) LoggedIn(request *http.Request) bool {
	session := useh.SessionHandler.GetSession(request)
	if session != nil {
		return true
	}
	return false
}
