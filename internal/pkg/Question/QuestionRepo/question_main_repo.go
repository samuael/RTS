package QuestionRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Question"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// QuestionRepo struct
type QuestionRepo struct {
	DB *gorm.DB
}

// NewQuestionRepo function
func NewQuestionRepo(db *gorm.DB) Question.QuestionRepo {
	return &QuestionRepo{
		DB: db,
	}
}

// CreateQuestion (question *entity.Question) (*entity.Question, error)
func (questionrepo *QuestionRepo) CreateQuestion(question *entity.Question) (*entity.Question, error) {
	newError := questionrepo.DB.Create(question).Error
	return question, newError
}

// SaveQuestion (question *entity.Question) (*entity.Question, error)
func (questionrepo *QuestionRepo) SaveQuestion(question *entity.Question) (*entity.Question, error) {
	newError := questionrepo.DB.Save(question).Error
	return question, newError
}

// DeleteQuestion (ID  uint) (*entity.Question, error)
func (questionrepo *QuestionRepo) DeleteQuestion(ID uint) error {
	newErra := questionrepo.DB.Delete(&entity.Question{}, ID).Error
	return newErra
}

// GetQuestions (StudentID uint) (*[]entity.Question, error)
func (questionrepo *QuestionRepo) GetQuestions(StudentID, Offset, Limit uint) (*[]entity.Question, error) {
	questions := &[]entity.Question{}
	askedquestion := &entity.AskedQuetion{}
	askedquestion.Studentid = StudentID
	newErra := questionrepo.DB.Table("asked_quetions").First(askedquestion, "studentid=?", StudentID).Error
	newErra = questionrepo.DB.Table("questions").Not([]int64(askedquestion.Questionsid)).Offset(Offset).Limit(Limit).Find(questions).Error //.Not(askedquestion.Questionsid)
	return questions, newErra
}

// GetStudentsAskedQuestionByID  method for Getting asked questions of a student by StudentsID ID
func (questionrepo *QuestionRepo) GetStudentsAskedQuestionByID(StudentID uint) (*entity.AskedQuetion, error) {
	askedQuestion := &entity.AskedQuetion{}
	era := questionrepo.DB.First(askedQuestion, "studentid=?", StudentID).Error
	return askedQuestion, era
}

// SaveAskedQuestion method for saving the Asked question of oa student
func (questionrepo *QuestionRepo) SaveAskedQuestion(askedQuestion *entity.AskedQuetion) (*entity.AskedQuetion, error) {
	erra := questionrepo.DB.Save(askedQuestion).Error
	return askedQuestion, erra
}

// GetAnswer (QuestionID, AnswerIndex uint) error
func (questionrepo *QuestionRepo) GetAnswer(QuestionID uint) (uint, error) {
	answerIndex := 0
	newErra := questionrepo.DB.Table("questions").Select("answerindex").Where("id=?", QuestionID).Row().Scan(&answerIndex)
	return uint(answerIndex), newErra
}

// GetGradeResult function  // returning Asked Questions Count and Answered Questions Count and error respectively
func (questionrepo *QuestionRepo) GetGradeResult(StudentID uint) (uint, uint, error) {
	var AskedQuestionCount, AnsweredQuestionCount uint
	newEra := questionrepo.DB.Table("students").Select("asked_questions_count , answered_question_count").Where("id=?", StudentID).Row().Scan(&AskedQuestionCount, &AnsweredQuestionCount)
	return AskedQuestionCount, AnsweredQuestionCount, newEra
}

// GetRank function to make a rank for the Students of the System to have grade it is Filterred by
// RoundID  of the System
// func (questionrepo *QuestionRepo) GetRank(RoundID, StudentID uint) (uint, error) {
// 	var Rank uint
// 	newError := questionrepo.DB.Table("students").Select("").Order("answered_question_count ASC").Where("round_refer=? and id=?", RoundID, StudentID).Row().Scan
// 	return Rank, newError
// }

// ResetResult (StudentID uint) error
func (questionrepo *QuestionRepo) ResetResult(StudentID uint) error {
	// Deleting asked Questions of student
	// Setting students AskedQuestionsCount AnsweredQuestionCount
	newErra := questionrepo.DB.Where("studentid=?", StudentID).Delete(&entity.AskedQuetion{}).Error
	newErra = questionrepo.DB.Table("students").Where("id=?", StudentID).Updates(map[string]uint{"answered_question_count": 0, "asked_questions_count": 0}).Error
	return newErra
}

// GetQuestionByID function
func (questionrepo *QuestionRepo) GetQuestionByID(ID uint) (*entity.Question, error) {
	question := &entity.Question{}
	era := questionrepo.DB.Table("questions").First(question, ID).Error
	return question, era
}
