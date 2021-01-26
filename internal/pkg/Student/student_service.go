//Package Student representing Student Struct
package Student

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

// StudentService representing repositoryService
type StudentService interface {
	GetStudent(student *entity.Student) *entity.Student
	GetStudentByID(ID uint) *entity.Student
	SaveStudent(student *entity.Student) *entity.Student
	SaveStudents(RoundID int, students *[]entity.Student) bool
	LogStudent(username, password string) *entity.Student
	GetStudentPaidMoreThanPaymentLimit(RoudID uint, PaymentLimit float64) *[]entity.Student
	GetStudentsOfCategory(BranchID, CategoryID, Offset, Limit uint, Active bool) *[]entity.Student
	ChangeImageURL(ID uint, val string) bool
	StudentsOfRound(RoundID uint) *[]entity.Student
	StudentsOfRoundWithPayment(round uint, amount float64) *[]entity.Student
}
