package Room

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

type RoomRepo interface {
	GetRoomByNumber(BranchID, RoomNumber uint) (*entity.Room, error)
	SaveRoom(room *entity.Room) (*entity.Room, error)
	GetRoomsOfABranch(BranchID uint) (*[]entity.Room, error)
	RoomsFreeInDate(BranchID uint, date *etc.Date) (*[]entity.Room, error)
}
