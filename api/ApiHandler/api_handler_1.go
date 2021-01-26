package ApiHandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// APIHandler struct
type APIHandler struct {
	VehicleService Vehicle.VehicleService
	Session        *session.Cookiehandler
	RoomService    Room.RoomService
	AdminService   Admin.AdminService
	StudentService Student.StudentService
	TeacherService Teacher.TeacherService
	TrainerService Trainer.TrainerService
}

// NewAPIHandler function
func NewAPIHandler(
	vehicleservice Vehicle.VehicleService,
	sess *session.Cookiehandler,
	AdminService Admin.AdminService,
	Studentser Student.StudentService,
	Teacherser Teacher.TeacherService,
	Trainerser Trainer.TrainerService,
	RoomService Room.RoomService) *APIHandler {
	return &APIHandler{
		VehicleService: vehicleservice,
		Session:        sess,
		AdminService:   AdminService,
		TeacherService: Teacherser,
		TrainerService: Trainerser,
		RoomService:    RoomService,
	}
}

// GetVehicles method returning the list of vehicles taking the CategoryId  as a list
func (apihandler *APIHandler) GetVehicles(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// Empty Json Return
	vehicles := []entity.Vehicle{}
	jsonReturn, erro := json.Marshal(vehicles)
	if erro != nil {
		fmt.Println("Error Has Happened But so Be it ")
	}
	jsonDecoder := json.NewDecoder(request.Body)
	// var categoryID uint
	categoryids := struct {
		CategoryID uint `json:"categoryid"`
	}{
		CategoryID: 0,
	}
	decodeEror := jsonDecoder.Decode(&categoryids)
	if decodeEror != nil {
		response.Write(jsonReturn)
		return
	}
	sessionFromHead := apihandler.Session.GetSession(request)
	vehicles = apihandler.VehicleService.GetVehicles(categoryids.CategoryID, sessionFromHead.BranchID)
	newJSONReturn, theErr := json.Marshal(vehicles)
	if theErr == nil {
		response.Write(newJSONReturn)
		return
	}
	response.Write(jsonReturn)
}

// ChangePassword method  having a csrf string for validation
func (apihandler *APIHandler) ChangePassword(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	SavedSession := apihandler.Session.GetSession(request)
	username := SavedSession.Username
	id := SavedSession.ID
	role := SavedSession.Role
	aboutUser := struct {
		Success         bool
		ID              uint
		Role            string
		Username        string
		NewPassword     string `json:"newpassword"`
		OldPassword     string `json:"oldpassword"`
		ConfirmPassword string `json:"confirmpassword"`
		Message         string
	}{
		Success:  false,
		ID:       id,
		Username: username,
	}

	newDecoder := json.NewDecoder(request.Body)
	decodeError := newDecoder.Decode(&aboutUser)
	aboutUser.ID = id
	aboutUser.Username = username
	aboutUser.Role = role
	passwordMatches := true
	if aboutUser.NewPassword != aboutUser.ConfirmPassword {
		aboutUser.Message = "  The New Password Has To Match "
		passwordMatches = false
	}
	jsonReturn, _ := json.Marshal(aboutUser)
	if decodeError != nil || !passwordMatches {
		response.Write(jsonReturn)
		return
	}
	switch role {
	case entity.SUPERADMIN,
		entity.SECRETART,
		entity.OWNER:
		{
			admin := apihandler.AdminService.GetAdmin(username, aboutUser.OldPassword)
			admin.Password = aboutUser.NewPassword
			admin = apihandler.AdminService.UpdateAdmin(admin)
			if admin != nil {
				aboutUser.Success = true
				aboutUser.Message = "Password Succesfully Changed "
				jsonReturn, _ := json.Marshal(aboutUser)
				response.Write(jsonReturn)
			} else {
				response.Write(jsonReturn)
			}
			break
		}
	case entity.TEACHER:
		{
			teacher := apihandler.TeacherService.GetTeacherByID(id)
			teacher.Password = aboutUser.NewPassword
			teacher = apihandler.TeacherService.SaveTeacher(teacher)
			if teacher != nil {
				aboutUser.Success = true
				aboutUser.Message = "Password Succesfully Changed "
				jsonReturn, _ := json.Marshal(aboutUser)
				response.Write(jsonReturn)
			} else {
				response.Write(jsonReturn)
			}
			break
		}
	case entity.FIELDMAN:
		{
			trainer := apihandler.TrainerService.GetTrainerByID(id)
			trainer.Password = aboutUser.NewPassword
			trainer = apihandler.TrainerService.SaveTrainer(trainer)
			if trainer != nil {
				aboutUser.Success = true
				aboutUser.Message = "Password Succesfully Changed "
				jsonReturn, _ := json.Marshal(aboutUser)
				response.Write(jsonReturn)
			} else {
				response.Write(jsonReturn)
			}
			break
		}
	case entity.STUDENT:
		{
			student := apihandler.StudentService.GetStudentByID(id)
			student.Password = aboutUser.NewPassword
			student = apihandler.StudentService.SaveStudent(student)
			if student != nil {
				aboutUser.Success = true
				aboutUser.Message = "Password Succesfully CHanged "
				jsonREturn, _ := json.Marshal(aboutUser)
				response.Write(jsonREturn)
			} else {
				response.Write(jsonReturn)
			}
			break
		}

	}

}

// RegisterRoom method
func (apihandler *APIHandler) RegisterRoom(response http.ResponseWriter, request *http.Request) {
	newRoom := &entity.Room{}
	response.Header().Add("Content-Type", "application/json")
	newDecoder := json.NewDecoder(request.Body)
	decodeError := newDecoder.Decode(newRoom)
	responseJSON := struct {
		Success bool   `json:"success,omitempty"`
		Message string `json:"message,omitempty"`
		Room    entity.Room
	}{
		Success: false,
		Message: "Invalid Request Body ",
	}
	if decodeError != nil {
		jsonReturn, _ := json.Marshal(responseJSON)
		response.Write(jsonReturn)
		return
	}
	sessionValue := apihandler.Session.GetSession(request)
	Branchid := sessionValue.BranchID
	newRoom.Branchid = Branchid
	newRoom.CreatedBy = sessionValue.Username
	newRoom.ReservedDates = []etc.Date{*etc.NewDate(0), *etc.NewDate(1)}
	// Check whether the Branch Capacity and number are taken Correctly
	if newRoom.Number <= 0 && newRoom.Capacity <= 0 {
		responseJSON.Message = "Please Enter Valid Number for Room Number and Room Capacity "
		jsonReturn, _ := json.Marshal(responseJSON)
		response.Write(jsonReturn)
		return
	}
	//  Check whether there is a room of BranchId and RoomNumber does Exist or not
	room := apihandler.RoomService.GetRoomByNumber(newRoom.Branchid, newRoom.Number)
	if room != nil {
		responseJSON.Success = false
		responseJSON.Message = "The Room Does Exit Please Enter Another Room Number "
		jsonReturn, _ := json.Marshal(responseJSON)
		response.Write(jsonReturn)
		return
	}
	// Saving the Room and if it is added Succesfuly Response Success aving the Branc Will Be true
	newRoom = apihandler.RoomService.SaveRoom(newRoom)
	if newRoom == nil {
		responseJSON.Message = "Sorry Can't Save The Room \n Please Try Again "
		jsonReturn, _ := json.Marshal(responseJSON)
		response.Write(jsonReturn)
		return
	}
	responseJSON.Success = true
	responseJSON.Message = " Succesfully Saved the room"
	responseJSON.Room = *newRoom
	jsonReturn, _ := json.Marshal(responseJSON)
	response.Write(jsonReturn)
	return
}
