package Question

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type QuestionRepo interface {
	CreateQuestion(question *entity.Question) (*entity.Question, error)
	DeleteQuestion(ID uint) error
	GetQuestions(StudentID, Offset, Limit uint) (*[]entity.Question, error)
	GetAnswer(QuestionID uint) (uint, error)
	ResetResult(StudentID uint) error
	GetGradeResult(StudentID uint) (uint, uint, error)
	GetQuestionByID(ID uint) (*entity.Question, error)
	SaveAskedQuestion(askedQuestion *entity.AskedQuetion) (*entity.AskedQuetion, error)
	GetStudentsAskedQuestionByID(StudentID uint) (*entity.AskedQuetion, error)
}
