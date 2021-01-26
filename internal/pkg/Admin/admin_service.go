package Admin

import (
	entity "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// AdminService interface representing DTIS Admin Service
type AdminService interface {
	RegisterAdmin(admin *entity.Admin) *entity.Admin
	GetAdmin(username, password string) *entity.Admin
	AdminsCount() uint
	UpdateAdmin(admin *entity.Admin) *entity.Admin
	GetAdminByID(ID uint, username string) *entity.Admin
	ChangeImageURL(ID uint, val string) bool
	GetAdminsOfSystem(BranchID int, Active int) *[]entity.Admin
	DeleteAdmin(AdminID uint) bool
	DeactivateAdmin(AdminID uint) bool
	ActivateAdmin(AdminID uint) bool
}
