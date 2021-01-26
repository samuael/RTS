package Course

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type CourseRepo interface {
	CreateCourse(course *entity.Course) (*entity.Course, error)
	DeleteCourse(course *entity.Course) error
	EditCourse(course *entity.Course) (*entity.Course, error)
	GetCourseByID(ID uint) (*entity.Course, error)
	SaveCourse(course *entity.Course) (*entity.Course, error)
	GetCourseofBranch(BranchID uint) (*[]entity.Course, error)
	GetCourseofBranchAndCategory(BranchID, CategoryID uint) (*[]entity.Course, error)
}
