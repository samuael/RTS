package TrainerService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// TrainerService struct
type TrainerService struct {
	TrainerRepo Trainer.TrainerRepo
}

// NewTrainerService function
func NewTrainerService(tr Trainer.TrainerRepo) Trainer.TrainerService {
	return &TrainerService{
		TrainerRepo: tr,
	}
}

// SaveTrainer (trainer *entity.FieldAssistant) *entity.FieldAssistant
func (trainser *TrainerService) SaveTrainer(trainer *entity.FieldAssistant) *entity.FieldAssistant {
	saveError := trainser.TrainerRepo.SaveTrainer(trainer)
	if saveError != nil {
		return nil
	}
	return trainer
}

// GetCount method returning the nimber of trainers in the database
func (trainser *TrainerService) GetCount() uint {
	count := trainser.TrainerRepo.GetCount()
	return count
}

// GetTrainerByID method
func (trainser *TrainerService) GetTrainerByID(ID uint) *entity.FieldAssistant {
	Trainer, newError := trainser.TrainerRepo.GetTrainerByID(ID)
	if newError != nil {
		return nil
	}
	return Trainer
}

// LogTrainer method
func (trainser *TrainerService) LogTrainer(username, password string) *entity.FieldAssistant {
	trainer, erra := trainser.TrainerRepo.LogTrainer(username, password)
	if erra != nil {
		return nil
	}
	return trainer
}

// TrainersOfCategoryID (ID uint) *[]entity.FieldAssistant
func (trainser *TrainerService) TrainersOfCategoryID(ID uint) *[]entity.FieldAssistant {
	trainers, newError := trainser.TrainerRepo.TrainersOfCategoryID(ID)
	if newError != nil {
		return nil
	}
	return trainers
}

// GetTrainersOfCategory (CategoryID uint) *[]entity.FieldAssistant
func (trainser TrainerService) GetTrainersOfCategory(CategoryID, Offset, Limit uint) *[]entity.FieldAssistant {
	trainers, era := trainser.TrainerRepo.GetTrainersOfCategory(CategoryID, Offset, Limit)
	if era != nil {
		return nil
	}
	return trainers
}

// GetFreeTrainers (CategoryID uint) (*[]entity.FieldAssistant, error)
func (trainser *TrainerService) GetFreeTrainers(CategoryID uint) *[]entity.FieldAssistant {
	trainers, era := trainser.TrainerRepo.GetFreeTrainers(CategoryID)
	if era != nil {
		return nil
	}
	return trainers
}

// ChangeImageURL method
func (trainser *TrainerService) ChangeImageURL(ID uint, ImageURL string) bool {
	erra := trainser.TrainerRepo.ChangeImageURL(ID, ImageURL)
	if erra != nil {
		return false
	}
	return true
}
