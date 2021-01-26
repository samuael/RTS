package PaymentRepo

import (
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Payment"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/jinzhu/gorm"
)

// GormPaymentRepo struct
type GormPaymentRepo struct {
	DB *gorm.DB
}

// NewGormPaymentRepo method
func NewGormPaymentRepo(db *gorm.DB) Payment.PaymentRepo {
	return &GormPaymentRepo{
		DB: db,
	}
}

// CreatePayment method returning new Instance of Payment
func (paymentrepo *GormPaymentRepo) CreatePayment(payment *entity.Payment) (*entity.Payment, error) {
	newErroe := paymentrepo.DB.Create(payment).Error
	return payment, newErroe
}

// PaymentsOfRound method returning the Payments of the Round Uding the ID Using RoundID  , RoundNumber  , and BranchNumber uint
func (paymentrepo *GormPaymentRepo) PaymentsOfRound(Roundid, BranchNumber uint) (*[]entity.Payment, error) {
	payments := &[]entity.Payment{}
	newErra := paymentrepo.DB.Model(&entity.Payment{}).Where("round_refer=? ", Roundid).Find(payments).Error
	for i := 0; i < len(*payments); i++ {
		paymentrepo.DB.Model((*payments)[i]).Related(&(*payments)[i].Round, "RoundRefer")
		// paymentrepo.DB.Model((*payments)[i]).Relat&ed((*payments)[i].Branch, "BranchRefer")
		paymentrepo.DB.Model((*payments)[i]).Related(&(*payments)[i].Admin, "AdminRefer")
		paymentrepo.DB.Model((*payments)[i]).Related(&(*payments)[i].Student, "StudentRefer")
	}
	return payments, newErra
}

// PaymentsOfStudent method
func (paymentrepo *GormPaymentRepo) PaymentsOfStudent(StudentID uint) (*[]entity.Payment, error) {
	payments := &[]entity.Payment{}
	newErroa := paymentrepo.DB.Where("student_refer=?", StudentID).Find(payments).Error
	for i := 0; i < len(*payments); i++ {
		paymentrepo.DB.Model((*payments)[i]).Related((*payments)[i].Round, "RoundRefer")
		// paymentrepo.DB.Model((*payments)[i]).Related((*payments)[i].Branch, "BranchRefer")
		// paymentrepo.DB.Model((*payments)[i]).Related((*payments)[i].Admin, "AdminRefer")
	}
	return payments, newErroa
}

// PaymentsOfSecretary Limited method returning limited amount of payments
func (paymentrepo *GormPaymentRepo) PaymentsOfSecretary(SecretaryID, limit uint, offser uint) (*[]entity.Payment, error) {
	payments := &[]entity.Payment{}
	errors := paymentrepo.DB.Where("admin_refer=?", SecretaryID).Limit(limit).Offset(offser).Find(payments).Error
	for i := 0; i < len(*payments); i++ {
		paymentrepo.DB.Model((*payments)[i]).Related(&(*payments)[i].Round, "RoundRefer")
		paymentrepo.DB.Model((*payments)[i]).Related(&(*payments)[i].Student, "StudentRefer")
	}
	return payments, errors
}

// StudentPaidAmount method
func (paymentrepo *GormPaymentRepo) StudentPaidAmount(StudentID uint) (float64, error) {
	var newAmount float64
	errors := paymentrepo.DB.Where("student_refer=?", StudentID).Select("sum(amount) as newAmount").Scan(&newAmount).Error
	return newAmount, errors
}

// AmountTake struct
type AmountTake struct {
	N float64
}

// RoundPaidAmount method
func (paymentrepo *GormPaymentRepo) RoundPaidAmount(Roundid uint) (float64, error) {
	amountTake := &AmountTake{}
	fmt.Println(Roundid, " ----This --- ")
	errorrs := paymentrepo.DB.Where("round_refer=?", Roundid).Select("sum(amount) as n ").Scan(amountTake).Error
	return amountTake.N, errorrs
}

// GetPaymentByID (PaymentID uint) (*entity.Payment, error)
func (paymentrepo *GormPaymentRepo) GetPaymentByID(PaymentID uint) (*entity.Payment, error) {
	payment := &entity.Payment{}
	newErra := paymentrepo.DB.Find(payment, PaymentID).Error
	payment.Round = &entity.Round{}
	payment.Date = &etc.Date{}
	paymentrepo.DB.Model(payment).Related(&payment.Round, "RoundRefer")
	paymentrepo.DB.Model(payment.Round).Related(&payment.Round.Category, "CategoryRefer")
	paymentrepo.DB.Model(payment).Related(&payment.Admin, "AdminRefer")
	paymentrepo.DB.Model(payment).Related(&payment.Student, "StudentRefer")
	paymentrepo.DB.Model(payment).Related(&payment.Date, "DateRefer")
	return payment, newErra
}

// GetBranchPaymentReport (Branch uint, date *etc.Date) (*[]entity.Payment, error)
func (paymentrepo *GormPaymentRepo) GetBranchPaymentReport(BranchID uint, date *etc.Date) (*[]entity.Payment, error) {
	payments := &[]entity.Payment{}
	newErro := paymentrepo.
		DB.
		Table("payments").
		Select("branch_refer , student_refer , round_refer , admin_refer , date_refer , amount ").
		Joins("left join dates on dates.id=payments.date_refer").
		Where("branch_refer=? && day=? && month=? year=?", BranchID, date.Day, date.Month, date.Year).
		Find(payments).Error
	for l := 0; l < len(*payments); l++ {
		payment := (*payments)[l]
		payment.Round = &entity.Round{}
		payment.Date = &etc.Date{}
		paymentrepo.DB.Model(payment).Related(&payment.Round, "RoundRefer")
		paymentrepo.DB.Model(payment.Round).Related(&payment.Round.Category, "CategoryRefer")
		paymentrepo.DB.Model(payment).Related(&payment.Admin, "AdminRefer")
		paymentrepo.DB.Model(payment).Related(&payment.Student, "StudentRefer")
		paymentrepo.DB.Model(payment).Related(&payment.Date, "DateRefer")
	}
	return payments, newErro
}

// // PaymentsOfEachSecretariesForMonth  returning the Payment sum where the secretory recieved in the Month
// func (paymentrepo *GormPaymentRepo) PaymentsOfEachSecretariesForMonth(BranchID, Year, Month uint) (map[*entity.Admin]float64, error) {
// 	themap := map[*entity.Admin]float64{}
// 	// First get the Secretories of the Branch
// 	secretories := &[]entity.Admin{}
// 	newEra := paymentrepo.DB.Find(secretories, "branch_refer=? and role=?", BranchID, entity.SECRETART).Error
// 	if newEra != nil {
// 		return themap, newEra
// 	}
// 	for i := 0; i < len(*secretories); i++ {
// 		admin := &(*secretories)[i]
// 		// paymentrepo.DB.Model(admin).
// 	}
// }

// GetDailyAdminsPaymentReport  method
func (paymentrepo *GormPaymentRepo) GetDailyAdminsPaymentReport(date etc.Date, BranchID uint) (map[*entity.Admin]float64, error) {
	// payments := &[]entity.Payment{}
	// payments := &[]entity.Payment{}
	rows, newError:= paymentrepo.DB.
		Table("payments").
		Select(" dates.day , dates.Year ,sum(payments.amount) as total group by dates.day").
		Joins("INNER JOIN dates ON payments.date_refer=dates.id").
		// Where("day=? and month=? and year=? and branch_refer=?", date.Day, date.Month, date.Year, BranchID).
		// Rows()
		Rows()
	fmt.Println(rows)
//	var day, year uint
//	var amount float64
	//fmt.Println(newError)
	//for newError == nil && rows.Next() {
	//	rows.Scan(&day, &year, &amount)
	//	fmt.Printf("Payment amount : %d : %d : %f  \n", day, year, amount)
	//}
	// fmt.Println("Result Printing Finished Payments Repo CLass")
	return nil, newError
}
