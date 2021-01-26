package QuestionService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Question"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// QuestionService struct
type QuestionService struct {
	Questionrepo Question.QuestionRepo
}

// NewQuestionService function
func NewQuestionService(Questionrepo Question.QuestionRepo) Question.QuestionService {
	return &QuestionService{
		Questionrepo: Questionrepo,
	}
}

// CreateQuestion (question *entity.Question) *entity.Question
func (quesser *QuestionService) CreateQuestion(question *entity.Question) *entity.Question {
	question, erra := quesser.Questionrepo.CreateQuestion(question)
	if erra != nil {
		return nil
	}
	return question
}

// DeleteQuestion (ID uint) bool
func (quesser *QuestionService) DeleteQuestion(ID uint) bool {
	erra := quesser.Questionrepo.DeleteQuestion(ID)
	if erra != nil {
		return false
	}
	return true
}

// GetQuestions (StudentID, Limit uint) []entity.Question
func (quesser *QuestionService) GetQuestions(StudentID, Offset, Limit uint) *[]entity.Question {
	questions, erra := quesser.Questionrepo.GetQuestions(StudentID, Offset, Limit)
	if erra != nil {
		return nil
	}
	return questions
}

// GetAnswer (QuestionID uint) uint
func (quesser *QuestionService) GetAnswer(QuestionID uint) int {
	value, Erra := quesser.Questionrepo.GetAnswer(QuestionID)
	if Erra != nil {
		return -1
	}
	return int(value)
}

// ResetResult (StudentID uint) bool
func (quesser *QuestionService) ResetResult(StudentID uint) bool {
	erra := quesser.Questionrepo.ResetResult(StudentID)
	if erra != nil {
		return false
	}
	return true
}

// GetGradeResult (StudentID uint) (uint, uint)  the bool to tell the result fetching is succesful or not
func (quesser *QuestionService) GetGradeResult(StudentID uint) (uint, uint, bool) {
	askedquestionCount, answeredCount, erra := quesser.Questionrepo.GetGradeResult(StudentID)
	if erra != nil {
		return 0, 0, false
	}
	return askedquestionCount, answeredCount, true
}

// GetQuestionByID function
func (quesser *QuestionService) GetQuestionByID(ID uint) *entity.Question {
	question, erra := quesser.Questionrepo.GetQuestionByID(ID)
	if erra != nil {
		return nil
	}
	return question
}

// GetStudentsAskedQuestionByID (StudentID uint) *entity.AskedQuetion
func (quesser *QuestionService) GetStudentsAskedQuestionByID(StudentID uint) *entity.AskedQuetion {
	asked, era := quesser.Questionrepo.GetStudentsAskedQuestionByID(StudentID)
	if era != nil {
		return nil
	}
	return asked
}

// SaveAskedQuestion (askedQuestion *entity.AskedQuetion) *entity.AskedQuetion
func (quesser *QuestionService) SaveAskedQuestion(askedQuestion *entity.AskedQuetion) *entity.AskedQuetion {
	askedQuestion, era := quesser.Questionrepo.SaveAskedQuestion(askedQuestion)
	if era != nil {
		return nil
	}
	return askedQuestion
}
