package Round

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type RoundService interface {
	GetRoundByID(ID uint) *entity.Round
	SaveRound(round *entity.Round) *entity.Round
	CreateRound(round *entity.Round) *entity.Round
	GetRounds(BranchID uint) *[]entity.Round
	DeleteScheduleDataUsingRoundID(RoundID uint) bool
	GetRoundByIDForSchedule(RoundID uint) *entity.Round
	UpdateToRegistration(RoundID uint) bool
	GetActiveRoundsOfCategory(CategoryID uint) *[]entity.Round
	GetRoundsOfCategory(RoundID uint) *[]entity.Round
	IsRoundNumberReseerved(BranchID, CategoryID, RoundNumber uint) bool
	UpdateStudentsCount(RoundID, StudentsCount uint) bool
}
