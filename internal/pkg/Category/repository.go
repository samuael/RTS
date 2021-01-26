package Category

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type CategoryRepo interface {
	GetCategories(branchID uint) ([]entity.Category, error)
	GetCategoryByID(ID uint) (*entity.Category, error)
	CreateCategory(category *entity.Category) (*entity.Category, error)
	DeleteCategory(category *entity.Category) (*entity.Category, error)
	DeleteCategoryByID(ID uint) error
	SaveCategory(category *entity.Category) (*entity.Category, error)
	UpdateCategoryLasrRoundCount(category *entity.Category, CategoryID, LastRoundNumber uint) error
}
