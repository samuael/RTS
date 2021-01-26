package StudentRepository

import (
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/jinzhu/gorm"
)

// StudentRepository in repository package of Student package
type StudentRepository struct {
	DB *gorm.DB
}

// NewStudentRepository y
func NewStudentRepository(db *gorm.DB) Student.StudentRepository {
	return &StudentRepository{
		DB: db,
	}
}

// SaveStudent ()
func (srepo *StudentRepository) SaveStudent(student *entity.Student) *entity.Student {
	erro := srepo.DB.Save(student).Error
	if erro != nil {
		return nil
	}
	return student
}

// SaveStudents ()
func (srepo *StudentRepository) SaveStudents(RoundID int, students *[]entity.Student) error {
	round := &entity.Round{}
	erro := srepo.DB.Table("rounds").Where("id=?", RoundID).Find(round).Error
	if erro != nil {
		return erro
	}
	round.Students = *students
	erro = srepo.DB.Table("rounds").Save(round).Error
	return erro
}

// GetStudentByID methid
func (srepo *StudentRepository) GetStudentByID(ID uint) (*entity.Student, error) {
	reciver := &entity.Student{}
	erro := srepo.DB.First(reciver, ID).Error
	if erro != nil || reciver.ID == 0 {
		return reciver, erro
	}
	reciver.BirthDate = &etc.Date{}
	reciver.Address = &entity.Address{}
	reciver.GuarantorAddress = &entity.Address{}
	reciver.Section = &entity.Section{}
	reciver.Round = &entity.Round{}
	srepo.DB.Model(reciver).Related(reciver.Address, "AddressRefer")
	srepo.DB.Model(reciver).Related(reciver.GuarantorAddress, "GuarantorAddressRefer")
	srepo.DB.Model(reciver).Related(reciver.BirthDate, "OwnerID")
	srepo.DB.Model(reciver).Related(&(reciver.Category), "CategoryID")
	srepo.DB.Model(reciver).Related(reciver.Round, "RoundRefer")
	srepo.DB.Model(reciver).Related(reciver.Section, "SectionRefer")
	if erro != nil {
		return nil, erro
	}
	return reciver, erro
}

// LogStudent method to signin a student
func (srepo *StudentRepository) LogStudent(username, password string) (*entity.Student, error) {
	student := &entity.Student{}
	student.Username = username
	student.Password = password
	newError := srepo.DB.First(student, "username=? and password=?", username, password).Error
	return student, newError
}

// GetStudentPaidMoreThanPaymentLimit method
func (srepo *StudentRepository) GetStudentPaidMoreThanPaymentLimit(RoundID uint, PaymentLimit float64) (*[]entity.Student, error) {
	students := &[]entity.Student{}
	newEra := srepo.DB.Table("students").Where("round_refer=? and paid_amount >=? ", RoundID, PaymentLimit).Find(students).Error
	for k := 0; k < len(*students); k++ {
		reciver := &(*students)[k]
		reciver.BirthDate = &etc.Date{}
		srepo.DB.Model(reciver).Related(&(reciver.Address), "AddressRefer")
		srepo.DB.Model(reciver).Related(&(reciver.GuarantorAddress), "GuarantorAddressRefer")
		srepo.DB.Model(reciver).Related((reciver.BirthDate), "OwnerID")
		srepo.DB.Model(reciver).Related(&(reciver.Category), "Category")
		srepo.DB.Model(reciver).Related(&(reciver.Round), "RoundRefer")
		srepo.DB.Model(reciver).Related(&(reciver.Section), "SectionRefer")
	}
	return students, newEra
}

// GetStudentsOfCategory method
func (srepo *StudentRepository) GetStudentsOfCategory(BranchID, CategoryID, Offset, Limit uint, Active bool) (*[]entity.Student, error) {
	students := &[]entity.Student{}
	newEra := srepo.DB.Table("students").Offset(Offset).Limit(Limit).Find(students, "category_id=? and branch_id=? and active=?", CategoryID, BranchID, Active).Error
	return students, newEra
}

// ChangeImageURL   (val string) error
func (srepo *StudentRepository) ChangeImageURL(ID uint, val string) error {
	era := srepo.DB.Table("students").Where("id=?", ID).Updates(map[string]string{"imageurl": val}).Error
	return era
}

// StudentsOfRound  method returning students of a round
func (srepo *StudentRepository) StudentsOfRound(RoundID uint) (*[]entity.Student, error) {
	students := &[]entity.Student{}
	erra := srepo.DB.Table("students").Find(students, "round_refer=?", RoundID).Error
	for k := 0; k < len(*students); k++ {
		student := &(*students)[0]
		student.Address = &entity.Address{}
		student.BirthAddress = &entity.Address{}
		student.GuarantorAddress = &entity.Address{}
		student.BirthDate = &etc.Date{}
		erra = srepo.DB.Model(student).Related(student.Address, "Address").Error
		erra = srepo.DB.Model(student).Related(student.BirthAddress, "BirthAddress").Error
		erra = srepo.DB.Model(student).Related(student.GuarantorAddress, "GuarantorAddress").Error
		erra = srepo.DB.Model(student).Related(student.BirthDate, "BirthDate").Error
	}
	return students, erra
}

// StudentsOfRoundWithPayment (round uint, amount float64) (*[]entity.Student, error)
func (srepo *StudentRepository) StudentsOfRoundWithPayment(round uint, amount float64) (*[]entity.Student, error) {
	students := &[]entity.Student{}
	erra := srepo.DB.Table("students").Find(students, "round_refer=? and paid_amount >=?", round, amount).Error
	for k := 0; k < len(*students); k++ {
		student := &(*students)[0]
		student.Address = &entity.Address{}
		student.BirthAddress = &entity.Address{}
		student.GuarantorAddress = &entity.Address{}
		student.BirthDate = &etc.Date{}
		// student.Category = &entity.Category{}
		srepo.DB.Model(student).Related(student.Address, "Address")
		srepo.DB.Model(student).Related(student.BirthAddress, "BirthAddress")
		srepo.DB.Model(student).Related(student.GuarantorAddress, "GuarantorAddress")
		srepo.DB.Model(student).Related(student.BirthDate, "BirthDate")
		srepo.DB.Model(student).Related(&student.Category, "Category")
	}
	fmt.Println(erra)
	return students, erra
}
