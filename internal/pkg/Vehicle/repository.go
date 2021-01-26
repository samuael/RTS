package Vehicle

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type VehicleRepo interface {
	GetVehicles(categoryID, branchID uint) ([]entity.Vehicle, error)
	GetVehicle(ID uint) (*entity.Vehicle, error)
	SaveVehicle(vehicle *entity.Vehicle) (*entity.Vehicle, error)
	IsVehicleReserved(VehicleID uint) error
	DeleteVehicle(VehicleID uint) error
	GetFreeVehiclesOfCategory(CategoryID uint) (*[]entity.Vehicle, error)
}
