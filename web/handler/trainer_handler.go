// Package handler for handling Trainer Related Routes
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

/*
* DetachVehicleFromFieldMan
*	ListOfTrainersOfCategory
*	ListOfFreeTrainers
*	AssignVehicleForTrainer
*	AssignVehicleForFieldMan
	DetachVehicleFromFieldMan
*/

// TrainerHandler struct
type TrainerHandler struct {
	TrainerService Trainer.TrainerService
	SessionService *session.Cookiehandler
	VehicleService Vehicle.VehicleService
}

// NewTrainerHandler function returning TrainerHandler
func NewTrainerHandler(TrainerService Trainer.TrainerService,
	SessionService *session.Cookiehandler,
	VehicleService Vehicle.VehicleService,
) *TrainerHandler {
	return &TrainerHandler{
		TrainerService: TrainerService,
		SessionService: SessionService,
		VehicleService: VehicleService,
	}
}

// ListOfTrainersOfCategory method
// Method GET
// PARAMETER category_id  , offset  , limit
func (trainerhandler *TrainerHandler) ListOfTrainersOfCategory(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	categoryID, era := strconv.Atoi(request.FormValue("category_id"))
	offset, era := strconv.Atoi(request.FormValue("offset"))
	limit, era := strconv.Atoi(request.FormValue("limit"))
	ds := struct {
		Success  bool
		Message  string
		Trainers []entity.FieldAssistant
	}{
		Success: false,
	}
	if era != nil || offset < 0 || limit <= 0 || categoryID <= 0 {
		ds.Message = "Error Input Values "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	trainers := trainerhandler.TrainerService.GetTrainersOfCategory(uint(categoryID), uint(offset), uint(limit))
	if trainers == nil {
		ds.Message = " No Trainer Record Found For the Category  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Message = " Successful"
	ds.Success = true
	ds.Trainers = *trainers
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// ListOfFreeTrainers method
// Method GET
// Variables category_id
func (trainerhandler *TrainerHandler) ListOfFreeTrainers(response http.ResponseWriter, request *http.Request) {
	requestersession := trainerhandler.SessionService.GetSession(request)
	if requestersession == nil {
		response.Write([]byte("<h1>  User Noot Authorized </h1>"))
		return
	}
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		Trainers []entity.FieldAssistant
	}{
		Success: false,
	}
	categoryidstring := request.FormValue("category_id")
	categoryID, era := strconv.Atoi(categoryidstring)
	if era != nil || categoryID <= 0 {
		ds.Message = "Invalid Category ID Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	trainers := trainerhandler.TrainerService.GetFreeTrainers(uint(categoryID))
	if trainers == nil {
		ds.Message = " not trainers for this category "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Succesful "
	ds.Trainers = *trainers
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}

// AssignVehicleForFieldMan method to add assign a field man for vehicle
// Method POST
// AUTHORIZATION SUPERADMIN
// Input JSON the datastructure of the json has to include
// trainer_id and vehicle_id
// remembering all the field man  , vehicle and admin has to be in the same branch
func (trainerhandler *TrainerHandler) AssignVehicleForFieldMan(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestersesion := trainerhandler.SessionService.GetSession(request)
	request.ParseForm()
	if requestersesion == nil {
		return
	}
	BranchID := requestersesion.BranchID
	ds := struct {
		Success        bool
		Message        string
		VehicleToAdmin entity.AddingVehicleToFieldMan
	}{
		Success:        false,
		VehicleToAdmin: entity.AddingVehicleToFieldMan{},
	}
	jsonDecoder := json.NewDecoder(request.Body)
	decodeError := jsonDecoder.Decode(&(ds.VehicleToAdmin))
	if decodeError != nil {
		ds.Message = "Invalid Input ...."
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if ds.VehicleToAdmin.VehicleID <= 0 || ds.VehicleToAdmin.FieldmanID <= 0 {
		fmt.Println(ds.VehicleToAdmin.VehicleID, ds.VehicleToAdmin.FieldmanID)
		ds.Message = " Value is Invalid "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	vehicle := trainerhandler.VehicleService.GetVehicleByID(uint(ds.VehicleToAdmin.VehicleID))
	trainer := trainerhandler.TrainerService.GetTrainerByID(uint(ds.VehicleToAdmin.FieldmanID))
	if vehicle == nil || trainer == nil {
		ds.Message = " Record Not Found "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if !(vehicle.BranchNo == trainer.BranchNumber && trainer.BranchNumber == BranchID) {
		ds.Message = " Invalid Action "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if vehicle.TrainerID > 0 {
		ds.Message = "The Vehicle Already oWner "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	} else if trainer.VehicleID > 0 {
		ds.Message = " The Trainer Already Has A Vehicle  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	trainer.Vehicle = *vehicle
	trainer.VehicleID = vehicle.ID
	vehicle.TrainerID = trainer.ID
	trainer = trainerhandler.TrainerService.SaveTrainer(trainer)
	vehicle = trainerhandler.VehicleService.SaveVehicle(vehicle)
	if trainer == nil || trainer == nil {
		ds.Message = " Internal Server Error Please Try Again "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	ds.Message = "success "
	ds.Success = true
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DetachVehicleFromFieldMan method for detaching a vehicle from a fieldman
// METHOD GET
// VAriable trainer_id
func (trainerhandler *TrainerHandler) DetachVehicleFromFieldMan(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestersession := trainerhandler.SessionService.GetSession(request)
	ds := struct {
		Success   bool
		Message   string
		TrainerID uint
	}{
		Success: false,
	}
	// tookAction := false
	BranchID := requestersession.BranchID
	trainerID, er := strconv.Atoi(request.FormValue("trainer_id"))
	if er != nil {
		ds.Message = "Invalid Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	trainer := trainerhandler.TrainerService.GetTrainerByID(uint(trainerID))
	ds.TrainerID = uint(trainerID)
	if trainer == nil {
		ds.Message = " No Record By This ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if trainer.BranchNumber != BranchID {
		ds.Message = " You Have No Authority For This Action " // when Accessing another branches Data Trainer
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	fmt.Println(trainer.VehicleID)
	if trainer.VehicleID == 0 {
		ds.Message = "The Trainer has no Vehicle Assigned to it "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	vehicle := trainerhandler.VehicleService.GetVehicleByID(trainer.VehicleID)
	if vehicle == nil {
		ds.Message = "Internal Server Error "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	vehicle.TrainerID = 0
	trainer.VehicleID = 0
	trainer.Vehicle = entity.Vehicle{}
	trainer = trainerhandler.TrainerService.SaveTrainer(trainer)
	vehicle = trainerhandler.VehicleService.SaveVehicle(vehicle)
	if trainer == nil || vehicle == nil {
		ds.Message = " Internal Server Error "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = " Successful "
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// GetTrainerByID method to get a json of trainer by ID
// Method GET
// trainer_id
// Authorization SECRETARY  , SUPERADMIN  , FIELDMAN , TEACHER
func (trainerhandler *TrainerHandler) GetTrainerByID(response http.ResponseWriter, request *http.Request) {
	requestersession := trainerhandler.SessionService.GetSession(request)
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success bool
		Message string
		Trainer entity.FieldAssistant
	}{
		Success: false,
	}
	trainerid, era := strconv.Atoi(request.FormValue("trainer_id"))
	if era != nil || trainerid <= 0 {
		ds.Message = "Invalid Request "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	trainer := trainerhandler.TrainerService.GetTrainerByID(uint(trainerid))
	if trainer == nil {
		ds.Message = " Record Not Found "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if trainer.BranchNumber != requestersession.BranchID {
		ds.Message = "You Are Not Authorized to access this Trainer "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = " Success "
	ds.Trainer = *trainer
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// AssignVehicleForTrainer method for adding a vehicle for a traienr
// Method POST
// AUTHORITY SUPERADMIN and SECRETARY
// func (trainerhandler *TrainerHandler) AssignVehicleForTrainer(response http.ResponseWriter, request *http.Request) {
// 	requestersesion := trainerhandler.SessionService.GetSession(request)
// 	if requestersesion == nil {
// 		response.Write([]byte("<h1> Un Authorized User  </h1>"))
// 		return
// 	}
// 	ds := struct {
// 		Success bool
// 		Message string
// 		Trainer entity.FieldAssistant
// 	}{
// 		Success: false,
// 	}
// 	response.Header().Add("Content-Type", "application/json")
// 	vehicleID, era := strconv.Atoi(request.FormValue("vehicle_id"))
// 	trainerID, era := strconv.Atoi(request.FormValue("vehicle_id"))
// 	if era != nil || vehicleID <= 0 || trainerID <= 0 {
// 		ds.Message = "Invalid Inout Value "
// 		jsonReturn, _ := json.Marshal(ds)
// 		response.Write(jsonReturn)
// 		return
// 	}
// 	vehicle := trainerhandler.VehicleService.GetVehicleByID(uint(vehicleID))
// 	if vehicle == nil {
// 		ds.Message = "No Vehicle By This ID "
// 		jsonReturn, _ := json.Marshal(ds)
// 		response.Write(jsonReturn)
// 		return
// 	}
// 	trainer := trainerhandler.TrainerService.GetTrainerByID(uint(trainerID))
// 	if trainer == nil {
// 		ds.Message = " No Trainer by this ID "
// 		jsonReturn, _ := json.Marshal(ds)
// 		response.Write(jsonReturn)
// 		return
// 	}
// 	trainer.VehicleID = vehicle.ID
// 	trainer.Vehicle = *vehicle
// 	trainer = trainerhandler.TrainerService.SaveTrainer(trainer)
// 	if trainer == nil {
// 		ds.Message = "Internal Server Error"
// 		jsonReturn, _ := json.Marshal(ds)
// 		response.Write(jsonReturn)
// 		return
// 	}
// 	ds.Success = true
// 	ds.Message = "Succesful"
// 	ds.Trainer = *trainer
// 	jsonReturn, _ := json.Marshal(ds)
// 	response.Write(jsonReturn)
// }

// PassTraining method for Trainer to pass a training
// Method GET
// Authorization is TRAINER (FieldAssistant)
// func (trainerhandler *TrainerHandler) PassTrainineg(response http.ResponseWriter, request *http.Request) {
// 	reqses := trainerhandler.SessionService.GetSession(request)
// }
