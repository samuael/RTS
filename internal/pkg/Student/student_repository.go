package Student

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// StudentRepository representing Student Repository
type StudentRepository interface {
	SaveStudent(student *entity.Student) *entity.Student
	SaveStudents(RoundID int, students *[]entity.Student) error
	GetStudentByID(ID uint) (*entity.Student, error)
	LogStudent(username, password string) (*entity.Student, error)
	GetStudentPaidMoreThanPaymentLimit(RoudID uint, PaymentLimit float64) (*[]entity.Student, error)
	GetStudentsOfCategory(BranchID, CategoryID, Offset, Limit uint, Active bool) (*[]entity.Student, error)
	ChangeImageURL(ID uint, val string) error
	StudentsOfRound(RoundID uint) (*[]entity.Student, error)
	StudentsOfRoundWithPayment(round uint, amount float64) (*[]entity.Student, error)
}
