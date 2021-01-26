package Lecture

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// LectureService struct
type LectureService interface {
	GetLectureByID(LectureID uint) *entity.Lecture
	SaveLecture(lecture *entity.Lecture) *entity.Lecture
	PopulateLectures(lectures *[]entity.Lecture) *[]entity.Lecture
	GetAllLecturesOfTeacherID(TeacherID uint) *[]entity.Lecture
	GetActiveLecturesOfTeacherID(TeacherID uint) *[]entity.Lecture
	GetActiveLecturesOfTeacherIDAndRoundID(TeacherID, RoundID uint) *[]entity.Lecture
	GetActiveLecturesOfRoundID(RoundID uint) *[]entity.Lecture
	PassLecture(LectureID uint) bool
}
