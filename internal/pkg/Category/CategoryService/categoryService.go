package CategoryService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// CategoryService struct
type CategoryService struct {
	CategoryRepo Category.CategoryRepo
}

// NewCategoryService  function
func NewCategoryService(repo Category.CategoryRepo) Category.CategoryService {
	return &CategoryService{
		CategoryRepo: repo,
	}
}

// GetCategories method
func (cateser *CategoryService) GetCategories(branchID uint) []entity.Category {
	categories, fetchError := cateser.CategoryRepo.GetCategories(branchID)
	if fetchError != nil {
		return nil
	}
	return categories
}

// GetCategoryByID method
func (cateser *CategoryService) GetCategoryByID(ID uint) *entity.Category {
	category, filterError := cateser.CategoryRepo.GetCategoryByID(ID)
	if filterError != nil {
		return nil
	}
	return category
}

// CreateCategory method
func (cateser *CategoryService) CreateCategory(category *entity.Category) *entity.Category {
	category, newError := cateser.CategoryRepo.CreateCategory(category)
	if newError != nil {
		return nil
	}
	return category
}

// DeleteCategory (category *entity.Category) (*entity.Category, error)
func (cateser *CategoryService) DeleteCategory(category *entity.Category) *entity.Category {
	category, erra := cateser.CategoryRepo.DeleteCategory(category)
	if erra != nil {
		return nil
	}
	return category
}

// DeleteCategoryByID (ID uint) (*entity.Category)
func (cateser *CategoryService) DeleteCategoryByID(ID uint) bool {
	erra := cateser.CategoryRepo.DeleteCategoryByID(ID)
	if erra != nil {
		return false
	}
	return true
}

// SaveCategory method to Update the Category
func (cateser *CategoryService) SaveCategory(category *entity.Category) *entity.Category {
	category, errros := cateser.CategoryRepo.SaveCategory(category)
	if errros != nil {
		return nil
	}
	return category
}

// UpdateCategoryLasrRoundCount method
func (cateser *CategoryService) UpdateCategoryLasrRoundCount(category *entity.Category, CategoryID, LastRoundNumber uint) bool {
	newErra := cateser.CategoryRepo.UpdateCategoryLasrRoundCount(category, CategoryID, LastRoundNumber)
	if newErra != nil {
		return false
	}
	return true
}
