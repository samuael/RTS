package Question

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// QuestionService inteerface
type QuestionService interface {
	CreateQuestion(question *entity.Question) *entity.Question
	DeleteQuestion(ID uint) bool
	GetQuestions(StudentID, Offset, Limit uint) *[]entity.Question
	GetGradeResult(StudentID uint) (uint, uint, bool)
	GetAnswer(QuestionID uint) int
	ResetResult(StudentID uint) bool
	GetQuestionByID(ID uint) *entity.Question
	GetStudentsAskedQuestionByID(StudentID uint) *entity.AskedQuetion
	SaveAskedQuestion(askedQuestion *entity.AskedQuetion) *entity.AskedQuetion
}
