package StudentService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// StudentService y
type StudentService struct {
	StudentRepo Student.StudentRepository
}

// GetStudent returning Student having username and Password
func (studentser *StudentService) GetStudent(student *entity.Student) *entity.Student {
	return nil
}

// NewStudentService y
func NewStudentService(studentrepo Student.StudentRepository) Student.StudentService {
	return &StudentService{
		StudentRepo: studentrepo,
	}
}

// GetStudentByID method
func (studentser *StudentService) GetStudentByID(ID uint) *entity.Student {
	student, erro := studentser.StudentRepo.GetStudentByID(ID)
	if erro != nil {
		return nil
	}
	return student
}

// SaveStudent method to save the student
func (studentser *StudentService) SaveStudent(student *entity.Student) *entity.Student {
	student = studentser.StudentRepo.SaveStudent(student)
	return student
}

// SaveStudents method to save the student
func (studentser *StudentService) SaveStudents(RoundID int, students *[]entity.Student) bool {
	erra := studentser.StudentRepo.SaveStudents(RoundID, students)
	if erra != nil {
		return false
	}
	return true
}

// LogStudent method
func (studentser *StudentService) LogStudent(username, password string) *entity.Student {
	student, newError := studentser.StudentRepo.LogStudent(username, password)
	if newError != nil {
		return nil
	}
	return student
}

// GetStudentPaidMoreThanPaymentLimit (RoudID uint, PaymentLimit float64) *[]entity.Student
func (studentser *StudentService) GetStudentPaidMoreThanPaymentLimit(RoundID uint, PaymentLimit float64) *[]entity.Student {
	students, era := studentser.StudentRepo.GetStudentPaidMoreThanPaymentLimit(RoundID, PaymentLimit)
	if era != nil {
		return nil
	}
	return students
}

// GetStudentsOfCategory   (BranchID, CategoryID, Offset, Limit uint) *[]entity.Student
func (studentser *StudentService) GetStudentsOfCategory(BranchID, CategoryID, Offset, Limit uint, Active bool) *[]entity.Student {
	students, era := studentser.StudentRepo.GetStudentsOfCategory(BranchID, CategoryID, Offset, Limit, Active)
	if era != nil {
		return nil
	}
	return students
}

// ChangeImageURL   (val string) bool this is the service for chane imaege of the students profule
func (studentser *StudentService) ChangeImageURL(ID uint, val string) bool {
	era := studentser.StudentRepo.ChangeImageURL(ID, val)
	if era != nil {
		return false
	}
	return true
}

// StudentsOfRound (RoundID uint) *[]entity.Student
func (studentser *StudentService) StudentsOfRound(RoundID uint) *[]entity.Student {
	students, ere := studentser.StudentRepo.StudentsOfRound(RoundID)
	if ere != nil {
		return nil
	}
	return students
}

// StudentsOfRoundWithPayment (round uint, amount float64) *[]entity.Student
func (studentser *StudentService) StudentsOfRoundWithPayment(round uint, amount float64) *[]entity.Student {
	students, era := studentser.StudentRepo.StudentsOfRoundWithPayment(round, amount)
	if era != nil {
		return nil
	}
	return students
}
