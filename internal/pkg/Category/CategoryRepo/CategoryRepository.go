package CategoryRepo

import (
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// CategoryRepo struct
type CategoryRepo struct {
	DB *gorm.DB
}

// NewCategoryRepo function
func NewCategoryRepo(db *gorm.DB) Category.CategoryRepo {
	return &CategoryRepo{
		DB: db,
	}
}

/*
	**********Methods************
	GetCategories(branchId uint )
	GetCategoryByID(ID uint) (*entity.Category, error)
	CreateCategory(category *entity.Category) (*entity.Category, error)

*/

// GetCategories method returning the Categories of Specific Branch
func (categoryrep *CategoryRepo) GetCategories(branchID uint) ([]entity.Category, error) {
	categories := []entity.Category{}
	fetchError := categoryrep.DB.Where(&entity.Category{Branchid: branchID}).Find(&categories).Error
	return categories, fetchError
}

// GetCategoryByID method returning the Categories of Specific Branch
func (categoryrep *CategoryRepo) GetCategoryByID(ID uint) (*entity.Category, error) {
	category := &entity.Category{}
	// category.ID = ID
	fetchError := categoryrep.DB.First(category, "id=?", ID).Error
	return category, fetchError
}

// CreateCategory method
func (categoryrep *CategoryRepo) CreateCategory(category *entity.Category) (*entity.Category, error) {
	newError := categoryrep.DB.Save(category).Error
	return category, newError
}

// DeleteCategory method
func (categoryrep *CategoryRepo) DeleteCategory(category *entity.Category) (*entity.Category, error) {
	erra := categoryrep.DB.Delete(category).Error
	return category, erra
}

// DeleteCategoryByID method for Deleting a category using the Category ID
func (categoryrep *CategoryRepo) DeleteCategoryByID(ID uint) error {
	category := &entity.Category{}
	category.ID = ID
	fmt.Println("Deleting Category Id  ", ID)
	newErra := categoryrep.DB.Where(&entity.Category{}, ID).Delete(category).Error
	// newErra = categoryrep.DB.Delete(category, "ID=?", ID).Error
	return newErra
}

// SaveCategory method to Update the Category
func (categoryrep *CategoryRepo) SaveCategory(category *entity.Category) (*entity.Category, error) {
	// sub := categoryrep.DB.Table("categories").Update(category).Debug()
	// fmt.Println(sub)
	errros := categoryrep.DB.Save(category).Debug().Error
	// fmt.Println(errros)
	return category, errros
}

// Adding Fieldman , Teacher  , Course to a round are Coded below

// AddFieldManToARound method
// func (roundhandler *)

// UpdateCategoryLasrRoundCount method
func (categoryrep *CategoryRepo) UpdateCategoryLasrRoundCount(category *entity.Category, CategoryID, LastRoundNumber uint) error {
	newErra := categoryrep.DB.Model(category).Save(category).Error // Updates(map[string]interface{}{"last_round_number": int(LastRoundNumber)}).Error
	fmt.Println(newErra)
	return newErra
}
