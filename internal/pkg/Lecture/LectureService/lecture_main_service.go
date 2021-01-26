package LectureService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Lecture"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// LectureService struct
type LectureService struct {
	LectureRepo Lecture.LectureRepo
}

// NewLectureService function
func NewLectureService(lecturerepo Lecture.LectureRepo) Lecture.LectureService {
	return &LectureService{
		LectureRepo: lecturerepo,
	}
}

// SaveLecture  (lecture *entity.Lecture) *entity.Lecture
func (lectureser *LectureService) SaveLecture(lecture *entity.Lecture) *entity.Lecture {
	lecture, era := lectureser.LectureRepo.SaveLecture(lecture)
	if era != nil {
		return nil
	}
	return lecture
}

// GetLectureByID (LectureID uint) *entity.Lecture
func (lectureser *LectureService) GetLectureByID(LectureID uint) *entity.Lecture {
	lecture, era := lectureser.LectureRepo.GetLectureByID(LectureID)
	if era != nil {
		return nil
	}
	return lecture
}

// PopulateLectures (lectures *[]entity.Lecture) *[]entity.Lecture
func (lectureser *LectureService) PopulateLectures(lectures *[]entity.Lecture) *[]entity.Lecture {
	lectures, era := lectureser.LectureRepo.PopulateLectures(lectures)
	if era != nil {
		return nil
	}
	return lectures
}

// GetAllLecturesOfTeacherID (TeacherID uint) *[]entity.Lecture
func (lectureser *LectureService) GetAllLecturesOfTeacherID(TeacherID uint) *[]entity.Lecture {
	lectures, era := lectureser.LectureRepo.GetAllLecturesOfTeacherID(TeacherID)
	if era != nil {
		return nil
	}
	return lectures
}

// GetActiveLecturesOfTeacherID (TeacherID uint) *[]entity.Lecture
// it may be any round
func (lectureser *LectureService) GetActiveLecturesOfTeacherID(TeacherID uint) *[]entity.Lecture {
	lectures, era := lectureser.LectureRepo.GetActiveLecturesOfTeacherID(TeacherID)
	if era != nil {
		return nil
	}
	return lectures
}

// GetActiveLecturesOfTeacherIDAndRoundID (TeacherID, RoundID uint) *[]entity.Lecture
// Because the teacher Mey have many rounds that he Teaches
func (lectureser *LectureService) GetActiveLecturesOfTeacherIDAndRoundID(TeacherID, RoundID uint) *[]entity.Lecture {
	lectures, erra := lectureser.LectureRepo.GetActiveLecturesOfTeacherIDAndRoundID(TeacherID, RoundID)
	if erra != nil {
		return nil
	}
	return lectures
}

// GetActiveLecturesOfRoundID (RoundID uint) *[]entity.Lecture
func (lectureser *LectureService) GetActiveLecturesOfRoundID(RoundID uint) *[]entity.Lecture {
	lectures, era := lectureser.LectureRepo.GetActiveLecturesOfRoundID(RoundID)
	if era != nil {
		return nil
	}
	return lectures
}

// PassLecture  to make the lecture as learned (LectureID uint) bool
func (lectureser *LectureService) PassLecture(LectureID uint) bool {
	era := lectureser.LectureRepo.PassLecture(LectureID)
	if era != nil {
		return false
	}
	return true
}
