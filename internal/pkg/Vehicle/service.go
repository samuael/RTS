package Vehicle

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// VehicleService interface
type VehicleService interface {
	GetVehicles(CategoryID, branchID uint) []entity.Vehicle
	GetVehicleByID(ID uint) *entity.Vehicle
	SaveVehicle(vehicle *entity.Vehicle) *entity.Vehicle
	IsVehicleReserved(VehicleID uint) bool
	DeleteVehicle(VehicleID uint) bool
	GetFreeVehiclesOfCategory(CategoryID uint) *[]entity.Vehicle
}
