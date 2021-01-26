package Information

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// InformationService interface
type InformationService interface {
	CreateInfo(info *entity.Information) *entity.Information
	SaveInfo(info *entity.Information) *entity.Information
	DeleteInfo(info *entity.Information) bool
	GetActiveInfos(BranchID uint) *[]entity.Information
	GetAllActiveInfos() *[]entity.Information
	GetAllInfos(BranchID uint) *[]entity.Information
	GetInfoByID(InfoID uint) *entity.Information
	GetUsersInfo(BranchID uint, Username string) *[]entity.Information // Not User Information
	ActivateInformation(InfoID uint) *entity.Information
	DeactivateInformation(InfoID uint) *entity.Information
}
