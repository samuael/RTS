package TrainerRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// TrainerRepo struct
type TrainerRepo struct {
	DB *gorm.DB
}

// NewTrainerRepo function
func NewTrainerRepo(db *gorm.DB) Trainer.TrainerRepo {
	return &TrainerRepo{
		DB: db,
	}
}

// SaveTrainer (trainer *entity.FieldAssistant) error
func (tr *TrainerRepo) SaveTrainer(trainer *entity.FieldAssistant) error {
	erro := tr.DB.Save(trainer).Error
	return erro
}

// GetCount method returning the nimber of trainers in the database
func (tr *TrainerRepo) GetCount() uint {
	var count uint
	tr.DB.Model(&entity.FieldAssistant{}).Count(&count)
	return count
}

// GetTrainerByID method
func (tr *TrainerRepo) GetTrainerByID(ID uint) (*entity.FieldAssistant, error) {
	trainer := &entity.FieldAssistant{}
	errors := tr.DB.Where(&entity.FieldAssistant{}).First(trainer, ID).Error
	return trainer, errors
}

// LogTrainer method
func (tr *TrainerRepo) LogTrainer(username, password string) (*entity.FieldAssistant, error) {
	trainer := &entity.FieldAssistant{}
	trainer.Username = username
	trainer.Password = password
	newErra := tr.DB.First(trainer, "username=? and password=?", username, password).Error
	return trainer, newErra
}

// TrainersOfCategoryID method returns Trainers Pointer having CategoryID of the parameter
func (tr *TrainerRepo) TrainersOfCategoryID(ID uint) (*[]entity.FieldAssistant, error) {
	trainers := &[]entity.FieldAssistant{}
	newError := tr.DB.Where("category_id=?", ID).Find(trainers).Error
	for l := 0; l < len(*trainers); l++ {
		trainer := &(*trainers)[l]
		newError = tr.DB.Model(trainer).Related(&trainer.Vehicle, "Vehicle").Error
	}
	return trainers, newError
}

// GetTrainersOfCategory method
func (tr *TrainerRepo) GetTrainersOfCategory(CategoryID, Offset, Limit uint) (*[]entity.FieldAssistant, error) {
	trainers := &[]entity.FieldAssistant{}
	eras := tr.DB.Table("field_assistants").Limit(Limit).Offset(Offset).Find(trainers, "category_refer=?", CategoryID).Error
	for j := 0; j < len(*trainers); j++ {
		trainer := &(*trainers)[j]
		tr.DB.Model(trainer).Related(&trainer.Vehicle, "Vehicle")
	}
	return trainers, eras
}

// GetFreeTrainers (CategoryID uint) (*[]entity.FieldAssistant, error)
func (tr *TrainerRepo) GetFreeTrainers(CategoryID uint) (*[]entity.FieldAssistant, error) {
	trainers := &[]entity.FieldAssistant{}
	newEra := tr.DB.Table("field_assistants").Find(trainers, "category_refer=? and vahicle_id >=? ", CategoryID, 1).Error
	for j := 0; j < len(*trainers); j++ {
		trainer := &(*trainers)[j]
		tr.DB.Model(trainer).Related(&trainer.Vehicle, "Vehicle")
	}
	return trainers, newEra
}

// ChangeImageURL method
func (tr *TrainerRepo) ChangeImageURL(ID uint, ImageURL string) error {
	erra := tr.DB.Table("field_sessions").Where("id=?", ID).Updates(map[string]string{"imageurl": ImageURL}).Error
	return erra
}
