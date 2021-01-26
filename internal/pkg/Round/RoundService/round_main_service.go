package RoundService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// RoundService struct
type RoundService struct {
	RoundRepo Round.RoundRepo
}

// NewRoundService function
func NewRoundService(roundRepo Round.RoundRepo) Round.RoundService {
	return &RoundService{
		RoundRepo: roundRepo,
	}
}

// GetRoundByID method returning fully generated round struct
func (roundser *RoundService) GetRoundByID(ID uint) *entity.Round {
	round, Xerror := roundser.RoundRepo.GetRoundByID(ID)
	if Xerror != nil {
		return nil
	}
	return round
}

// Using the Field Man ID i can Create A field Man for Specific Round

// SaveRound method
func (roundser *RoundService) SaveRound(round *entity.Round) *entity.Round {
	newRound, erroa := roundser.RoundRepo.SaveRound(round)
	if erroa != nil {
		return nil
	}
	return newRound
}

// CreateRound method
func (roundser *RoundService) CreateRound(round *entity.Round) *entity.Round {
	newRound, erroa := roundser.RoundRepo.SaveRound(round)
	if erroa != nil {
		return nil
	}
	return newRound
}

// GetRounds method
func (roundser *RoundService) GetRounds(BranchID uint) *[]entity.Round {
	round, erro := roundser.RoundRepo.GetRounds(BranchID)
	if erro != nil {
		return nil
	}
	return round
}

// DeleteScheduleDataUsingRoundID (RoundID uint) bool
func (roundser *RoundService) DeleteScheduleDataUsingRoundID(RoundID uint) bool {
	newErra := roundser.RoundRepo.DeleteFieldDates(RoundID)
	if newErra != nil {
		return false
	}
	newErra = roundser.RoundRepo.DeleteLectureDates(RoundID)
	if newErra != nil {
		return false
	}
	newErra = roundser.RoundRepo.DeleteFieldSessionAndItsStudents(RoundID)
	if newErra != nil {
		return false
	}
	newErra = roundser.RoundRepo.DeleteLecturesAndItsStudents(RoundID)
	if newErra != nil {
		return false
	}
	newErra = roundser.RoundRepo.DeleteRoundRelatedRoomsDates(RoundID)
	if newErra != nil {
		return false
	}
	return true
}

// GetRoundByIDForSchedule (RoundID uint) *entity.Round
func (roundser *RoundService) GetRoundByIDForSchedule(RoundID uint) *entity.Round {
	round, newErro := roundser.RoundRepo.GetRoundByIDForSchedule(RoundID)
	if newErro != nil {
		return nil
	}
	return round
}

// UpdateToRegistration  (RoundID uint) bool
func (roundser *RoundService) UpdateToRegistration(RoundID uint) bool {
	era := roundser.RoundRepo.UpdateToRegistration(RoundID)
	if era != nil {
		return false
	}
	return true
}

// GetActiveRoundsOfCategory   ( CategoryID uint) *[]entity.Round
func (roundser *RoundService) GetActiveRoundsOfCategory(CategoryID uint) *[]entity.Round {
	rounds, era := roundser.RoundRepo.GetActiveRoundsOfCategory(CategoryID)
	if era != nil {
		return nil
	}
	return rounds
}

// GetRoundsOfCategory   ( CategoryID uint) *[]entity.Round
func (roundser *RoundService) GetRoundsOfCategory(CategoryID uint) *[]entity.Round {
	rounds, era := roundser.RoundRepo.GetRoundsOfCategory(CategoryID)
	if era != nil {
		return nil
	}
	return rounds
}

// IsRoundNumberReseerved (BranchID, CategoryID, RoundNumber uint) bool
func (roundser *RoundService) IsRoundNumberReseerved(BranchID, CategoryID, RoundNumber uint) bool {
	count := roundser.RoundRepo.IsRoundNumberReseerved(BranchID, CategoryID, RoundNumber)
	if count >= 1 {
		return true
	}
	return false
}

// UpdateStudentsCount  (RoundID, Studentscount uint) error
func (roundser *RoundService) UpdateStudentsCount(RoundID, Studentscount uint) bool {
	era := roundser.RoundRepo.UpdateStudentsCount(RoundID, Studentscount)
	if era != nil {
		return false
	}
	return true
}
