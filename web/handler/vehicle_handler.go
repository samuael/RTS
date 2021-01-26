// Package handler packahe for handling the Routing
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
)

// VehicleHandler struct
type VehicleHandler struct {
	VehicleService Vehicle.VehicleService
	SessionService *session.Cookiehandler
	TrainerService Trainer.TrainerService
}

// NewVehicleHandler function
func NewVehicleHandler(
	VehicleService Vehicle.VehicleService,
	SessionService *session.Cookiehandler,
	TrainerService Trainer.TrainerService,
) *VehicleHandler {
	return &VehicleHandler{
		VehicleService: VehicleService,
		SessionService: SessionService,
		TrainerService: TrainerService,
	}
}

// CreateVehicle method
// Method POST
// Input MULTIPART FORM REQUEST
func (vehiclehandler *VehicleHandler) CreateVehicle(response http.ResponseWriter, request *http.Request) {
	requestersession := vehiclehandler.SessionService.GetSession(request)
	if requestersession == nil {
		response.Write([]byte("<h1>  UnAuthorized User  </h1>"))
		return
	}
	era := request.ParseMultipartForm(123456789)
	if era != nil {
		response.Write([]byte("<h1>  Request Error  </h1>"))
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	ds := struct {
		Success bool
		Message string
		Vehicle entity.Vehicle
	}{
		Success: false,
	}
	response.Header().Add("Content-Type", "application/json")
	boardNumberstring := request.FormValue("board_number")
	categoryidstring := request.FormValue("category_id")
	categoryID, errors := strconv.Atoi(categoryidstring)
	if errors != nil || boardNumberstring == "" {
		ds.Message = " Input Values Error "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	image, header, erroe := request.FormFile("images")
	hasImage := true
	var imageDirectory string
	if erroe != nil {
		hasImage = false
	} else {
		defer image.Close()
	}
	var newImage *os.File
	var eroo error
	if hasImage {
		randomname := Helper.GenerateRandomString(5, Helper.NUMBERS)
		newImage, eroo = os.Create(entity.PATHToVehiclesImageFromMain + randomname + header.Filename)
		if eroo != nil {
			fmt.Println("Error Whille Creating a File ")
			newImage = nil
			imageDirectory = ""
			hasImage = false
		} else {
			newImage.Close()
		}
		filename := strings.TrimPrefix(entity.PATHToVehiclesImageFromMain+randomname+header.Filename, "../..")
		imageDirectory = filename
	}
	if hasImage {
		io.Copy(newImage, image)
	}
	vehicle := &entity.Vehicle{
		BranchNo:    requestersession.BranchID,
		CategoryID:  uint(categoryID),
		BoardNumber: boardNumberstring,
	}
	if hasImage {
		vehicle.Imageurl = imageDirectory
	}
	vehicle = vehiclehandler.VehicleService.SaveVehicle(vehicle)
	if vehicle == nil {
		ds.Message = " Error While saving the vehicle"
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Success ful "
	ds.Vehicle = *vehicle
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DeleteVehicle method to delete UnReserved Vehicles
// Method POST
// AUTHORIZATION SUPERADMIN   , SECCRETARY
//
func (vehiclehandler *VehicleHandler) DeleteVehicle(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	requestersession := vehiclehandler.SessionService.GetSession(request)
	if requestersession == nil {
		response.Write([]byte("<h1>  UnAuthorized User  </h1>"))
		return
	}
	response.Header().Add("Content-Type", "application/json")
	vehicleidstring := request.FormValue("vehicle_id")
	ds := struct {
		Success   bool
		Message   string
		VehicleID uint
	}{
		Success: false,
	}
	vehicleID, era := strconv.Atoi(vehicleidstring)
	if era != nil {
		ds.Message = "Invalid Data Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}

	isReserved := vehiclehandler.VehicleService.IsVehicleReserved(uint(vehicleID))
	if isReserved {
		ds.Message = "The Vehicle is Still In Use First Please Detach The Veihcle "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	success := vehiclehandler.VehicleService.DeleteVehicle(uint(vehicleID))
	if !success {
		ds.Message = " Error While Deleting the Vehicle "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Succes fully Deleted the Vehicle "
	ds.VehicleID = uint(vehicleID)
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// FreeVehiclesOfCategory method
// METHOD GET
// variable category_id
func (vehiclehandler *VehicleHandler) FreeVehiclesOfCategory(response http.ResponseWriter, request *http.Request) {
	requestersession := vehiclehandler.SessionService.GetSession(request)
	if requestersession == nil {
		response.Write([]byte("<h1> Not Authorized </h1>"))
		return
	}
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		Vehicles []entity.Vehicle
	}{
		Success: false,
	}
	categoryID, era := strconv.Atoi(request.FormValue("category_id"))
	if era != nil {
		ds.Message = " Invalid Input "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	vehicles := vehiclehandler.VehicleService.GetFreeVehiclesOfCategory(uint(categoryID))
	if vehicles == nil {
		ds.Message = "No Record vehicle Record found For This Category "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Vehicles = *vehicles
	ds.Message = "Succesfull "
	ds.Success = true
	jsonRturn, _ := json.Marshal(ds)
	response.Write(jsonRturn)
}
