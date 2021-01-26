package Trainer

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type TrainerService interface {
	SaveTrainer(trainer *entity.FieldAssistant) *entity.FieldAssistant
	GetCount() uint
	GetTrainerByID(ID uint) *entity.FieldAssistant
	LogTrainer(username, password string) *entity.FieldAssistant
	TrainersOfCategoryID(ID uint) *[]entity.FieldAssistant
	GetTrainersOfCategory(CategoryID, Offset, Limit uint) *[]entity.FieldAssistant
	GetFreeTrainers(CategoryID uint) *[]entity.FieldAssistant
	ChangeImageURL(ID uint, value string) bool
}
