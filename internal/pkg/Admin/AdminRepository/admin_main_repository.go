package AdminRepository

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Admin"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// AdminRepo representing Admin Related Repository
type AdminRepo struct {
	db *gorm.DB
}

// NewAdminRepo y
func NewAdminRepo(dbs *gorm.DB) Admin.AdminRepo {
	return &AdminRepo{db: dbs}
}

// RegisterStudent yea
func (adminr *AdminRepo) RegisterStudent(student *entity.Student) error {
	eror := adminr.db.Save(student).Error
	return eror
}

// ReigisterAdmin method
func (adminr *AdminRepo) ReigisterAdmin(admin *entity.Admin) error {
	dberror := adminr.db.Save(admin).Error
	return dberror
}

// GetCount method
func (adminr *AdminRepo) GetCount() (count uint) {
	adminr.db.Model(&entity.Admin{}).Count(&count)
	return count
}

// GetAdmin  method
func (adminr *AdminRepo) GetAdmin(admin *entity.Admin) error {
	erro := adminr.db.Where("username=? and password=?", admin.Username, admin.Password).First(admin).Error
	branch := entity.Branch{}
	adminr.db.Model(admin).Related(&branch, "BranchRefer")
	admin.Branch = branch
	address := &entity.Address{}
	adminr.db.Model(&branch).Related(address, "AddressRefer")
	branch.Address = *address
	return erro
}

// UpdateAdmin method
func (adminr *AdminRepo) UpdateAdmin(admin *entity.Admin) error {
	erro := adminr.db.Save(admin).Error
	return erro
}

// GetAdminByID method
func (adminr *AdminRepo) GetAdminByID(ID uint) (*entity.Admin, error) {
	admin := &entity.Admin{}
	err := adminr.db.Where("id=?", ID).First(admin).Error
	return admin, err
}

// ChangeImageURL (val string) error  this is the repository
func (adminr *AdminRepo) ChangeImageURL(ID uint, ImageURL string) error {
	era := adminr.db.Table("admins").Where("id=?", ID).Updates(map[string]string{"imageurl": ImageURL}).Error
	return era
}

// GetAdminsOfSystem (BranchID uint, Active byte) (*[]entity.Admin, error)
func (adminr *AdminRepo) GetAdminsOfSystem(BranchID int, Active int) (*[]entity.Admin, error) {
	admins := &[]entity.Admin{}
	var era error
	if BranchID >= 0 {
		switch Active {
		case entity.All:
			{
				era = adminr.db.Table("admins").Find(admins, "branch_refer=? and role =?", BranchID, entity.SUPERADMIN).Error
				break
			}
		case entity.Active:
			{
				era = adminr.db.Table("admins").Find(admins, "branch_refer=? and active=? and role =?", BranchID, true, entity.SUPERADMIN).Error
				break
			}
		case entity.Passive:
			{
				era = adminr.db.Table("admins").Find(admins, "branch_refer=? and active=? and role =?", BranchID, false, entity.SUPERADMIN).Error
				break
			}
		}
	} else {
		switch Active {
		case entity.All:
			{
				era = adminr.db.Table("admins").Find(admins, "role =?", entity.SUPERADMIN).Error
				break
			}
		case entity.Active:
			{
				era = adminr.db.Table("admins").Find(admins, "active=? and role =?", BranchID, true, entity.SUPERADMIN).Error
				break
			}
		case entity.Passive:
			{
				era = adminr.db.Table("admins").Find(admins, "active=?  and role =? ", BranchID, false, entity.SUPERADMIN).Error
				break
			}
		}
	}
	return admins, era
}

// DeleteAdmin method to delete an admin
// the imagefile of the admin should be deleted from the File Directory While deleting the admin
func (adminr *AdminRepo) DeleteAdmin(AdminID uint) error {
	newEra := adminr.db.Table("admins").Where("id=?", AdminID).Delete(nil).Error
	return newEra
}

// DeactivateAdmin method to make the admin passive
func (adminr *AdminRepo) DeactivateAdmin(AdminID uint) error {
	era := adminr.db.Table("admins").Where("id=?", AdminID).Updates(map[string]interface{}{"active": false}).Error
	return era
}

// ActivateAdmin method to make the admin active
func (adminr *AdminRepo) ActivateAdmin(AdminID uint) error {
	era := adminr.db.Table("admins").Where("id=?", AdminID).Updates(map[string]interface{}{"active": true}).Error
	return era
}
