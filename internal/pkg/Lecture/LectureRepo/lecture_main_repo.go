package LectureRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Lecture"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// LectureRepo struct
type LectureRepo struct {
	DB *gorm.DB
}

// NewLectureRepo function
func NewLectureRepo(db *gorm.DB) Lecture.LectureRepo {
	return &LectureRepo{
		DB: db,
	}
}

// PopulateLectures method
func (lecturerepo *LectureRepo) PopulateLectures(lectures *[]entity.Lecture) (*[]entity.Lecture, error) {
	var newError error
	for i := 0; i < len(*lectures); i++ {
		lecture := &(*lectures)[i]
		newError = lecturerepo.DB.Model(lecture).Related(&(lecture.Course), "CourseRefer").Error
		newError = lecturerepo.DB.Model(lecture).Related(&lecture.StartDate, "StartDateRefer").Error
		newError = lecturerepo.DB.Model(lecture).Related(&lecture.EndDate, "EndDateRefer").Error
	}
	return lectures, newError
}

// GetAllLecturesOfTeacherID mthod
func (lecturerepo *LectureRepo) GetAllLecturesOfTeacherID(TeacherID uint) (*[]entity.Lecture, error) {
	lectures := &[]entity.Lecture{}
	newError := lecturerepo.DB.Where("teacher_refer=?", TeacherID).Find(lectures).Error
	// Populating the Start Datetime and the Ending Date Time of a lecture
	lectures, newError = lecturerepo.PopulateLectures(lectures)
	return lectures, newError
}

// GetActiveLecturesOfTeacherID mthod
func (lecturerepo *LectureRepo) GetActiveLecturesOfTeacherID(TeacherID uint) (*[]entity.Lecture, error) {
	lectures := &[]entity.Lecture{}
	newError := lecturerepo.DB.Where("teacher_refer=? and passed=?", TeacherID, false).Find(lectures).Error
	// Populating the Start Datetime and the Ending Date Time of a lecture
	lectures, newError = lecturerepo.PopulateLectures(lectures)
	return lectures, newError
}

// GetActiveLecturesOfTeacherIDAndRoundID method
func (lecturerepo *LectureRepo) GetActiveLecturesOfTeacherIDAndRoundID(TeacherID, RoundID uint) (*[]entity.Lecture, error) {
	lectures := &[]entity.Lecture{}
	newError := lecturerepo.DB.Where("teacher_refer=? and ", TeacherID).Find(lectures).Error
	// Populating the Start Datetime and the Ending Date Time of a lecture
	lectures, newError = lecturerepo.PopulateLectures(lectures)
	return lectures, newError
}

// GetActiveLecturesOfRoundID method
func (lecturerepo *LectureRepo) GetActiveLecturesOfRoundID(RoundID uint) (*[]entity.Lecture, error) {
	lectures := &[]entity.Lecture{}
	newError := lecturerepo.DB.Where("roundid=? and ", RoundID).Find(lectures).Error
	// Populating the Start Datetime and the Ending Date Time of a lecture
	lectures, newError = lecturerepo.PopulateLectures(lectures)
	return lectures, newError
}

// PassLecture method
func (lecturerepo *LectureRepo) PassLecture(LectureID uint) error {
	newError := lecturerepo.DB.Table("lectures").Where("id=?", LectureID).Updates("passed=?", true).Error
	return newError
}

// GetLectureByID method
func (lecturerepo *LectureRepo) GetLectureByID(LectureID uint) (*entity.Lecture, error) {
	lecture := &entity.Lecture{}
	newError := lecturerepo.DB.Table("lectures").Where("id=?", LectureID).First(lecture).Error
	return lecture, newError
}

// SaveLecture method
func (lecturerepo *LectureRepo) SaveLecture(lecture *entity.Lecture) (*entity.Lecture, error) {
	newEra := lecturerepo.DB.Save(lecture).Error
	return lecture, newEra
}
