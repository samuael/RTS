package RoomRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/jinzhu/gorm"
)

// RoomRepo struct
type RoomRepo struct {
	DB *gorm.DB
}

// NewRoomRepo method
func NewRoomRepo(db *gorm.DB) Room.RoomRepo {
	return &RoomRepo{
		DB: db,
	}
}

// GetRoomByNumber method
func (roomrepo *RoomRepo) GetRoomByNumber(BranchID, RoomNumber uint) (*entity.Room, error) {
	rooms := &entity.Room{
		Number:   RoomNumber,
		Branchid: BranchID,
	}
	newRoomError := roomrepo.DB.Where(rooms).First(rooms).Error

	// POPULATING THE room and some data types
	roomrepo.DB.Model(rooms).Related(&rooms.ReservedDates, "ReservedDates")
	return rooms, newRoomError
}

// SaveRoom method to save the room
func (roomrepo *RoomRepo) SaveRoom(room *entity.Room) (*entity.Room, error) {
	newErrr := roomrepo.DB.Save(room).Error
	return room, newErrr
}

// GetRoomsOfABranch method
func (roomrepo *RoomRepo) GetRoomsOfABranch(BranchID uint) (*[]entity.Room, error) {
	rooms := &[]entity.Room{}
	newError := roomrepo.DB.Where("branchid=?", BranchID).Find(rooms).Error
	for k := 0; k < len(*rooms); k++ {
		newError = roomrepo.DB.Model(&(*rooms)[k]).Related(&(&(*rooms)[k]).ReservedDates, "ReservedDates").Error
	}
	return rooms, newError
}

// RoomsFreeInDate method
func (roomrepo *RoomRepo) RoomsFreeInDate(BranchID uint, date *etc.Date) (*[]entity.Room, error) {
	rooms := &[]entity.Room{}
	errors := roomrepo.DB.Where("branchid=?", BranchID).Find(rooms).Error
	// Rooms related Datas
	for i := 0; i <= len(*rooms); i++ {
		room := (*rooms)[i]
		roomrepo.DB.Model(&room).Related(&room.ReservedDates, "ReservedDates")
	}
	selectedrooms := []entity.Room{}
	for i := 0; i < len(*rooms); i++ {
		room := (*rooms)[i]
		reserved := false
		for _, dat := range room.ReservedDates {
			if date.Day == dat.Day && date.Month == dat.Month && dat.Year == date.Year && dat.Day == date.Day {
				reserved = true
			}
		}
		if !reserved {
			selectedrooms = append(selectedrooms, room)
		}
	}
	*rooms = selectedrooms
	return rooms, errors
}
