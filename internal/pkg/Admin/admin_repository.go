package Admin

import entity "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type AdminRepo interface {
	RegisterStudent(student *entity.Student) error
	GetAdmin(admin *entity.Admin) error
	GetCount() (count uint)
	ReigisterAdmin(admin *entity.Admin) error
	UpdateAdmin(admin *entity.Admin) error
	GetAdminByID(ID uint) (*entity.Admin, error)
	ChangeImageURL(ID uint, val string) error
	GetAdminsOfSystem(BranchID int, Active int) (*[]entity.Admin, error)
	DeleteAdmin(AdminID uint) error
	DeactivateAdmin(AdminID uint) error
	ActivateAdmin(AdminID uint) error
}
