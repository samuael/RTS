package Branch

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// BranchService interface
type BranchService interface {
	GetDefaultBranch() *entity.Branch
	GetBranchs() *[]entity.Branch
	GetBranchByID(id uint) *entity.Branch
	CreateBranch(branch *entity.Branch) *entity.Branch
	DeleteBranch(BranchID int) bool
	ChangeEmail(BranchID int, Email string) bool
	ChangePhones(BranchID int, Phones []string) bool
	DeleteDate(dateRefer int64) bool
}
