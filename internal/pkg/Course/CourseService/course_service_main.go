package CourseService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Course"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// GormCourseService struct representing
type GormCourseService struct {
	CourseRepo Course.CourseRepo
}

// NewGormCourseService function
func NewGormCourseService(courserepo Course.CourseRepo) Course.CourseService {
	return &GormCourseService{
		CourseRepo: courserepo,
	}
}

// CreateCourse method representing to create a course
func (courseser *GormCourseService) CreateCourse(course *entity.Course) *entity.Course {
	course, newEroro := courseser.CourseRepo.CreateCourse(course)
	if newEroro != nil {
		return nil
	}
	return course
}

// DeleteCourse method to delete a course
func (courseser *GormCourseService) DeleteCourse(course *entity.Course) bool {
	newEroro := courseser.CourseRepo.DeleteCourse(course)
	if newEroro != nil {
		return false
	}
	return true
}

// EditCourse method for editing th ecourse
func (courseser *GormCourseService) EditCourse(course *entity.Course) *entity.Course {
	course, newEroro := courseser.CourseRepo.EditCourse(course)
	if newEroro != nil {
		return nil
	}
	return course
}

// GetCourseByID method
func (courseser *GormCourseService) GetCourseByID(ID uint) *entity.Course {
	course, era := courseser.CourseRepo.GetCourseByID(ID)
	if era != nil {
		return nil
	}
	return course
}

// SaveCourse (course *entity.Course) (*entity.Course, error)
func (courseser *GormCourseService) SaveCourse(course *entity.Course) *entity.Course {
	course, er := courseser.CourseRepo.SaveCourse(course)
	if er != nil {
		return nil
	}
	return course

}

// GetCourseofBranchAndCategory

// GetCourseofBranch method
func (courseser *GormCourseService) GetCourseofBranch(BranchID uint) *[]entity.Course {
	courses, erra := courseser.CourseRepo.GetCourseofBranch(BranchID)
	if erra != nil {
		return nil
	}
	return courses
}

// GetCourseofBranchAndCategory method
func (courseser *GormCourseService) GetCourseofBranchAndCategory(BranchID, CategoryID uint) *[]entity.Course {
	courses, erra := courseser.CourseRepo.GetCourseofBranchAndCategory(BranchID, CategoryID)
	if erra != nil {
		return nil
	}
	return courses
}
