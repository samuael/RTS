package CourseRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Course"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// GormCourseRepo struct
type GormCourseRepo struct {
	DB *gorm.DB
}

// NewGormCourseRepo method returning new Instance of Gorm Repo for Course
func NewGormCourseRepo(db *gorm.DB) Course.CourseRepo {
	return &GormCourseRepo{
		DB: db,
	}
}

// CreateCourse method representing to create a course
func (courserepo *GormCourseRepo) CreateCourse(course *entity.Course) (*entity.Course, error) {
	newError := courserepo.DB.Create(course).Error
	return course, newError
}

// DeleteCourse method to delete a course
func (courserepo *GormCourseRepo) DeleteCourse(course *entity.Course) error {
	newError := courserepo.DB.Delete(course).Error
	return newError
}

// EditCourse method for editing th ecourse
func (courserepo *GormCourseRepo) EditCourse(course *entity.Course) (*entity.Course, error) {
	newError := courserepo.DB.Update(course).Error
	return course, newError
}

// GetCourseByID method
func (courserepo *GormCourseRepo) GetCourseByID(ID uint) (*entity.Course, error) {
	course := &entity.Course{}
	newErra := courserepo.DB.Where("id=?", ID).First(course).Error
	return course, newErra
}

// SaveCourse (course *entity.Course) (*entity.Course, error)
func (courserepo *GormCourseRepo) SaveCourse(course *entity.Course) (*entity.Course, error) {
	newEra := courserepo.DB.Save(course).Error
	return course, newEra
}

// GetCourseofBranchAndCategory

// GetCourseofBranch  (BranchID uint) *[]entity.Course
func (courserepo *GormCourseRepo) GetCourseofBranch(BranchID uint) (*[]entity.Course, error) {
	courses := &[]entity.Course{}
	newErra := courserepo.DB.Where("branch_id=?", BranchID).Find(courses).Error
	return courses, newErra
}

// GetCourseofBranchAndCategory  (BranchID uint) *[]entity.Course
func (courserepo *GormCourseRepo) GetCourseofBranchAndCategory(BranchID, CategoryID uint) (*[]entity.Course, error) {
	courses := &[]entity.Course{}
	newErra := courserepo.DB.Where("branch_id=? and  categoryid=? ", BranchID, CategoryID).Find(courses).Error
	return courses, newErra
}
