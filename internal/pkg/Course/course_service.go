package Course

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type CourseService interface {
	CreateCourse(course *entity.Course) *entity.Course
	DeleteCourse(course *entity.Course) bool
	EditCourse(course *entity.Course) *entity.Course
	GetCourseByID(ID uint) *entity.Course
	SaveCourse(course *entity.Course) *entity.Course
	GetCourseofBranch(BranchID uint) *[]entity.Course
	GetCourseofBranchAndCategory(BranchID, CategoryID uint) *[]entity.Course
}
