package Trainer

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type TrainerRepo interface {
	SaveTrainer(trainer *entity.FieldAssistant) error
	GetCount() uint
	GetTrainerByID(ID uint) (*entity.FieldAssistant, error)
	LogTrainer(username, password string) (*entity.FieldAssistant, error)
	TrainersOfCategoryID(ID uint) (*[]entity.FieldAssistant, error)
	GetTrainersOfCategory(CategoryID, Offset, Limit uint) (*[]entity.FieldAssistant, error)
	GetFreeTrainers(CategoryID uint) (*[]entity.FieldAssistant, error)
	ChangeImageURL(ID uint, value string) error
}
