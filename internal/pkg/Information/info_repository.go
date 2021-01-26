package Information

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// InformationRepo interface
type InformationRepo interface {
	CreateInfo(info *entity.Information) (*entity.Information, error)
	SaveInfo(info *entity.Information) (*entity.Information, error)
	DeleteInfo(info *entity.Information) error
	GetActiveInfos(BranchID uint) (*[]entity.Information, error)
	GetAllInfos(BranchID uint) (*[]entity.Information, error)
	GetInfoByID(InfoID uint) (*entity.Information, error)
	GetUsersInfo(BranchID uint, Username string) (*[]entity.Information, error)
	ActivateInformation(InfoID uint) (*entity.Information, error)
	DeactivateInformation(InfoID uint) (*entity.Information, error)
	GetAllActiveInfos() (*[]entity.Information, error)
}
