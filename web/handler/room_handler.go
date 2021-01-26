package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// RoomHandler strunct
type RoomHandler struct {
	RoomService    Room.RoomService
	SessionHandler *session.Cookiehandler
}

// NewRoomHandler method
func NewRoomHandler(roomservice Room.RoomService, Session *session.Cookiehandler) *RoomHandler {
	return &RoomHandler{
		RoomService:    roomservice,
		SessionHandler: Session,
	}
}

// GetRoomsOfBranch method   Method Get
// variable branch_id
// authorization  for any one
func (roomhandler *RoomHandler) GetRoomsOfBranch(response http.ResponseWriter, request *http.Request) {
	// Data Dtructure representing rooms of a branch
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		Rooms    []entity.Room
		BranchID uint
	}{
		Success: false,
		Rooms:   []entity.Room{},
	}
	branchidstring := request.FormValue("branch_id")
	branchid, parseerra := strconv.Atoi(branchidstring)
	if branchidstring == "" || parseerra != nil || branchid <= 0 {
		ds.Message = " Invalid Input || Branch ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.BranchID = uint(branchid)
	rooms := roomhandler.RoomService.GetRoomsOfABranch(ds.BranchID)
	if rooms == nil {
		ds.Message = " Record Not Found  "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Rooms = *rooms
	ds.Success = true
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
	return
}
