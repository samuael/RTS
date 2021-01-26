package VehicleService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// VehicleService struct
type VehicleService struct {
	VehicleRepo Vehicle.VehicleRepo
}

// NewVehicleService function
func NewVehicleService(repo Vehicle.VehicleRepo) Vehicle.VehicleService {
	return &VehicleService{
		VehicleRepo: repo,
	}
}

// SaveVehicle method
func (vehicleser *VehicleService) SaveVehicle(vehicle *entity.Vehicle) *entity.Vehicle {
	vehicle, errors := vehicleser.VehicleRepo.SaveVehicle(vehicle)
	if errors != nil {
		return nil
	}
	return vehicle
}

// GetVehicles method
func (vehicleser *VehicleService) GetVehicles(CategoryID, branchID uint) []entity.Vehicle {
	vehicles, fetchError := vehicleser.VehicleRepo.GetVehicles(CategoryID, branchID)
	if fetchError != nil {
		return nil
	}
	return vehicles
}

// GetVehicleByID method
func (vehicleser *VehicleService) GetVehicleByID(ID uint) *entity.Vehicle {
	vehicle, theError := vehicleser.VehicleRepo.GetVehicle(ID)
	if theError != nil {
		return nil
	}
	return vehicle
}

// IsVehicleReserved  (VehicleID uint) error
func (vehicleser *VehicleService) IsVehicleReserved(VehicleID uint) bool {
	newEera := vehicleser.VehicleRepo.IsVehicleReserved(VehicleID)
	if newEera != nil {
		return false
	}
	return true
}

// DeleteVehicle  (VehicleID uint) error
func (vehicleser *VehicleService) DeleteVehicle(VehicleID uint) bool {
	newEera := vehicleser.VehicleRepo.DeleteVehicle(VehicleID)
	if newEera != nil {
		return false
	}
	return true
}

// GetFreeVehiclesOfCategory  (CategoryID uint) *[]entity.Vehicle
func (vehicleser *VehicleService) GetFreeVehiclesOfCategory(CategoryID uint) *[]entity.Vehicle {
	vehicles, era := vehicleser.VehicleRepo.GetFreeVehiclesOfCategory(CategoryID)
	if era != nil {
		return nil
	}
	return vehicles
}
