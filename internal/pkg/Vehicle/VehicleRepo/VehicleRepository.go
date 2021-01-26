package VehicleRepo

import (
	"errors"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// VehicleRepo struct
type VehicleRepo struct {
	DB *gorm.DB
}

// NewVehicleRepo function
func NewVehicleRepo(db *gorm.DB) Vehicle.VehicleRepo {
	return &VehicleRepo{
		DB: db,
	}
}

// SaveVehicle method
func (vehiclerepo *VehicleRepo) SaveVehicle(vehicle *entity.Vehicle) (*entity.Vehicle, error) {
	era := vehiclerepo.DB.Save(vehicle).Error
	return vehicle, era
}

// GetVehicles method returning the Categories of Specific Branch
func (vehiclerepo *VehicleRepo) GetVehicles(categoryID, branchID uint) ([]entity.Vehicle, error) {
	vehicles := []entity.Vehicle{}
	fetchError := vehiclerepo.DB.Where(&entity.Vehicle{CategoryID: categoryID, BranchNo: branchID}).Find(&vehicles).Error
	return vehicles, fetchError
}

// GetVehicle method
func (vehiclerepo *VehicleRepo) GetVehicle(ID uint) (*entity.Vehicle, error) {
	vehicle := &entity.Vehicle{}
	fetchError := vehiclerepo.DB.Where(&entity.Vehicle{}, ID).First(vehicle).Error
	vehiclerepo.DB.Model(vehicle).Related(&vehicle.Category, "Category")
	return vehicle, fetchError
}

// IsVehicleReserved  (VehicleID uint) error
func (vehiclerepo *VehicleRepo) IsVehicleReserved(VehicleID uint) error {
	vehicle := &entity.Vehicle{}
	newEra := vehiclerepo.DB.First(vehicle, "id=?", VehicleID).Error
	if vehicle.Reserved || newEra != nil {
		return errors.New("The Thing is Reserved ")
	}
	return nil
}

// DeleteVehicle method to delete a vehicle Using it's Vehicle ID only
func (vehiclerepo *VehicleRepo) DeleteVehicle(VehicleID uint) error {
	newEra := vehiclerepo.DB.Delete(&entity.Vehicle{}, "id=?", VehicleID).Error
	return newEra
}

// GetFreeVehiclesOfCategory (CategoryID uint) (*[]entity.Vehicle, error)
func (vehiclerepo *VehicleRepo) GetFreeVehiclesOfCategory(CategoryID uint) (*[]entity.Vehicle, error) {
	vehicles := &[]entity.Vehicle{}
	newEra := vehiclerepo.DB.Table("vehicles").Find(vehicles, "category_refer=? && reserved= ?", CategoryID, false).Error
	return vehicles, newEra
}
