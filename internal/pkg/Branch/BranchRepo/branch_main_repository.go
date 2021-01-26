package BranchRepo

import (
	"log"

	"github.com/lib/pq"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// BranchRepo struct
type BranchRepo struct {
	DB *gorm.DB
}

// NewBranchRepo function
func NewBranchRepo(db *gorm.DB) Branch.BranchRepo {
	return &BranchRepo{
		DB: db,
	}
}

// GetDefaultBranch method
func (brepo *BranchRepo) GetDefaultBranch(branch *entity.Branch) error {
	erro := brepo.DB.First(branch, 1).Error
	if erro != nil {
		return erro
	}
	address := &entity.Address{}
	erro = brepo.DB.Model(branch).Related(address, "AddressRefer").Error
	if erro != nil {
		log.Println(erro.Error())
	}
	branch.Address = *address
	return erro
}

// GetBranchs method
func (brepo *BranchRepo) GetBranchs(branch *[]entity.Branch) error {
	erro := brepo.DB.Find(branch).Error
	// Find their Address for each Branch
	for index, branc := range *branch {
		address := &entity.Address{}
		brepo.DB.Model(branc).Related(address, "AddressRefer")
		((*branch)[index]).Address = *address
	}
	return erro
}

// GetBranchByID returning the Branch taking the banch id as a parameter
func (brepo *BranchRepo) GetBranchByID(id uint) (*entity.Branch, error) {
	branc := &entity.Branch{}
	erro := brepo.DB.Table("branches").Where("id=?", id).First(branc).Error
	if erro != nil {
		return nil, erro
	}
	brepo.DB.Model(branc).Related(&branc.Address, "AddressRefer")
	brepo.DB.Model(branc).Related(&branc.LicenceGivenDate, "LicenceGivenDateRefer")
	return branc, erro
}

// CreateBranch (branch *entity.Branch) *entity.Branch
func (brepo *BranchRepo) CreateBranch(branch *entity.Branch) (*entity.Branch, error) {
	newEra := brepo.DB.Save(branch).Error
	return branch, newEra
}

// ChangeEmail (BranchID int, Email string) error
func (brepo *BranchRepo) ChangeEmail(BranchID int, Email string) error {
	era := brepo.DB.Table("branches").Where("id=?", BranchID).Updates(map[string]string{"email": Email}).Error
	return era
}

// ChangePhones  (Phones []string ) bool
func (brepo *BranchRepo) ChangePhones(BranchID int, Phones []string) error {
	era := brepo.DB.Table("branches").Where("id=?", BranchID).Updates(map[string]interface{}{"phone_numbers": pq.StringArray(Phones)}).Error
	return era
}

// DeleteDate for saving the memory
func (brepo *BranchRepo) DeleteDate(dateRefer int64) error {
	era := brepo.DB.Table("dates").Where("id=?", dateRefer).Delete(nil, "id=?", dateRefer).Error
	return era
}
