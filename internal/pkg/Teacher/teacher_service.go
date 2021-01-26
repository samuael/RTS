package Teacher

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// TeacherService interface holding the method of TeacherService
type TeacherService interface {
	SaveTeacher(teacher *entity.Teacher) *entity.Teacher
	TeachersCount() uint
	GetTeacherByID(ID uint) *entity.Teacher
	LogTeacher(username, password string) *entity.Teacher
	TeachersOfBranchByID(BranchID uint) *[]entity.Teacher
	GetActiveLectures(BranchID, TeacherID uint) *entity.ActiveLecturesListDataStructure
	ChangeImageURL(ID uint, value string) bool
}
