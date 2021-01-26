package RoomService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// RoomService struct
type RoomService struct {
	RoomRepo Room.RoomRepo
}

// NewRoomService function
func NewRoomService(RoomRepo Room.RoomRepo) Room.RoomService {
	return &RoomService{
		RoomRepo: RoomRepo,
	}
}

// GetRoomByNumber method
func (roomser *RoomService) GetRoomByNumber(BranchID, RoomNumber uint) *entity.Room {
	room, erro := roomser.RoomRepo.GetRoomByNumber(BranchID, RoomNumber)
	if erro != nil {
		return nil
	}
	return room
}

// SaveRoom method to save a room
func (roomser *RoomService) SaveRoom(room *entity.Room) *entity.Room {
	room, erroa := roomser.RoomRepo.SaveRoom(room)
	if erroa != nil {
		return nil
	}
	return room
}

// GetRoomsOfABranch (BranchID uint) *[]entity.Room
func (roomser *RoomService) GetRoomsOfABranch(BranchID uint) *[]entity.Room {
	rooms, era := roomser.RoomRepo.GetRoomsOfABranch(BranchID)
	if era != nil {
		return nil
	}
	return rooms
}

// RoomsFreeInDate method
func (roomser *RoomService) RoomsFreeInDate(BranchID uint, date *etc.Date) *[]entity.Room {
	rooms, erra := roomser.RoomRepo.RoomsFreeInDate(BranchID, date)
	if erra != nil {
		return nil
	}
	return rooms
}
