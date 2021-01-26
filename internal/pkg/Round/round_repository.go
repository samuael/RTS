package Round

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type RoundRepo interface {
	GetRoundByID(ID uint) (*entity.Round, error)
	SaveRound(round *entity.Round) (*entity.Round, error)
	CreateRound(round *entity.Round) (*entity.Round, error)
	GetRounds(BranchID uint) (*[]entity.Round, error)
	DeleteFieldDates(RoundID uint) error
	DeleteLectureDates(RoundID uint) error
	DeleteFieldSessionAndItsStudents(RoundID uint) error
	DeleteLecturesAndItsStudents(RoundID uint) error
	DeleteRoundRelatedRoomsDates(RoundID uint) error
	GetRoundByIDForSchedule(RoundID uint) (*entity.Round, error)
	UpdateToRegistration(RoundID uint) error
	GetActiveRoundsOfCategory(CategoryID uint) (*[]entity.Round, error)
	GetRoundsOfCategory(CategoryID uint) (*[]entity.Round, error)
	IsRoundNumberReseerved(BranchID uint, CategoryID uint, RoundNumber uint) uint
	UpdateStudentsCount(RoundID, StudentsCount uint) error
}
