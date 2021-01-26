package Branch

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// BranchRepo representing repository interface
type BranchRepo interface {
	GetDefaultBranch(branch *entity.Branch) error
	GetBranchs(branch *[]entity.Branch) error
	GetBranchByID(id uint) (*entity.Branch, error)
	CreateBranch(branch *entity.Branch) (*entity.Branch, error)
	// DeleteBranch(BranchID uint) error
	ChangeEmail(BranchID int, Email string) error
	ChangePhones(BranchID int, Phones []string) error
	DeleteDate(dateRefer int64) error
}
