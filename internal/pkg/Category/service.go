package Category

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type CategoryService interface {
	GetCategories(branchID uint) []entity.Category
	GetCategoryByID(ID uint) *entity.Category
	CreateCategory(category *entity.Category) *entity.Category
	DeleteCategory(category *entity.Category) *entity.Category
	DeleteCategoryByID(ID uint) bool
	SaveCategory(category *entity.Category) *entity.Category
	UpdateCategoryLasrRoundCount(category *entity.Category, CategoryID, LastRoundNumber uint) bool
}
