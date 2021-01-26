package Lecture

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// LectureRepo
type LectureRepo interface {
	GetLectureByID(LectureID uint) (*entity.Lecture, error)
	SaveLecture(lecture *entity.Lecture) (*entity.Lecture, error)
	PopulateLectures(lectures *[]entity.Lecture) (*[]entity.Lecture, error)
	GetAllLecturesOfTeacherID(TeacherID uint) (*[]entity.Lecture, error)
	GetActiveLecturesOfTeacherID(TeacherID uint) (*[]entity.Lecture, error)
	GetActiveLecturesOfTeacherIDAndRoundID(TeacherID, RoundID uint) (*[]entity.Lecture, error)
	GetActiveLecturesOfRoundID(RoundID uint) (*[]entity.Lecture, error)
	PassLecture(LectureID uint) error
}
