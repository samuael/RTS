package PaymentService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Payment"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// GormPaymentService struct
type GormPaymentService struct {
	PaymentRepo Payment.PaymentRepo
}

// NewGormPaymentService function
func NewGormPaymentService(paymentrepo Payment.PaymentRepo) Payment.PaymentService {
	return &GormPaymentService{
		PaymentRepo: paymentrepo,
	}
}

// CreatePayment method
func (paymentser *GormPaymentService) CreatePayment(payment *entity.Payment) *entity.Payment {
	payment, erra := paymentser.PaymentRepo.CreatePayment(payment)
	if erra != nil {
		return nil
	}
	return payment
}

// PaymentsOfRound (Roundid, BranchNumber uint) *[]entity.Payment
func (paymentser *GormPaymentService) PaymentsOfRound(Roundid, BranchNumber uint) *[]entity.Payment {
	payments, era := paymentser.PaymentRepo.PaymentsOfRound(Roundid, BranchNumber)
	if era != nil {
		return nil
	}
	return payments
}

// PaymentsOfStudent (StudentID uint) *[]entity.Payment
func (paymentser *GormPaymentService) PaymentsOfStudent(StudentID uint) *[]entity.Payment {
	payments, erra := paymentser.PaymentRepo.PaymentsOfStudent(StudentID)
	if erra != nil {
		return nil
	}
	return payments
}

// PaymentsOfSecretary (SecretaryID, limit uint) *[]entity.Payment
func (paymentser *GormPaymentService) PaymentsOfSecretary(SecretaryID, limit, offset uint) *[]entity.Payment {
	payments, newErrro := paymentser.PaymentRepo.PaymentsOfSecretary(SecretaryID, limit, offset)
	if newErrro != nil {
		return nil
	}
	return payments
}

// StudentPaidAmount (StudentID uint) float64
func (paymentser *GormPaymentService) StudentPaidAmount(StudentID uint) float64 {
	amount, errra := paymentser.PaymentRepo.StudentPaidAmount(StudentID)
	if errra != nil {
		return amount
	}
	return amount
}

// RoundPaidAmount (Roundid uint) float64
func (paymentser *GormPaymentService) RoundPaidAmount(Roundid uint) float64 {
	amount, era := paymentser.PaymentRepo.RoundPaidAmount(Roundid)
	if era != nil {
		return amount
	}
	return amount
}

// GetSinglePaymentReciptData method generatinf payment and Related Datas of the Payment ID
func (paymentser *GormPaymentService) GetSinglePaymentReciptData(PaymentID uint) *entity.SinglePaymentDataStructure {
	singlePaymentdatastructure := &entity.SinglePaymentDataStructure{}
	payment, era := paymentser.PaymentRepo.GetPaymentByID(PaymentID)
	if era != nil || payment == nil {
		singlePaymentdatastructure.Message = "No Payment Found For This Request "
		singlePaymentdatastructure.Success = false
		return singlePaymentdatastructure
	}
	singlePaymentdatastructure.Success = true
	singlePaymentdatastructure.Payment = *payment
	singlePaymentdatastructure.Message = " Succesful"
	return singlePaymentdatastructure
}

// GetBranchPaymentReport  (Branch uint, date *etc.Date) *[]entity.Payment
func (paymentser *GormPaymentService) GetBranchPaymentReport(BranchID uint, date *etc.Date) *[]entity.Payment {
	payments, era := paymentser.PaymentRepo.GetBranchPaymentReport(BranchID, date)
	if era != nil {
		return nil
	}
	return payments
}

// DPR func
func (paymentser *GormPaymentService) DPR(date etc.Date, BranchID int) {
	paymentser.PaymentRepo.GetDailyAdminsPaymentReport(date, uint(BranchID))
}

// GetDailyPaymentReport (date etc.Date, BranchID uint) *entity.DailyPaymentResult
func (paymentser *GormPaymentService) GetDailyPaymentReport(date etc.Date, BranchID uint) *entity.DailyPaymentResult {

	return nil
}
