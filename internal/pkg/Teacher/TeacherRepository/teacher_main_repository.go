//Package TeacherRepo
package TeacherRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/jinzhu/gorm"
)

// TeacherRepo struct
type TeacherRepo struct {
	DB *gorm.DB
}

// NewTeacherRepo function
func NewTeacherRepo(db *gorm.DB) Teacher.TeacherRepo {
	return &TeacherRepo{
		DB: db,
	}
}

// SaveTeacher method
func (tr *TeacherRepo) SaveTeacher(teacher *entity.Teacher) error {
	saveError := tr.DB.Save(teacher).Error
	return saveError
}

// TeachersCount method
func (tr *TeacherRepo) TeachersCount(count *uint) error {
	erro := tr.DB.Model(&entity.Teacher{}).Count(count).Error
	return erro
}

// GetTeacherByID method
func (tr *TeacherRepo) GetTeacherByID(ID uint) (*entity.Teacher, error) {
	teacher := &entity.Teacher{}
	newError := tr.DB.Where("id=?", ID).First(teacher).Error
	return teacher, newError
}

// LogTeacher method
func (tr *TeacherRepo) LogTeacher(username, password string) (*entity.Teacher, error) {
	teacher := &entity.Teacher{}
	teacher.Username = username
	teacher.Password = password
	theError := tr.DB.First(teacher, "username=? and password=?", username, password).Error
	return teacher, theError
}

// TeachersOfBranchByID method to find the Teachers of specified Branch
func (tr *TeacherRepo) TeachersOfBranchByID(BranchID uint) (*[]entity.Teacher, error) {
	teachers := &[]entity.Teacher{}
	newError := tr.DB.Where("branch_number=?", BranchID).Find(teachers).Error
	return teachers, newError
}

// GetActiveLectures BranchID, TeacherID uint) (*entity.ActiveLecturesListDataStructure, error)
func (tr *TeacherRepo) GetActiveLectures(BranchID, TeacherID uint) (*[]entity.Lecture, error) {
	lectures := &[]entity.Lecture{}
	newErra := tr.
		DB.
		Table("lectures").
		Joins("left join rounds on rounds.id=lectures.roundid").
		Where("branch_refer=? and teacher_refer=?", BranchID, TeacherID).
		Find(lectures).
		Error
	for l := 0; l < len(*lectures); l++ {
		lecture := &(*lectures)[l]
		newErra = tr.DB.Model(lecture).Related(&lecture.Course, "Course").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.Teacher, "TeacherRefer").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.StartDate, "StartDate").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.EndDate, "EndDate").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.Round, "Roundid").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.Round.Category, "Category").Error
	}
	// Populating the Datas
	return lectures, newErra
}

// GetTodaysLecture method
// Fetching from lectures Criteria date Has to Be Simmilar with today
func (tr *TeacherRepo) GetTodaysLecture(BranchID, TeacherID uint) (*[]entity.Lecture, error) {
	lectures := &[]entity.Lecture{}
	today := etc.NewDate(0)
	newErra := tr.
		DB.
		Table("lectures").
		Joins("left join dates on dates.id=lectures.start_date_refer").
		Where("branch_refer=? and teacher_refer=? && day=? and month= ? and year=? and actuve=?",
			BranchID,
			TeacherID,
			today.Day,
			today.Month,
			today.Year,
			true).
		Find(lectures).
		Error
	for l := 0; l < len(*lectures); l++ {
		lecture := &(*lectures)[l]
		newErra = tr.DB.Model(lecture).Related(&lecture.Course, "Course").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.Teacher, "TeacherRefer").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.StartDate, "StartDate").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.EndDate, "EndDate").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.Round, "Roundid").Error
		newErra = tr.DB.Model(lecture).Related(&lecture.Round.Category, "Category").Error
	}
	return lectures, newErra
}

// ChangeImageURL  method for saving the new Image Url
func (tr *TeacherRepo) ChangeImageURL(ID uint, ImageURL string) error {
	era := tr.DB.Table("teachers").Where("id=?", ID).Updates(map[string]string{"imageurl": ImageURL}).Error
	return era
}
