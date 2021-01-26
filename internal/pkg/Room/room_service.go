package Room

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

type RoomService interface {
	GetRoomByNumber(BranchID, RoomNumber uint) *entity.Room
	SaveRoom(room *entity.Room) *entity.Room
	GetRoomsOfABranch(BranchID uint) *[]entity.Room
	RoomsFreeInDate(BranchID uint, date *etc.Date) *[]entity.Room
}
