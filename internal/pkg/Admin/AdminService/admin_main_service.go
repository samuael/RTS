//Package AdminService for serving admins related Services
package AdminService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// AdminService struct
type AdminService struct {
	Adminrepo Admin.AdminRepo
}

//NewAdminService  rrr
func NewAdminService(adminrepo Admin.AdminRepo) Admin.AdminService {
	return &AdminService{Adminrepo: adminrepo}
}

// AdminsCount method
func (adminser *AdminService) AdminsCount() uint {
	count := adminser.Adminrepo.GetCount()
	return count
}

// RegisterAdmin method
func (adminser *AdminService) RegisterAdmin(admin *entity.Admin) *entity.Admin {
	registrationError := adminser.Adminrepo.ReigisterAdmin(admin)
	if registrationError != nil {
		return nil
	}
	return admin
}

// GetAdmin method
func (adminser *AdminService) GetAdmin(username, password string) *entity.Admin {
	admin := &entity.Admin{}
	admin.Username = username
	admin.Password = password

	erro := adminser.Adminrepo.GetAdmin(admin)
	if erro != nil {
		return nil
	}
	return admin
}

// GetAdminByID you  can leav the Username entry "" if you want but the ID is meandatory
// returns an admin pointer populated
func (adminser *AdminService) GetAdminByID(ID uint, username string) *entity.Admin {
	admin := &entity.Admin{}
	admin.ID = ID
	if username != "" {
		admin.Username = username
	}
	admin, erro := adminser.Adminrepo.GetAdminByID(admin.ID)
	if erro != nil {
		return nil
	}
	return admin
}

// UpdateAdmin method
func (adminser *AdminService) UpdateAdmin(admin *entity.Admin) *entity.Admin {
	newErr := adminser.Adminrepo.UpdateAdmin(admin)
	if newErr != nil {
		return nil
	}
	return admin
}

// ChangeImageURL (val string) error  method
func (adminser *AdminService) ChangeImageURL(ID uint, ImageURL string) bool {
	newEra := adminser.Adminrepo.ChangeImageURL(ID, ImageURL)
	if newEra != nil {
		return false
	}
	return true
}

// GetAdminsOfSystem method
func (adminser *AdminService) GetAdminsOfSystem(BranchID int, Active int) *[]entity.Admin {
	admins, erro := adminser.Adminrepo.GetAdminsOfSystem(BranchID, Active)
	if erro != nil {
		return nil
	}
	return admins
}

// DeleteAdmin (AdminID uint) error
func (adminser *AdminService) DeleteAdmin(AdminID uint) bool {
	era := adminser.Adminrepo.DeleteAdmin(AdminID)
	if era != nil {
		return false
	}
	return true
}

// DeactivateAdmin ( AdminID uint) bool
func (adminser *AdminService) DeactivateAdmin(AdminID uint) bool {
	era := adminser.Adminrepo.DeactivateAdmin(AdminID)
	if era != nil {
		return false
	}
	return true
}

// ActivateAdmin (AdminID uint) bool
func (adminser *AdminService) ActivateAdmin(AdminID uint) bool {
	era := adminser.Adminrepo.ActivateAdmin(AdminID)
	if era != nil {
		return false
	}
	return true
}
