package InformationRepo

import (
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Information"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// InfoRepo struct
type InfoRepo struct {
	DB *gorm.DB
}

// NewInfoRepo function
func NewInfoRepo(db *gorm.DB) Information.InformationRepo {
	return &InfoRepo{
		DB: db,
	}
}
// CreateInfo (info *entity.Information) (*entity.Information, error)
func (inforepo *InfoRepo) CreateInfo(info *entity.Information) (*entity.Information, error) {
	newErro := inforepo.DB.Create(info).Error
	return info, newErro
}

// SaveInfo (info *entity.Information) (*entity.Information, error)
func (inforepo *InfoRepo) SaveInfo(info *entity.Information) (*entity.Information, error) {
	newError := inforepo.DB.Save(info).Error
	return info, newError
}

// DeleteInfo (info *entity.Information) (*entity.Information, error)
func (inforepo *InfoRepo) DeleteInfo(info *entity.Information) error {
	erra := inforepo.DB.Delete(info).Error
	return erra
}

// GetActiveInfos (BranchID uint) (*entity.Information, error)
func (inforepo *InfoRepo) GetActiveInfos(BranchID uint) (*[]entity.Information, error) {
	infos := &[]entity.Information{}
	erra := inforepo.DB.Where("active=? AND branch_id=?", true, BranchID).Find(infos).Error
	return infos, erra
}

// GetAllInfos (BranchID uint) (*entity.Information, error)
func (inforepo *InfoRepo) GetAllInfos(BranchID uint) (*[]entity.Information, error) {
	infos := &[]entity.Information{}
	erras := inforepo.DB.Where("branch_id=?", BranchID).Find(infos).Error
	return infos, erras
}

// GetInfoByID (InfoID uint) (*entity.Information, error)
func (inforepo *InfoRepo) GetInfoByID(InfoID uint) (*entity.Information, error) {
	info := &entity.Information{}
	erra := inforepo.DB.Where("id=?", InfoID).First(info).Error
	return info, erra
}

// GetUsersInfo (BranchID uint, Username string) (*entity.Information, error)
func (inforepo *InfoRepo) GetUsersInfo(BranchID uint, Username string) (*[]entity.Information, error) {
	info := &[]entity.Information{}
	erra := inforepo.DB.Where("branch_id=? and username=?", BranchID, Username).Find(info).Error
	return info, erra
}

// ActivateInformation (InfoID uint) (*entity.Information, error)
func (inforepo *InfoRepo) ActivateInformation(InfoID uint) (*entity.Information, error) {
	info := &entity.Information{}
	erra := inforepo.DB.Where("id=?", InfoID).Find(info).Error
	info.Active = true
	erra = inforepo.DB.Where("id=?", InfoID).Save(info).Error
	return info, erra
}

// DeactivateInformation (InfoID uint) (*entity.Information, error)
func (inforepo *InfoRepo) DeactivateInformation(InfoID uint) (*entity.Information, error) {
	info := &entity.Information{}
	fmt.Println(InfoID)
	erra := inforepo.DB.Where("id=?", InfoID).Find(info).Error
	info.Active = false
	erra = inforepo.DB.Where("id=?", InfoID).Save(info).Error
	return info, erra
}

// GetAllActiveInfos method
func (inforepo *InfoRepo) GetAllActiveInfos() (*[]entity.Information, error) {
	infos := &[]entity.Information{}
	erras := inforepo.DB.Where("active=?", true).Find(infos).Error
	return infos, erras
}
