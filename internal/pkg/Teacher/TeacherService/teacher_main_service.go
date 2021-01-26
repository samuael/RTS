// Package TeacherService package
package TeacherService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// TeacherService struct
type TeacherService struct {
	TeacherRepo Teacher.TeacherRepo
}

// NewTeacherService function
func NewTeacherService(teacha Teacher.TeacherRepo) Teacher.TeacherService {
	return &TeacherService{
		TeacherRepo: teacha,
	}
}

// SaveTeacher method
func (teacherser *TeacherService) SaveTeacher(teacher *entity.Teacher) *entity.Teacher {
	errors := teacherser.TeacherRepo.SaveTeacher(teacher)
	if errors != nil {
		return nil
	}
	return teacher
}

// TeachersCount method
func (teacherser *TeacherService) TeachersCount() uint {
	var count uint
	count = 0
	erro := teacherser.TeacherRepo.TeachersCount(&count)
	if erro != nil {
		return 0
	}
	return count
}

// GetTeacherByID method returning a Teacher struct pointer
func (teacherser *TeacherService) GetTeacherByID(ID uint) *entity.Teacher {
	teacher, errors := teacherser.TeacherRepo.GetTeacherByID(ID)
	if errors != nil {
		return nil
	}
	return teacher

}

// LogTeacher method
func (teacherser *TeacherService) LogTeacher(username, password string) *entity.Teacher {
	teacher, erra := teacherser.TeacherRepo.LogTeacher(username, password)
	if erra != nil {
		return nil
	}
	return teacher
}

// TeachersOfBranchByID (BranchID uint) *[]entity.Teacher returning Slice of Teachers Pointer of a branch
func (teacherser *TeacherService) TeachersOfBranchByID(BranchID uint) *[]entity.Teacher {
	teachers, erra := teacherser.TeacherRepo.TeachersOfBranchByID(BranchID)
	if erra != nil {
		return nil
	}
	return teachers
}

// GetActiveLectures  (BranchID, TeacherID uint) *entity.ActiveLecturesListDataStructure
func (teacherser *TeacherService) GetActiveLectures(BranchID, TeacherID uint) *entity.ActiveLecturesListDataStructure {
	activeLecturesListds := &entity.ActiveLecturesListDataStructure{
		Success: false,
	}
	lectures, era := teacherser.TeacherRepo.GetActiveLectures(BranchID, TeacherID)
	if era != nil || lectures == nil {
		activeLecturesListds.Message = "No Active Lecture For The User "
		return activeLecturesListds
	}
	activeLecturesListds.Success = true
	activeLecturesListds.Message = "Succesful"
	activeLecturesListds.Lectures = *lectures
	return activeLecturesListds
}

// GetTodaysLectures  (BranchID, TeacherID uint) *entity.ActiveLecturesListDataStructure
func (teacherser *TeacherService) GetTodaysLectures(BranchID, TeacherID uint) *entity.ActiveLecturesListDataStructure {
	activeLecturesListds := &entity.ActiveLecturesListDataStructure{
		Success: false,
	}
	lectures, era := teacherser.TeacherRepo.GetTodaysLecture(BranchID, TeacherID)
	if era != nil || lectures == nil {
		activeLecturesListds.Message = "No Lecture For Today!"
		return activeLecturesListds
	}
	activeLecturesListds.Success = true
	activeLecturesListds.Message = "Lectures of Today "
	activeLecturesListds.Lectures = *lectures
	return activeLecturesListds
}

// ChangeImageURL  method in the service
func (teacherser *TeacherService) ChangeImageURL(ID uint, ImageURL string) bool {
	era := teacherser.TeacherRepo.ChangeImageURL(ID, ImageURL)
	if era != nil {
		return false
	}
	return true
}
