package BranchService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// BranchService Struct
type BranchService struct {
	BranchRepo Branch.BranchRepo
}

// NewBranchService function
func NewBranchService(repo Branch.BranchRepo) Branch.BranchService {
	return &BranchService{
		BranchRepo: repo,
	}
}

// GetDefaultBranch method
func (branchservice *BranchService) GetDefaultBranch() *entity.Branch {
	branch := &entity.Branch{}
	branchservice.BranchRepo.GetDefaultBranch(branch)
	return branch
}

// GetBranchs method
func (branchservice *BranchService) GetBranchs() *[]entity.Branch {
	branch := &[]entity.Branch{}
	erro := branchservice.BranchRepo.GetBranchs(branch)
	if erro != nil {
		return nil
	}
	return branch
}

// GetBranchByID method
func (branchservice *BranchService) GetBranchByID(id uint) *entity.Branch {
	branch, erro := branchservice.BranchRepo.GetBranchByID(id)
	if erro != nil {
		return nil
	}
	return branch
}

// CreateBranch  (branch *entity.Branch) *entity.Branch
func (branchservice *BranchService) CreateBranch(branch *entity.Branch) *entity.Branch {
	branch, era := branchservice.BranchRepo.CreateBranch(branch)
	if era != nil {
		return nil
	}
	return branch
}

// DeleteBranch  for deletign branch and related Resources
func (branchservice *BranchService) DeleteBranch(BranchID int) bool {
	// Deleting Dates related To This Branch
	// err := branchservice.BranchRepo.DeleteDatesOfBranch()
	return false
}

// ChangeEmail   service method
func (branchservice *BranchService) ChangeEmail(BranchID int, Email string) bool {
	era := branchservice.BranchRepo.ChangeEmail(BranchID, Email)
	if era != nil {
		return false
	}
	return true
}

// ChangePhones (Phones []string ) bool
func (branchservice *BranchService) ChangePhones(BranchID int, Phones []string) bool {
	era := branchservice.BranchRepo.ChangePhones(BranchID, Phones)
	if era != nil {
		return false
	}
	return true
}

// DeleteDate ( dateRefer int64) bool
func (branchservice *BranchService) DeleteDate(dateRefer int64) bool {
	era := branchservice.BranchRepo.DeleteDate(dateRefer)
	if era != nil {
		return false
	}
	return true
}
