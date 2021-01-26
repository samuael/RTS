package InformationService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Information"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// InfoService struct
type InfoService struct {
	InfoRepo Information.InformationRepo
}

// NewInformationService function
func NewInformationService(inforepo Information.InformationRepo) Information.InformationService {
	return &InfoService{
		InfoRepo: inforepo,
	}
}

// CreateInfo (info *entity.Information) *entity.Information
func (infoser *InfoService) CreateInfo(info *entity.Information) *entity.Information {
	info, erra := infoser.InfoRepo.CreateInfo(info)
	if erra != nil {
		return nil
	}
	return info
}

// SaveInfo (info *entity.Information) *entity.Information
func (infoser *InfoService) SaveInfo(info *entity.Information) *entity.Information {
	info, era := infoser.InfoRepo.SaveInfo(info)
	if era != nil {
		return info
	}
	return info
}

// DeleteInfo (info *entity.Information) *entity.Information
func (infoser *InfoService) DeleteInfo(info *entity.Information) bool {
	era := infoser.InfoRepo.DeleteInfo(info)
	if era != nil {
		return false
	}
	return true
}

// GetActiveInfos (BranchID uint) *entity.Information
func (infoser *InfoService) GetActiveInfos(BranchID uint) *[]entity.Information {
	infos, era := infoser.InfoRepo.GetActiveInfos(BranchID)
	if era != nil {
		return nil
	}
	return infos
}

// GetAllInfos (BranchID uint) *entity.Information
func (infoser *InfoService) GetAllInfos(BranchID uint) *[]entity.Information {
	infos, era := infoser.InfoRepo.GetAllInfos(BranchID)
	if era != nil {
		return nil
	}
	return infos
}

// GetInfoByID (InfoID uint) *entity.Information
func (infoser *InfoService) GetInfoByID(InfoID uint) *entity.Information {
	info, erra := infoser.InfoRepo.GetInfoByID(InfoID)
	if erra != nil {
		return nil
	}
	return info
}

// GetUsersInfo (BranchID uint, Username string) *entity.Information
func (infoser *InfoService) GetUsersInfo(BranchID uint, Username string) *[]entity.Information {
	infos, era := infoser.InfoRepo.GetUsersInfo(BranchID, Username)
	if era != nil {
		return nil
	}
	return infos
}

// ActivateInformation method to activate an Information
func (infoser *InfoService) ActivateInformation(InfoID uint) *entity.Information {
	info, erra := infoser.InfoRepo.ActivateInformation(InfoID)
	if erra != nil {
		return nil
	}
	return info
}

// DeactivateInformation method to activate an Information
func (infoser *InfoService) DeactivateInformation(InfoID uint) *entity.Information {
	info, erra := infoser.InfoRepo.DeactivateInformation(InfoID)
	if erra != nil {
		return nil
	}
	return info
}

// GetAllActiveInfos method
func (infoser *InfoService) GetAllActiveInfos() *[]entity.Information {
	infos, errra := infoser.InfoRepo.GetAllActiveInfos()
	if errra != nil {
		return nil
	}
	return infos
}
